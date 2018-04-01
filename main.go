package main

import (
	"fmt"
	"blockchain/core"
	"github.com/boltdb/bolt"
)

var targetInt uint = 254
var dbPath = "dbchain.bolt"

func main() {
	fmt.Println("hello, blockchain !!!")
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		fmt.Printf("open db error: %s \n", err)
		return
	}
	defer db.Close()

	blockchain, err := core.NewBlockchain(targetInt, db)
	if err != nil {
		fmt.Printf("init block chain error: %s\n", err)
		return
	}
	//err = addTestBlock(blockchain)
	//if err != nil {
	//	fmt.Printf("add test block error: %s\n", err)
	//	return
	//}

	blockchain.PrintBlocks()
}

func addTestBlock(blockchain *core.Blockchain) error {
	fmt.Println("--------------- add block 1 ---------------")
	err := blockchain.Add([]byte("block 1"))
	if err != nil {
		fmt.Printf("add block failed: %s\n", err)
		return err
	}
	fmt.Println("--------------- add block 2 ---------------")
	err = blockchain.Add([]byte("block 2"))
	if err != nil {
		fmt.Printf("add block failed: %s\n", err)
		return err
	}
	fmt.Println("--------------- add block OK ---------------")
	return nil
}
