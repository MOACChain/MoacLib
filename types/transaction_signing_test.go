// Copyright 2017-2020  The MOAC Foundation
// This file is part of the moac-lib library.
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
	"math/big"
	"testing"

	"github.com/MOACChain/MoacLib/common"
	"github.com/MOACChain/MoacLib/crypto"
	"github.com/MOACChain/MoacLib/log"
	"github.com/MOACChain/MoacLib/rlp"
)

var emptyAddress = common.StringToAddress("")

func TestEIP155Signing(t *testing.T) {
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)

	signer := NewPanguSigner(big.NewInt(18))
	tx, err := SignTx(NewTransaction(0, addr, new(big.Int), new(big.Int), new(big.Int), 0, nil), signer, key)
	if err != nil {
		t.Fatal(err)
	}

	from, err := Sender(signer, tx)
	if err != nil {
		t.Fatal(err)
	}
	if from != addr {
		t.Errorf("exected from and address to be equal. Got %x want %x", from, addr)
	}
}

func TestEIP155ChainId(t *testing.T) {
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)

	//Set the signer with chainID = 18
	signer := NewPanguSigner(big.NewInt(18))
	tx, err := SignTx(NewTransaction(0, addr, new(big.Int), new(big.Int), new(big.Int), 0, nil), signer, key)
	if err != nil {
		t.Fatal(err)
	}
	if !tx.Protected() {
		t.Fatal("expected tx to be protected")
	}

	if tx.ChainId().Cmp(signer.chainId) != 0 {
		t.Error("expected chainId to be", signer.chainId, "got", tx.ChainId())
	}

	tx = NewTransaction(0, addr, new(big.Int), new(big.Int), new(big.Int), 0, nil)
	tx, err = SignTx(tx, PanguSigner{}, key)
	if err != nil {
		t.Fatal(err)
	}

	if tx.Protected() {
		t.Error("didn't expect tx to be protected")
	}

	if tx.ChainId().Sign() != 0 {
		t.Error("expected chain id to be 0 got", tx.ChainId())
	}
}

/*
 * Use MOAC signature algorithm
 * to sign and decode the signature
 * Can be used to verify Chain3 libs.
 */
func TestEIP155MoacSigning(t *testing.T) {
	// Customized test cases
	for i, test := range []struct {
		chainid     int64
		txRlp, addr string
	}{
		//chainid = 100
		{100, "f8713c80850ba43b7400834c4b4094d814f2ac2c4ca49b33066582e4e97ebae02f2ab9888ac7230489e8000000808081eca043d6fa3c6f3b75356ad034d118a8e17d660c95a9490445080a2ae6990b5c24d0a02ab20ac0039a2dceec7117ecb689525dae3d9fb5a7198f3e686369901ecabf5a", "0x7312F4B8A4457a36827f185325Fd6B66a3f8BB8B"},
		//testnet 101,
		{101, "f8718080850ba43b7400834c4b4094d814f2ac2c4ca49b33066582e4e97ebae02f2ab9888ac7230489e8000000808081eea04fdeb0c315ce473e08303e7705d75b5a743b96b1abf4ffaf19f3e367f252648aa07eaa4c5d75c90398c4f016a17754e99e6d88f8629b98245bce57fa2b3f5af3ea", "0x7312F4B8A4457a36827f185325Fd6B66a3f8BB8B"},
		//testnet 106,
		{106, "f86f80808504a817c8008252089432d6f648a651c5e458315641863a386914adb74788016345785d8a000000808081f8a085d45752bf46afc756cb3f1719fe8d8c05ecf8dfb3534551ad340f723d4563fa9ff7689ed34ea42fd5d8f72449bcf0a919866534ce7b7b9d2c95118281cc5160", "0x05a729a0B7965dBAaad6e4Ef9566ca96DE1E0d27"},
	} {
		//Setup the signer with chainID
		//99- mainnet, 100 - devnet, 101 - testnet,
		signer := NewPanguSigner(big.NewInt(test.chainid))

		var tx *Transaction
		err := rlp.DecodeBytes(common.Hex2Bytes(test.txRlp), &tx)
		log.Info("[core/tx_pool.go->TxPool.add] tx=%v", tx.String())
		if err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}

		from, err := Sender(signer, tx)

		if err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}

		addr := common.HexToAddress(test.addr)
		if from != addr {
			t.Errorf("%d: expected %x got %x", i, addr, from)
		}
	}
}

func TestChainId(t *testing.T) {
	key, _ := defaultTestKey()

	tx := NewTransaction(0, common.Address{}, new(big.Int), new(big.Int), new(big.Int), 0, nil)

	var err error
	tx, err = SignTx(tx, NewPanguSigner(big.NewInt(101)), key)
	if err != nil {
		t.Fatal(err)
	}

	_, err = Sender(NewPanguSigner(big.NewInt(100)), tx)
	if err != ErrInvalidChainId {
		t.Error("expected error:", ErrInvalidChainId)
	}

	_, err = Sender(NewPanguSigner(big.NewInt(101)), tx)
	if err != nil {
		t.Error("expected no error")
	}
}
