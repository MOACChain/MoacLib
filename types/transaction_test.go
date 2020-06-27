// Copyright 2017  The MOAC-lib Authors
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

package types

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/MOACChain/MoacLib/common"
	"github.com/MOACChain/MoacLib/crypto"
	"github.com/MOACChain/MoacLib/rlp"
)

// The values in those tests are for the Transaction Tests

var (
	emptyTx = NewTransaction(
		0,
		common.HexToAddress("7312F4B8A4457a36827f185325Fd6B66a3f8BB8B"),
		big.NewInt(0), big.NewInt(0), big.NewInt(0), 0,
		nil,
	)

	rightvrsTx, _ = NewTransaction(
		22,
		common.HexToAddress("7312F4B8A4457a36827f185325Fd6B66a3f8BB8B"),
		big.NewInt(400000000000),
		big.NewInt(40000000000),
		big.NewInt(2100000),
		0,
		nil,
	).WithSignature(
		PanguSigner{},
		common.Hex2Bytes("98ff921201554726367d2be8c804a7ff89ccf285ebc57dff8ae4c44b9c19ac4a8887321be575c8095f789dd4c743dfe42c1820f9231f98a962b210e3ac2452a301"),
	)

	panguTx = NewTransaction(
		42,
		common.HexToAddress("0000000000000000000000000000000000000065"),
		big.NewInt(0), big.NewInt(0), big.NewInt(20000000000), 0,
		nil,
	)
)

func TestTransactionSigHash(t *testing.T) {

	outHash := emptyTx.SigHash(PanguSigner{})
	// fmt.Printf("Signed:%v\n", len(outHash))
	// fmt.Printf("Signed:%v\n", outHash)
	// fmt.Printf("Signed:%v\n", len(emptyTx.Hash()))
	// fmt.Printf("Signed:%v\n", common.HexToHash("52e8d2bc1372f0c6672f54b1173b302d1204ff5687f797c5f2bf5cca504d9e9a"))

	if outHash != common.HexToHash("33f1c77c68175ffede28ddd0361239251164093901ea11a7baa40d64d820005f") {
		t.Errorf("empty transaction hash mismatch, got %x", emptyTx.Hash())
	}

	if rightvrsTx.SigHash(PanguSigner{}) != common.HexToHash("394ce9ac0a875696de603110f976b6c4f9c09a74783b2d35b51d5d398e7af6c4") {
		t.Errorf("RightVRS transaction hash mismatch, got %x", rightvrsTx.SigHash(PanguSigner{}))
	}
	outHash = panguTx.SigHash(PanguSigner{})
	//093d41e8ce16bcaea82d67ba7b90f899b03085e9d7c44a6781f70bdd30265c0c, with common.Address
	//d2e3b81804a788fb52f7117a359a554bd13b1c58d14ab5f1542f3e920808d28a, interface
	if outHash != common.HexToHash("093d41e8ce16bcaea82d67ba7b90f899b03085e9d7c44a6781f70bdd30265c0c") {
		t.Errorf("Pangu transaction hash mismatch, got %x", panguTx.SigHash(PanguSigner{}))
	}

}

func TestTransactionEncode(t *testing.T) {
	txb, err := rlp.EncodeToBytes(rightvrsTx)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}
	should := common.FromHex("f86d168083200b208509502f9000947312f4b8a4457a36827f185325fd6b66a3f8bb8b855d21dba0008080801ca098ff921201554726367d2be8c804a7ff89ccf285ebc57dff8ae4c44b9c19ac4aa08887321be575c8095f789dd4c743dfe42c1820f9231f98a962b210e3ac2452a3")
	if !bytes.Equal(txb, should) {
		t.Errorf("encoded RLP mismatch, got %x", txb)
	}
}

// GETH TX decoder
func decodeTx(data []byte) (*Transaction, error) {
	var tx Transaction
	t, err := &tx, rlp.Decode(bytes.NewReader(data), &tx)

	fmt.Printf("decoded tx: %s\n", t.String())

	return t, err
}

func defaultTestKey() (*ecdsa.PrivateKey, common.Address) {
	key, _ := crypto.HexToECDSA("c75a5f85ef779dcf95c651612efb3c3b9a6dfafb1bb5375905454d9fc8be8a6b")
	addr := crypto.PubkeyToAddress(key.PublicKey)
	return key, addr
}

/*
 * Test the decoding of MOAC format transaction
 *
 { nonce: '0x04',
  gasPrice: '9c40',
  gasLimit: '07d0',
  to: '0x7312F4B8A4457a36827f185325Fd6B66a3f8BB8B',
  value: '01100000000000000000',
  data: '0x00',
  shardingFlag: 0,
 }
*/
func TestDecodeMoacTx(t *testing.T) {
	addFrom := "0x7312F4B8A4457a36827f185325Fd6B66a3f8BB8B"
	addTo := "0xD814F2ac2c4cA49b33066582E4e97EBae02F2aB9"
	addNonce := uint64(60)

	// notice the cmd string should not have '0x' as prefix
	cmd := "f8713c80850ba43b7400834c4b4094d814f2ac2c4ca49b33066582e4e97ebae02f2ab9888ac7230489e8000000808081eca043d6fa3c6f3b75356ad034d118a8e17d660c95a9490445080a2ae6990b5c24d0a02ab20ac0039a2dceec7117ecb689525dae3d9fb5a7198f3e686369901ecabf5a"
	tx, err := decodeTx(common.Hex2Bytes(cmd))
	if err != nil {
		fmt.Printf("rec add: %v\n", tx.TxData.Recipient)
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("TX nonce: %v\n", tx.TxData.AccountNonce)
	fmt.Printf("TX amount: %v\n", tx.TxData.Amount)
	fmt.Printf("TX gasLimit: %v\n", tx.TxData.GasLimit)
	fmt.Printf("TX ShardingFlag: %v\n", tx.GetShardingFlag())
	fmt.Printf("TX system: %v\n", tx.GetSystemFlag())

	fmt.Printf("Tx src:%v\n%v\n", tx.GetSender(), addFrom)

	if strings.ToLower(tx.TxData.Recipient.Hex()) != strings.ToLower(addTo) {
		t.Error("Derived address doesn't match")
		fmt.Printf("Get:%s, want %s\n", tx.TxData.Recipient.Hex(), addTo)
	}
	if addNonce != tx.TxData.AccountNonce {
		t.Error("Derived nonce doesn't match")
		fmt.Printf("Get:%d, want %d\n", tx.TxData.AccountNonce, addNonce)
	}
	//Test if the signature is valid
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

}

/*
 * Test the system and sharding flags set
 */

func TestControlFlags(t *testing.T) {
	var f1 uint64 = 1
	//Set system flag and clear it
	emptyTx.SetSystemFlag(f1)

	if emptyTx.GetSystemFlag() != f1 {
		fmt.Printf("Set system flag Error: %v\n", emptyTx.TxData.ShardingFlag)
	}
	emptyTx.SetShardingFlag(f1)

	if emptyTx.GetShardingFlag() != f1 {
		fmt.Printf("Set sharding flag 0 Error: %v\n", emptyTx.TxData.ShardingFlag)
	}
	f1 = 0
	emptyTx.SetSystemFlag(f1)

	if emptyTx.GetSystemFlag() != f1 {
		fmt.Printf("Set system flag 0 Error: %v\n", emptyTx.TxData.ShardingFlag)
	}

	emptyTx.SetShardingFlag(f1)

	if emptyTx.GetShardingFlag() != f1 {
		fmt.Printf("Set sharding flag 0 Error: %v\n", emptyTx.TxData.ShardingFlag)
	}

}

// MOAC didn't support HomeStead signer or Frontier Signer due to security reasons
func TestRecipientEmpty(t *testing.T) {
	_, addr := defaultTestKey()
	tx, err := decodeTx(common.Hex2Bytes("f901603d808561c9f368008301c1a08080b90109608060405234801561001057600080fd5b5060ea8061001f6000396000f30060806040526004361060525763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166306661abd81146057578063771602f714607b578063a87d942c146093575b600080fd5b348015606257600080fd5b50606960a5565b60408051918252519081900360200190f35b348015608657600080fd5b50606960043560243560ab565b348015609e57600080fd5b50606960b8565b60005481565b6000805460010190550190565b600054905600a165627a7a7230582021c8fccfce143160093bf7a7678648bb60fd456b95752067a6c81042f70eab970029808081eca0efb1986ea2a905dbbee47ebb3603e3099b989c603a84cd7346402696abcada45a00b675c05feab7815f8d83b78c80055c26cb75bfc6fe586b2299e0ff7c3009590"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// For PanguSigner, need to be initialized with chainId
	psigner := NewPanguSigner(big.NewInt(100))
	from, err := Sender(psigner, tx)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if addr != from {
		t.Error("derived address doesn't match")
	}
}

func TestRecipientNormal(t *testing.T) {
	_, addr := defaultTestKey()
	tx, err := decodeTx(common.Hex2Bytes("f8713d80850ba43b7400834c4b4094d814f2ac2c4ca49b33066582e4e97ebae02f2ab9888ac7230489e8000000808081eca0dbb8aa849c8ac42f174b38253ec498ca1c28cb524f6e1527f145045d253567e4a07a2a6f9878518ba8f4f9ffc79845b6d74f88f6c0601aebea68b18b239f8d84c3"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// For PanguSigner, need to be initialized with chainId
	psigner := NewPanguSigner(big.NewInt(100))
	from, err := Sender(psigner, tx)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if addr != from {
		t.Error("derived address doesn't match")
	}
}

// Tests that transactions can be correctly sorted according to their price in
// decreasing order, but at the same time with increasing nonces when issued by
// the same account.
func TestTransactionPriceNonceSort(t *testing.T) {
	// Generate a batch of accounts to start with
	keys := make([]*ecdsa.PrivateKey, 25)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = crypto.GenerateKey()
	}

	signer := PanguSigner{}
	// Generate a batch of transactions with overlapping values, but shifted nonces
	groups := map[common.Address]Transactions{}
	for start, key := range keys {
		addr := crypto.PubkeyToAddress(key.PublicKey)
		for i := 0; i < 25; i++ {
			tx, _ := SignTx(NewTransaction(uint64(start+i), common.Address{}, big.NewInt(100), big.NewInt(100), big.NewInt(int64(start+i)), 0, nil), signer, key)
			groups[addr] = append(groups[addr], tx)
		}
	}
	// Sort the transactions and cross check the nonce ordering
	txset := NewTransactionsByPriceAndNonce(signer, groups)

	txs := Transactions{}
	for {
		if tx := txset.Peek(); tx != nil {
			txs = append(txs, tx)
			txset.Shift()
		}
		break
	}
	for i, txi := range txs {
		fromi, _ := Sender(signer, txi)

		// Make sure the nonce order is valid
		for j, txj := range txs[i+1:] {
			fromj, _ := Sender(signer, txj)

			if fromi == fromj && txi.Nonce() > txj.Nonce() {
				t.Errorf("invalid nonce ordering: tx #%d (A=%x N=%v) < tx #%d (A=%x N=%v)", i, fromi[:4], txi.Nonce(), i+j, fromj[:4], txj.Nonce())
			}
		}
		// Find the previous and next nonce of this account
		prev, next := i-1, i+1
		for j := i - 1; j >= 0; j-- {
			if fromj, _ := Sender(signer, txs[j]); fromi == fromj {
				prev = j
				break
			}
		}
		for j := i + 1; j < len(txs); j++ {
			if fromj, _ := Sender(signer, txs[j]); fromi == fromj {
				next = j
				break
			}
		}
		// Make sure that in between the neighbor nonces, the transaction is correctly positioned price wise
		for j := prev + 1; j < next; j++ {
			fromj, _ := Sender(signer, txs[j])
			if j < i && txs[j].GasPrice().Cmp(txi.GasPrice()) < 0 {
				t.Errorf("invalid gasprice ordering: tx #%d (A=%x P=%v) < tx #%d (A=%x P=%v)", j, fromj[:4], txs[j].GasPrice(), i, fromi[:4], txi.GasPrice())
			}
			if j > i && txs[j].GasPrice().Cmp(txi.GasPrice()) > 0 {
				t.Errorf("invalid gasprice ordering: tx #%d (A=%x P=%v) > tx #%d (A=%x P=%v)", j, fromj[:4], txs[j].GasPrice(), i, fromi[:4], txi.GasPrice())
			}
		}
	}
}

// TestTransactionJSON tests serializing/de-serializing to/from JSON.
func TestTransactionJSON(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("could not generate key: %v", err)
	}
	signer := NewPanguSigner(common.Big1)

	for i := uint64(0); i < 25; i++ {
		var tx *Transaction
		switch i % 2 {
		case 0:
			tx = NewTransaction(i, common.Address{1}, common.Big0, common.Big1, common.Big2, 0, []byte("abcdef"))
		case 1:
			tx = NewContractCreation(i, common.Big0, common.Big1, common.Big2, 0, []byte("abcdef"))
		}

		tx, err := SignTx(tx, signer, key)
		if err != nil {
			t.Fatalf("could not sign transaction: %v", err)
		}

		data, err := json.Marshal(tx)
		if err != nil {
			t.Errorf("json.Marshal failed: %v", err)
		}

		var parsedTx *Transaction
		if err := json.Unmarshal(data, &parsedTx); err != nil {
			t.Errorf("json.Unmarshal failed: %v", err)
		}

		// compare nonce, price, gaslimit, recipient, amount, payload, V, R, S
		if tx.Hash() != parsedTx.Hash() {
			t.Errorf("parsed tx differs from original tx, want %v, got %v", tx, parsedTx)
		}
		if tx.ChainId().Cmp(parsedTx.ChainId()) != 0 {
			t.Errorf("invalid chain id, want %d, got %d", tx.ChainId(), parsedTx.ChainId())
		}
	}
}
