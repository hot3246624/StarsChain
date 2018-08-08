package plasma_core

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"crypto/ecdsa"
	"github.com/hot3246624/StarsChain/plasma_core/utils"
	"github.com/hot3246624/StarsChain/plasma_core/utils/merkle"
	"fmt"
)

type Block struct {
	TransactionSet []Transaction
	Number *big.Int
	Sig []byte
	SpentUTXOs map[int64]bool
}

func (bl *Block)Hash() common.Hash {
	return utils.RlpHash(bl)
}

func (bl *Block)Singer() common.Address {
	return utils.GetSigner(bl.Hash().Bytes(), bl.Sig)
}

func (bl *Block)Merkle() *merkle.FixedMerkle {
	var hashedTransactions [][]byte
	for _, transaction := range bl.TransactionSet{
		hashedTransactions = append(hashedTransactions, transaction.MerkleHash().Bytes())
	}
	return &merkle.FixedMerkle{Depth: 16, Leaves: hashedTransactions, Hashed: true}
}

func (bl *Block)Root() []byte{
	return bl.Merkle().Root
}

func (bl *Block)IsDepositBlock() bool{
	return len(bl.TransactionSet) == 1 && bl.TransactionSet[0].IsDepositTransaction()
}

//TODO set sig = nil
func (bl *Block)Encoded() common.Hash{
	bl.Sig = nil
	return utils.RlpHash(bl)
}

func (bl *Block)Sign(key *ecdsa.PrivateKey){
	sig, err := utils.SignData(bl.Hash().Bytes(), key)
	if err != nil {
		fmt.Println("cannot be sign!x2")
	}
	bl.Sig = sig
}

func (bl *Block)AddTransaction(tx Transaction) {
	bl.TransactionSet = append(bl.TransactionSet, tx)
	inputs := [][]interface{}{{tx.blknum1, tx.txindex1, tx.oindex1}, {tx.blknum2, tx.txindex2, tx.oindex2}}
	for _, i := range inputs {
		inputId := utils.EncodeUTXOID(i[0].(*big.Int), i[1].(int64), i[2].(int64))
		bl.SpentUTXOs[inputId] = true
	}
}

