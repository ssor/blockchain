package corelib

import "testing"

func TestNewWallet(t *testing.T) {
    w, err := NewWallet()
    if err != nil {
        t.Fatal(err)
    }
    err = w.SaveToFile()
    if err != nil {
        t.Fatal(err)
    }
}
