package db

import (
    "fmt"
    "github.com/boltdb/bolt"
)

var (
    LocalDB *bolt.DB
)

func Init(args ...string) error {
    if args == nil {
        return fmt.Errorf("no args")
    }
    if len(args) > 0 {
        dbPath := args[0]
        db, err := initLocalDB(dbPath)
        if err != nil {
            return err
        }
        LocalDB = db
    }
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
