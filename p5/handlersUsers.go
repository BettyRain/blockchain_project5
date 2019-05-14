package p5

import (
	"../p3/data"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	//"../p1"
	"../p2"
	"../p3"
	"./dataPr5"
)

var SBC data.SyncBlockChain
var PatientList dataPr5.PatientList
var DoctorList dataPr5.DoctorList

var ifStarted bool
var kv map[string]string
var ID string
var Peers data.PeerList
var FIRST_NODE_ADDR = "http://localhost:8813"
var SELF_ADDR = "http://localhost:" + os.Args[1]

func init() {
	SBC = data.NewBlockChain()
	kv = make(map[string]string)
	IDnew, _ := strconv.ParseInt(os.Args[1], 10, 64)
	Peers = data.NewPeerList(int32(IDnew), 32)
	ID = os.Args[1]

	//id, _ := strconv.ParseInt(os.Args[1], 10, 64)
}

func Patient(w http.ResponseWriter, r *http.Request) {
	//View dataPr5 by personal code

	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "docInfo.html")
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		docID := r.FormValue("id")
		fmt.Println("HERE IS ID OF DOC", docID)
		//TODO: Verify DocID
		//TODO: if patients id = show
		respData := p3.GetSBC()
		SBC.UpdateEntireBlockChain(string(respData))
		var latestBlocks []p2.Block
		latestBlocks = SBC.GetLatestBlocks()
		var length int32
		res := ""
		mapData := ""
		num := 0
		isData := false
		if len(latestBlocks) > 0 {
			length = latestBlocks[0].GetHeight()
			if len(latestBlocks) > 1 {
				num = len(latestBlocks) - 1
			} else {
				num = len(latestBlocks)
			}

			for j := 0; j < num; j++ {
				res += "Patient results" + "\n"
				res += ID + "\n"
				latestBlock := latestBlocks[j]
				for i := length - 1; i >= 0; i-- {
					for k, v := range latestBlock.Value.GetKeyValue() {
						if k != "" {
							//TODO: data verificaction
							decryptedMap := PatientList.DecryptPatInfo(v)
							fmt.Println(decryptedMap)
							for ke, va := range decryptedMap {
								if ke == ID {
									mapData += "Patient ID: " + ke + ", Patient Data = " + va + "\n"
									isData = true
								}
							}

						}
					}
					if isData == true {
						res += latestBlock.ShowBlockData()
						res += mapData
						res += "\n"
					}
					parentBlock := p2.Block{}
					parentBlock = SBC.GetParentBlock(latestBlock)
					latestBlock = parentBlock
					isData = false
					mapData = ""
				}
				res += "\n"
			}
		}
		fmt.Fprintf(w, "%s\n", res)
	}
}

func Show(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n%s", Peers.Show())
}

func Patients(w http.ResponseWriter, r *http.Request) {
	//TODO: change doc ID
	//docID := "123"
	respData := p3.GetSBC()
	SBC.UpdateEntireBlockChain(string(respData))
	var latestBlocks []p2.Block
	latestBlocks = SBC.GetLatestBlocks()
	var length int32
	res := ""
	mapData := ""
	num := 0
	isData := false
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
				for k, v := range latestBlock.Value.GetKeyValue() {
					//(doc.VerifyDocSign(string(signature), patInfo, doc.selfId))
					//fmt.Println("VERIFY")
					/*if len(k)>0{
						fmt.Println("^^^^^^^^^^^^^^^^^^^^6")
						fmt.Println(k)
						fmt.Println(v)
						fmt.Println("^^^^^^^^^^^^^^^^^^^^6")
					}
						fmt.Println("ID = ", ID)
					fmt.Println(DoctorList.VerifyDocSign(k, v, ID))
					fmt.Println("VERIFY")
					if DoctorList.VerifyDocSign(k, v, ID) {
						mapData += "Patient ID: " + k + ", Patient Data = " + v + "\n"
						isData = true
					}*/

					if k != "nil" {
						isVerified := DoctorList.VerifyDocSign(k, v, ID)
						fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^")
						fmt.Println(isVerified)
						fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^")
						//TODO: data verificaction
						fmt.Println("IT IS NOT")
						decryptedMap := PatientList.DecryptPatInfo(v)
						fmt.Println(decryptedMap)
						for ke, va := range decryptedMap {
							mapData += "Patient ID: " + ke + ", Patient Data = " + va + "\n"
							isData = true
						}

					}
				}
				if isData == true {
					res += latestBlock.ShowBlockData()
					res += mapData
					res += "\n"
				}
				parentBlock := p2.Block{}
				parentBlock = SBC.GetParentBlock(latestBlock)
				latestBlock = parentBlock
				isData = false
				mapData = ""
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

		prMap := PatientList.GetPublicKeys()
		_, isExist := prMap[id]
		if isExist == false {
			PatientList.Register(id)
		}

		dataPool := dataPr5.AddToPool(ID, id, info, PatientList, DoctorList)

		p3.ForwardNewData(dataPool)

		//fmt.Println(kv)
		fmt.Fprintf(w, "Sent to miners\n")
		fmt.Fprintf(w, "Patient ID = %s\n", id)
		fmt.Fprintf(w, "Patient Information = %s\n", info)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

//Create public-private keys for patients
//Registation
func StartPat(w http.ResponseWriter, r *http.Request) {
	patId := os.Args[1]
	PatientList = dataPr5.NewPatientList(patId)
	ifStarted = true
	fmt.Fprintf(w, "%s\n", "You are registered, you can use /patient to see your information from a special doctor")
}

//Create public-private keys for doctors
//Registation
func StartDoc(w http.ResponseWriter, r *http.Request) {
	docId := os.Args[1]
	DoctorList = dataPr5.NewDoctorList(docId)
	ifStarted = true
	FirstHeartBeat()
	go StartHeartBeat()
	fmt.Fprintf(w, "%s\n", "You are registered, you can use /add to add new information or /patiens to see added information")
}

func ForwardPatientList(list dataPr5.PatientList) {
	//TODO: send patients list to doctors
	//patientListJson, _ := json.Marshal(list)
	//for key, _ := range Peers.Copy() {
	//	http.Post(key+"/heartbeat/receive", "application/json", bytes.NewBuffer(heartBeatJson))
	//}
}

func StartHeartBeat() {
	for true {
		fmt.Println("MY HEARTBEAT")
		duration := time.Duration(10) * time.Second
		time.Sleep(duration)

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
			for key, _ := range Peers.Copy() {
				http.Post(key+"/heartbeat/receive", "application/json", bytes.NewBuffer(heartBeatJson))
			}
		}
	}
}

func HeartBeatReceive(w http.ResponseWriter, r *http.Request) {
	fmt.Println("I GOT IT")
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

	newHops := heartBeatDataNew.Hops - 1
	heartBeatDataNew.Hops = newHops
	if heartBeatDataNew.Hops > 0 {
		ForwardHeartBeat(heartBeatDataNew)
	}
	w.WriteHeader(http.StatusOK)
}

func ForwardHeartBeat(heartBeatData data.HeartBeatData) {
	heartBeatJson, _ := json.Marshal(heartBeatData)
	for key, _ := range Peers.Copy() {
		http.Post(key+"/heartbeat/receive", "application/json", bytes.NewBuffer(heartBeatJson))
	}
}

func FirstHeartBeat() {
	peerMapJson, err := Peers.PeerMapToJson()
	if err != nil {
		println(err)
	}
	if SELF_ADDR != FIRST_NODE_ADDR {
		fmt.Println("MY FIRST HEARTBEAT")
		heartBeatFirst := data.HeartBeatData{false, Peers.GetSelfId(), "", peerMapJson, SELF_ADDR, 0}
		heartBeatJsonFirst, _ := json.Marshal(heartBeatFirst)

		http.Post(FIRST_NODE_ADDR+"/heartbeat/receive", "application/json", bytes.NewBuffer(heartBeatJsonFirst))
	}
}
