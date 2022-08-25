package main

import (
	"os"

	"github.com/vijaynyaya/tensor_programming_golang_blockchain/cli"
)

// Example Usage:
// $ go run main.go print
// $ go run main.go add "The Block's Data"

func main() {
	defer os.Exit(0)
	// Run 🏃‍♂️
	cli := cli.CommandLine{}
	cli.Run()
}
