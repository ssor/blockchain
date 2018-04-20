package main

import (
    "fmt"
    "io/ioutil"
    "gopkg.in/yaml.v2"
    "flag"
    "blockchain/core"
)

var targetInt uint = 254
//var dbPath = "my.db"

const (
    nodeTypeCore   = "corelib"
    nodeTypeWallet = "wallet"
    nodeTypeMiner  = "miner"
)

var (
    node = flag.String("core", "corelib", "core type to startup")
)

type NodeConfig struct {
    PortsMax int `yaml:"port_max"`
    PortsMin int `yaml:"port_min"`
}

type Config struct {
    Core   *NodeConfig
    Wallet *NodeConfig
    Miner  *NodeConfig
    Dbpath string
}

func readConfig() (*Config, error) {
    bs, err := ioutil.ReadFile("config.yaml")
    if err != nil {
        fmt.Println(err)
        return nil, err
    }
    var config Config
    err = yaml.Unmarshal(bs, &config)
    if err != nil {
        fmt.Println(err)
        return nil, err
    }
    return &config, nil
}

func main() {
    flag.Parse()
    if flag.Parsed() == false {
        flag.PrintDefaults()
        return
    }
    fmt.Println("core: ", *node)

    config, err := readConfig()
    if err != nil {
        return
    }

    fmt.Println("db: ", config.Dbpath)
    switch *node {
    case nodeTypeCore:
        fmt.Println("ports: ", config.Core.PortsMin, "-", config.Core.PortsMax)
        core.StartNode(config.Core.PortsMin)
    case nodeTypeMiner:
        fmt.Println("ports: ", config.Miner.PortsMin, "-", config.Miner.PortsMax)
    case nodeTypeWallet:
        fmt.Println("ports: ", config.Wallet.PortsMin, "-", config.Wallet.PortsMax)
    default:
        fmt.Println("unknown core type")
        return
    }

    //log.Infof("hello, blockchain !!!")
    //blockchain, db, err := initBlockchain()
    //if err != nil {
    //	panic(err)
    //}
    //defer func() {
    //	if db != nil {
    //		db.Close()
    //	}
    //}()
    //
    //switch os.Args[1] {
    //case "printchain":
    //err = printBlocksCmd.Parse(os.Args[2:])
    //case "addblock":
    //	err = addBlockCmd.Parse(os.Args[2:])
    //default:
    //	flag.PrintDefaults()
    //	os.Exit(1)
    //}
    //
    //if printBlocksCmd.Parsed() {
    //blockchain.PrintBlocks()
    //return
    //}
    //if addBlockCmd.Parsed() {
    //	if len(*blockData) > 0 {
    //		err = blockchain.Add([]byte(*blockData))
    //		if err != nil {
    //			fmt.Printf("add block failed\n")
    //			return
    //		} else {
    //			fmt.Printf("add block success\n")
    //			return
    //		}
    //	} else {
    //		fmt.Printf("no data for block\n")
    //	}
    //} else {
    //	flag.PrintDefaults()
    //	return
    //}

}

//
//func initBlockchain(blockchainCreator string) (*corelib.Blockchain, *bolt.DB, error) {
//    db, err := bolt.Open(dbPath, 0600, nil)
//    if err != nil {
//        fmt.Printf("open db error: %s \n", err)
//        return nil, nil, err
//    }
//
//    blockchain, err := corelib.NewBlockchain(targetInt, db, blockchainCreator)
//    if err != nil {
//        fmt.Printf("init block chain error: %s\n", err)
//        return nil, db, err
//    }
//
//    return blockchain, db, nil
//}
