package main

import (
    "time"
    "bytes"
    "encoding/gob"
    "log"
    "crypto/sha256"
)

type Block struct {
    Timestamp    int64
    Transactions []*Transaction
    PreviousHash []byte
    Hash         []byte
    Nonce        int
}

func NewBlock(transactions []*Transaction, previousBlockHash []byte) *Block {
    block := &Block{time.Now().Unix(), transactions, previousBlockHash, []byte{}, 0}
    proofOfWork := NewProofOfWork(block)
    nonce, hash := proofOfWork.Work()
    block.Nonce = nonce
    block.Hash = hash[:]

    return block
}

func GenesisBlock(transaction *Transaction) *Block {
    return NewBlock([]*Transaction{transaction}, []byte{})
}

func (block *Block) Serialize() []byte {
    var serialized bytes.Buffer
    encoder := gob.NewEncoder(&serialized)
    err := encoder.Encode(block)
    if err != nil {
        log.Panic(err)
    }

    return serialized.Bytes()
}

func DeserializeBlock(serialized []byte) *Block {
    var block Block
    decoder := gob.NewDecoder(bytes.NewReader(serialized))
    err := decoder.Decode(&block)
    if err != nil {
        log.Panic(err)
    }
    return &block
}

func(block *Block) HashTransactions() []byte {
    var txHashes [][]byte
    var txHash [32]byte

    for _, tx := range block.Transactions {
        txHashes = append(txHashes, tx.ID)
    }
    txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

    return txHash[:]
}
