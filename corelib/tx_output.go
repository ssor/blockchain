package corelib

import (
	"bytes"
	"crypto/sha256"
)

type TXOutput struct {
	Value        int    //0.00000001 BTC
	ScriptPubKey string // used to lock value to address(someone)
}

func (out *TXOutput) canBeUnlockedWith(unlockData string) bool {
	return out.ScriptPubKey == unlockData
}

func (out *TXOutput) Hash() []byte {
	data := bytes.Join([][]byte{
		IntToHex(int64(out.Value)), []byte(out.ScriptPubKey),
	}, []byte{})
	hash := sha256.Sum256(data)
	return hash[:]
}

func NewTxOutput(v int, key string) *TXOutput {
	output := &TXOutput{
		Value:        v,
		ScriptPubKey: key,
	}
	return output
}
