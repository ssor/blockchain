package core

import (
	"time"
	"bytes"
	"crypto/sha256"
	"math/big"
	"strconv"
	"encoding/gob"
	"fmt"
)

const (
	genesisBlockData = "genesis block"
)

type Block struct {
	Timestamp    int64
	Data         []byte
	Hash         []byte
	PreBlockHash []byte
	Nonce        int64
}

func NewBlock(preBlockHash, data []byte) *Block {
	b := &Block{
		Timestamp:    time.Now().Unix(),
		PreBlockHash: preBlockHash,
		Data:         data,
	}
	return b
}

func NewGenesisBlock() *Block {
	b := NewBlock([]byte{}, []byte(genesisBlockData))
	b.GenerateHash()
	return b
}

func (b *Block) serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		fmt.Printf("### encode error: %s\n", err)
		return nil, fmt.Errorf("cannot serialize block")
	}
	return result.Bytes(), nil
}

func (b *Block) isGenesisBlock() bool {
	if string(b.Data) == genesisBlockData {
		return true
	} else {
		return false
	}
}

func deserializeBlock(data []byte) (*Block, error) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	var b Block
	err := decoder.Decode(&b)
	if err != nil {
		fmt.Printf("### decode error: %s\n", err)
		return nil, fmt.Errorf("cannot decode block")
	} else {
		return &b, nil
	}
}

func (b *Block) validate(targetBits uint) bool {
	if b.isGenesisBlock() {
		return true
	}
	data := b.prepareHashData(targetBits, b.Nonce)
	hash := sha256.Sum256(data)

	var hashInt big.Int
	hashInt.SetBytes(hash[:])
	if hashInt.Cmp(GetTarget(targetBits)) < 0 {
		return true
	} else {
		return false
	}
}

func (b *Block) prepareHashData(targetBits uint, nonce int64) []byte {
	data := bytes.Join([][]byte{
		b.PreBlockHash, b.Data, IntToHex(int64(targetBits)), IntToHex(nonce),
	}, []byte{})
	return data[:]
}

func (b *Block) GenerateHash() {
	timestamp := strconv.FormatInt(b.Timestamp, 10)
	headersRaw := bytes.Join([][]byte{[]byte(timestamp), b.Data, b.PreBlockHash}, []byte{})
	result := sha256.Sum256(headersRaw)
	b.Hash = result[:]
}
