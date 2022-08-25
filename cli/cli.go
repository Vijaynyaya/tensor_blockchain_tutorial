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
type CommandLine struct{}

// Prints the command line interface's usage
func (cli *CommandLine) PrintUsage() {
	fmt.Println(
		"Usage:\n" +
			" getbalance -address ADDRESS \t get the balance for an address\n" +
			" createblockchain -address ADDRESS \t create a new blockchain with a genesis block mined with the given address\n" +
			" send -from FROM -to TO -amount AMOUNT \t sends the AMOUNT from FROM address to TO address\n" +
			" printchain \t prints out the blocks in the blockchain",
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

func (cli *CommandLine) send(from, to string, amount int) {
	chain := blockchain.ContinueBlockchian()
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Success!")
}

func (cli *CommandLine) createBlockChain(address string) {
	chain := blockchain.InitBlockChain(address)
	chain.Database.Close()
	fmt.Println("Finished!")
}

func (cli *CommandLine) getBalance(address string) {
	chain := blockchain.ContinueBlockchian()
	defer chain.Database.Close()

	balance := 0
	UTXOs := chain.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

// Prints all the blocks in the blockchain
func (cli *CommandLine) PrintChain() {
	chain := blockchain.ContinueBlockchian()
	defer chain.Database.Close()
	iter := chain.Iterator()

	// traverse the blocks in the blockchain
	for {
		// Get the block
		block := iter.Next()
		// Print the block
		strBlock := fmt.Sprintf(
			"Block %x {\n"+
				"\tPrevious Hash: %x\n"+
				"\tNonce: %v\n"+
				"}\n",
			block.Hash, block.PrevHash, block.Nonce,
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

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "wallet address")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "Genesis address")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.String("amount", "", "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	default:
		cli.PrintUsage()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed() {
		// 🤔 How's addBlockData initalized by addBlockCmd.String() ?
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockChain(*createBlockchainAddress)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount == "" {
			sendCmd.Usage()
			runtime.Goexit()
		}
		amount, err := strconv.Atoi(*sendAmount)
		if err != nil {
			sendCmd.Usage()
			runtime.Goexit()
		}
		cli.send(*sendFrom, *sendTo, amount)
	}

	if printChainCmd.Parsed() {
		cli.PrintChain()
	}
}
