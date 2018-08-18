package client

import (
	"github.com/regcostajr/go-web3/providers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hot3246624/StarsChain/plasma_core"
	"math/big"
	"crypto/ecdsa"
	"github.com/hot3246624/StarsChain/plasma/child_chain"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/fabric/core/cclifecycle"
)

type Client struct {
	rootChainProvide providers.HTTPProvider
	childChainUrl string
	rootChain common.Address
	childChain ChildChainService
}

func NewClient(rootChainProvider providers.HTTPProvider, childChainUrl string) Client {
	//TODO deploy the rootchain
	client := Client{childChain:childChainUrl}
	client.rootChainProvide = providers.HTTPProvider{address:"http://localhost:8545"}
	if childChainUrl == "" {
		childChainUrl = "http://localhost:8546/jsonrpc"
	}
	client.childChain = ChildChainService{url: childChainUrl}
	return client
}

func (cli *Client)createTransaction(blknum1 big.Int, txindex1 int64, oindex1 int64, blknum2 big.Int, txindex2 int64, oindex2 int64, newowner1 common.Address, amount1 big.Int, newowner2 common.Address, amount2 big.Int, cur12 common.Address, fee big.Int) plasma_core.Transaction{
	return plasma_core.Transaction{blknum1:blknum1, txindex1:txindex1, oindex1:oindex1,
		blknum2:blknum2, txindex2:txindex2, oindex2:oindex2,
		cur12:cur12,
		newowner1:newowner1, amount1:amount1,
		newowner2:newowner2, amount2:amount2,
		}
}

func (cli *Client)signTransaction(transaction plasma_core.Transaction, key1 *ecdsa.PrivateKey, key2 *ecdsa.PrivateKey) plasma_core.Transaction{
	if key1.D.Int64() == 0{
		transaction.Sign1(key2)
	}
	if key2.D.Int64() == 0 {
		transaction.Sign1(key1)
	}
	return transaction
}

func (cli *Client)Deposit(amount big.Int, owner common.Address) {
	//TODO deposit
	cli.rootChain.deposit()
}

func (cli *Client)applyTransaction(transaction plasma_core.Transaction){
	cli.childChain.ApplyTransaction(transaction)
}

func (cli *Client)submitBlock(block plasma_core.Block){
	cli.childChain.SubmitBlock(block)
}

func (cli *Client)withdraw(blknum big.Int, txindex int64, oindex int64, tx plasma_core.Transaction, proof []byte, sigs []byte){
	utxoPos := blknum.Int64() * 1000000000 + txindex * 10000 + oindex * 1
	encodedTransaction, err := encodeToString(tx)
	if err != nil {

	}
	//TODO
	cli.rootChain.startExit(utxoPos, encodedTransaction, proof, sigs, transact={'from': '0x' + tx.newowner1.hex()})

}

func (cli *Client)withdrawDeposit(owner common.Address, depositPos []byte, amount big.Int){
	//TODO
	cli.rootChain.startDepositExit(deposit_pos, amount, transact={'from': owner})
}

func (cli *Client)getTransaction(blknum big.Int, txindex int64) plasma_core.Transaction{
	//TODO unarchive transaction from RLP
	encodedTransaction, err := cli.childChain.GetTransaction(blknum, txindex)
	if err != nil {

	}
	return rlp.DecodeBytes(utils.decode_hex(encoded_transaction), Transaction)
}

func (cli *Client)getCurrentBlock(){
	encodedBlock, err := cli.childChain.GetCurrentBlock()
	if err != nil {

	}
	//TODO what is usage of rlp
	return rlp.decode(utils.decode_hex(encoded_block), Block)
}

func (cli *Client)getBlock(blknum big.Int){
	encodedBlock, err := cli.childChain.GetBlock(blknum)
	if err := nil{

	}
	return rlp.decode(utils.decode_hex(encodedBlock), Block)
}

func (cli *Client)getCurrentBlockNum() big.Int{
	//TODO why we need simplejson, compare to rlp
	num, err := cli.childChain.GetCurrentBlockNum()
	if err != nil {

	}
	return num["num"]
}
