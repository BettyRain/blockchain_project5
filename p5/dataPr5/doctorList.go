package dataPr5

import (
	"crypto/rsa"
	"fmt"
	"os"
	"sync"
)

type DoctorList struct {
	selfId string
	pubMap map[string]rsa.PublicKey
	//id -> public key
	prMap map[string]rsa.PrivateKey
	//id -> private key
	mux sync.Mutex
}

func NewDoctorList(id string) DoctorList {
	docList := DoctorList{}
	docList.Register(id)
	docList.mux = sync.Mutex{}
	return docList
}

func (doc *DoctorList) Register(id string) {
	doc.mux.Lock()
	defer doc.mux.Unlock()
	//pat.patID = id
	if os.Args[1] == "8813" {
		doc.pubMap = make(map[string]rsa.PublicKey)
		doc.prMap = make(map[string]rsa.PrivateKey)
	}
	privateKey, publicKey := GenerateKeys()
	doc.prMap[id] = *privateKey
	doc.pubMap[id] = *publicKey
	doc.selfId = id

	fmt.Println("---- NEW DOCTOR ----")
	fmt.Printf("privateKey=%v\n", privateKey)
	fmt.Printf("publicKey=%v\n", publicKey)
	fmt.Printf("SelfIdPat=%v\n", id)
	fmt.Println("---- NEW DOCTOR ----")
}
