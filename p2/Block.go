package p2

import (
	"../p1"
	"encoding/hex"
	"encoding/json"
	"golang.org/x/crypto/sha3"
	"strconv"
	"time"
)

//  Block.Header.Nonce: Add a field "Nonce" to Block.Header.
//  Update BlockJson structure, EncodeToJson() and DecodeFromJson() functions accordingly.

type Header struct {
	height     int32
	timestamp  int64 //The value must be in the UNIX timestamp format such as 1550013938
	hash       string
	parentHash string
	size       int32 //You have a mpt, you convert it to byte array, then size = len(byteArray)
	nonce      string
}

type Block struct {
	//Block{Header{Height, Timestamp, Hash, ParentHash, Size}, Value}
	Header Header
	Value  p1.MerklePatriciaTrie
}

type jsonBlockSt struct {
	Hash       string
	TimeStamp  int64
	Height     int32
	ParentHash string
	Size       int32
	Nonce      string
	Mpt        map[string]string
}

func Initial(new_height int32, new_parentHash string, new_value p1.MerklePatriciaTrie, new_nonce string) Block {
	//This function takes arguments(such as height, parentHash, and value of MPT type) and forms a block.
	time_str := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	new_timestamp, err := strconv.ParseInt(time_str, 10, 64)
	new_size := int32(len([]byte(new_value.String())))
	new_hash := strconv.FormatInt(int64(new_height), 10) + time_str + new_parentHash + new_value.GetRoot() + strconv.FormatInt(int64(new_size), 10)
	sum := sha3.Sum256([]byte(new_hash))
	new_header := Header{height: new_height, timestamp: new_timestamp, hash: hex.EncodeToString(sum[:]), parentHash: new_parentHash, size: new_size, nonce: new_nonce}
	new_block := Block{Value: new_value, Header: new_header}
	if err != nil {
		return Block{}
	}
	return new_block
}

func (decodedBlock *Block) DecodeFromJson(jsonString string) {
	//This function takes a string that represents the JSON value of a block as an input,
	//and decodes the input string back to a block instance.
	byt := []byte(jsonString)
	blocks := jsonBlockSt{}
	err := json.Unmarshal(byt, &blocks)
	println(err)
	decodedBlock.DecodeFromStruct(blocks)
}

func (blc *Block) EncodeToJSON() (string, error) {
	//This function encodes a block instance into a JSON format string.
	//Note that the block's value is an MPT, and you have to record all of the (key, value) pairs
	//that have been inserted into the MPT in your JSON string.

	group := &jsonBlockSt{
		Hash:       blc.Header.hash,
		TimeStamp:  blc.Header.timestamp,
		Height:     blc.Header.height,
		ParentHash: blc.Header.parentHash,
		Size:       blc.Header.size,
		Nonce:      blc.Header.nonce,
		Mpt:        blc.Value.GetKeyValue(),
	}
	b, err := json.Marshal(group)

	if err != nil {
		return "", err
	}

	return string(b), err
}

func (decodedBlock *Block) DecodeFromStruct(blocks jsonBlockSt) {
	//This function creates a Block from json Block structure
	//Create new mpt
	mptDecoded := p1.MerklePatriciaTrie{}
	for key, value := range blocks.Mpt {
		mptDecoded.Insert(key, value)
	}
	decodedBlock.Header.hash = blocks.Hash
	decodedBlock.Header.parentHash = blocks.ParentHash
	decodedBlock.Header.height = blocks.Height
	decodedBlock.Header.size = blocks.Size
	decodedBlock.Header.timestamp = blocks.TimeStamp
	decodedBlock.Header.nonce = blocks.Nonce
	decodedBlock.Value = mptDecoded
}

func (block *Block) GetParentHash() string {
	return block.Header.parentHash
}

func (block *Block) GetHash() string {
	return block.Header.hash
}

func (block *Block) GetHeight() int32 {
	return block.Header.height
}

func (block *Block) GetNonce() string {
	return block.Header.nonce
}

func (block *Block) GetMPTRoot() string {
	return block.Value.GetRoot()
}

func (block *Block) ShowBlock() string {
	res := "Height = " + strconv.FormatInt(int64(block.Header.height), 10)
	res += ", Timestamp = " + strconv.FormatInt(block.Header.timestamp, 10)
	res += ", Hash = " + block.Header.hash
	res += ", ParentHash = " + block.Header.parentHash
	res += ", MPT root = " + block.GetMPTRoot()
	res += ", Size = " + strconv.FormatInt(int64(block.Header.size), 10) + "\n"
	return res
}

func (block *Block) ShowMap() string {
	res := "Block № " + strconv.FormatInt(int64(block.Header.height), 10) + "\n"
	res += "Timestamp =" + strconv.FormatInt(block.Header.timestamp, 10) + "\n"
	for k, v := range block.Value.GetKeyValue() {
		res += "Patient ID: " + k + ", Patient Data = " + v + "\n"
	}
	res += "\n"
	return res
}

func (block *Block) ShowBlockData() string {
	res := "Block № " + strconv.FormatInt(int64(block.Header.height), 10) + "\n"

	//tm := time.Unix(block.Header.timestamp, 0)
	//TODO: change data format
	res += "Timestamp = " + strconv.FormatInt(block.Header.timestamp, 10) + "\n"
	//time.Now().UTC().UnixNano()
	//st := tm.UTC().String()
	//res += "Timestamp = " + st + "\n"
	return res
}
