package main

import "github.com/boltdb/bolt"

type BlockChainIterator struct {
    currentHash []byte
    db *bolt.DB
}


func (iterator *BlockChainIterator) Next() *Block {
    var block *Block
    iterator.db.View(func(tx *bolt.Tx) error {
        bucket := tx.Bucket([]byte(blocksBucket))
        serialized := bucket.Get(iterator.currentHash)
        block = DeserializeBlock(serialized)
        return nil
    })

    iterator.currentHash = block.PreviousHash

    return block
}
