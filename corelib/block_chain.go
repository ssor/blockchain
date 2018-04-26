package corelib

import (
    "fmt"
    log "github.com/sirupsen/logrus"
    "blockchain/db"
)

const (
    defaultTargetInt    uint = 254
    genesisCoinbaseData      = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
)

var (
    ErrorNonceTooLarge   = fmt.Errorf("nonce is larger than setting")
    ErrorInvalidateBlock = fmt.Errorf("block is invalidate")

    blockchainLogger = log.WithFields(log.Fields{
        "file_name": "block_chain",
    })
)

type Blockchain struct {
    targetBits uint
}

func NewBlockchain(address string) (*Blockchain, error) {
    bc := &Blockchain{}
    err := bc.init(address)
    return bc, err
}

func BlockchainInitialized() bool {
    status := false
    err := db.GetDb().View(func(tx boltTx) error {
        bucket := tx.Bucket([]byte(db.BucketBlocks))
        if bucket == nil {
            return nil
        }
        v := bucket.Get([]byte(db.BucketLastBlockKey))
        if v == nil {
            return nil
        }
        status = true
        return nil
    })
    if err != nil {
        return false
    }
    return status
}

func (bc *Blockchain) init(address string) error {
    err := db.GetDb().Update(func(tx boltTx) error {
        if BlockchainInitialized() == true {
            return nil
        }
        //bucket, err := tx.CreateBucketIfNotExists([]byte(db.BucketBlocks))
        //if err != nil {
        //   blockchainLogger.Error("### create bucket error: %s\n", err)
        //   return fmt.Errorf("cannot create db bucket")
        //}

        cbtx := NewCoinbaseTX(address, genesisBlockData)
        genesis := NewGenesisBlock(cbtx)
        return bc.addBlockToDB(genesis)
    })
    if err != nil {
        return err
    }
    return nil
}

func (bc *Blockchain) addBlockToDB(block *Block) error {
    if block.validate(bc.targetBits) == false {
        return ErrorInvalidateBlock
    }

    db.GetDb().Batch(func(tx boltTx) error {
        bucket := tx.Bucket([]byte(db.BucketBlocks))
        if bucket == nil {
            blockchainLogger.Errorf("### bucket [%s] should exists, \n", string(db.BucketBlocks))
            return fmt.Errorf("cannot find bucket")
        }

        bs := bucket.Get(block.Hash)
        if bs != nil {
            blockchainLogger.Info("block already exists")
            return nil
        }

        bs, err := block.serialize()
        if err != nil {
            blockchainLogger.Errorf("### block serialize error: %s\n", err)
            return fmt.Errorf("cannot serialize block")
        }
        if err := bucket.Put(block.Hash, bs); err != nil {
            blockchainLogger.Errorf("### bucket put error: %s\n", err)
            return fmt.Errorf("cannot store block data")
        }
        if err := bucket.Put([]byte(db.BucketLastBlockKey), block.Hash); err != nil {
            blockchainLogger.Errorf("### bucket put error: %s\n", err)
            return fmt.Errorf("cannot store last block data")
        }
        return nil
    })
    blockchainLogger.Infof("current block change to -> %x\n", block.Hash)
    return nil
}

func (bc *Blockchain) getBlock(hash []byte) (*Block, error) {
    var blockRaw []byte
    err := db.GetDb().View(func(tx boltTx) error {
        bucket := tx.Bucket([]byte(db.BucketBlocks))
        if bucket == nil {
            blockchainLogger.Errorf("### bucket [%s] should exists, \n", string(db.BucketBlocks))
            return fmt.Errorf("cannot find bucket")
        } else {
            blockRaw = bucket.Get(hash)
            if blockRaw == nil {
                return fmt.Errorf("block key not found, %x\n", hash)
            }
            return nil
        }
        return nil
    })
    if err != nil {
        return nil, err
    } else {
        return deserializeBlock(blockRaw)
    }
}

func (bc *Blockchain) getCurrentBlockHash() ([]byte, error) {
    var hash []byte
    err := db.GetDb().View(func(tx boltTx) error {
        bucket := tx.Bucket([]byte(db.BucketBlocks))
        if bucket == nil {
            blockchainLogger.Errorf("### bucket [%s] should exists\n", string(db.BucketBlocks))
            return fmt.Errorf("cannot find bucket")
        }

        hash = bucket.Get([]byte(db.BucketLastBlockKey))
        if hash == nil {
            return fmt.Errorf("last block key not found")
        }
        return nil
    })

    if err != nil {
        return nil, err
    } else {
        return hash[:], nil
    }
}

func (bc *Blockchain) Add(transactions []*Transaction) error {
    currentHash, err := bc.getCurrentBlockHash()
    if err != nil {
        return err
    }
    if currentHash == nil || len(currentHash) <= 0 {
        return fmt.Errorf("current block not found")
    }
    newblock := NewBlock(currentHash, transactions)
    pow := NewProofOfWork(newblock, bc.targetBits)
    nonce, hash := pow.Mine()
    if nonce >= 0 {
        newblock.Nonce = nonce
        newblock.Hash = hash[:]
        return bc.addBlockToDB(newblock)
    } else {
        return ErrorNonceTooLarge
    }
}

//
//func (bc *Blockchain) saveToDB(block *Block) error {
//    err := db.GetDb().Update(func(tx boltTx) error {
//        bucket, err := tx.CreateBucketIfNotExists([]byte(db.BucketBlocks))
//        if err != nil {
//            blockchainLogger.Errorf("### create bucket error: %s\n", err)
//            return fmt.Errorf("cannot create db bucket")
//        } else {
//            return bc.addBlockToDB(block)
//        }
//        return nil
//    })
//    if err != nil {
//        return err
//    }
//    return nil
//}

func (bc *Blockchain) PrintBlocks() {
    currentHash, err := bc.getCurrentBlockHash()
    if err != nil {
        fmt.Printf("last block hash not found")
        return
    }
    index := 0
    hash := currentHash
    for {
        if v, err := bc.getBlock(hash); err == nil {
            blockchainLogger.Infof("%d: \n", index)
            blockchainLogger.Infof("	hash:     %x\n", string(v.Hash))
            blockchainLogger.Infof("	pre.hash: %x\n", string(v.PreBlockHash))
            //fmt.Printf("	data:     %s\n", v.Data)
            hash = v.PreBlockHash
            index++
        } else {
            blockchainLogger.Infof("%d blocks in chain\n", index)
            break
        }
    }
}
