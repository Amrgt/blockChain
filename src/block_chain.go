package src

type BlockChain struct {
    blocks []*Block
}

func (blockChain *BlockChain) AddBlock(data string) {
    previousBlock := blockChain.blocks[len(blockChain.blocks) - 1]
    nextBlock := NewBlock(data, previousBlock.Hash)
    blockChain.blocks = append(blockChain.blocks, nextBlock)
}

func NewBlockChain() *BlockChain {
    return &BlockChain{[]*Block{GenesisBlock()}}
}
