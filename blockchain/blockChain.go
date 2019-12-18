package main

import (
	"fmt"
	"os"

	bolt "go.etcd.io/bbolt"
)

const dbFile = "chain.db"
const blockBucket = "bucket"
const lastHashKey = "key"
const genesisInfo = "Boeing to halt 737 Max production in January"

// BlockChain struct
type BlockChain struct {
	db   *bolt.DB
	tail []byte
}

func isDBExist() bool {
	_, err := os.Stat(dbFile)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// InitBlockChain func
func InitBlockChain(address string) *BlockChain {
	if isDBExist() {
		fmt.Println("Blockchain exist already, need not to create!")
		os.Exit(1)
	}

	db, err := bolt.Open(dbFile, 0600, nil)
	CheckErr("NewBlockChain position #1", err)

	var lastHash []byte

	db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucket([]byte(blockBucket))
		CheckErr("InitBlockChain position #1", err)
		genesis := NewGenesisBlock(NewCoinbase(address, genesisInfo))
		bucket.Put(genesis.Hash, genesis.Serialize())
		CheckErr("InitBlockChain position #2", err)
		bucket.Put([]byte(lastHashKey), genesis.Hash)
		CheckErr("InitBlockChain position #3", err)
		lastHash = genesis.Hash

		return nil
	})

	return &BlockChain{db, lastHash}
}

// GetBlockChainHandler func
func GetBlockChainHandler() *BlockChain {
	if !isDBExist() {
		fmt.Println("Please create blockchain first!")
		os.Exit(1)
	}

	db, err := bolt.Open(dbFile, 0600, nil)
	CheckErr("GetBlockChainHandler1", err)

	var lastHash []byte

	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket != nil {
			lastHash = bucket.Get([]byte(lastHashKey))
		}

		return nil
	})

	return &BlockChain{db, lastHash}
}

// AddBlock func
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	var prevBlockHash []byte
	err := bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))

		if bucket == nil {
			os.Exit(1)
		}

		prevBlockHash = bucket.Get([]byte(lastHashKey))
		return nil
	})

	block := NewBlock(txs, prevBlockHash)
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

// BlockChainIterator struct
type BlockChainIterator struct {
	currHash []byte
	db       *bolt.DB
}

// NewIterator func
func (bc *BlockChain) NewIterator() *BlockChainIterator {
	return &BlockChainIterator{currHash: bc.tail, db: bc.db}
}

// Next func
func (it *BlockChainIterator) Next() *Block {
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

func (bc *BlockChain) FindUTXOTransactions(address string) []Transaction {
	var UTXOTransanctions []Transaction

	spentUTXO := make(map[string][]int64)

	it := bc.NewIterator()
	for {
		block := it.Next()

		for _, tx := range block.Transactions {
			if !tx.IsCoinbase() {
				for _, input := range tx.TXInputs {
					if input.CanUnlockUTXOWith(address) {
						spentUTXO[string(input.TXID)] = append(spentUTXO[string(input.TXID)], input.Vout)
					}
				}
			}

		OUTPUTS:
			for currIndex, output := range tx.TXOutputs {
				if spentUTXO[string(tx.TXID)] != nil {
					indexes := spentUTXO[string(tx.TXID)]
					for _, index := range indexes {
						if int64(currIndex) == index {
							continue OUTPUTS
						}
					}
				}
				if output.CanBeUnlockedUTXOWith(address) {
					UTXOTransanctions = append(UTXOTransanctions, *tx)
				}
			}

		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return UTXOTransanctions
}

func (bc *BlockChain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput
	txs := bc.FindUTXOTransactions(address)

	for _, tx := range txs {
		for _, utxo := range tx.TXOutputs {
			if utxo.CanBeUnlockedUTXOWith(address) {
				UTXOs = append(UTXOs, utxo)
			}
		}
	}
	return UTXOs
}

func (bc *BlockChain) FindSuitableUTXOs(address string, amount float64) (map[string][]int64, float64) {
	txs := bc.FindUTXOTransactions(address)
	validUTXOs := make(map[string][]int64)
	var total float64

FIND:
	for _, tx := range txs {
		outputs := tx.TXOutputs

		for index, output := range outputs {
			if output.CanBeUnlockedUTXOWith(address) && total < amount {
				total += output.Value
				validUTXOs[string(tx.TXID)] = append(validUTXOs[string(tx.TXID)], int64(index))
			} else {
				break FIND
			}
		}
	}
	return validUTXOs, total
}
