package p2

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/sha3"
	"sort"
)

type BlockChain struct {
	Chain map[int32][]Block //This is a map which maps a block height to a list of blocks
	//MPT value is an instance of MerklePatriciaTrie. You would create a mpt like P1 does, then pass it to Initial().
	Length int32 //Length equals to the highest block height
}

func (blch *BlockChain) Get(height int32) ([]Block, bool) {
	//This function takes a height as the argument,
	//returns the list of blocks stored in that height or None if the height doesn't exist.
	getList, ok := blch.Chain[height]
	if ok == false {
		return nil, false
	} else {
		return getList, true
	}
}

func (blch *BlockChain) GetBlock(height int32, hash string) (Block, bool) {
	//This function takes a height as the argument,
	//returns the list of blocks stored in that height or None if the height doesn't exist.
	getList, ok := blch.Chain[height]
	if ok == false {
		return Block{}, false
	} else {
		for _, block := range getList {
			if block.GetHash() == hash {
				return block, true
			}
		}
		return Block{}, false
	}
}

func (blch *BlockChain) Insert(blc Block) {
	//This function takes a block as the argument, insert that block to the BlockChain.Chain map.
	heightNewBlock := blc.GetHeight()
	if len(blch.Chain) == 0 {
		//chain is null
		blockList := []Block{blc}
		blch.Chain = make(map[int32][]Block)
		blch.Chain[heightNewBlock] = blockList
		blch.Length = heightNewBlock
	} else {
		blockList, _ := blch.Get(heightNewBlock)
		if len(blockList) == 0 {
			//no such blocks with that height
			blockList = []Block{blc}
			if heightNewBlock > blch.Length {
				blch.Length = heightNewBlock
			}
		} else {
			//if block already exists in blockchain
			isEqual := false
			for _, block := range blockList {
				if block.GetHash() == blc.GetHash() {
					isEqual = true
				}
			}
			if !isEqual {
				//add to the end of the list
				blockList = append(blockList, blc)
			}
		}
		blch.Chain[heightNewBlock] = blockList
		if blch.Length < blc.GetHeight() {
			blch.Length = blc.GetHeight()
		}
	}
}

func (blch *BlockChain) EncodeToJSON() (string, error) {
	//This function iterates over all the blocks,
	//generate blocks' JsonString by the function you implemented previously,
	//and return the list of those JsonStritgns.
	jsonString := "[\n"
	for _, list_block := range blch.Chain {
		for _, block := range list_block {
			blockJson, err := block.EncodeToJSON()
			jsonString += blockJson
			jsonString += ","
			if err != nil {
				return "", err
			}
		}
	}
	leng := len(jsonString)
	jsonString = jsonString[:leng-1]
	jsonString += "]"
	return jsonString, nil
}

func DecodeFromJSON(jsonString string) BlockChain {
	//Description: This function is called upon a blockchain instance.
	//It takes a blockchain JSON string as input,
	//decodes the JSON string back to a list of block JSON strings,
	//decodes each block JSON string back to a block instance,
	//and inserts every block into the blockchain.
	blockch := NewBlockChain()
	decodedBlock := Block{}
	blocks := make([]jsonBlockSt, 0)
	err := json.Unmarshal([]byte(jsonString), &blocks)
	if err != nil {
		println(err)
	}
	for _, block := range blocks {
		decodedBlock = Block{}
		decodedBlock.DecodeFromStruct(block)
		blockch.Insert(decodedBlock)
	}
	if err != nil {
		return BlockChain{}
	}
	return blockch
}

func NewBlockChain() BlockChain {
	blockch := BlockChain{make(map[int32][]Block), 0}
	return blockch
}

func (bc *BlockChain) Show() string {
	rs := ""
	var idList []int
	for id := range bc.Chain {
		idList = append(idList, int(id))
	}
	sort.Ints(idList)
	for _, id := range idList {
		var hashs []string
		for _, block := range bc.Chain[int32(id)] {
			hashs = append(hashs, block.GetHash()+"<="+block.GetParentHash())
		}
		sort.Strings(hashs)
		rs += fmt.Sprintf("%v: ", id)
		for _, h := range hashs {
			rs += fmt.Sprintf("%s, ", h)
		}
		rs += "\n"
	}
	sum := sha3.Sum256([]byte(rs))
	rs = fmt.Sprintf("This is the BlockChain: %s\n", hex.EncodeToString(sum[:])) + rs
	return rs
}

func (bc *BlockChain) GetLatestBlocks() []Block {
	//This function returns the list of blocks of height "BlockChain.length".
	blocks, _ := bc.Get(bc.Length)
	return blocks
}

func (bc *BlockChain) GetParentBlock(blc Block) Block {
	//This function takes a block as the parameter, and returns its parent block.
	parentHash := blc.GetParentHash()
	height := blc.GetHeight()
	parentHeight := height - 1
	parentBlock, _ := bc.GetBlock(parentHeight, parentHash)
	return parentBlock
}
