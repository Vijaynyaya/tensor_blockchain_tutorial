package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/vijaynyaya/tensor_programming_golang_blockchain/blockchain"
)

// CommandLine allows the user to interact with the blockchain through the comman line
type CommandLine struct {
	Blockchain *blockchain.BlockChain
}

// Prints the command line interface's usage
func (cli *CommandLine) PrintUsage() {
	fmt.Println(
		"Usage:\n" +
			" add -block <BLOCK_DATA> \t adds a block to the blockchain\n" +
			" print \t prints out the blocks in the blockchain",
	)

}

// Validates the arguments passed through the command line interface
func (cli *CommandLine) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		// It's important to allow BadgerDB enough time to garbage collect and prevent the corruption of data.
		// Therefore, use runtime.Goexit()
		// unlike os.Exit(), runtime.Goexit() exits the application by shutting down the goroutine.
		runtime.Goexit()
	}
}

// Adds a block with the supplied data to the bockchain
func (cli *CommandLine) AddBlock(data string) {
	cli.Blockchain.AddBlock(data)
	fmt.Println("Added Block!")
}

// Prints all the blocks in the blockchain
func (cli *CommandLine) PrintChain() {
	iter := cli.Blockchain.Iterator()

	// traverse the blocks in the blockchain
	for {
		// Get the block
		block := iter.Next()
		// Print the block
		strBlock := fmt.Sprintf(
			"Block %x {\n"+
				"\tPrevious Hash: %x\n"+
				"\tData: %s\n"+
				"\tNonce: %v\n"+
				"}\n",
			block.Hash, block.PrevHash, block.Data, block.Nonce,
		)
		fmt.Print(strBlock)
		// Validate the block 👩‍⚖️
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		// Break the loop on genesis block
		if len(block.PrevHash) == 0 {
			break
		}
	}

}

func (cli *CommandLine) Run() {
	cli.ValidateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	default:
		cli.PrintUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		// 🤔 How's addBlockData initalized by addBlockCmd.String() ?
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.AddBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.PrintChain()
	}
}
