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

import "errors"

var (
	ErrOutOfGas                 = errors.New("out of gas")
	ErrCodeStoreOutOfGas        = errors.New("contract creation code storage out of gas")
	ErrDepth                    = errors.New("max call depth exceeded")
	ErrTraceLimitReached        = errors.New("the number of logs reached the specified limit")
	ErrInsufficientBalance      = errors.New("insufficient balance for transfer")
	ErrContractAddressCollision = errors.New("contract address collision")
	ErrEmptyCode                = errors.New("code empty")
	ErrGRPCService              = errors.New("grpc service error")
	ErrInputFormat              = errors.New("Input format error")
	ErrInvalidCode              = errors.New("Invalid opcode")
	ErrReturnDataOutOfBounds    = errors.New("return data out of bounds")
	ErrReturnStackExceeded      = errors.New("return stack limit reached")
	ErrInvalidRetsub            = errors.New("invalid retsub")
	ErrInvalidJump              = errors.New("invalid jump destination")
	ErrInvalidSubroutineEntry   = errors.New("invalid subroutine entry")
	ErrExecutionReverted        = errors.New("execution reverted")
)
