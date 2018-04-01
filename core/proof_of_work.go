package core

import (
	"math/big"
	"crypto/sha256"
	"fmt"
)

const (
	maxNonce = 1000 * 1000 * 1000
)

type ProofOfWork struct {
	block      *Block
	target     *big.Int
	targetBits uint
}

func NewProofOfWork(block *Block, targetBits uint) *ProofOfWork {
	pw := &ProofOfWork{
		block:      block,
		targetBits: targetBits,
	}
	pw.target = GetTarget(targetBits)
	return pw
}


func (pow *ProofOfWork) Hash(nonce int64) []byte {
	data := pow.block.prepareHashData(pow.targetBits, nonce)
	hash := sha256.Sum256(data)
	fmt.Printf("nonce: %d  Hash: %x\n", nonce, hash)
	return hash[:]
}

func (pow *ProofOfWork) Mine() (int64, []byte) {
	var nonce int64 = 0
	//var hash []byte
	var hashInt big.Int

	for nonce < maxNonce {
		hashTemp := pow.Hash(nonce)
		hashInt.SetBytes(hashTemp[:])
		if hashInt.Cmp(pow.target) < 0 {
			return nonce, hashTemp[:]
		} else {
			nonce++
		}
	}
	return -1, nil
}
