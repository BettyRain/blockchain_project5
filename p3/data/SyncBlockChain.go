package data

import (
	"../../p1"
	"../../p2"
	"sync"
)

type SyncBlockChain struct {
	bc  p2.BlockChain
	mux sync.Mutex
}

func NewBlockChain() SyncBlockChain {
	return SyncBlockChain{bc: p2.NewBlockChain()}
}

func (sbc *SyncBlockChain) Get(height int32) ([]p2.Block, bool) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.Get(height)
}

func (sbc *SyncBlockChain) GetBlock(height int32, hash string) (p2.Block, bool) {

	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.GetBlock(height, hash)
}

func (sbc *SyncBlockChain) Insert(block p2.Block) {
	sbc.mux.Lock()
	if block.GetHeight() > 0 && block.GetHash() != "" {
		sbc.bc.Insert(block)
	}
	sbc.mux.Unlock()
}

//CheckParentHash() is used to check if the parent hash(or parent block) exist in the current blockchain
//when you want to insert a new block sent by others.
//For example, if you received a block of height 7,
//you should check if its parent block(of height 6) exist in your blockchain.
//If not, you should ask others to download that parent block of height 6 before inserting the block
//of height 7.

func (sbc *SyncBlockChain) CheckParentHash(insertBlock p2.Block) bool {
	sbc.mux.Lock()
	height := insertBlock.GetHeight()
	parentHeight := height - 1
	blockList, _ := sbc.bc.Get(parentHeight)
	sbc.mux.Unlock()
	for _, block := range blockList {
		if block.GetHash() == insertBlock.GetParentHash() {
			return true
		}
	}
	return false
}

func (sbc *SyncBlockChain) UpdateEntireBlockChain(blockChainJson string) {
	sbc.mux.Lock()
	sbc.bc = p2.DecodeFromJSON(blockChainJson)
	sbc.mux.Unlock()
}

func (sbc *SyncBlockChain) BlockChainToJson() (string, error) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.EncodeToJSON()
}

func (sbc *SyncBlockChain) GenBlock(mpt p1.MerklePatriciaTrie) p2.Block {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	height := sbc.bc.Length
	blockList, _ := sbc.bc.Get(height)
	if len(blockList) == 0 {
		blockList, _ = sbc.bc.Get(height - 1)
	}
	parentBlock := blockList[0]
	parentHash := parentBlock.GetHash()
	return p2.Initial(height+1, parentHash, mpt, "")
}

func (sbc *SyncBlockChain) Show() string {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.Show()
}

func (sbc *SyncBlockChain) GetLatestBlocks() []p2.Block {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.GetLatestBlocks()
}

func (sbc *SyncBlockChain) GetParentBlock(blc p2.Block) p2.Block {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.GetParentBlock(blc)
}

func (sbc *SyncBlockChain) ShowBlock() string {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.Show()
}

func (sbc *SyncBlockChain) GetLength() int32 {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.Length
}
