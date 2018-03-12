package main

import (
    "flag"
    "os"
    "fmt"
    "strconv"
)

const (
    addBlock = "add_block"
    printChain = "print_chain"
)

type CLI struct {
    blockChain *BlockChain
}

func (cli *CLI) Run() {
    addBlockCmd := flag.NewFlagSet(addBlock, flag.ExitOnError)
    printChainCmd := flag.NewFlagSet(printChain, flag.ExitOnError)

    addBlockData := addBlockCmd.String("data", "", "Block data")

    switch os.Args[1] {
    case addBlock:
        addBlockCmd.Parse(os.Args[2:])
    case printChain:
        printChainCmd.Parse(os.Args[2:])
    default:
        cli.PrintUsage()
        os.Exit(1)
    }

    if addBlockCmd.Parsed() {
        if *addBlockData == "" {
            addBlockCmd.Usage()
            os.Exit(1)
        }
        cli.AddBlock(*addBlockData)
    }

    if printChainCmd.Parsed() {
        cli.PrintChain()
    }
}

func (cli *CLI) AddBlock(data string) {
    cli.blockChain.AddBlock(data)
    fmt.Println("Added block: " + data)
}

func (cli *CLI) GetChain() []*Block {
    iterator := cli.blockChain.Iterator()
    var blocks []*Block

    for {
        block := iterator.Next()
        blocks = append(blocks, block)

        if len(block.PreviousHash) == 0 {
            break
        }
    }

    return blocks
}

func (cli *CLI) PrintChain() {
    iterator := cli.blockChain.Iterator()

    for {
        block := iterator.Next()

        fmt.Printf("Previous hash: %x\n", block.PreviousHash)
        fmt.Printf("Data: %s\n", block.Data)
        fmt.Printf("Hash: %x\n", block.Hash)
        pow := NewProofOfWork(block)
        fmt.Printf("Proof of Work: %s\n", strconv.FormatBool(pow.Validate()))
        fmt.Println()

        if len(block.PreviousHash) == 0 {
            break
        }
    }
}

func (cli *CLI) PrintUsage() {
    fmt.Println("cli <command> [option]")
    fmt.Println("commands:")
    fmt.Println("   " + printChain + ": print the block chain")
    fmt.Println("   " + addBlock + " -data <block name> : add the block to the block chain")
    fmt.Println("")
}
