/*
Blockchain: a public database distributed across different peers:
- all peers or nodes are not required to be trustworthy.
- the system works as long as a majority of nodes can be trusted.
- allows the creation of crypto-currencies and smart contracts.

// implementation details
- hashes uniquely identify each block.
- hashes are generated from the data and the hash of the previous block
- hashes can be compared to establish the autheticity of data.
*/
package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

/*
Blockchain is composed of blocks. 🧱
Each block contains the data to be persisted to the database and a hash associated with the block.
*/
type Block struct {
	Hash     []byte // #️⃣ Hash of the block
	Data     []byte // 📄 Data inside the block
	PrevHash []byte // #️⃣🖇 Last block's hash (back linked list).
}

// Generates a hash based on the previous hash and the data.
func (b *Block) DeriveHash() {
	// 🤔 two dimensional slice of bytes?
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})

	// TODO: implement a hashing function that's more secure than sha256.Sum256()
	hash := sha256.Sum256(info) // 🤔 checksum?
	b.Hash = hash[:]
}

// Creates a block 🏭
func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash}
	block.DeriveHash()

	return block
}

// Blockchain ⛓
type BlockChain struct {
	blocks []*Block // TODO: rather than stroing blocks in slice, reference blocks by their hash or value.
}

// Adds a block to the blockchain 👶
func (chain *BlockChain) AddBlock(data string) {
	indexLastBlock := len(chain.blocks) - 1
	prevBlock := chain.blocks[indexLastBlock]

	newBlock := CreateBlock(data, prevBlock.Hash)
	chain.blocks = append(chain.blocks, newBlock)
}

// Creates the genesis block 🎅🧬
func Genesis() *Block {
	return CreateBlock("Genesis: The First Block", []byte{})
}

// Initializes a new blockchain 🤶
func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}

func main() {
	// create a blockchain
	chain := InitBlockChain()

	// add some blocks
	// 👩‍🔬 Try changing the data in any one of the blocks, it'll generate a different hash for that block.
	chain.AddBlock("ONE")
	chain.AddBlock("TWO")
	chain.AddBlock("THREE")

	// traverse the blocks in the blockchain
	for i, block := range chain.blocks {
		// print the formatted string representation of the block
		strBlock := fmt.Sprintf(
			"Block %d %x {\n"+
				"\tPrevious Hash: %x\n"+
				"\tData: %s\n"+
				"}\n",
			i, block.Hash, block.PrevHash, block.Data,
		)
		fmt.Print(strBlock)
	}

}
