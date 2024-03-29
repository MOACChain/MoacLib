// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package vm

import (
	"sync/atomic"

	"github.com/MOACChain/MoacLib/common"
	"github.com/MOACChain/MoacLib/common/math"
	"github.com/MOACChain/MoacLib/crypto"
	"github.com/MOACChain/MoacLib/log"
	"github.com/MOACChain/MoacLib/params"
)

// Config are the configuration options for the Interpreter
type Config struct {
	// Debug enabled debugging Interpreter options
	Debug bool
	// EnableJit enabled the JIT VM
	EnableJit bool
	// ForceJit forces the JIT VM
	ForceJit bool
	// Tracer is the op code logger
	Tracer Tracer
	// NoRecursion disabled Interpreter call, callcode,
	// delegate call and create.
	NoRecursion bool
	// Disable gas metering
	DisableGasMetering bool
	// Enable recording of SHA3/keccak preimages
	EnablePreimageRecording bool
	// JumpTable contains the EVM instruction table. This
	// may be left uninitialised and will be set to the default
	// table.
	JumpTable [256]operation
}

// Interpreter is used to run MoacNode based contracts and will utilise the
// passed evmironment to query external sources for state information.
// The Interpreter will run the byte code VM or JIT VM based on the passed
// configuration.
type Interpreter struct {
	evm      *EVM
	Cfg      Config
	gasTable params.GasTable

	readOnly   bool   // Whether to throw on stateful modifications
	returnData []byte // Last CALL's return data for subsequent reuse
}

// NewInterpreter returns a new instance of the Interpreter.
func NewInterpreter(evm *EVM, cfg Config) *Interpreter {
	// We use the STOP instruction whether to see
	// the jump table was initialised. If it was not
	// we'll set the default jump table.
	if !cfg.JumpTable[STOP].valid {
		var jt [256]operation
		switch {
		case evm.chainRules.IsFuxi:
			jt = fuxiInstructionSet
		default:
			jt = panguInstructionSet
		}

		cfg.JumpTable = jt
	}

	return &Interpreter{
		evm:      evm,
		Cfg:      cfg,
		gasTable: evm.ChainConfig().GasTable(evm.BlockNumber),
	}
}

// MOAC enforce with the IsByzantium
func (in *Interpreter) enforceRestrictions(op OpCode, operation operation, stack *Stack) error {
	if in.readOnly {
		// If the interpreter is operating in readonly mode, make sure no
		// state-modifying operation is performed. The 3rd stack item
		// for a call operation is the value. Transferring value from one
		// account to the others means the state is modified and should also
		// return with an error.
		if operation.writes || (op == CALL && stack.Back(2).BitLen() > 0) {
			return errWriteProtection
		}
	}
	return nil
}

// Run loops and evaluates the contract's code with the given input data and returns
// the return byte-slice and an error if one occurred.
//
// It's important to note that any errors returned by the interpreter should be
// considered a revert-and-consume-all-gas operation. No error specific checks
// should be handled to reduce complexity and errors further down the in.
func (in *Interpreter) Run(snapshot int, contract *Contract, input []byte, precompiledContracts ContractsInterface, msgHash *common.Hash) (ret []byte, err error) {
	log.Debugf("[core/vm/interpreter.go->Run] contract %s input size in bytes: %d", contract.Address().String(), len(input))
	// Increment the call depth which is restricted to 1024
	in.evm.depth++
	defer func() { in.evm.depth-- }()

	// Reset the previous call's return data. It's unimportant to preserve the old buffer
	// as every returning call will return new data anyway.
	in.returnData = nil

	// Don't bother with the execution if there's no code.
	if len(contract.Code) == 0 {
		log.Debugf("len(contract.Code) == 0 err %v", err)
		return nil, nil
	}

	codehash := contract.CodeHash // codehash is used when doing jump dest caching
	if codehash == (common.Hash{}) {
		codehash = crypto.Keccak256Hash(contract.Code)
	}

	var (
		op     OpCode             // current opcode
		mem    = NewMemory()      // bound memory
		stack  = newstack()       // local stack
		rstack = newReturnStack() // return stack
		// For optimisation reason we're using uint64 as the program counter.
		// It's theoretically possible to go above 2^64. The YP defines the PC
		// to be uint256. Practically much less so feasible.
		pc   = uint64(0) // program counter
		cost uint64
	)
	contract.Input = input

	defer func() {
		if err != nil && in.Cfg.Debug {
			in.Cfg.Tracer.CaptureState(in.evm, pc, op, contract.GasRemaining, cost, mem, stack, contract, in.evm.depth, err)
		}
	}()
	log.Debugf("[core/vm/interpreter.go->Run] gasRemaining: %d", contract.GasRemaining)
	// The Interpreter main run loop (contextual). This loop runs until either an
	// explicit STOP, RETURN or SELFDESTRUCT is executed, an error occurred during
	// the execution of one of the operations or until the done flag is set by the
	// parent context.
	for atomic.LoadInt32(&in.evm.abort) == 0 {
		// Get the memory location of pc
		op = contract.GetOp(pc)

		// get the operation from the jump table matching the opcode
		operation := in.Cfg.JumpTable[op]
		if err := in.enforceRestrictions(op, operation, stack); err != nil {
			log.Debugf("enforceRestrictions err %v", err)
			return nil, err
		}

		// if the op is invalid abort the process and return an error
		if !operation.valid {
			log.Debugf("!operation.valid[%s] err %v", op, err)
			return nil, ErrInvalidCode
		}

		// validate the stack and make sure there enough stack items available
		// to perform the operation
		if err := operation.validateStack(stack); err != nil {
			log.Debugf("operation.validateStack(stack) err %v", err)
			return nil, err
		}

		var memorySize uint64
		// calculate the new memory size and expand the memory to fit
		// the operation
		if operation.memorySize != nil {
			memSize, overflow := bigUint64(operation.memorySize(stack))
			if overflow {
				log.Debugf("overflow err on op: %v pc: %v, stack: %v, %v", op, pc, stack, operation)
				return nil, errGasUintOverflow
			}
			// memory is expanded in words of 32 bytes. GasRemaining
			// is also calculated in words.
			if memorySize, overflow = math.SafeMul(toWordSize(memSize), 32); overflow {
				log.Debugf("overflow 2 err %v", err)
				return nil, errGasUintOverflow
			}
		}

		if !in.Cfg.DisableGasMetering {
			// consume the gas and return an error if not enough gas is available.
			// cost is explicitly set so that the capture state defer method cas get the proper cost
			cost, err = operation.gasCost(in.gasTable, in.evm, contract, stack, mem, memorySize)
			if err != nil || !contract.UseGas(cost) {
				log.Debugf("!contract.UseGas(cost) err %v", err)
				return nil, ErrOutOfGas
			}
		}
		if memorySize > 0 {
			mem.Resize(memorySize)
		}

		if in.Cfg.Debug {
			in.Cfg.Tracer.CaptureState(in.evm, pc, op, contract.GasRemaining, cost, mem, stack, contract, in.evm.depth, err)
		}

		res, err := operation.execute(&pc, in.evm, contract, mem, stack, rstack, precompiledContracts, msgHash)

		// if the operation clears the return data (e.g. it has returning data)
		// set the last return to the result of the operation.
		if operation.returns {
			in.returnData = common.CopyBytes(res)
		}

		switch {
		case err != nil:
			return nil, err
		case operation.reverts:
			return res, errExecutionReverted

		case operation.halts:
			return res, nil
		case !operation.jumps:
			pc++
		}
	}
	return nil, nil
}
