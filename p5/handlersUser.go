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
var addPage = " <!DOCTYPE html> <html> <head> <meta charset='UTF-8' /> </head> <body> <div> <form method='POST' action='/add'><label>Name</label><input name='name' type='text' value='' /><label>Address</label><input name='address' type='text' value='' /><input type='submit' value='submit' /></form></div></body></html>"
var kv map[string]string

func init() {
	SBC = data.NewBlockChain()
	kv = make(map[string]string)
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
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "addInfo.html")
	case "POST":
		id := r.FormValue("id")
		info := r.FormValue("info")

		kv[id] = info
		http.ServeFile(w, r, "postInfo.html")
		//if err := r.ParseForm(); err != nil {
		//	fmt.Fprintf(w, "ParseForm() err: %v", err)
		//	return
		//}
		//fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		//fmt.Println(r.PostForm)
		//		fmt.Println(r)
		//		name := r.FormValue("id")
		//		address := r.FormValue("info")
		//fmt.Fprintf(w, "Patient ID = %s\n", name)
		//		fmt.Fprintf(w, "Patient Information = %s\n", address)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func SendToMiners(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Send to Miners")
	fmt.Println(kv)

	//send hashmap to miners
	//so they will add in into block
}

//Сделаем две кнопки, одна будет снова вызывать форму добавления (что будет идти в хэшмеп,
// другая отправлять в блок и майнерам
