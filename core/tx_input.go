package core

import (
	"bytes"
	"crypto/sha256"
)

type TXInput struct {
	TxId      []byte // output of previous tx
	Vout      int
	ScriptSig string // signature of belonging
}

func (in *TXInput) canBeUnlockedWith(unlockData string) bool {
	return in.ScriptSig == unlockData
}

func (in *TXInput) Hash() []byte {
	data := bytes.Join([][]byte{
		in.TxId, IntToHex(int64(in.Vout)), []byte(in.ScriptSig),
	}, []byte{})
	hash := sha256.Sum256(data)
	return hash[:]
}

func NewTxInput(id []byte, index int, sig string) *TXInput {
	input := &TXInput{
		TxId:      id,
		Vout:      index,
		ScriptSig: sig,
	}
	return input
}

func NewEmptyInput(sig string) *TXInput {
	return NewTxInput([]byte{}, -1, sig)
}
