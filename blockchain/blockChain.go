package main

import (
	"os"
	bolt "go.etcd.io/bbolt"
)

const dbFile = "achain.db"
const blockBucket = "bucket"
const lastHashKey = "key"

type BlockChain struct {
	db *bolt.DB
	tail []byte
}

func NewBlockChain() *BlockChain {
	db, err := bolt.Open(dbFile, 0600, nil)
	CheckErr("NewBlockChain1", err)
	
	var lastHash []byte

	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket != nil {
			lastHash = bucket.Get([]byte(lastHashKey))
		} else {
			genesis := NewGenesisBlock()
			bucket, err := tx.CreateBucket([]byte(blockBucket))
			CheckErr("NewBlockChain2", err)
			bucket.Put(genesis.Hash, genesis.Serialize())
			CheckErr("NewBlockChain3", err)
			bucket.Put([]byte(lastHashKey), genesis.Hash)
			CheckErr("NewBlockChain4", err)
			lastHash = genesis.Hash
		}

		return nil
	})

	return &BlockChain{db, lastHash}
}

func (bc *BlockChain)AddBlock(data string) {
	var prevBlockHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))

		if bucket == nil {
			os.Exit(1)
		}

		prevBlockHash = bucket.Get([]byte(lastHashKey))
		return nil
	})
	
	block := NewBlock(data, prevBlockHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))

		if bucket == nil {
			os.Exit(1)
		}

		err := bucket.Put(block.Hash, block.Serialize())
		CheckErr("AddBlock1", err)
		err = bucket.Put([]byte(lastHashKey), block.Hash)
		CheckErr("AddBlock2", err)
		bc.tail = block.Hash

		return nil
	})

	CheckErr("AddBlock3", err)

}

type BlockChainIterator struct {
	currHash []byte
	db *bolt.DB
}

func (bc *BlockChain)NewIterator() *BlockChainIterator {
	return &BlockChainIterator{currHash:bc.tail, db:bc.db}
}

func (it *BlockChainIterator)Next() *Block {
	var block *Block

	err := it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			return nil
		}

		data := bucket.Get(it.currHash)
		block = Deserialize(data)
		it.currHash = block.PrevBlockHash
		return nil
	})	
	CheckErr("Next()", err)
	return block
}
