package db

import (
    "fmt"
    "github.com/boltdb/bolt"
    "github.com/sirupsen/logrus"
)

const (
    BucketWallet       = "wallet"
    BucketBlocks       = "blocks"
    BucketLastBlockKey = "l"
    //GenesisBlockHash   = "genesis"
)

var (
    localDB *bolt.DB
)

func GetDb() *bolt.DB {
    return localDB
}

func Init(args ...string) error {
    if args == nil {
        return fmt.Errorf("init db failed, no args")
    }
    if len(args) > 0 {
        dbPath := args[0]
        db, err := initLocalDB(dbPath)
        if err != nil {
            logrus.Error("init db error: ", err)
            return err
        }
        localDB = db
    }
    logrus.Info("init db success")
    return nil
}

func initLocalDB(dbPath string) (*bolt.DB, error) {
    db, err := bolt.Open(dbPath, 0600, nil)
    if err != nil {
        fmt.Printf("open db error: %s \n", err)
        return nil, err
    }
    return db, nil
}
