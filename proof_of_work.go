package main

import (
    "math/big"
    "bytes"
    "math"
    "crypto/sha256"
    "fmt"
    "encoding/binary"
    "log"
)

const minBits = 18
const maxNonce = math.MaxInt64

type ProofOfWork struct {
    block *Block
    goal  *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
    goal := big.NewInt(1)
    goal.Lsh(goal, uint(256 - minBits))

    proofOfWork := &ProofOfWork{block, goal}
    return proofOfWork
}

func (proofOfWork *ProofOfWork) OrganizeData(nonce int) []byte {
    data := bytes.Join([][]byte{
        proofOfWork.block.PreviousHash,
        proofOfWork.block.HashTransactions(),
        IntToHex(proofOfWork.block.Timestamp),
        IntToHex(int64(minBits)),
        IntToHex(int64(nonce)),
    }, []byte{},)

    return data
}

func (proofOfWork *ProofOfWork) Work() (int, []byte) {
    var hashInt big.Int
    var hash [32]byte
    nonce := 0

    fmt.Printf("Mining the new block")

    for nonce < maxNonce  {
        data := proofOfWork.OrganizeData(nonce)
        hash := sha256.Sum256(data)

        hashInt.SetBytes(hash[:])
        if hashInt.Cmp(proofOfWork.goal) != -1 {
            nonce++
        } else {
            fmt.Println(hash)
            break
        }
    }

    return nonce, hash[:]
}

func (proofOfWork *ProofOfWork) Validate() bool {
    var hashInt big.Int

    data := proofOfWork.OrganizeData(proofOfWork.block.Nonce)
    hash := sha256.Sum256(data)
    hashInt.SetBytes(hash[:])

    isValid := hashInt.Cmp(proofOfWork.goal) == -1

    return isValid
}

func IntToHex(num int64) []byte {
    buff := new(bytes.Buffer)
    err := binary.Write(buff, binary.BigEndian, num)
    if err != nil {
        log.Panic(err)
    }

    return buff.Bytes()
}
