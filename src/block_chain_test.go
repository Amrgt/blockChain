package src

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "regexp"
)

func TestNewBlockChain(t *testing.T) {
    blockChain := NewBlockChain()
    blockChain.AddBlock("FirstBlock")
    blockChain.AddBlock("SecondBlock")

    for i := 0; i < len(blockChain.blocks)-1; i++ {
        block := blockChain.blocks[i]
        proofOfWork := NewProofOfWork(block)

        if i == 0 {
            assert.Equal(t, []byte("Genesis Block"), block.Data, "Genesis Block badly initialized")
            assert.Empty(t, block.PreviousHash, "Genesis Block badly initialized")
        } else {
            assert.Equal(t, blockChain.blocks[i-1].Hash, block.PreviousHash, "Hash should follow next block's PreviousHash")
            assert.Regexp(t, regexp.MustCompile("^.*Block$"), regexp.MustCompile(string(block.Data[:])), "Data should be consistent")
        }
        assert.True(t, proofOfWork.Validate())
    }
}
