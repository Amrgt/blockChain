package src

import (
    "strconv"
    "crypto/sha256"
    "bytes"
    "time"
)

type Block struct {
    Timestamp int64
    Data []byte
    PreviousHash []byte
    Hash []byte
}

func (block *Block) calcHash() {
    timestamp := []byte(strconv.FormatInt(block.Timestamp, 10))
    headers := bytes.Join([][]byte{block.PreviousHash, block.Data, timestamp}, []byte{})
    hash := sha256.Sum256(headers)
    block.Hash = hash[:]
}

func NewBlock(data string, previousBlockHash []byte) *Block {
    block := &Block{time.Now().Unix(), []byte(data), previousBlockHash, []byte{}}
    block.calcHash()
    return block
}

func GenesisBlock() *Block {
    return NewBlock("Genesis Block", []byte{})
}
