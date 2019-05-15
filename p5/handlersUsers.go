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

//All data initialization
func init() {
	SBC = data.NewBlockChain()
	kv = make(map[string]string)
	IDnew, _ := strconv.ParseInt(os.Args[1], 10, 64)
	Peers = data.NewPeerList(int32(IDnew), 32)
	ID = os.Args[1]
	PatientList = dataPr5.PatientList{}
	PatientList.PrMap = make(map[string][]byte)
	PatientList.PubMap = make(map[string][]byte)
	DoctorList = dataPr5.DoctorList{}
	DoctorList.PubMap = make(map[string][]byte)
}

//Shows the patient's data from a definite doctor from canonical chain
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
				res += "My ID: " + ID + "\n"
				latestBlock := latestBlocks[j]
				for i := length - 1; i >= 0; i-- {
					for k, v := range latestBlock.Value.GetKeyValue() {

						if k != "nil" {
							signature := latestBlock.Value.GetSignature()[k]
							if k == docID {
								isVerified := DoctorList.VerifyDocSignForPatient(k, v, signature)
								if isVerified {
									decryptedMap := PatientList.DecryptPatInfo(v)
									for ke, va := range decryptedMap {
										mapData += "Patient ID: " + ke + ", Patient Data = " + va + "\n"
										isData = true
									}
								}
							}
						}
					}
					if isData == true {
						res += latestBlock.ShowBlockData()
						res += mapData
						res += "Doctor ID: " + docID + "\n"
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

//Shows doctors peers and patients
func Show(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n%s", Peers.Show(), PatientList.Show())
}

//Shows all patients of one doctor from canonical chain
func Patients(w http.ResponseWriter, r *http.Request) {
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
			res += "Doctor ID: " + ID + "\n"
			latestBlock := latestBlocks[j]
			for i := length - 1; i >= 0; i-- {
				for k, v := range latestBlock.Value.GetKeyValue() {
					if k != "nil" {
						signature := latestBlock.Value.GetSignature()[k]
						isVerified := DoctorList.VerifyDocSign(k, v, ID, signature)
						if isVerified {
							decryptedMap := PatientList.DecryptPatInfo(v)
							for ke, va := range decryptedMap {
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

//Adds patient's data, calls hasher, sends hashed and signed data to the queue
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

		prMap := PatientList.PubMap
		_, isExist := prMap[id]
		if isExist == false {
			PatientList.Register(id)
		}

		dataPool := dataPr5.AddToPool(ID, id, info, PatientList, DoctorList)
		p3.ForwardNewData(dataPool)

		fmt.Fprintf(w, "Sent to miners\n")
		fmt.Fprintf(w, "Patient ID = %s\n", id)
		fmt.Fprintf(w, "Patient Information = %s\n", info)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

//Create public-private keys for patients
//Special patient's registation
func StartPat(w http.ResponseWriter, r *http.Request) {
	patId := os.Args[1]
	PatientList = dataPr5.NewPatientList(patId)
	ifStarted = true
	jsonPublicMap, _ := json.Marshal(PatientList.PubMap)
	jsonPrivateMap, _ := json.Marshal(PatientList.PrMap)
	patientMessage := dataPr5.PatientMessage{JsonPubMap: jsonPublicMap, JsonPrMap: jsonPrivateMap, Hops: 3}
	FirstForwardPatientList(patientMessage)
	fmt.Fprintf(w, "%s\n", "You are registered, you can use /patient to see your information from a special doctor")
}

//Create public-private keys for doctors
//Special doctor's registation
func StartDoc(w http.ResponseWriter, r *http.Request) {
	docId := os.Args[1]
	DoctorList = dataPr5.NewDoctorList(docId)
	ifStarted = true
	FirstHeartBeat()
	go StartHeartBeat()
	fmt.Fprintf(w, "%s\n", "You are registered, you can use /add to add new information or /patiens to see added information")
}

//Send Patient List to all doctors
func ForwardPatientList(list dataPr5.PatientMessage) {
	patientListJson, _ := json.Marshal(list)
	for key, _ := range Peers.Copy() {
		http.Post(key+"/patientlist/receive", "application/json", bytes.NewBuffer(patientListJson))
	}
}

//Receive patient list and add info to each doc's patient list
func PatientListReceive(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		println(err)
	}
	patientData := dataPr5.PatientMessage{}
	err = json.Unmarshal([]byte(string(body)), &patientData)
	if err != nil {
		println(err)
	}
	privateMap := make(map[string][]byte)
	publicMap := make(map[string][]byte)
	err = json.Unmarshal(patientData.JsonPrMap, &privateMap)
	err = json.Unmarshal(patientData.JsonPubMap, &publicMap)
	PatientList.AddNewPatient(publicMap, privateMap)
	newHops := patientData.Hops - 1
	patientData.Hops = newHops

	if patientData.Hops > 0 {
		ForwardPatientList(patientData)
	}
	w.WriteHeader(http.StatusOK)
}

//Send patient list to first node
func FirstForwardPatientList(list dataPr5.PatientMessage) {
	patientListJson, _ := json.Marshal(list)
	http.Post(FIRST_NODE_ADDR+"/patientlist/receive", "application/json", bytes.NewBuffer(patientListJson))
}

//Heart beat from doctors to doctors: send their ips and patients lists
func StartHeartBeat() {
	for true {
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
		if len(PatientList.PubMap) > 0 {
			jsonPublicMap, _ := json.Marshal(DoctorList.PubMap)
			for key := range PatientList.PubMap {
				realKey := "http://localhost:" + key
				http.Post(realKey+"/doctorlist/receive", "application/json", bytes.NewBuffer(jsonPublicMap))
			}
		}
	}
}

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

	newHops := heartBeatDataNew.Hops - 1
	heartBeatDataNew.Hops = newHops
	if heartBeatDataNew.Hops > 0 {
		ForwardHeartBeat(heartBeatDataNew)
	}
	w.WriteHeader(http.StatusOK)
}

//Receive doctor list by patients, so patients will know doc's public keys
func DoctorListReceive(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		println(err)
	}
	var m map[string][]byte
	err = json.Unmarshal([]byte(string(body)), &m)
	if err != nil {
		println(err)
	}
	for key, value := range m {
		DoctorList.PubMap[key] = value
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
		heartBeatFirst := data.HeartBeatData{false, Peers.GetSelfId(), "", peerMapJson, SELF_ADDR, 0}
		heartBeatJsonFirst, _ := json.Marshal(heartBeatFirst)
		http.Post(FIRST_NODE_ADDR+"/heartbeat/receive", "application/json", bytes.NewBuffer(heartBeatJsonFirst))
	}
}
