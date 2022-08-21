package main

import (
	"fmt"
	"strconv"

	"github.com/vijaynyaya/tensor_programming_golang_blockchain/blockchain"
)

func main() {
	// create a blockchain
	chain := blockchain.InitBlockChain()

	// add some blocks
	// 👩‍🔬 Try changing the data in any one of the blocks, it'll generate a different hash for that block.
	chain.AddBlock("ONE")
	chain.AddBlock("TWO")
	chain.AddBlock("THREE")

	// traverse the blocks in the blockchain
	for i, block := range chain.Blocks {
		// print the formatted string representation of the block
		strBlock := fmt.Sprintf(
			"Block %d %x {\n"+
				"\tPrevious Hash: %x\n"+
				"\tData: %s\n"+
				"\tNonce: %v\n"+
				"}\n",
			i, block.Hash, block.PrevHash, block.Data, block.Nonce,
		)
		fmt.Print(strBlock)

		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}

}
