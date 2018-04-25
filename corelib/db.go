package corelib

import "github.com/boltdb/bolt"

type Store interface {
    Put([]byte, []byte) error
    Get([]byte) []byte
}

type boltDB = *bolt.DB
type boltTx = *bolt.Tx
