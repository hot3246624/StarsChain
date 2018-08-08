package merkle

//import "github.com/ethereum/go-ethereum/common"

type Node struct {
	data []byte
	left *Node
	right *Node
}

func NewNode(data []byte) *Node{
	return &Node{data,nil,nil}
}
