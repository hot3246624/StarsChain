package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"crypto/ecdsa"
	"github.com/sirupsen/logrus"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
)

func SigToPub(hash []byte, sig []byte) (pubkey *ecdsa.PublicKey) {
	pubkey, err := crypto.SigToPub(hash, sig)
	if err != nil {
		logrus.Error("cannot get pubkey!")
		return nil
	}
	return pubkey
}

func SignData(data []byte, key *ecdsa.PrivateKey) ([]byte, error){
	var sig []byte
	sig, err := crypto.Sign(data, key)
	if err != nil {
		logrus.Error("cannot be sign!")
		return nil, err
	}
	return sig, nil
}

func RlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

func GetSigner(hash []byte, sig []byte) common.Address{
	pubkey := SigToPub(hash, sig)
	return crypto.PubkeyToAddress(*pubkey)
}