package main

import "github.com/boltdb/bolt"

const dbFile  = "blockChain.db"
const blocksBucket = "bucket"

type BlockChain struct {
    tip []byte
    db *bolt.DB
}

func (blockChain *BlockChain) AddBlock(data string) {
    var previousHash []byte

    blockChain.db.View(func(tx *bolt.Tx) error {
        block := tx.Bucket([]byte(blocksBucket))
        previousHash = block.Get([]byte("1"))
        return nil
    })

    // Mine the new block
    newBlock := NewBlock(data, previousHash)

    blockChain.db.Update(func(tx *bolt.Tx) error {
        bucket := tx.Bucket([]byte(blocksBucket))
        bucket.Put(newBlock.Hash, newBlock.Serialize())
        bucket.Put([]byte("1"), newBlock.Hash)
        blockChain.tip = newBlock.Hash

        return nil
    })
}

func NewBlockChain() *BlockChain {
    var tip []byte

    db, err := bolt.Open(dbFile, 0600, nil)
    err = db.Update(func(tx *bolt.Tx) error {

        bucket := tx.Bucket([]byte(blocksBucket))

        if bucket == nil { // Create a new BlockChain
            genesis := GenesisBlock()
            bucket, err = tx.CreateBucket([]byte(blocksBucket))
            err = bucket.Put(genesis.Hash, genesis.Serialize())
            err = bucket.Put([]byte("1"), genesis.Hash)
            tip = genesis.Hash
        } else {
            tip = bucket.Get([]byte("1"))
        }

        return nil
    })

    blockChain := BlockChain{tip, db}
    return &blockChain
}

func (blockChain *BlockChain) Iterator() *BlockChainIterator {
    iterator := &BlockChainIterator{blockChain.tip, blockChain.db}
    return iterator
}
