package child_chain

import (
	"net/rpc"
	"github.com/hot3246624/StarsChain/plasma_core"
	"github.com/ethereum/go-ethereum/common"
	"net/rpc/jsonrpc"
	"fmt"
	"os"
	"net"
)

//TODO transfrom deployer
var rootChain = make([]byte, 1)
var childChain = ChildChain{operator:common.HexToAddress(plasma_core.AUTHORITY["address"][2:]), rootChain: rootChain }


type RPCResponse int

//TODO temporary way
func (rpc *RPCResponse)SubmitBlock(block plasma_core.Block, reply *int) error{
	childChain.submitBlock(block)
	*reply = 1
	return nil
}

func (rpc *RPCResponse)ApplyTransaction(block plasma_core.Block, reply *int) error{
	childChain.submitBlock(block)
	*reply = 1
	return nil
}

func (rpc *RPCResponse)GetTransaction(block plasma_core.Block, reply *int) error{
	childChain.submitBlock(block)
	*reply = 1
	return nil
}

func (rpc *RPCResponse)GetCurrentBlock(block plasma_core.Block, reply *int) error{
	childChain.submitBlock(block)
	*reply = 1
	return nil
}

func (rpc *RPCResponse)GetCurrentBlockNum(block plasma_core.Block, reply *int) error{
	childChain.submitBlock(block)
	*reply = 1
	return nil
}

func (rpc *RPCResponse)GetBlock(block plasma_core.Block, reply *int) error{
	childChain.submitBlock(block)
	*reply = 1
	return nil
}


func Application(portAddr string) {
	rpcRes := new(RPCResponse)
	rpc.Register(rpcRes)
	tcpAddr, err := net.ResolveTCPAddr("tcp", portAddr)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		jsonrpc.ServeConn(conn)
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
