package main

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "regexp"
)

func TestNewBlockChain(t *testing.T) {
    blockChain := NewBlockChain()

    cli := CLI{blockChain}
    cli.AddBlock("FirstBlock")
    cli.AddBlock("SecondBlock")
    cli.PrintChain()
    blocks := cli.GetChain()

    for i := 0; i < len(blocks)-1; i++ {
        block := blocks[i]
        proofOfWork := NewProofOfWork(block)

        if i == 0 {
            assert.Equal(t, []byte("Genesis Block"), block.Data, "Genesis Block badly initialized")
            assert.Empty(t, block.PreviousHash, "Genesis Block badly initialized")
        } else {
            assert.Equal(t, blocks[i-1].Hash, block.PreviousHash, "Hash should follow next block's PreviousHash")
            assert.Regexp(t, regexp.MustCompile("^.*Block$"), regexp.MustCompile(string(block.Data[:])), "Data should be consistent")
        }
        assert.True(t, proofOfWork.Validate())
    }
}
