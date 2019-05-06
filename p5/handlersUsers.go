package p5

import (
	"../p3/data"
	"fmt"
	"net/http"
	//"../p1"
	"../p2"
	"../p3"
	"./dataPr5"
)

var SBC data.SyncBlockChain
var Peers data.PeerList
var ifStarted bool
var kv map[string]string

func init() {
	SBC = data.NewBlockChain()
	kv = make(map[string]string)
	//id, _ := strconv.ParseInt(os.Args[1], 10, 64)
}

func Patient(w http.ResponseWriter, r *http.Request) {
	//View dataPr5 by personal code
	fmt.Fprintf(w, "View dataPr5 by patient\n")
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
	//Add patient dataPr5
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "addInfo.html")
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		id := r.FormValue("id")
		info := r.FormValue("info")
		kv[id] = info
		dataPr5.AddToPool("123", id, info)
		//fmt.Println(kv)
		fmt.Fprintf(w, "Sent to miners\n")
		fmt.Fprintf(w, "Patient ID = %s\n", id)
		fmt.Fprintf(w, "Patient Information = %s\n", info)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func VerifyHash(hash string) bool {
	//Compare Hash with new generated Hash of <PatID, [PatInfo]PKpat>
	// (Verify that dataPr5 hasn't been changed)
	return false
}
