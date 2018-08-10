package plasma_core

import (
	//"log"
	"github.com/sirupsen/logrus"
	"math/big"
	"github.com/hot3246624/StarsChain/plasma_core/utils"
	"errors"
	"encoding/hex"
)

type Chain struct {
	Operator interface{}
	Blocks []Block
	ParentQueue map[int64][]Block
	ChildBlockInterval int64
	NextChildBlock int64
	NextDepositBlock int64
}

var (
	isNextChildBlock bool
)

func (ch *Chain)AddBlock(bl Block) bool{
	isNextChildBlock = bl.Number.Int64() == ch.NextChildBlock
	if isNextChildBlock || bl.Number.Int64() == ch.NextDepositBlock {
		ch.validateBlock(bl)
		ch.applyBlock(bl)
		if isNextChildBlock {
			ch.NextDepositBlock = ch.NextChildBlock + 1
			ch.NextChildBlock += ch.ChildBlockInterval
		}else{
			ch.NextDepositBlock += 1
		}
	}else if bl.Number.Int64() > ch.NextDepositBlock {
		parentBlockNumber := bl.Number.Int64() - 1
		if _, ok := ch.ParentQueue[parentBlockNumber]; !ok {
			ch.ParentQueue[parentBlockNumber] = []Block{}
		}
		ch.ParentQueue[parentBlockNumber] = append(ch.ParentQueue[parentBlockNumber], bl)
		return false
	}else{
		return false
	}
	if _, ok := ch.ParentQueue[bl.Number.Int64()]; ok {
		for _, blk := range ch.ParentQueue[bl.Number.Int64()] {
			ch.AddBlock(blk)
		}
		delete(ch.ParentQueue, bl.Number.Int64())
	}
	return true
}

func (ch *Chain)validateBlock(blk Block) error{
	if !blk.IsDepositBlock() && (hex.EncodeToString(blk.Sig) == NULL_SIGNATURE || utils.AddressToHex(blk.Singer()) != ch.Operator){
		logrus.Print("the block is invalid!")
		return errors.New("the block is invalid!")
	}
	for _, tx := range blk.TransactionSet {
		ch.validateTransaction(tx)
	}
	return nil
}

func (ch *Chain)applyBlock(blk Block){
	for _, tx := range blk.TransactionSet {
		ch.applyTransaction(tx)
	}
	ch.Blocks[blk.Number.Int64()] = blk
}

func (ch *Chain)applyTransaction(tx Transaction){
	inputs := [][]interface{}{{tx.blknum1, tx.txindex1, tx.oindex1}, {tx.blknum2, tx.txindex2, tx.oindex2}}
	for _, i := range inputs{
		blkNum := i[0]
		if blkNum.(big.Int).Int64() == 0 {
			continue
		}
		inputId := utils.EncodeUTXOID(i[0].(*big.Int), i[1].(int64), i[2].(int64))
		ch.MarkUXTOSpent(inputId)
	}
}

func (ch *Chain)MarkUXTOSpent(utxoID int64) {
	_, _, oindex := utils.DecodeUTXOID(utxoID)
	tx := ch.getTransaction(utxoID)
	if oindex == 0 {
		tx.spent1 = true
	}else{
		tx.spent2 = true
	}
}

func (ch *Chain)GetBlock(blknum big.Int) Block{
	for i, blk := range ch.Blocks{
		if blk.Number.Int64() == blknum.Int64(){
			return ch.Blocks[i]
		}
	}
	return Block{}
}

func (ch *Chain)GetTransaction(transactionID int64) Transaction{
	blknum, txindex, _ := utils.DecodeUTXOID(transactionID)
	return ch.Blocks[blknum.Int64()].TransactionSet[txindex]
}

func (ch *Chain)ValidateTransaction(tx Transaction) error{
	inputAmount := int64(0)
	outputAmount := tx.amount1.Int64() + tx.amount2.Int64()
	var validSignature  = false
	var spent  = false
	tempSpent := make(map[int64]Transaction)
	inputs := [][]interface{}{{tx.blknum1, tx.txindex1, tx.oindex1}, {tx.blknum2, tx.txindex2, tx.oindex2}}
	for _, i := range inputs {
		blkNum, txindex, oindex := i[0].(big.Int), i[1].(int64), i[2].(int64)
		if blkNum.Int64() == 0 {
			continue
		}

		inputTX := ch.Blocks[blkNum.Int64()].TransactionSet[txindex]

		if oindex == 0 {
			validSignature = hex.EncodeToString(tx.sig1) != NULL_SIGNATURE && inputTX.newowner1 == tx.Sender1()
			spent = inputTX.spent1
			inputAmount = inputTX.amount1.Int64()
		} else {
			validSignature = hex.EncodeToString(tx.sig2) != NULL_SIGNATURE && inputTX.newowner2 == tx.Sender2()
			spent = inputTX.spent2
			inputAmount = inputTX.amount2.Int64()
		}
		utxoID := utils.EncodeUTXOID(&blkNum,txindex,oindex)
		if _, ok := tempSpent[utxoID]; !ok || spent {
			return errors.New("failed to validate tx")
		}

		if !validSignature {
			return errors.New("failed to validate tx")
		}
	}

	if !tx.IsDepositTransaction() && inputAmount < outputAmount {
		return errors.New("failed to validate tx")
	}
	return nil
}