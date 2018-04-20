package corelib

import (
	"math/big"
	"encoding/binary"
)

// larger targetBits means larger int
func GetTarget(targetBits uint) *big.Int {
	target := big.NewInt(1)
	target.Lsh(target, targetBits)
	return target
}


func IntToHex(i int64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buf, i)
	return buf[:n]
}

