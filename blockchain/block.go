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

// Blockchain ⛓
type BlockChain struct {
	Blocks []*Block // TODO: rather than stroing blocks in slice, reference blocks by their hash or value.
}

// Adds a block to the blockchain 👶
func (chain *BlockChain) AddBlock(data string) {
	indexLastBlock := len(chain.Blocks) - 1
	prevBlock := chain.Blocks[indexLastBlock]

	newBlock := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, newBlock)
}

// Creates the genesis block 🎅🧬
func Genesis() *Block {
	return CreateBlock("Genesis: The First Block", []byte{})
}

// Initializes a new blockchain 🤶
func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}
