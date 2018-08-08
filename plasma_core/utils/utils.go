package utils

import (
	"github.com/hot3246624/StarsChain/plasma_core/utils/merkle"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"github.com/hot3246624/StarsChain/plasma_core"
	"github.com/ethereum/go-ethereum/crypto"
	"strconv"
	"crypto/ecdsa"
)

func GetEmptyMerkleTreeHash(depth int) []byte{
	zerosHash := []byte(plasma_core.NULL_HASH)
	for i := 0; i < depth; i ++ {
		zerosHash = crypto.Keccak256(zerosHash, zerosHash)
	}
	return zerosHash
}

func GetMerkleOfLeaves(depth int, leaves [][]byte) merkle.FixedMerkle {
	return merkle.FixedMerkle{Depth:depth, Leaves: leaves}
}

//TODO：inp是什么类型待定
//int/float类型转变为[]byte, 直接先转成str，比如通过FormatInt/Float等;而[]byte转化为字符串推荐用hex.encodeToString.
func BytesFillLeft(inp []byte, length int) []byte {
	return append([]byte(strconv.FormatInt(int64(length - len(inp)), 10)), inp...)
}


func GetDepositHash(owner common.Address, token []byte, value big.Int) []byte {
	return crypto.Keccak256(owner.Bytes(), token, []byte(plasma_core.CreateRepeatBytes("\x00", 31)), []byte(value.String()))
}

func ConfirmTX(tx plasma_core.Transaction, root []byte, key *ecdsa.PrivateKey) ([]byte, error){
	sig, err := crypto.Sign(crypto.Keccak256(tx.Hash().Bytes(), root), key)
	if err != nil {
		return nil, err
	}
	return sig, nil
}