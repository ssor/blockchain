package core

import (
	"fmt"
	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
)

const (
	blocksBucket = "blocks"
	lastBlockKey = "l"
)

var (
	ErrorNonceTooLarge   = fmt.Errorf("nonce is larger than setting")
	ErrorInvalidateBlock = fmt.Errorf("block is invalidate")

	blockchainLogger = log.WithFields(log.Fields{
		"file_name": "block_chain",
	})

)

type BlockchainDB interface {
	Update(func(tx *bolt.Tx) error) error
	View(func(tx *bolt.Tx) error) error
}

type Store interface {
	Put([]byte, []byte) error
	Get([]byte) []byte
}

type Blockchain struct {
	targetBits uint
	db         BlockchainDB
}

func NewBlockchain(targetBits uint, db BlockchainDB, address string) (*Blockchain, error) {
	bc := &Blockchain{
		targetBits: targetBits,
		db:         db,
	}
	err := bc.init(address)
	return bc, err
}

func (bc *Blockchain) init(address string) error {
	err := bc.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(blocksBucket))
		if err != nil {
			blockchainLogger.Error("### create bucket error: %s\n", err)
			return fmt.Errorf("cannot create db bucket")
		}

		if bc.alreayInit(bucket) {
			return nil
		}
		cbtx := NewCoinbaseTX(address, genesisBlockData)
		genesis := NewGenesisBlock(cbtx)
		return bc.addBlockToDB(genesis, bucket)
	})
	if err != nil {
		return err
	}
	return nil
}

func (bc *Blockchain) alreayInit(store Store) bool {
	hash := store.Get([]byte(lastBlockKey))
	if hash != nil {
		return true
	}
	return false
}

func (bc *Blockchain) addBlockToDB(block *Block, store Store) error {
	if block.validate(bc.targetBits) {
		bs := store.Get(block.Hash)
		if bs != nil {
			blockchainLogger.Info("block already exists")
			return nil
		}

		bs, err := block.serialize()
		if err != nil {
			blockchainLogger.Errorf("### block serialize error: %s\n", err)
			return fmt.Errorf("cannot serialize block")
		}
		if err := store.Put(block.Hash, bs); err != nil {
			blockchainLogger.Errorf("### bucket put error: %s\n", err)
			return fmt.Errorf("cannot store block data")
		}
		if err := store.Put([]byte(lastBlockKey), block.Hash); err != nil {
			blockchainLogger.Errorf("### bucket put error: %s\n", err)
			return fmt.Errorf("cannot store last block data")
		}
		blockchainLogger.Infof("current block change to -> %x\n", block.Hash)
		return nil
	} else {
		return ErrorInvalidateBlock
	}
}

func (bc *Blockchain) getBlock(hash []byte) (*Block, error) {
	var blockRaw []byte
	err := bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		if bucket == nil {
			blockchainLogger.Errorf("### bucket [%s] should exists, \n", string(blocksBucket))
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
		block, err := deserializeBlock(blockRaw)
		if err != nil {
			return nil, err
		} else {
			return block, nil
		}
	}

}
func (bc *Blockchain) getCurrentBlockHash() ([]byte, error) {
	var hash []byte
	err := bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		if bucket == nil {
			blockchainLogger.Errorf("### bucket [%s] should exists\n", string(blocksBucket))
			return fmt.Errorf("cannot find bucket")
		} else {
			hash = bucket.Get([]byte(lastBlockKey))
			if hash == nil {
				return fmt.Errorf("last block key not found")
			}
			return nil
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
		return bc.saveToDB(newblock)
	} else {
		return ErrorNonceTooLarge
	}
}

func (bc *Blockchain) saveToDB(block *Block) error {
	err := bc.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(blocksBucket))
		if err != nil {
			blockchainLogger.Errorf("### create bucket error: %s\n", err)
			return fmt.Errorf("cannot create db bucket")
		} else {
			return bc.addBlockToDB(block, bucket)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

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
