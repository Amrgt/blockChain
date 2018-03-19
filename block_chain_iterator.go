package main

import (
    "github.com/boltdb/bolt"
    "log"
    "fmt"
)

type BlockChainIterator struct {
    currentHash []byte
    db          *bolt.DB
}


func (iterator *BlockChainIterator) Next() *Block {
    var block *Block
    err := iterator.db.View(func(tx *bolt.Tx) error {
        bucket := tx.Bucket([]byte(blocksBucket))
        serialized := bucket.Get(iterator.currentHash)
        block = DeserializeBlock(serialized)

        return nil
    })

    if err != nil {
        log.Panic(err)
    }

    iterator.currentHash = block.PreviousHash

    return block
}
