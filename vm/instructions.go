// Copyright 2015 The go-ethereum Authors
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
	"errors"
	"fmt"
	"math/big"

	"github.com/holiman/uint256"
	"github.com/MOACChain/MoacLib/common"
	"github.com/MOACChain/MoacLib/common/math"
	"github.com/MOACChain/MoacLib/crypto"
	"github.com/MOACChain/MoacLib/log"
	"github.com/MOACChain/MoacLib/types"
	"github.com/MOACChain/MoacLib/params"
)

var (
	bigZero                  = new(big.Int)
	errWriteProtection       = errors.New("evm: write protection")
	errReturnDataOutOfBounds = errors.New("evm: return data out of bounds")
	errExecutionReverted     = errors.New("evm: execution reverted")
	errMaxCodeSizeExceeded   = errors.New("evm: max code size exceeded")
)

func opAdd(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	y.Add(&x, y)
	return nil, nil
}

func opSub(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	y.Sub(&x, y)
	return nil, nil
}

func opMul(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	y.Mul(&x, y)
	return nil, nil
}

func opDiv(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	y.Div(&x, y)
	return nil, nil
}

func opSdiv(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	y.SDiv(&x, y)
	return nil, nil
}

func opMod(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	y.Mod(&x, y)
	return nil, nil
}

func opSmod(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	y.SMod(&x, y)
	return nil, nil
}

func opExp(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	base, exponent := stack.pop(), stack.peek()
	exponent.Exp(&base, exponent)
	return nil, nil
}

func opSignExtend(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	back, num := stack.pop(), stack.peek()
	num.ExtendSign(num, &back)
	return nil, nil
}

func opNot(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	x := stack.peek()
	x.Not(x)
	return nil, nil
}

func opLt(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	if x.Lt(y) {
		y.SetOne()
	} else {
		y.Clear()
	}
	return nil, nil
}

func opGt(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	if x.Gt(y) {
		y.SetOne()
	} else {
		y.Clear()
	}
	return nil, nil
}

func opSlt(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	if x.Slt(y) {
		y.SetOne()
	} else {
		y.Clear()
	}
	return nil, nil
}

func opSgt(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	if x.Sgt(y) {
		y.SetOne()
	} else {
		y.Clear()
	}
	return nil, nil
}

func opEq(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	if x.Eq(y) {
		y.SetOne()
	} else {
		y.Clear()
	}
	return nil, nil
}

func opIszero(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	x := stack.peek()
	if x.IsZero() {
		x.SetOne()
	} else {
		x.Clear()
	}
	return nil, nil
}

func opAnd(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	y.And(&x, y)
	return nil, nil
}

func opOr(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	y.Or(&x, y)
	return nil, nil
}

func opXor(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	y.Xor(&x, y)
	return nil, nil
}

func opByte(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	th, val := stack.pop(), stack.peek()
	val.Byte(&th)
	return nil, nil
}

func opAddmod(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	x, y, z := stack.pop(), stack.pop(), stack.peek()
	if z.IsZero() {
		z.Clear()
	} else {
		z.AddMod(&x, &y, z)
	}
	return nil, nil
}

func opMulmod(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	x, y, z := stack.pop(), stack.pop(), stack.peek()
	z.MulMod(&x, &y, z)
	return nil, nil
}

// opSHL implements Shift Left
// The SHL instruction (shift left) pops 2 values from the stack, first arg1 and then arg2,
// and pushes on the stack arg2 shifted to the left by arg1 number of bits.
func opSHL(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	// Note, second operand is left in the stack; accumulate result into it, and no need to push it afterwards
	shift, value := stack.pop(), stack.peek()
	if shift.LtUint64(256) {
		value.Lsh(value, uint(shift.Uint64()))
	} else {
		value.Clear()
	}
	return nil, nil
}

// opSHR implements Logical Shift Right
// The SHR instruction (logical shift right) pops 2 values from the stack, first arg1 and then arg2,
// and pushes on the stack arg2 shifted to the right by arg1 number of bits with zero fill.
func opSHR(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	// Note, second operand is left in the stack; accumulate result into it, and no need to push it afterwards
	shift, value := stack.pop(), stack.peek()
	if shift.LtUint64(256) {
		value.Rsh(value, uint(shift.Uint64()))
	} else {
		value.Clear()
	}
	return nil, nil
}

// opSAR implements Arithmetic Shift Right
// The SAR instruction (arithmetic shift right) pops 2 values from the stack, first arg1 and then arg2,
// and pushes on the stack arg2 shifted to the right by arg1 number of bits with sign extension.
func opSAR(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	shift, value := stack.pop(), stack.peek()
	if shift.GtUint64(256) {
		if value.Sign() >= 0 {
			value.Clear()
		} else {
			// Max negative shift: all bits set
			value.SetAllOne()
		}
		return nil, nil
	}
	n := uint(shift.Uint64())
	value.SRsh(value, n)
	return nil, nil
}

func opSelfBalance(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	balance, _ := uint256.FromBig(evm.StateDB.GetBalance(contract.Address()))
	stack.push(balance)
	return nil, nil
}

// opChainID implements CHAINID opcode
func opChainID(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	chainId, _ := uint256.FromBig(evm.chainConfig.ChainId)
	stack.push(chainId)
	return nil, nil
}

func opBeginSub(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	return nil, ErrInvalidSubroutineEntry
}

func opJumpSub(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	if len(rstack.data) >= 1023 {
		return nil, ErrReturnStackExceeded
	}
	pos := stack.pop()
	if !pos.IsUint64() {
		return nil, ErrInvalidJump
	}

	posU64 := pos.Uint64()
	if !contract.jumpdests.has(contract.CodeHash, contract.Code, pos.ToBig()) {
		return nil, ErrInvalidJump
	}

	rstack.push(uint32(*pc))
	*pc = posU64 + 1
	return nil, nil
}

func opReturnSub(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	if len(rstack.data) == 0 {
		return nil, ErrInvalidRetsub
	}
	// Other than the check that the return stack is not empty, there is no
	// need to validate the pc from 'returns', since we only ever push valid
	//values onto it via jumpsub.
	v := rstack.pop()
	*pc = uint64(v) + 1
	return nil, nil
}

// opExtCodeHash returns the code hash of a specified account.
// There are several cases when the function is called, while we can relay everything
// to `state.GetCodeHash` function to ensure the correctness.
//   (1) Caller tries to get the code hash of a normal contract account, state
// should return the relative code hash and set it as the result.
//
//   (2) Caller tries to get the code hash of a non-existent account, state should
// return common.Hash{} and zero will be set as the result.
//
//   (3) Caller tries to get the code hash for an account without contract code,
// state should return emptyCodeHash(0xc5d246...) as the result.
//
//   (4) Caller tries to get the code hash of a precompiled account, the result
// should be zero or emptyCodeHash.
//
// It is worth noting that in order to avoid unnecessary create and clean,
// all precompile accounts on mainnet have been transferred 1 wei, so the return
// here should be emptyCodeHash.
// If the precompile account is not transferred any amount on a private or
// customized chain, the return value will be zero.
//
//   (5) Caller tries to get the code hash for an account which is marked as suicided
// in the current transaction, the code hash of this account should be returned.
//
//   (6) Caller tries to get the code hash for an account which is marked as deleted,
// this account should be regarded as a non-existent account and zero should be returned.
func opExtCodeHash(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	slot := stack.peek()
	address := common.Address(slot.Bytes20())
	if evm.StateDB.Empty(address) {
		slot.Clear()
	} else {
		slot.SetBytes(evm.StateDB.GetCodeHash(address).Bytes())
	}
	return nil, nil
}

func opCreate2(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	var (
		endowment    = stack.pop()
		offset, size = stack.pop(), stack.pop()
		salt         = stack.pop()
		input        = memory.GetCopy(int64(offset.Uint64()), int64(size.Uint64()))
		gas          = contract.GasRemaining
	)

	// Apply EIP150
	gas -= gas / 64
	contract.UseGas(gas)
	// reuse size int for stackvalue
	stackvalue := size
	bigEndowment := big.NewInt(0)
	if !endowment.IsZero() {
		bigEndowment = endowment.ToBig()
	}
	res, addr, returnGas, suberr := evm.Create2(contract, input, gas, bigEndowment, 0, precompiledContracts, msgHash, &salt)
	// Push item on the stack based on the returned error.
	if suberr != nil {
		stackvalue.Clear()
	} else {
		stackvalue.SetBytes(addr.Bytes())
	}
	stack.push(&stackvalue)
	contract.GasRemaining += returnGas

	if suberr == ErrExecutionReverted {
		return res, nil
	}
	return nil, nil
}

func opSha3(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	offset, size := stack.pop(), stack.pop()
	data := memory.Get(int64(offset.Uint64()), int64(size.Uint64()))
	hash := crypto.Keccak256(data)

	if evm.VmConfig.EnablePreimageRecording {
		evm.StateDB.AddPreimage(common.BytesToHash(hash), data)
	}

	stack.push(new(uint256.Int).SetBytes(hash))

	return nil, nil
}

func opAddress(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	stack.push(new(uint256.Int).SetBytes(contract.Address().Bytes()))
	return nil, nil
}

func opBalance(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	slot := stack.peek()
	address := common.Address(slot.Bytes20())
	slot.SetFromBig(evm.StateDB.GetBalance(address))
	return nil, nil
}

func opOrigin(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	stack.push(new(uint256.Int).SetBytes(evm.Origin.Bytes()))
	return nil, nil
}

func opCaller(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	stack.push(new(uint256.Int).SetBytes(contract.Caller().Bytes()))
	return nil, nil
}

func opCallValue(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	v, _ := uint256.FromBig(contract.value)
	stack.push(v)
	return nil, nil
}

func opCallDataLoad(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	v := stack.pop()
	stack.push(new(uint256.Int).SetBytes(getDataBig(contract.Input, v.ToBig(), common.Big32)))
	return nil, nil
}

func opCallDataSize(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	stack.push(new(uint256.Int).SetUint64(uint64(len(contract.Input))))
	return nil, nil
}

func opCallDataCopy(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	var (
		memOffset  = stack.pop()
		dataOffset = stack.pop()
		length     = stack.pop()
	)
	memory.Set(memOffset.Uint64(), length.Uint64(), getDataBig(contract.Input, dataOffset.ToBig(), length.ToBig()))
	return nil, nil
}

func opReturnDataSize(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	stack.push(new(uint256.Int).SetUint64(uint64(len(evm.interpreter.returnData))))
	return nil, nil
}

func opReturnDataCopy(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	var (
		memOffset  = stack.pop()
		dataOffset = stack.pop()
		length     = stack.pop()
	)

	offset64, overflow := dataOffset.Uint64WithOverflow()
	if overflow {
		return nil, ErrReturnDataOutOfBounds
	}
	// we can reuse dataOffset now (aliasing it for clarity)
	var end = dataOffset
	end.Add(&dataOffset, &length)
	end64, overflow := end.Uint64WithOverflow()
	if overflow || uint64(len(evm.interpreter.returnData)) < end64 {
		return nil, ErrReturnDataOutOfBounds
	}
	memory.Set(memOffset.Uint64(), length.Uint64(), evm.interpreter.returnData[offset64:end64])
	return nil, nil
}

func opExtCodeSize(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	a := stack.pop()
	addr := common.Uint256ToAddress(&a)
	a.SetUint64(uint64(evm.StateDB.GetCodeSize(addr)))
	log.Debugf("opExtCodeSize addr %v code size %v", addr.Hex(), a)
	stack.push(&a)

	return nil, nil
}

func opCodeSize(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	l := new(uint256.Int)
	l.SetUint64(uint64(len(contract.Code)))
	stack.push(l)
	return nil, nil
}

func opCodeCopy(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	var (
		memOffset  = stack.pop()
		codeOffset = stack.pop()
		length     = stack.pop()
	)
	codeCopy := getDataBig(contract.Code, codeOffset.ToBig(), length.ToBig())
	memory.Set(memOffset.Uint64(), length.Uint64(), codeCopy)
	return nil, nil
}

func opExtCodeCopy(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	v := stack.pop()
	var (
		addr       = common.Uint256ToAddress(&v)
		memOffset  = stack.pop()
		codeOffset = stack.pop()
		length     = stack.pop()
	)
	codeCopy := getDataBig(evm.StateDB.GetCode(addr), codeOffset.ToBig(), length.ToBig())
	memory.Set(memOffset.Uint64(), length.Uint64(), codeCopy)
	return nil, nil
}

func opGasprice(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	v, _ := uint256.FromBig(evm.GasPrice)
	stack.push(v)
	return nil, nil
}

func opBlockhash(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	num := stack.peek()
	num64, overflow := num.Uint64WithOverflow()
	if overflow {
		num.Clear()
		return nil, nil
	}
	var upper, lower uint64
	upper = evm.BlockNumber.Uint64()
	if upper < 257 {
		lower = 0
	} else {
		lower = upper - 256
	}
	if num64 >= lower && num64 < upper {
		num.SetBytes(evm.GetHash(num64).Bytes())
	} else {
		num.Clear()
	}
	return nil, nil
}

func opCoinbase(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	stack.push(new(uint256.Int).SetBytes(evm.Coinbase.Bytes()))
	return nil, nil
}

func opTimestamp(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	v, _ := uint256.FromBig(evm.Time)
	stack.push(v)
	return nil, nil
}

func opNumber(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	v, _ := uint256.FromBig(evm.BlockNumber)
	stack.push(v)
	return nil, nil
}

func opDifficulty(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	v, _ := uint256.FromBig(evm.Difficulty)
	stack.push(v)
	return nil, nil
}

func opGasLimit(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	v, _ := uint256.FromBig(evm.GasLimit)
	stack.push(v)
	return nil, nil
}

func opPop(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	stack.pop()
	return nil, nil
}

func opMload(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	offset := stack.pop()
	val := new(uint256.Int).SetBytes(memory.Get(int64(offset.Uint64()), 32))
	stack.push(val)

	return nil, nil
}

func opMstore(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	// pop value of the stack
	mStart, val := stack.pop(), stack.pop()
	memory.Set(mStart.Uint64(), 32, math.PaddedBigBytes(val.ToBig(), 32))
	return nil, nil
}

func opMstore8(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	off, val := stack.pop(), stack.pop()
	memory.store[off.Uint64()] = byte(val.Uint64())

	return nil, nil
}

func opSload(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	loc := stack.peek()
	hash := common.Hash(loc.Bytes32())
	val := evm.StateDB.GetState(contract.Address(), hash)
	loc.SetBytes(val.Bytes())
	return nil, nil
}

func opSstore(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	loc := stack.pop()
	val := stack.pop()
	evm.StateDB.SetState(contract.Address(), common.Hash(loc.Bytes32()), common.Hash(val.Bytes32()))
	return nil, nil
}

func opJump(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	pos := stack.pop()
	if !contract.jumpdests.has(contract.CodeHash, contract.Code, pos.ToBig()) {
		nop := contract.GetOp(pos.Uint64())
		return nil, fmt.Errorf("invalid jump destination (%v) %v", nop, pos)
	}
	*pc = pos.Uint64()

	return nil, nil
}

func opJumpi(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	pos, cond := stack.pop(), stack.pop()
	if cond.Sign() != 0 {
		if !contract.jumpdests.has(contract.CodeHash, contract.Code, pos.ToBig()) {
			nop := contract.GetOp(pos.Uint64())
			return nil, fmt.Errorf("invalid jump destination (%v) %v", nop, pos)
		}
		*pc = pos.Uint64()
	} else {
		*pc++
	}

	return nil, nil
}

func opJumpdest(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	return nil, nil
}

func opPc(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	stack.push(new(uint256.Int).SetUint64(*pc))
	return nil, nil
}

func opMsize(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {

	stack.push(new(uint256.Int).SetUint64(uint64(memory.Len())))
	return nil, nil
}

func opGas(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	stack.push(new(uint256.Int).SetUint64(contract.GasRemaining))
	return nil, nil
}

func opCreate(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	var (
		value        = stack.pop()
		offset, size = stack.pop(), stack.pop()
		input        = memory.Get(int64(offset.Uint64()), int64(size.Uint64()))
		gas          = contract.GasRemaining
	)
	gas -= gas / 64

	contract.UseGas(gas)
	res, addr, returnGas, suberr := evm.Create(contract, input, gas, value.ToBig(), 0, precompiledContracts, msgHash) //TODO not unknow, but what
	// Push item on the stack based on the returned error. If the ruleset is
	// pangu we must check for CodeStoreOutOfGasError (pangu only
	// rule) and treat as an error, if the ruleset is frontier we must
	// ignore this error and pretend the operation was successful.
	if evm.ChainConfig().IsPangu(evm.BlockNumber) && suberr == ErrCodeStoreOutOfGas {
		stack.push(new(uint256.Int))
	} else if suberr != nil && suberr != ErrCodeStoreOutOfGas {
		stack.push(new(uint256.Int))
	} else {
		v, _ := uint256.FromBig(addr.Big())
		stack.push(v)
	}
	contract.GasRemaining += returnGas

	if suberr == errExecutionReverted {
		return res, nil
	}
	return nil, nil
}

func opCall(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	g := stack.pop()
	gas := g.Uint64()
	// pop gas and value of the stack.
	addr, v := stack.pop(), stack.pop()
	value := math.U256(v.ToBig())
	// pop input size and offset
	inOffset, inSize := stack.pop(), stack.pop()
	// pop return size and offset
	retOffset, retSize := stack.pop(), stack.pop()

	address := common.Uint256ToAddress(&addr)

	// Get the arguments from the memory
	args := memory.Get(int64(inOffset.Uint64()), int64(inSize.Uint64()))

	if value.Sign() != 0 {
		gas += params.CallStipend
	}
	ret, returnGas, err := evm.Call(contract, address, args, gas, value, false, 0, precompiledContracts, msgHash)
	if err != nil {
		stack.push(new(uint256.Int))
	} else {
		stack.push(uint256.NewInt().SetUint64(1))
	}
	if err == nil || err == errExecutionReverted {
		memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
	}
	contract.GasRemaining += returnGas

	return ret, nil
}

func opCallCode(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	g := stack.pop()
	gas := g.Uint64()
	// pop gas and value of the stack.
	addr, v := stack.pop(), stack.pop()
	value := math.U256(v.ToBig())
	// pop input size and offset
	inOffset, inSize := stack.pop(), stack.pop()
	// pop return size and offset
	retOffset, retSize := stack.pop(), stack.pop()

	address := common.Uint256ToAddress(&addr)

	// Get the arguments from the memory
	args := memory.Get(int64(inOffset.Uint64()), int64(inSize.Uint64()))

	if value.Sign() != 0 {
		gas += params.CallStipend
	}

	ret, returnGas, err := evm.CallCode(contract, address, args, gas, value, precompiledContracts, msgHash)
	if err != nil {
		stack.push(new(uint256.Int))
	} else {
		stack.push(uint256.NewInt().SetUint64(1))
	}
	if err == nil || err == errExecutionReverted {
		memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
	}
	contract.GasRemaining += returnGas

	return ret, nil
}

func opDelegateCall(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	g, to, inOffset, inSize, outOffset, outSize := stack.pop(), stack.pop(), stack.pop(), stack.pop(), stack.pop(), stack.pop()
	gas := g.Uint64()

	toAddr := common.Uint256ToAddress(&to)
	args := memory.Get(int64(inOffset.Uint64()), int64(inSize.Uint64()))

	ret, returnGas, err := evm.DelegateCall(contract, toAddr, args, gas, precompiledContracts, msgHash)
	if err != nil {
		stack.push(new(uint256.Int))
	} else {
		stack.push(uint256.NewInt().SetUint64(1))
	}
	if err == nil || err == errExecutionReverted {
		memory.Set(outOffset.Uint64(), outSize.Uint64(), ret)
	}
	contract.GasRemaining += returnGas
	return ret, nil
}

func opStaticCall(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	// pop gas
	g := stack.pop()
	gas := g.Uint64()
	// pop address
	addr := stack.pop()
	// pop input size and offset
	inOffset, inSize := stack.pop(), stack.pop()
	// pop return size and offset
	retOffset, retSize := stack.pop(), stack.pop()

	address := common.Uint256ToAddress(&addr)

	// Get the arguments from the memory
	args := memory.Get(int64(inOffset.Uint64()), int64(inSize.Uint64()))

	ret, returnGas, err := evm.StaticCall(contract, address, args, gas, precompiledContracts, msgHash)
	if err != nil {
		stack.push(new(uint256.Int))
	} else {
		stack.push(uint256.NewInt().SetUint64(1))
	}
	if err == nil || err == errExecutionReverted {
		memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
	}
	contract.GasRemaining += returnGas
	return ret, nil
}

func opReturn(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	offset, size := stack.pop(), stack.pop()
	ret := memory.GetPtr(int64(offset.Uint64()), int64(size.Uint64()))
	return ret, nil
}

func opRevert(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	offset, size := stack.pop(), stack.pop()
	ret := memory.GetPtr(int64(offset.Uint64()), int64(size.Uint64()))
	return ret, nil
}

func opStop(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	return nil, nil
}

func opSuicide(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	balance := evm.StateDB.GetBalance(contract.Address())
	addr := stack.pop()
	evm.StateDB.AddBalance(common.Uint256ToAddress(&addr), balance)

	evm.StateDB.Suicide(contract.Address())
	return nil, nil
}

// following functions are used by the instruction jump  table

// make log instruction function
func makeLog(size int) executionFunc {
	return func(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
		topics := make([]common.Hash, size)
		mStart, mSize := stack.pop(), stack.pop()
		for i := 0; i < size; i++ {
			t := stack.pop()
			topics[i] = common.BigToHash(t.ToBig())
		}

		d := memory.Get(int64(mStart.Uint64()), int64(mSize.Uint64()))
		evm.StateDB.AddLog(&types.Log{
			Address: contract.Address(),
			Topics:  topics,
			Data:    d,
			// This is a non-consensus field, but assigned here because
			// core/state doesn't know the current block number.
			BlockNumber: evm.BlockNumber.Uint64(),
		})

		return nil, nil
	}
}

// make push instruction function
func makePush(size uint64, pushByteSize int) executionFunc {
	return func(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
		codeLen := len(contract.Code)

		startMin := codeLen
		if int(*pc+1) < startMin {
			startMin = int(*pc + 1)
		}

		endMin := codeLen
		if startMin+pushByteSize < endMin {
			endMin = startMin + pushByteSize
		}

		integer := uint256.NewInt()
		stack.push(integer.SetBytes(common.RightPadBytes(contract.Code[startMin:endMin], pushByteSize)))

		*pc += size
		return nil, nil
	}
}

// make push instruction function
func makeDup(size int64) executionFunc {
	return func(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
		stack.dup(int(size))
		return nil, nil
	}
}

// make swap instruction function
func makeSwap(size int64) executionFunc {
	// switch n + 1 otherwise n would be swapped with n
	size += 1
	return func(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, rstack *ReturnStack, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
		stack.swap(int(size))
		return nil, nil
	}
}
