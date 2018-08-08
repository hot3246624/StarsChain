package plasma_core

import (
	"math/big"
	"github.com/ethereum/go-ethereum/common"
	"crypto/ecdsa"
	"github.com/sirupsen/logrus"
	"github.com/hot3246624/StarsChain/plasma_core/utils"
)

type Transaction struct {
	//blocknumber 1
	blknum1 big.Int
	txindex1 int64
	oindex1 int64
	blknum2 big.Int
	txindex2 int64
	oindex2 int64
	//token address
	cur12 common.Address
	newowner1 common.Address
	amount1 big.Int
	newowner2 common.Address
	amount2 big.Int
	sig1 []byte
	sig2 []byte
	spent1 bool
	spent2 bool
}

func (tx *Transaction)Hash() common.Hash {
	return utils.RlpHash(tx)
}

func (tx *Transaction)MerkleHash() common.Hash {
	hash := tx.Hash()
	temp := append(hash.Bytes(), tx.sig1...)
	temp = append(temp, tx.sig2...)
	return utils.RlpHash(temp)
}

func (tx *Transaction)IsSingleUTXO() bool {
	return tx.blknum2.Int64() == 0
}

func (tx *Transaction)IsDepositTransaction() bool{
	return tx.blknum1.Int64() == 0 && tx.blknum2.Int64() == 0
}

func (tx *Transaction)Sender1() common.Address{
	return utils.GetSigner(tx.Hash().Bytes(), tx.sig1)
}

func (tx *Transaction)Sender2() common.Address{
	return utils.GetSigner(tx.Hash().Bytes(), tx.sig2)
}

func (tx *Transaction)Encoded() common.Hash{
	tx.sig1, tx.sig2 = nil, nil
	return utils.RlpHash(tx)
}

func (tx *Transaction)Sign1(key *ecdsa.PrivateKey){
	sig1, err := utils.SignData(tx.Hash().Bytes(), key)
	if err != nil {
		logrus.Error("sign1 failed")
		return
	}
	tx.sig1 = sig1
}

func (tx *Transaction)Sign2(key *ecdsa.PrivateKey){
	sig2, err := utils.SignData(tx.Hash().Bytes(), key)
	if err != nil {
		logrus.Error("sign2 failed")
		return
	}
	tx.sig2 = sig2
}