package main

import (
    "flag"
    "os"
    "fmt"
    "strconv"
    "log"
)

const (
    send = "send"
    printChain = "print_chain"
    createBlockChain = "create_block_chain"
    getBalance = "get_balance"
)

type CLI struct {}

func (cli *CLI) PrintUsage() {
    fmt.Println("Usage:")
    fmt.Println("   " + getBalance + " -address ADDRESS - Get balance of ADDRESS")
    fmt.Println("  " + createBlockChain + " -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
    fmt.Println("  " + printChain + "  - Print all the blocks of the blockchain")
    fmt.Println("  " + send + " -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")
}

func (cli *CLI) ValidateArgs() {
    if len(os.Args) < 2 {
        cli.PrintUsage()
        os.Exit(1)
    }
}

func (cli *CLI) getBalance(address string) {
    blockChain := NewBlockChain(address)
    defer blockChain.db.Close()

    balance := 0
    UTXOs := blockChain.FindUTXO(address)

    for _, out := range UTXOs {
        balance += out.Value
    }

    fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func (cli *CLI) send(from, to string, amount int) {
    blockChain := NewBlockChain(from)
    defer blockChain.db.Close()

    tx := NewUTXOTransaction(from, to, amount, blockChain)
    blockChain.AddBlock([]*Transaction{tx})
    fmt.Println("Success!")
}

func (cli *CLI) CreateBlockChain(address string) {
    blockChain := CreateBlockChain(address)
    blockChain.db.Close()
    fmt.Println("BlockChain created!")
}

func (cli *CLI) PrintChain() {

    blockChain := NewBlockChain("")
    defer blockChain.db.Close()

    iterator := blockChain.Iterator()

    for {
        block := iterator.Next()

        fmt.Printf("Previous hash: %x\n", block.PreviousHash)
        fmt.Printf("Hash: %x\n", block.Hash)
        pow := NewProofOfWork(block)
        fmt.Printf("Proof of Work: %s\n", strconv.FormatBool(pow.Validate()))
        fmt.Println()

        if len(block.PreviousHash) == 0 {
            break
        }
    }
}

func (cli *CLI) Run() {
    cli.ValidateArgs()

    getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
    createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
    sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
    printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

    getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
    createBlockChainAddress := createBlockChainCmd.String("address", "", "The address to send genesis block reward to")
    sendFrom := sendCmd.String("from", "", "Source wallet address")
    sendTo := sendCmd.String("to", "", "Destination wallet address")
    sendAmount := sendCmd.Int("amount", 0, "Amount to send")

    switch os.Args[1] {
    case getBalance:
        err := getBalanceCmd.Parse(os.Args[2:])
        if err != nil {
            log.Panic(err)
        }
    case createBlockChain:
        err := createBlockChainCmd.Parse(os.Args[2:])
        if err != nil {
            log.Panic(err)
        }
    case printChain:
        err := printChainCmd.Parse(os.Args[2:])
        if err != nil {
            log.Panic(err)
        }
    case send:
        err := sendCmd.Parse(os.Args[2:])
        if err != nil {
            log.Panic(err)
        }
    default:
        cli.PrintUsage()
        os.Exit(1)
    }

    if getBalanceCmd.Parsed() {
        if *getBalanceAddress == "" {
            getBalanceCmd.Usage()
            os.Exit(1)
        }
        cli.getBalance(*getBalanceAddress)
    }

    if createBlockChainCmd.Parsed() {
        if *createBlockChainAddress == "" {
            createBlockChainCmd.Usage()
            os.Exit(1)
        }
        cli.CreateBlockChain(*createBlockChainAddress)
    }

    if printChainCmd.Parsed() {
        cli.PrintChain()
    }

    if sendCmd.Parsed() {
        if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
            sendCmd.Usage()
            os.Exit(1)
        }

        cli.send(*sendFrom, *sendTo, *sendAmount)
    }
}
