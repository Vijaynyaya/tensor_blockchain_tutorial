/*
Blockchain is an open and public database.
Sensitive information such as accounts, balances, addresses, claims, senders and receivers should not be stored in it.
Thus, everything is derived from the inputs and outputs stored in the blockhain.
*/
package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

type Transaction struct {
	ID      []byte // hash of the transaction
	Inputs  []TxInput
	Outputs []TxOutput
}

// Hashes the transaction and sets that as the transaction's ID
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	// byte encode the transaction
	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	Handle(err)

	// hash the transaction
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// Output is indivisible.
type TxOutput struct {
	Value  int    // Value in tokens
	PubKey string // used to unlock the Value
}

// An output can only be unlocked by the address it pertains to.
func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}

// Input references a previous output.
type TxInput struct {
	ID  []byte // Transaction of the Output
	Out int    // Output's index
	Sig string // the data to be used with the Output's PubKey
}

// Like output, an input can only be unlocked by the address it pertains to.
func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}

// Creates a transaction transferring a given amount from one address to another
func NewTransaction(from, to string, amount int, chain *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	acc, validOutputs := chain.FindSpendableOutputs(from, amount)

	// Insufficient accumulated balance
	if amount > acc {
		log.Panic("Error: not enough funds")
	}

	// remember: validOutputs is map[string][]int
	for txid, outs := range validOutputs {
		// convert string to []byte
		txID, err := hex.DecodeString(txid)
		Handle(err)

		// spend outputs (debit)
		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	// credit the receiving address
	outputs = append(outputs, TxOutput{amount, to})

	// remember: outputs are indivisible
	if acc > amount {
		// credit the leftover
		outputs = append(outputs, TxOutput{acc - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}

// Coinbase Transaction: The first transaction. (Genesis Block)
// It has a single input that references an empty output with some arbitrary data for the signature.
// It has a single output.
// It has an attached reward(subsidy) which is released to the node that mines the coinbase.
func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	txin := TxInput{[]byte{}, -1, data}
	txout := TxOutput{100, to}

	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	tx.SetID()

	return &tx
}

// Checks if the transaction is a coinbase transaction
func (tx *Transaction) IsCoinbase() bool {
	hasSingleInput := len(tx.Inputs) == 1
	singleInputHasEmptyID := len(tx.Inputs[0].ID) == 0
	singleInputHasInvalidOutputIndex := tx.Inputs[0].Out == -1

	return hasSingleInput && singleInputHasEmptyID && singleInputHasInvalidOutputIndex
}
