package dataPr5

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"os"
	"sync"
)

type DoctorList struct {
	selfId string
	pubMap map[string][]byte
	//id -> public key
	prMap map[string][]byte
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
	//pat.patID = id
	if len(doc.pubMap) < 1 {
		doc.pubMap = make(map[string][]byte)
		doc.prMap = make(map[string][]byte)
	}

	privateKey, publicKey := GenerateKeys()
	doc.prMap[id] = privateKey
	doc.pubMap[id] = publicKey
	doc.selfId = id

	fmt.Println("---- NEW DOCTOR ----")
	fmt.Printf("privateKey=%v\n", privateKey)
	fmt.Printf("publicKey=%v\n", publicKey)
	fmt.Printf("SelfIdDoc=%v\n", id)
	fmt.Println("---- NEW DOCTOR ----")
	defer doc.mux.Unlock()
}

func (doc *DoctorList) SignByDoc(patInfo string, docID string) []byte {
	//preparing data for signatures
	rng := rand.Reader
	message := []byte(patInfo)
	hashed := sha256.Sum256(message)

	//get doctor's private key
	value, _ := doc.prMap[docID]
	key := BytesToPrivateKey(value)

	//get the signature
	signature, err := rsa.SignPKCS1v15(rng, key, crypto.SHA256, hashed[:])
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(doc.VerifyDocSign(string(signature), patInfo, doc.selfId))
	return signature
}

func (doc *DoctorList) VerifyDocSign(pastID string, patInfo string, docID string, sign []byte) bool {
	message := []byte(patInfo)
	hashed := sha256.Sum256(message)
	value, _ := doc.prMap[docID]
	key := BytesToPrivateKey(value)

	err := rsa.VerifyPKCS1v15(&key.PublicKey, crypto.SHA256, hashed[:], sign)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from verification: %s\n", err)
		return false
	} else {
		return true
	}
}
