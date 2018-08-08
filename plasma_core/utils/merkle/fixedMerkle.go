package merkle

import (
	//"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
	"errors"
	"github.com/hot3246624/StarsChain/plasma_core"
	"math"
	"encoding/hex"
)

type FixedMerkle struct {
	Depth int
	LeafCount int
	Hashed bool
	Leaves [][]byte
	Tree [][]Node
	Root []byte
}

func (fm *FixedMerkle)createNodes(leaves [][]byte) (nodes []Node) {
	for _, leaf := range leaves{
		nodes = append(nodes, Node{data:leaf})
	}
	return
}

//TODO 此处叶子节点集合暂定为leaf
func (fm *FixedMerkle)createTree(leaves []Node) {
	if len(leaves) == 1 {
		fm.Root = leaves[0].data
	}
	nextLevel := len(leaves)
	treeLevel := make([]Node, nextLevel)
	for i := 0; i < nextLevel; i += 2{
		combined := crypto.Keccak256(leaves[i].data, leaves[i+1].data)
		nextNode := Node{combined,NewNode(leaves[i].data),NewNode(leaves[i+1].data)}
		treeLevel = append(treeLevel, nextNode)
	}
	fm.Tree = append(fm.Tree, treeLevel)
	fm.createTree(treeLevel)
}

//根据之前的方法，leaves中应该存放的是node，这里暂定选择其data部分或者rlp之后的部分
func (fm *FixedMerkle)checkMembership(leaf []byte, index uint, proof []byte) bool{
	if !fm.Hashed {
		leaf = crypto.Keccak256(leaf)
	}
	computedHash := leaf
	for i := 0; i < int(fm.Depth) * 32; i += 32 {
		segment := proof[i:i+32]
		if index % 2 == 0 {
			computedHash = crypto.Keccak256(computedHash, segment)
		}else {
			computedHash = crypto.Keccak256(segment, computedHash)
		}
		index /= 2
	}
	return hex.EncodeToString(computedHash) == hex.EncodeToString(fm.Root)
}

func (fm *FixedMerkle)createMembershipProof(leaf []byte) ([]byte, error){
	if !fm.Hashed {
		leaf = crypto.Keccak256(leaf)
	}
	if !fm.isMember(leaf) {
		return nil, &MemberNotExistException{"leaf is not in the merkle tree!"}
	}
	id := 0
	for index, value := range fm.Leaves {
		if hex.EncodeToString(value) == hex.EncodeToString(leaf) {
			id = index
		}
	}
	var siblingIndex int
	var proof []byte
	for i := 0; i < int(fm.Depth); i ++ {
		if id % 2 == 0 {
			siblingIndex = id + 1
		}else {
			siblingIndex = id - 1
		}
		id /= 2
		nodes := fm.Tree[i]
		proof = append(proof, nodes[siblingIndex].data...)
	}
	return proof, nil
}

func (fm *FixedMerkle)isMember(leaf []byte) bool {
	for _, value := range fm.Leaves {
		if hex.EncodeToString(value) == hex.EncodeToString(leaf){
			return true
		}
	}
	return false
}

func (fm *FixedMerkle)isNotMember(leaf []byte) bool {
	for _, value := range fm.Leaves {
		if hex.EncodeToString(value) == hex.EncodeToString(leaf){
			return false
		}
	}
	return true
}

func NewFixedMerkle(depth int, leaves [][]byte, hashed bool) (*FixedMerkle, error) {
	if depth < 1{
		logrus.Error("depth should be at least 1")
		return nil, errors.New("depth should be at least 1")
	}

	fixedMerkle := &FixedMerkle{Depth:depth, Leaves:leaves, Hashed:hashed}
	fixedMerkle.LeafCount = int(math.Pow(2, float64(depth)))
	if len(leaves) > fixedMerkle.LeafCount {
		return nil, errors.New("num of leaves exceed max avaiable num with the depth")
	}

	if !hashed {
		var tempLeaves [][]byte
		for _, leaf := range leaves {
			tempLeaves = append(tempLeaves, crypto.Keccak256(leaf))
		}
		leaves = tempLeaves
	}
	for i := 0; i < 2 * int(depth) - len(leaves); i++ {
		leaves = append(leaves, []byte(plasma_core.NULL_HASH))
	}
	fixedMerkle.Leaves = leaves
	fixedMerkle.Tree = append(fixedMerkle.Tree, fixedMerkle.createNodes(leaves))
	fixedMerkle.createTree(fixedMerkle.Tree[0])
	return fixedMerkle, nil
}