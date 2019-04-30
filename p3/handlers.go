package p3

import (
	"../p1"
	"../p2"
	//"../p5"
	"./data"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/sha3"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var TA_SERVER = "http://localhost:6688"
var REGISTER_SERVER = TA_SERVER + "/peer"
var BC_DOWNLOAD_SERVER = TA_SERVER + "/upload"
var SELF_ADDR = "http://localhost:" + os.Args[1]
var FIRST_NODE_ADDR = "http://localhost:6686"
var BLOCKCHAIN_JSON = "[{\"hash\": \"3ff3b4efe9177f705550231079c2459ba54a22d340a517e84ec5261a0d74ca48\", \"timeStamp\": 1234567890, \"height\": 1, \"parentHash\": \"genesis\", \"size\": 1174, \"mpt\": {\"hello\": \"world\", \"charles\": \"ge\"}}, {\"hash\": \"24cf2c336f02ccd526a03683b522bfca8c3c19aed8a1bed1bbc23c33cd8d1159\", \"timeStamp\": 1234567890, \"height\": 2, \"parentHash\": \"3ff3b4efe9177f705550231079c2459ba54a22d340a517e84ec5261a0d74ca48\", \"size\": 1231, \"mpt\": {\"hello\": \"world\", \"charles\": \"ge\"}}]"

var SBC data.SyncBlockChain
var Peers data.PeerList
var ifStarted bool

//Create SyncBlockChain and PeerList instances.
func init() {
	// This function will be executed before everything else.
	// Do some initialization here.
	SBC = data.NewBlockChain()
	// When server works
	// Peers = data.NewPeerList(Register(), 32)
	// While server doesn't work -> use port as id
	id, _ := strconv.ParseInt(os.Args[1], 10, 64)
	Peers = data.NewPeerList(int32(id), 32)
	ifStarted = true
}

// Register ID, download BlockChain, start HeartBea
func Start(w http.ResponseWriter, r *http.Request) {
	if os.Args[1] == "6686" {
		SBC.UpdateEntireBlockChain(BLOCKCHAIN_JSON)
		Upload(w, r)
	} else {
		Download()
	}
	FirstHeatdBeat()
	go StartHeartBeat()
	go StartTryingNonces()
}

//Display peerList and sbc
//Shows the PeerMap and the BlockChain.
func Show(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n%s", Peers.Show(), SBC.Show())
}

// Register to TA's server, get an ID
func Register() int32 {
	resp, err := http.Get(REGISTER_SERVER)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer resp.Body.Close()
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			fmt.Printf("%s", err2)
			os.Exit(1)
		}
		content, _ := strconv.ParseInt(string(bodyBytes), 10, 64)
		return int32(content)
	}
	return 0
}

// Download blockchain from TA server
// Download the current BlockChain from your own first node(can be hardcoded).
// It's ok to use this function only after launching a new node.
// You may not need it after node starts heartBeats.
func Download() {
	//call upload on get from first node (self address port)
	resp, err := http.Get(FIRST_NODE_ADDR + "/upload")

	if err != nil {
		println(err)
		os.Exit(1)
	}
	if resp.StatusCode == 200 {
		respData, erro := ioutil.ReadAll(resp.Body)
		if erro != nil {
			println(err)
			os.Exit(1)
		}
		SBC.UpdateEntireBlockChain(string(respData))
	}

}

func GetSBC() string {
	//call upload on get from first node (self address port)
	resp, err := http.Get(FIRST_NODE_ADDR + "/upload")

	if err != nil {
		println(err)
		os.Exit(1)
	}
	if resp.StatusCode == 200 {
		respData, erro := ioutil.ReadAll(resp.Body)
		if erro != nil {
			println(err)
			os.Exit(1)
		}
		return string(respData)
		//SBC.UpdateEntireBlockChain(string(respData))
	}
	return ""
}

// Upload blockchain to whoever called this method, return jsonStr
// Return the BlockChain's JSON. And add the remote peer into the PeerMap.
func Upload(w http.ResponseWriter, r *http.Request) {
	blockChainJson, err := SBC.BlockChainToJson()
	if err != nil {
		//	data.PrintError(err, "Upload")
		fmt.Printf("%s", err)
	}
	fmt.Fprint(w, blockChainJson)
}

// Upload a block to whoever called this method, return jsonStr
// Return the Block's JSON.
func UploadBlock(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")
	height, err := strconv.ParseInt(paths[2], 10, 64)
	if err != nil {
		fmt.Fprint(w, http.StatusInternalServerError)
	}
	hash := paths[3]
	//if you don't have the block, return HTTP 204: StatusNoContent;
	//if there's an error, return HTTP 500: InternalServerError.
	block := p2.Block{}
	exist := true
	block, exist = SBC.GetBlock(int32(height), hash)
	jsonBlock := ""
	if exist {
		if block.GetHeight() > 0 {
			jsonBlock, err = block.EncodeToJSON()
		}
		if err != nil {
			fmt.Fprint(w, http.StatusInternalServerError)
		} else {
			fmt.Fprint(w, jsonBlock)
		}

	} else {
		fmt.Fprint(w, http.StatusNoContent)
	}
}

// Received a heartbeat
// Add the remote address, and the PeerMapJSON into local PeerMap.
// Then check if the HeartBeatData contains a new block.
// If so, do these: (1) check if the parent block exists.
// If not, call AskForBlock() to download the parent block.
// (2) it verifies the nonce, insert the new block from HeartBeatData.
// (3) HeartBeatData.hops minus one, and if it's still bigger than 0,
// //call ForwardHeartBeat() to forward this heartBeat to all peers.
//For the respond of HeartBeats, return 200 is enough
func HeartBeatReceive(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		println(err)
	}
	heartBeatData := data.HeartBeatData{}
	err = json.Unmarshal([]byte(string(body)), &heartBeatData)
	if err != nil {
		println(err)
	}
	heartBeatDataNew := heartBeatData
	Peers.InjectPeerMapJson(heartBeatDataNew.PeerMapJson, SELF_ADDR)
	Peers.Add(heartBeatDataNew.Addr, heartBeatDataNew.Id)

	if heartBeatDataNew.IfNewBlock {
		fmt.Println("RECEIVE A BLOCK")
		newBlock := p2.Block{}
		newBlock.DecodeFromJson(heartBeatDataNew.BlockJson)

		if VerifyNonce(newBlock.GetNonce(), newBlock.GetParentHash(), newBlock.GetMPTRoot()) {
			fmt.Println("NONCE VERIFIED")
			if !SBC.CheckParentHash(newBlock) {
				AskForBlock(newBlock.GetHeight()-1, newBlock.GetParentHash())
			}
			SBC.Insert(newBlock)
		} else {
			fmt.Println("NONCE NOT VERIFIED")
		}
	}
	newHops := heartBeatDataNew.Hops - 1
	heartBeatDataNew.Hops = newHops
	if heartBeatDataNew.Hops > 0 {
		ForwardHeartBeat(heartBeatDataNew)
	}
	w.WriteHeader(http.StatusOK)
}

// Ask another server to return a block of certain height and hash
// Loop through all peers in local PeerMap to download a block. As soon as one peer returns the block, stop the loop.
func AskForBlock(height int32, hash string) {
	Peers.Rebalance()

	for key, _ := range Peers.Copy() {
		resp, err := http.Get(key + "/block/" + string(height) + "/" + hash)
		if err != nil {
			println(err)
			os.Exit(1)
		}
		if resp.StatusCode == 200 {
			respData, erro := ioutil.ReadAll(resp.Body)
			if erro != nil {
				println(err)
				os.Exit(1)
			}
			respBlc := p2.Block{}
			respBlc.DecodeFromJson(string(respData))
			if !SBC.CheckParentHash(respBlc) {
				AskForBlock(respBlc.GetHeight()-1, respBlc.GetParentHash())
			}
			SBC.Insert(respBlc)
			break
		}
	}

}

//Send the HeartBeatData to all peers in local PeerMap.
//send heartbeat to all local nodes in peermap, I  SHOULD rebalance before sending
func ForwardHeartBeat(heartBeatData data.HeartBeatData) {
	Peers.Rebalance()
	heartBeatJson, _ := json.Marshal(heartBeatData)
	for key, _ := range Peers.Copy() {
		http.Post(key+"/heartbeat/receive", "application/json", bytes.NewBuffer(heartBeatJson))
	}
}

func FirstHeatdBeat() {
	peerMapJson, err := Peers.PeerMapToJson()
	if err != nil {
		println(err)
	}
	if SELF_ADDR != FIRST_NODE_ADDR {
		heartBeatFirst := data.HeartBeatData{false, Peers.GetSelfId(), "", peerMapJson, SELF_ADDR, 0}
		heartBeatJsonFirst, _ := json.Marshal(heartBeatFirst)
		http.Post(FIRST_NODE_ADDR+"/heartbeat/receive", "application/json", bytes.NewBuffer(heartBeatJsonFirst))
	}
}

func StartHeartBeat() {
	for true {
		duration := time.Duration(5) * time.Second
		time.Sleep(duration)
		Peers.Rebalance()

		peerMapJson, err := Peers.PeerMapToJson()
		if err != nil {
			println(err)
		}

		if len(Peers.Copy()) > 0 {
			heartBeat := data.PrepareHeartBeatData(Peers.GetSelfId(), peerMapJson, SELF_ADDR)
			heartBeatJson, err := json.Marshal(heartBeat)
			if err != nil {
				println(err)
			}
			Peers.Rebalance()
			for key, _ := range Peers.Copy() {
				http.Post(key+"/heartbeat/receive", "application/json", bytes.NewBuffer(heartBeatJson))
			}
		}
	}
}

//This function prints the current canonical chain, and chains of all forks if there are forks.
func Canonical(w http.ResponseWriter, r *http.Request) {
	var latestBlocks []p2.Block
	latestBlocks = SBC.GetLatestBlocks()
	var length int32
	res := ""
	if len(latestBlocks) > 0 {
		length = latestBlocks[0].GetHeight()

		if len(latestBlocks) > 1 {
			for j := 0; j < len(latestBlocks); j++ {
				res += "BlockChain #" + strconv.Itoa(j+1) + "\n"
				latestBlock := latestBlocks[j]
				for i := length - 1; i >= 0; i-- {
					res += latestBlock.ShowBlock()
					parentBlock := p2.Block{}
					parentBlock = SBC.GetParentBlock(latestBlock)
					latestBlock = parentBlock
				}
				res += "\n"
			}
		} else {
			res = "BlockChain #1" + "\n"
			latestBlock := latestBlocks[0]
			for i := length - 1; i >= 0; i-- {
				res += latestBlock.ShowBlock()
				parentBlock := p2.Block{}
				parentBlock = SBC.GetParentBlock(latestBlock)
				latestBlock = parentBlock
			}
		}
	}
	fmt.Fprintf(w, "%s\n", res)
}

func StartTryingNonces() {
start:
	newMPT := p1.MerklePatriciaTrie{}
	newMPT = data.GenerateRandomMPT()
	//TODO: new mpt should be with special data
	for true {
		blocks := SBC.GetLatestBlocks()
		blocksCount := len(blocks)
		parentBlock := p2.Block{}
		if blocksCount > 1 {
			num := rand.Intn(blocksCount)
			parentBlock = blocks[num]
		} else {
			parentBlock = blocks[0]
		}
		nonce := GenerateNonce()
		for !VerifyNonce(nonce, parentBlock.GetHash(), newMPT.GetRoot()) {
			if parentBlock.GetHeight() < SBC.GetLength() {
				goto start
			}
			nonce = GenerateNonce()
		}

		newHeight := parentBlock.GetHeight() + 1
		newBlock := p2.Initial(newHeight, parentBlock.GetHash(), newMPT, nonce)
		blockJson, _ := newBlock.EncodeToJSON()
		peerMapJson, _ := Peers.PeerMapToJson()
		heartBeat := data.NewHeartBeatData(true, Peers.GetSelfId(), blockJson, peerMapJson, SELF_ADDR)
		ForwardHeartBeat(heartBeat)
		fmt.Println("IT IS MY BLOCK WOHOO")
	}
}

func GenerateNonce() string {
	nonce := make([]byte, 8)
	if _, err := rand.Read(nonce); err != nil {
		return ""
	}
	return hex.EncodeToString(nonce)
}

func VerifyNonce(nonce string, hash string, root string) bool {
	concat := nonce + hash + root
	hasher := sha3.Sum256([]byte(concat))
	sha := hex.EncodeToString(hasher[:])
	if strings.HasPrefix(sha, "00000") {
		return true
	} else {
		return false
	}
}
