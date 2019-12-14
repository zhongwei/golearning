package main

import (
	"time"
	"bytes"
	"encoding/gob"
)

type Block struct {
	Version int64
	PrevBlockHash []byte
	Hash []byte
	MerkelRoot []byte
	TimeStamp int64
	Bits int64
	Nonce int64
	Data []byte
}

func (block *Block)Serialize() []byte {
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

func NewBlock(data string, prevBlockHash []byte) *Block {
	var block Block

	block = Block {
		Version: 1,
		PrevBlockHash: prevBlockHash,
		MerkelRoot: []byte{},
		TimeStamp: time.Now().Unix(),
		Bits: targetBits,
		Nonce: 0,
		Data: []byte(data),
	}

	pow := NewProofOfWork(&block)
	nonce, hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash

	return &block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block!", []byte{})
}
