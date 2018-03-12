package main

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestBlock_Serialization(t *testing.T) {
    block := GenesisBlock()
    serialized := block.Serialize()
    deSerialized := DeserializeBlock(serialized)
    assert.Equal(t, block.Nonce, deSerialized.Nonce)
    assert.Equal(t, block.Hash, deSerialized.Hash)
    assert.Equal(t, block.Timestamp, deSerialized.Timestamp)
    assert.Equal(t, block.Data, deSerialized.Data)
    assert.Equal(t, block.PreviousHash, deSerialized.PreviousHash)
}
