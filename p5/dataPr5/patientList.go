package dataPr5

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"strconv"
	"sync"
)

//Public Key, Private Key, PatID, lock
type PatientList struct {
	patID  string
	pubMap map[string]int
	//id -> public key
	prMap map[string]int
	//id -> private key
	mux sync.Mutex
}

func NewPatientList(id string, maxLength int32) PatientList {
	peerList := PatientList{}
	peerList.Register(id)
	peerList.mux = sync.Mutex{}
	return peerList
}

func (pat *PatientList) CreateKey(id string) {
	pat.mux.Lock()
	//create keys
	privateKey := 0
	publicKey := 0
	pat.pubMap[id] = publicKey
	pat.prMap[id] = privateKey
	pat.mux.Unlock()
}

func (pat *PatientList) Register(id string) {
	pat.mux.Lock()
	defer pat.mux.Unlock()
	pat.patID = id
	fmt.Printf("SelfId=%v\n", id)
}

func (pat *PatientList) GetSelfId() string {
	pat.mux.Lock()
	defer pat.mux.Unlock()
	return pat.patID
}

func (pat *PatientList) Show() string {
	pat.mux.Lock()
	defer pat.mux.Unlock()
	res := "This is a Patients Map: \n"
	for key, value := range pat.prMap {
		res += "ID = " + key + ", Public Key = " + strconv.Itoa(int(value)) + "\n"
	}
	return res
}

func GenerateKeys() {
	reader := rand.Reader
	bitSize := 2048

	key, _ := rsa.GenerateKey(reader, bitSize)
	publicKey := key.PublicKey

	fmt.Println("Public Key")
	fmt.Println(publicKey)
	fmt.Println("Private Key")
	fmt.Println(key)
}
