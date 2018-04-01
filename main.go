package main

import (
	"fmt"
	"blockchain/core"
	"github.com/boltdb/bolt"
	"flag"
	"os"
)

var targetInt uint = 254
var dbPath = "dbchain.bolt"

var (
	printBlocksCmd = flag.NewFlagSet("printchain", flag.ExitOnError)
	addBlockCmd    = flag.NewFlagSet("addblock", flag.ExitOnError)
	blockData      = addBlockCmd.String("data", "", "data of new block")
)

func main() {
	flag.Parse()

	fmt.Println("hello, blockchain !!!")
	blockchain, db, err := initBlockchain()
	if err != nil {
		panic(err)
	}
	defer func() {
		if db != nil {
			db.Close()
		}
	}()

	switch os.Args[1] {
	case "printchain":
		err = printBlocksCmd.Parse(os.Args[2:])
	case "addblock":
		err = addBlockCmd.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if printBlocksCmd.Parsed() {
		blockchain.PrintBlocks()
		return
	}
	if addBlockCmd.Parsed() {
		if len(*blockData) > 0 {
			err = blockchain.Add([]byte(*blockData))
			if err != nil {
				fmt.Printf("add block failed\n")
				return
			} else {
				fmt.Printf("add block success\n")
				return
			}
		} else {
			fmt.Printf("no data for block\n")
		}
	} else {
		flag.PrintDefaults()
		return
	}

}

func initBlockchain() (*core.Blockchain, *bolt.DB, error) {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		fmt.Printf("open db error: %s \n", err)
		return nil, nil, err
	}

	blockchain, err := core.NewBlockchain(targetInt, db)
	if err != nil {
		fmt.Printf("init block chain error: %s\n", err)
		return nil, db, err
	}
	//err = addTestBlock(blockchain)
	//if err != nil {
	//	fmt.Printf("add test block error: %s\n", err)
	//	return
	//}

	return blockchain, db, nil
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
