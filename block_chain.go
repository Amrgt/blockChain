package main

import (
    "github.com/boltdb/bolt"
    "log"
)

const dbFile  = "blockChain.db"
const blocksBucket = "bucket"

type BlockChain struct {
    tip []byte
    db  *bolt.DB
}

func (blockChain *BlockChain) AddBlock(data string) {
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
    newBlock := NewBlock(data, previousHash)

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

func NewBlockChain() *BlockChain {
    var tip []byte

    db, err := bolt.Open(dbFile, 0600, nil)
    if err != nil {
        log.Panic(err)
    }

    err = db.Update(func(tx *bolt.Tx) error {

        bucket := tx.Bucket([]byte(blocksBucket))

        if bucket == nil { // Create a new BlockChain
            genesis := GenesisBlock()
            bucket, err := tx.CreateBucket([]byte(blocksBucket))
            if err != nil {
                log.Panic(err)
            }

            err = bucket.Put(genesis.Hash, genesis.Serialize())
            if err != nil {
                log.Panic(err)
            }

            err = bucket.Put([]byte("1"), genesis.Hash)
            if err != nil {
                log.Panic(err)
            }

            tip = genesis.Hash
        } else {
            tip = bucket.Get([]byte("1"))
        }

        return nil
    })

    if err != nil {
        log.Panic(err)
    }

    blockChain := BlockChain{tip, db}
    return &blockChain
}

func (blockChain *BlockChain) Iterator() *BlockChainIterator {
    iterator := &BlockChainIterator{blockChain.tip, blockChain.db}
    return iterator
}
