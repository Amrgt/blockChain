package main

import (
    "fmt"
    "bytes"
    "encoding/gob"
    "log"
    "crypto/sha256"
    "encoding/hex"
)

const reward = 1

type TXInput struct {
    TxId      []byte
    Output    int
    ScriptSig string
}

type TXOutput struct {
    Value        int
    ScriptPubKey string
}

type Transaction struct {
    ID     []byte
    Input  []TXInput
    Output []TXOutput
}

func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
    return in.ScriptSig == unlockingData
}

func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
    return out.ScriptPubKey == unlockingData
}

func (tx Transaction) IsCoinbase() bool {
    return len(tx.Input) == 1 && len(tx.Input[0].TxId) == 0 && tx.Input[0].Output == -1
}

func NewCoinBaseTransaction(to, data string) *Transaction {
    if data == "" {
        data = fmt.Sprintf("Reward to %s", to)
    }

    txIn := TXInput{[]byte{}, -1, data}
    txOut := TXOutput{reward, to}
    tx := Transaction{nil, []TXInput{txIn}, []TXOutput{txOut}}
    tx.SetID()

    return &tx
}

func (tx *Transaction) SetID() {
    var encoded bytes.Buffer
    var hash [32]byte

    encoder := gob.NewEncoder(&encoded)
    err := encoder.Encode(tx)
    if err != nil {
        log.Panic(err)
    }

    hash = sha256.Sum256(encoded.Bytes())
    tx.ID = hash[:]
}

func (blockChain *BlockChain) FindUnspentTransactions(address string) []Transaction {
    var unspentTXs []Transaction
    spentTXOs := make(map[string][]int)
    bci := blockChain.Iterator()

    for {
        block := bci.Next()

        for _, tx := range block.Transactions {
            txID := hex.EncodeToString(tx.ID)

        Outputs:
            for outIdx, out := range tx.Output {
                // Was the output spent?
                if spentTXOs[txID] != nil {
                    for _, spentOut := range spentTXOs[txID] {
                        if spentOut == outIdx {
                            continue Outputs
                        }
                    }
                }

                if out.CanBeUnlockedWith(address) {
                    unspentTXs = append(unspentTXs, *tx)
                }
            }

            if tx.IsCoinbase() == false {
                for _, in := range tx.Input {
                    if in.CanUnlockOutputWith(address) {
                        inTxID := hex.EncodeToString(in.TxId)
                        spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Output)
                    }
                }
            }
        }

        if len(block.PreviousHash) == 0 {
            break
        }
    }

    return unspentTXs
}

func (blockChain *BlockChain) FindUTXO(address string) []TXOutput {
    var UTXOs []TXOutput
    unspentTransactions := blockChain.FindUnspentTransactions(address)

    for _, tx := range unspentTransactions {
        for _, out := range tx.Output {
            if out.CanBeUnlockedWith(address) {
                UTXOs = append(UTXOs, out)
            }
        }
    }

    return UTXOs
}

func NewUTXOTransaction(from, to string, amount int, blockChain *BlockChain) *Transaction {
    var inputs []TXInput
    var outputs []TXOutput

    acc, validOutputs := blockChain.FindSpendableOutputs(from, amount)

    if acc < amount {
        log.Panic("ERROR: Not enough funds")
    }

    // Build a list of inputs
    for txId, outs := range validOutputs {
        txID, _ := hex.DecodeString(txId)

        for _, out := range outs {
            input := TXInput{txID, out, from}
            inputs = append(inputs, input)
        }
    }

    // Build a list of outputs
    outputs = append(outputs, TXOutput{amount, to})
    if acc > amount {
        outputs = append(outputs, TXOutput{acc - amount, from}) // a change
    }

    tx := Transaction{nil, inputs, outputs}
    tx.SetID()

    return &tx
}

func (blockChain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
    unspentOutputs := make(map[string][]int)
    unspentTXs := blockChain.FindUnspentTransactions(address)
    accumulated := 0

Work:
    for _, tx := range unspentTXs {
        txID := hex.EncodeToString(tx.ID)

        for outIdx, out := range tx.Output {
            if out.CanBeUnlockedWith(address) && accumulated < amount {
                accumulated += out.Value
                unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

                if accumulated >= amount {
                    break Work
                }
            }
        }
    }

    return accumulated, unspentOutputs
}

