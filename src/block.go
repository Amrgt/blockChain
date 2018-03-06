package src

import (
    "time"
)

type Block struct {
    Timestamp int64
    Data []byte
    PreviousHash []byte
    Hash []byte
    Nonce int
}

func NewBlock(data string, previousBlockHash []byte) *Block {
    block := &Block{time.Now().Unix(), []byte(data), previousBlockHash, []byte{}, 0}
    proofOfWork := NewProofOfWork(block)
    nonce, hash := proofOfWork.Work()
    block.Nonce = nonce
    block.Hash = hash[:]
    return block
}

func GenesisBlock() *Block {
    return NewBlock("Genesis Block", []byte{})
}
