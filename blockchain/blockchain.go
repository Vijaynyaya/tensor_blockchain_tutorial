package blockchain

import (
	"encoding/hex"
	"fmt"
	"os"
	"runtime"

	"github.com/dgraph-io/badger"
)

const (
	// Path to the database
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "First transaction from genesis"
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
func InitBlockChain(address string) *BlockChain {
	// Exit if a blockchain already exists
	if DBexists() {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}

	// New blockchain
	var lastHash []byte

	// Initialize blockcahin database
	opts := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		// Initialize a new blockchain database
		fmt.Println("No existing blockchain found")
		// Add genesis
		cbtx := CoinbaseTx(address, genesisData)
		genesis := Genesis(cbtx)
		fmt.Println("Genesis created")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		Handle(err)
		// Set last hash (lh)
		err = txn.Set([]byte("lh"), genesis.Hash)
		lastHash = genesis.Hash

		return err
	})

	Handle(err)

	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

// Checks if a blockchain database already exists
func DBexists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func ContinueBlockchian() *BlockChain {
	// Exit if the blockchain does not exist
	if !DBexists() {
		fmt.Println("No existing blockchian found, create one!")
		runtime.Goexit()
	}

	var lastHash []byte

	// Access blockchain database
	opts := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opts)
	Handle(err)

	// Get hash of the last block in the blockchain (lh)
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		return err
	})
	Handle(err)

	chain := BlockChain{lastHash, db}

	return &chain

}

// Adds a block to the blockchain 👶
func (chain *BlockChain) AddBlock(transactions []*Transaction) {
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
	newBlock := CreateBlock(transactions, lastHash)

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

// Returns a sturct to iterate over the blockchain
func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}

	return iter
}

// Returns the next block in the blockchain
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

// Returns all the transactions which include unspent outputs with the given address
func (chain *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction
	// Maps transactions to output indexes referenced in inputs.
	spentTXOs := make(map[string][]int)

	// traverse the blockchain
	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

			// Traverse the inputs to find the spent outputs
			if !tx.IsCoinbase() {
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
					}
				}
			}

			// Traverse the outputs
		Outputs:
			for outIdx, out := range tx.Outputs {
				// Is spent
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				// Is unspent
				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}

		}

		// Stop at the genesis block
		if len(block.PrevHash) == 0 {
			break
		}
	}

	return unspentTxs
}

// Returns all the unspent transaction outputs for the given address
func (chain *BlockChain) FindUTXO(address string) []TxOutput {
	var UTXOs []TxOutput
	unspentTransactions := chain.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

// Returns the unspent transactions that will allow the spending of the given amount for the given address.
func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := chain.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOuts
}
