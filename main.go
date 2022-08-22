package main

import (
	"os"

	"github.com/vijaynyaya/tensor_programming_golang_blockchain/blockchain"
	"github.com/vijaynyaya/tensor_programming_golang_blockchain/cli"
)

// Example Usage:
// $ go run main.go print
// $ go run main.go add "The Block's Data"

func main() {
	defer os.Exit(0)
	// create a blockchain
	chain := blockchain.InitBlockChain()
	// Properly close the database before the main() function exits
	defer chain.Database.Close()

	// Run 🏃‍♂️
	cli := cli.CommandLine{Blockchain: chain}
	cli.Run()
}
