type HistoriaClinica struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Paciente     string `json:"paciente"`
	CreationDate string `json:"creation_date"`
	Historia     string `json:"historia"`
	IsGenesis    bool   `json:"is_genesis"`
}

type Block struct {
	Pos       int
	Data      HistoriaClinica
	Timestamp string
	Hash      string
	PrevHash  string
}

func (b *Block) generateHash() {
	// get string val of the Data
	bytes, _ := json.Marshal(b.Data)
	// concatenate the dataset
	data := string(b.Pos) + b.Timestamp + string(bytes) + b.PrevHash
	hash := sha256.New()
	hash.Write([]byte(data))
	b.Hash = hex.EncodeToString(hash.Sum(nil))
}

type Blockchain struct {
	blocks []*Block
}

var BlockChain *Blockchain

func CreateBlock(prevBlock *Block, Item HistoriaClinica) *Block {
	block := &Block{}
	block.Pos = prevBlock.Pos + 1
	block.Timestamp = time.Now().String()
	block.Data = Item
	block.PrevHash = prevBlock.Hash
	block.generateHash()

	return block
}

func (bc *Blockchain) AddBlock(data HistoriaClinica) {
	// get previous block
	prevBlock := bc.blocks[len(bc.blocks)-1]
	// create new block
	block := CreateBlock(prevBlock, data)
	bc.blocks = append(bc.blocks, block)
}

func GenesisBlock() *Block {
	return CreateBlock(&Block{}, HistoriaClinica{IsGenesis: true})

}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{GenesisBlock()}}
}

func validBlock(block, prevBlock *Block) bool {

	if prevBlock.Hash != block.PrevHash {
		return false
	}

	if !block.validateHash(block.Hash) {
		return false
	}

	if prevBlock.Pos+1 != block.Pos {
		return false
	}
	return true
}

func (b *Block) validateHash(hash string) bool {
	b.generateHash()
	if b.Hash != hash {
		return false
	}
	return true
}

func getBlockchain(w http.ResponseWriter, r *http.Request) {
	jbytes, err := json.MarshalIndent(BlockChain.blocks, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	// write JSON string
	io.WriteString(w, string(jbytes))
}