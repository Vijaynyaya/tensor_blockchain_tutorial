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
package blockchain

import (
	"bytes"
	"encoding/gob"
)

/*
Blockchain is composed of blocks. 🧱
Each block contains the data to be persisted to the database and a hash associated with the block.
*/
type Block struct {
	Hash     []byte // #️⃣ Hash of the block
	Data     []byte // 📄 Data inside the block
	PrevHash []byte // #️⃣🖇 Last block's hash (back linked list)
	Nonce    int
}

// Creates a block 🏭
func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	// 👩‍⚖️ define proof of work (requirements)
	pow := NewProof(block)

	// 🏃‍♂️ work to meet the requirements (mine 👷‍♀️) (generate proof of work)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// Creates the genesis block 🎅🧬
func Genesis() *Block {
	return CreateBlock("Genesis: The First Block", []byte{})
}

// Serializes a block to a bytes
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	Handle(err)

	return res.Bytes()
}

// Deserailizes a block from bytes
func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	Handle(err)

	return &block
}
