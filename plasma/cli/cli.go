package cli

import (
	"github.com/hot3246624/StarsChain/plasma/client"
	"fmt"
	"math/big"
	"github.com/ethereum/go-ethereum/common"
	"crypto/ecdsa"
)

var ContextSetting = map[string]interface{}{"help_option_names":[...]string{"-h","--help"}}


func cli(ctx interface{}) {
	ctx.obj = client.NewClient()
}

func clientCall(fn func([]interface{}) (interface{}, error), args []interface{}, successmessage string) interface{}{
	output,err := fn(args)
	if err != nil {
		//TODO
	}
	if successmessage != ""{
		fmt.Println(successmessage)
	}
	return output
}

func deposit(client client.Client, amount big.Int, address common.Address) {
	client.Deposit(amount,address)
	fmt.Printf("Deposit %v to %v", amount, address)
}

func sendtx(){
	//TODO
}

func submitBlock(client client.Client, key *ecdsa.PrivateKey) {
	//TODO
}

func withDraw(client client.Client, blknum big.Int, txindex uint64, oindex uint64, key1 *ecdsa.PrivateKey, key2 *ecdsa.PrivateKey){
	//TODO
}

func withdrawdeposit(client client.Client, owner common.Address, blknum big.Int, amount big.Int) {
	//TODO
}