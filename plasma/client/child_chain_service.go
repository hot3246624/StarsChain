package client

import (
	"github.com/hot3246624/StarsChain/plasma/child_chain"
	"net/http"
	"github.com/bitly/go-simplejson"
	"strings"
	"io/ioutil"
	"github.com/hot3246624/StarsChain/plasma_core"
	"math/big"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"encoding/hex"
)

type ChildChainService struct {
	url string
	methods []interface{}
	child_chain.ChildChain
}

var errorJson,_ = simplejson.NewJson([]byte(`"result":"-1"`))

func NewCCS(url string) *ChildChainService {
	ccs := &ChildChainService{url:url, methods:nil, ChildChain: child_chain.ChildChain{}}

	//ccs.methods
	return ccs
}

func (ccs *ChildChainService)sendRequest(method string, args []interface{}) (*simplejson.Json, error){
	payload := fmt.Sprintf(`{
		"data":{
		"method": %s,
		"params": %v
		"jsonrpc": "2.0",
		"id": 0,
		}
	}`, method, args)
	resp, err := http.Post(ccs.url,"json", strings.NewReader(payload))

	if err != nil{
		return errorJson, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return errorJson, err
	}
	js, err := simplejson.NewJson(body)
	result := js.Get("result")
	return result, nil
}

func (ccs *ChildChainService)ApplyTransaction(transaction plasma_core.Transaction) (*simplejson.Json, error) {
	dataStr, err := encodeToString(transaction)
	if err != nil {
		return errorJson, err
	}
	return ccs.sendRequest("apply_transaction", []interface{}{dataStr})
}

func (ccs *ChildChainService)SubmitBlock(block plasma_core.Block) (*simplejson.Json, error){
	dataStr, err := encodeToString(block)
	if err != nil {
		return errorJson, err
	}
	return ccs.sendRequest("submit_block", []interface{}{dataStr})
}

func (ccs *ChildChainService)GetTransaction(blknum big.Int, txindex int) (*simplejson.Json, error){
	return ccs.sendRequest("get_transaction", []interface{}{blknum, txindex})
}

func (ccs *ChildChainService)GetCurrentBlock() (*simplejson.Json, error){
	return ccs.sendRequest("get_current_block", []interface{}{})
}

func (ccs *ChildChainService)GetBlock(blknum big.Int) (*simplejson.Json, error){
	return ccs.sendRequest("get_block", []interface{}{blknum})
}

func (ccs *ChildChainService)GetCurrentBlockNum() (*simplejson.Json, error){
	return ccs.sendRequest("get_current_block_num", []interface{}{})
}

func encodeToString(data interface{}) (string, error) {
	content, err := rlp.EncodeToBytes(data)
	if err != nil {
		return "", err
	}
	contentStr := hex.EncodeToString(content)
	return contentStr, nil
}