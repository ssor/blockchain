package core

import (
	"fmt"
	"bytes"
	"crypto/sha256"
)

const (
	/*
	subsidy is the amount of reward.
	In Bitcoin, this number is not stored anywhere and calculated based only on the total number of blocks:
	the number of blocks is divided by 210000.
	Mining the genesis block produced 50 BTC, and every 210000 blocks the reward is halved.
	*/
	subsidy = 50
)

type Transaction struct {
	ID   []byte
	Vin  []*TXInput
	Vout []*TXOutput
}

func (tx *Transaction) setID() {
	buf := bytes.NewBuffer([]byte{})
	for _, in := range tx.Vin {
		buf.Write(in.Hash())
	}
	for _, out := range tx.Vout {
		buf.Write(out.Hash())
	}
	hash := sha256.Sum256(buf.Bytes())
	tx.ID = hash[:]
}

func NewCoinbaseTX(to, data string) *Transaction {
	if len(data) <= 0 {
		data = fmt.Sprintf("reward to %s", to)
	}
	txIn := NewEmptyInput(data)
	txOut := NewTxOutput(subsidy, to)
	tx := &Transaction{
		Vin:  []*TXInput{txIn},
		Vout: []*TXOutput{txOut},
	}
	tx.setID()
	return tx
}
