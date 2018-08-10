package utils

import (
	"github.com/hot3246624/StarsChain/plasma_core/transaction"
	"math/big"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hot3246624/StarsChain/plasma_core"
)

const BLKNUM_OFFSET  = 1000000000
const TXINDEX_OFFSET  = 10000

func EncodeUTXOID(blknum *big.Int, txinde, oindex int64) int64{
	return blknum.Int64()*BLKNUM_OFFSET + txinde * TXINDEX_OFFSET + oindex
}

func DecodeUTXOID(utxoID int64) (blknum *big.Int, txindex, oindex int64) {
	blknum =  big.NewInt(int64(utxoID / BLKNUM_OFFSET))
	txindex = (utxoID % BLKNUM_OFFSET) / BLKNUM_OFFSET
	oindex = utxoID - blknum.Int64() * BLKNUM_OFFSET - txindex * TXINDEX_OFFSET
	return blknum, txindex, oindex
}

func GetDepositTX(owner common.Address, amount big.Int) *plasma_core.Transaction{
	return &plasma_core.Transaction{blknum1:0,txindex1:0,oindex1:0,blknum2:0,txindex2:0,oindex2:0,cur12:plasma_core.NULL_ADDRESS,newowner1:owner,amout1:amount,newowner2:plasma_core.NULL_ADDRESS,amount2:0}
}


