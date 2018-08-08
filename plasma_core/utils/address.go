package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"bytes"
)

func AddressToHex(address common.Address) string{
	return "0x" + address.Hex()
}

func AddressToByte(address []byte) []byte{
	return address[2:]
}