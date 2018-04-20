package corelib

import (
	"time"
	"bytes"
	"crypto/sha256"
	"math/big"
	"encoding/gob"
	"fmt"
	log "github.com/sirupsen/logrus"
)

const (
	genesisBlockData = "genesis block"
)

var (
	blockLogger = log.WithFields(log.Fields{
		"file_name": "block",
	})
)

type Block struct {
	Timestamp    int64
	Transactions []*Transaction
	Hash         []byte
	PreBlockHash []byte
	Nonce        int64
}

func NewBlock(preBlockHash []byte, transactions []*Transaction) *Block {
	b := &Block{
		Timestamp:    time.Now().Unix(),
		PreBlockHash: preBlockHash,
		Transactions: transactions,
	}
	return b
}

func NewGenesisBlock(coinbase *Transaction) *Block {
	b := NewBlock([]byte{}, []*Transaction{coinbase})
	b.Hash = []byte(genesisBlockData)
	//b.GenerateHash()
	return b
}

func (b *Block) serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		blockLogger.Errorf("### encode error: %s\n", err)
		return nil, fmt.Errorf("cannot serialize block")
	}
	return result.Bytes(), nil
}

func (b *Block) isGenesisBlock() bool {
	if string(b.Hash) == genesisBlockData {
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
		blockLogger.Errorf("### decode error: %s\n", err)
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

func (b *Block) hashTransactions() []byte {
	buf := bytes.NewBuffer([]byte{})
	for _, tx := range b.Transactions {
		buf.Write(tx.ID)
	}
	hash := sha256.Sum256(buf.Bytes())
	return hash[:]
}

func (b *Block) prepareHashData(targetBits uint, nonce int64) []byte {
	data := bytes.Join([][]byte{
		b.PreBlockHash, b.hashTransactions(), IntToHex(int64(targetBits)), IntToHex(nonce),
	}, []byte{})
	return data[:]
}

//
//func (b *Block) GenerateHash() {
//	timestamp := strconv.FormatInt(b.Timestamp, 10)
//	headersRaw := bytes.Join([][]byte{[]byte(timestamp), b.Data, b.PreBlockHash}, []byte{})
//	result := sha256.Sum256(headersRaw)
//	b.Hash = result[:]
//}
