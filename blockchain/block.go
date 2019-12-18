package main

import (
	"bytes"
	"encoding/gob"
	"time"
)

type Block struct {
	Version       int64
	PrevBlockHash []byte
	Hash          []byte
	MerkelRoot    []byte
	TimeStamp     int64
	Bits          int64
	Nonce         int64
	Transactions  []*Transaction
}

func (block *Block) Serialize() []byte {
	/*
		fmt.Printf("It's here2! Version = %v \n", block.Version)
		fmt.Printf("It's here2! PrevBlockHash = %v \n", block.PrevBlockHash)
		fmt.Printf("It's here2! Hash = %v \n", block.Hash)
		fmt.Printf("It's here2! MerkelRoot = %v \n", block.MerkelRoot)
		fmt.Printf("It's here2! TimeStamp = %v \n", block.TimeStamp)
		fmt.Printf("It's here2! Bits = %v \n", block.Bits)
		fmt.Printf("It's here2! Nonce = %v \n", block.Nonce)
		fmt.Printf("It's here2! Transactions = %v \n", block.Transactions)
	*/
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(block)
	CheckErr("Serialize", err)
	return buffer.Bytes()
}

func Deserialize(data []byte) *Block {
	if len(data) == 0 {
		return nil
	}

	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	CheckErr("Deserialize", err)
	return &block
}

func NewBlock(txs []*Transaction, prevBlockHash []byte) *Block {
	var block Block

	block = Block{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		MerkelRoot:    []byte{},
		TimeStamp:     time.Now().Unix(),
		Bits:          targetBits,
		Nonce:         0,
		Transactions:  txs,
	}

	pow := NewProofOfWork(&block)
	nonce, hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash

	return &block
}

func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}
