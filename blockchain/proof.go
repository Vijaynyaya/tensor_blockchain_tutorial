/*
Consensus Algorithms:
- Proof Of Work: "A solution that is difficult to find but easy to verify."
*/

// Requirements:
// Given data A, find a number x (nonce) such as that the hash of x appended to A results is a number less than B (Target).
package blockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

const Difficulty = 12 // the greater the difficulty the smaller the (Target)

type ProofOfWork struct {
	Block  *Block
	Target *big.Int // The generated hash should be smaller than the target hash
}

// Initializes proof of work for a given block 👩‍🍳 by defining the target hash
func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1) // 1
	// Left shift the target ⬅
	// In binary, the target becomes 1 followed by 256 - Difficulty zeroes
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{b, target}
	return pow // 🍞
}

// 🧪👨‍🎨 Tests different unique values of nonce to arrive at the hash that satisfies 😌 the proof of work requirements
func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	// create a counter (nonce) which starts at 0
	nonce := 0

	for nonce < math.MaxInt64 {
		// Take the block and create a hash of the block's data plus 💑 the counter (nonce)
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])

		// check the hash to see if it meets a set of requirements (hash is smaller than the target)
		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Println()

	return nonce, hash[:]
}

// 🔎📐 Checks if the hash generated meets the requirements (Target)
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	// generate the hash with the proposed nonce
	data := pow.InitData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	// verify that the resulting hash is smaller 🤏 than the target hash
	return intHash.Cmp(pow.Target) == -1
}

// Concatenates the block data 👨, the previous hash 👩, the difficulty 👧 and the given nonce 👦 to a byte array
func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)
	return data // 👨‍👩‍👦‍👦
}
