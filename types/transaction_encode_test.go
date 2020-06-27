// Copyright 2017  The MOAC Foundation
// This file is modified from the go-ethereum library.
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

package types

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"strings"
	"testing"

	"github.com/MOACChain/MoacLib/common"
	"github.com/MOACChain/MoacLib/rlp"
)

// transaction format in Nuwa version
type txnvwa struct {
	AccountNonce   uint64          `json:"nonce"    gencodec:"required"`
	SystemContract uint64          `json:"syscnt" gencodec:"required"`
	Price          *big.Int        `json:"gasPrice" gencodec:"required"`
	GasLimit       *big.Int        `json:"gas"      gencodec:"required"`
	Recipient      *common.Address `json:"to"       rlp:"nil"` // nil means contract creation
	Amount         *big.Int        `json:"value"    gencodec:"required"`
	Payload        []byte          `json:"input"    gencodec:"required"`
	ShardingFlag   uint64          `json:"shardingFlag" gencodec:"required"`
	Via            interface{}     `json:"via"       rlp:"nil"`

	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash *common.Hash `json:"hash" rlp:"-"`
}

// Empty TX with signature fields
var emptyTxWithSignature = txnvwa{
	AccountNonce:   0,
	SystemContract: 0,
	ShardingFlag:   0,
	Via:            nil,
	Payload:        nil,
	Amount:         big.NewInt(0),
	GasLimit:       big.NewInt(0),
	Price:          big.NewInt(0),
	V:              big.NewInt(0),
	R:              big.NewInt(0),
	S:              big.NewInt(0),
}

// Recipient:      nil,
func (e *txnvwa) EncodeRLP(w io.Writer) error {
	if e == nil {
		w.Write([]byte{0, 0, 0, 0})
	} else {
		w.Write([]byte{0, 1, 0, 1, 0, 1, 0, 1, 0, 1})
	}
	return nil
}

type encTTest struct {
	val           interface{}
	output, error string
}

type simplestruct struct {
	A uint
	B string
}

var encTxTests = []encTTest{

	// structs
	{val: simplestruct{}, output: "C28080"},
	{val: simplestruct{A: 3, B: "foo"}, output: "C50383666F6F"},

	{val: (*[]struct{ uint })(nil), output: "C0"},
	{val: (*interface{})(nil), output: "C0"},

	// interfaces
	// {val: []io.Reader{reader}, output: "C3C20102"}, // the contained value is a struct

	// Encoder
	{val: (*txnvwa)(nil), output: "00000000"},
	{val: &txnvwa{}, output: "00010001000100010001"},
	// {val: &txnvwa{0}, error: "test error"},
	// verify that pointer method txnvwa.EncodeRLP is called for
	// addressable non-pointer values.
	{val: &struct{ TE txnvwa }{txnvwa{}}, output: "CA00010001000100010001"},
	// {val: &struct{ TE txnvwa }{emptyTxWithSignature}, error: "test error"},
	// verify the error for non-addressable non-pointer Encoder
	{val: txnvwa{}, error: "rlp: game over: unadressable value of type types.txnvwa, EncodeRLP is pointer method"},
	// verify the special case for []byte
	// {val: []byteEncoder{0, 1, 2, 3, 4}, output: "C5C0C0C0C0C0"},
}

func unhex2(str string) []byte {
	b, err := hex.DecodeString(strings.Replace(str, " ", "", -1))
	if err != nil {
		panic(fmt.Sprintf("invalid hex string: %q", str))
	}
	return b
}

func runEncTxTests(t *testing.T, f func(val interface{}) ([]byte, error)) {
	for i, test := range encTxTests {
		output, err := f(test.val)
		if err != nil && test.error == "" {
			t.Errorf("test %d: unexpected error: %v\nvalue %#v\ntype %T",
				i, err, test.val, test.val)
			continue
		}
		if test.error != "" && fmt.Sprint(err) != test.error {
			t.Errorf("test %d: error mismatch\ngot   %v\nwant  %v\nvalue %#v\ntype  %T",
				i, err, test.error, test.val, test.val)
			continue
		}
		if err == nil && !bytes.Equal(output, unhex2(test.output)) {
			t.Errorf("test %d: output mismatch:\ngot   %X\nwant  %s\nvalue %#v\ntype  %T",
				i, output, test.output, test.val, test.val)
		}
	}
}

//TestEncodeTx: test the TX transaction format
func TestEncodeTx(t *testing.T) {
	fmt.Println("Test the Encoding function run...")
	runEncTxTests(t, func(val interface{}) ([]byte, error) {
		b := new(bytes.Buffer)
		err := rlp.Encode(b, val)
		return b.Bytes(), err
	})
}
