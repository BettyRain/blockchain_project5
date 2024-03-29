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
	PubMap map[string][]byte
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

//Doctor's registration, public-private keys computations
func (doc *DoctorList) Register(id string) {
	doc.mux.Lock()
	if len(doc.PubMap) < 1 {
		doc.PubMap = make(map[string][]byte)
		doc.prMap = make(map[string][]byte)
	}
	privateKey, publicKey := GenerateKeys()
	doc.prMap[id] = privateKey
	doc.PubMap[id] = publicKey
	doc.selfId = id

	fmt.Println("---- NEW DOCTOR ----")
	fmt.Printf("privateKey=%v\n", privateKey)
	fmt.Printf("publicKey=%v\n", publicKey)
	fmt.Printf("SelfIdDoc=%v\n", id)
	fmt.Println("---- NEW DOCTOR ----")
	defer doc.mux.Unlock()
}

//Doctor's signature creation
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
	return signature
}

//Doctor verifies his signature (to have an ability to access the data)
func (doc *DoctorList) VerifyDocSign(pastID string, patInfo string, docID string, sign []byte) bool {
	message := []byte(patInfo)
	hashed := sha256.Sum256(message)
	value, _ := doc.PubMap[docID]
	key := BytesToPublicKey(value)

	err := rsa.VerifyPKCS1v15(key, crypto.SHA256, hashed[:], sign)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from verification: %s\n", err)
		return false
	} else {
		return true
	}
}

//Patient verifies his doctor's signature
func (doc *DoctorList) VerifyDocSignForPatient(docID string, patInfo string, sign []byte) bool {
	message := []byte(patInfo)
	hashed := sha256.Sum256(message)
	value, exist := doc.PubMap[docID]
	if !exist {
		return false
	}
	key := BytesToPublicKey(value)

	err := rsa.VerifyPKCS1v15(key, crypto.SHA256, hashed[:], sign)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from verification: %s\n", err)
		return false
	} else {
		return true
	}
}
