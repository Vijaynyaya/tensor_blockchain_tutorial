package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

const (
	// Path to the database
	dbPath = "./tmp/blocks"
)

// Blockchain ⛓
type BlockChain struct {
	LastHash []byte // back linked list
	Database *badger.DB
}

// Usefull for iterating block-by-block over the blockchain.
type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

// Initializes a new blockchain 🤶
func InitBlockChain() *BlockChain {
	var lastHash []byte

	// Access the blockcahin database
	opts := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		// Try fetching the hash of the last block in the blockchain (lh)
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			// Initialize a new blockchain database
			fmt.Println("No existing blockchain found")
			// Add genesis block
			genesis := Genesis()
			fmt.Println("Genesis proved")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)
			// Set last hash (lh)
			err = txn.Set([]byte("lh"), genesis.Hash)
			lastHash = genesis.Hash

			return err
		} else {
			// Get hash of the last block in the blockchain (lh)
			item, err := txn.Get([]byte("lh"))
			Handle(err)
			err = item.Value(func(val []byte) error {
				lastHash = val
				return nil
			})
			return err
		}
	})

	Handle(err)

	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

// Adds a block to the blockchain 👶
func (chain *BlockChain) AddBlock(data string) {
	// Get last block's hash
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})

		return err
	})
	Handle(err)

	// Create new block
	newBlock := CreateBlock(data, lastHash)

	// Add new block to the database
	err = chain.Database.Update(func(txn *badger.Txn) error {
		// add new block
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		// update last hash to be the new block's hash
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
	Handle(err)
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}

	return iter
}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	// Get the block from the database
	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)
		err = item.Value(func(byteEncodedVal []byte) error {
			block = Deserialize(byteEncodedVal)
			return nil
		})

		return err
	})
	Handle(err)

	// Hash of next block
	iter.CurrentHash = block.PrevHash

	return block
}
