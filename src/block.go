package main

import (
    "time"
    "bytes"
    "encoding/gob"
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

func (block *Block) Serialize() []byte {
    var serialized bytes.Buffer
    encoder := gob.NewEncoder(&serialized)
    encoder.Encode(block)
    return serialized.Bytes()
}

func DeserializeBlock(serialized []byte) *Block {
    var block Block
    decoder := gob.NewDecoder(bytes.NewReader(serialized))
    decoder.Decode(&block)
    if len(block.PreviousHash) == 0 {
        block.PreviousHash = []byte{}
    }
    return &block
}
