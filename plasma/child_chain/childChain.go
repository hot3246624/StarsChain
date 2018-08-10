package child_chain

import (
	"github.com/hot3246624/StarsChain/plasma_core"
	"math/big"
	"github.com/hot3246624/StarsChain/plasma_core/utils"
	"github.com/ethereum/go-ethereum/common"
)

type ChildChain struct {
	operator common.Address
	rootChain []byte
	chain plasma_core.Chain
	currentBlock plasma_core.Block
	eventLister RootEventListener

}

func NewCC(operator common.Address, rootChain []byte) *ChildChain{
	cc := &ChildChain{operator: operator, rootChain:rootChain}
	cc.chain = plasma_core.Chain{Operator:operator}
	cc.currentBlock = plasma_core.Block{Number:big.NewInt(cc.chain.NextChildBlock)}
	cc.eventLister = RootEventListener{rootChain:rootChain, confirmations:0}
	cc.eventLister.on("Deposit", cc.applyDeposit)
	cc.eventLister.on("ExitStarted", cc.applyExit)
}

func (cc *ChildChain)applyExit(event map[string]map[string]interface{}){
	eventArgs := event["arg"]
	utxoID := eventArgs["utxoPos"]
	cc.chain.MarkUXTOSpent(utxoID.(int64))
}


func (cc *ChildChain)applyDeposit(event map[string]map[string]interface{}){
	eventArgs := event["arg"]
	owner := eventArgs["depositor"]
	amount := eventArgs["amount"]
	blknum := eventArgs["depositBlock"]

	depositTX := utils.GetDepositTX(owner.(common.Address), amount.(big.Int))
	depositBlock := plasma_core.Block{TransactionSet:[]plasma_core.Transaction{*depositTX}, Number:big.NewInt(blknum.(int64))}
	cc.chain.AddBlock(depositBlock)
}

func (cc *ChildChain)applyTransaction(tx plasma_core.Transaction) (int64, error){
	err := cc.chain.ValidateTransaction(tx)
	if err != nil {
		return 0, err
	}
	cc.currentBlock.AddTransaction(tx)
	return utils.EncodeUTXOID(cc.currentBlock.Number, int64(len(cc.currentBlock.TransactionSet) - 1), 0), nil
}

func (cc *ChildChain)submitBlock(block plasma_core.Block){
	cc.chain.AddBlock(block)
	//TODO for the different between web3.js with go-web3
	cc.rootChain.tranact()
	cc.currentBlock = plasma_core.Block{Number:big.NewInt(cc.chain.NextChildBlock)}
}

func (cc *ChildChain)getTransaction(txID int64) plasma_core.Transaction{
	return cc.chain.GetTransaction(txID)
}

func (cc *ChildChain)getBlock(blknum big.Int) plasma_core.Block{
	return cc.chain.GetBlock(blknum)
}

func (cc *ChildChain)getCurrentBlock() plasma_core.Block{
	return cc.currentBlock
}