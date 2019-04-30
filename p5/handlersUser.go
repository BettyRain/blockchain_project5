package p5

import (
	"../p3/data"
	"fmt"
	"net/http"

	//"../p1"
	"../p2"
	"../p3"
)

var SBC data.SyncBlockChain
var Peers data.PeerList
var ifStarted bool

func init() {
	SBC = data.NewBlockChain()
	//id, _ := strconv.ParseInt(os.Args[1], 10, 64)
}

func Patient(w http.ResponseWriter, r *http.Request) {
	//View data by personal code
}

func Patients(w http.ResponseWriter, r *http.Request) {
	respData := p3.GetSBC()
	SBC.UpdateEntireBlockChain(string(respData))
	var latestBlocks []p2.Block
	latestBlocks = SBC.GetLatestBlocks()
	var length int32
	res := ""
	num := 0
	if len(latestBlocks) > 0 {
		length = latestBlocks[0].GetHeight()
		if len(latestBlocks) > 1 {
			num = len(latestBlocks) - 1
		} else {
			num = len(latestBlocks)
		}

		for j := 0; j < num; j++ {
			res += "Patient results" + "\n"
			latestBlock := latestBlocks[j]
			for i := length - 1; i >= 0; i-- {
				res += latestBlock.ShowMap()
				parentBlock := p2.Block{}
				parentBlock = SBC.GetParentBlock(latestBlock)
				latestBlock = parentBlock
			}
			res += "\n"
		}
	}
	fmt.Fprintf(w, "%s\n", res)
}

func AddData(w http.ResponseWriter, r *http.Request) {
	//Add patient data
}
