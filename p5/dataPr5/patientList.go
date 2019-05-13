package dataPr5

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"
	"strconv"
	"sync"
)

//Public Key, Private Key, PatID, lock
type PatientList struct {
	//patID  string
	pubMap map[string]rsa.PublicKey
	//id -> public key
	prMap map[string]rsa.PrivateKey
	//id -> private key
	mux sync.Mutex
}

func NewPatientList(id string) PatientList {
	patList := PatientList{}
	patList.Register(id)
	patList.mux = sync.Mutex{}
	return patList
}

func (pat *PatientList) Register(id string) {
	pat.mux.Lock()
	defer pat.mux.Unlock()
	//pat.patID = id
	if os.Args[1] == "9913" {
		pat.pubMap = make(map[string]rsa.PublicKey)
		pat.prMap = make(map[string]rsa.PrivateKey)
	}
	privateKey, publicKey := GenerateKeys()
	pat.prMap[id] = privateKey
	pat.pubMap[id] = publicKey

	fmt.Println("---- NEW PATIENT ----")
	fmt.Printf("privateKey=%v\n", privateKey)
	fmt.Printf("publicKey=%v\n", publicKey)
	fmt.Printf("SelfIdPat=%v\n", id)
	fmt.Println("---- NEW PATIENT ----")
}

func (pat *PatientList) Show() string {
	pat.mux.Lock()
	defer pat.mux.Unlock()
	res := "This is a Patients Map: \n"
	for key, value := range pat.pubMap {
		res += "ID = " + key + ", Public Key = " + strconv.Itoa(int(value.E)) + "\n"
	}
	return res
}

func GenerateKeys() (rsa.PrivateKey, rsa.PublicKey) {
	reader := rand.Reader
	bitSize := 2048
	privateKey, _ := rsa.GenerateKey(reader, bitSize)
	publicKey := privateKey.PublicKey
	return *privateKey, publicKey
}
