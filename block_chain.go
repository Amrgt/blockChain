package main

import (
    "github.com/boltdb/bolt"
    "log"
    "fmt"
    "os"
)

const dbFile  = "blockChain.db"
const blocksBucket = "bucket"
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type BlockChain struct {
    tip []byte
    db  *bolt.DB
}

func (blockChain *BlockChain) AddBlock(transactions []*Transaction) {
    var previousHash []byte

    err := blockChain.db.View(func(tx *bolt.Tx) error {
        block := tx.Bucket([]byte(blocksBucket))
        previousHash = block.Get([]byte("1"))
        return nil
    })

    if err != nil {
        log.Panic(err)
    }

    // Mine the new block
    newBlock := NewBlock(transactions, previousHash)

    err = blockChain.db.Update(func(tx *bolt.Tx) error {
        bucket := tx.Bucket([]byte(blocksBucket))
        err := bucket.Put(newBlock.Hash, newBlock.Serialize())
        if err != nil {
            log.Panic(err)
        }

        err = bucket.Put([]byte("1"), newBlock.Hash)
        if err != nil {
            log.Panic(err)
        }
        blockChain.tip = newBlock.Hash

        return nil
    })

    if err != nil {
        log.Panic(err)
    }
}

func NewBlockChain(address string) *BlockChain {
    if dbExists() == false {
        fmt.Println("No existing blockchain found. Create one first.")
        os.Exit(1)
    }

    var tip []byte
    db, err := bolt.Open(dbFile, 0600, nil)
    if err != nil {
        log.Panic(err)
    }

    err = db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(blocksBucket))
        tip = b.Get([]byte("l"))

        return nil
    })

    if err != nil {
        log.Panic(err)
    }

    bc := BlockChain{tip, db}

    return &bc
}


func CreateBlockChain(address string) *BlockChain {
    if dbExists() {
        fmt.Println("Blockchain already exists.")
        os.Exit(1)
    }

    var tip []byte
    db, err := bolt.Open(dbFile, 0600, nil)
    if err != nil {
        log.Panic(err)
    }

    err = db.Update(func(tx *bolt.Tx) error {
        cbTx := NewCoinBaseTransaction(address, genesisCoinbaseData)
        genesis := GenesisBlock(cbTx)

        b, err := tx.CreateBucket([]byte(blocksBucket))
        if err != nil {
            log.Panic(err)
        }

        err = b.Put(genesis.Hash, genesis.Serialize())
        if err != nil {
            log.Panic(err)
        }

        err = b.Put([]byte("l"), genesis.Hash)
        if err != nil {
            log.Panic(err)
        }
        tip = genesis.Hash

        return nil
    })

    if err != nil {
        log.Panic(err)
    }

    bc := BlockChain{tip, db}

    return &bc
}

func dbExists()    bool {
        if _, err := os.Stat(dbFile); os.IsNotExist(err) {
        return false
    }

    return true
}

func (blockChain *BlockChain) Iterator() *BlockChainIterator {
    iterator := &BlockChainIterator{blockChain.tip, blockChain.db}
    return iterator
}
