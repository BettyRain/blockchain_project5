package dataPr5

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"os"
	"reflect"
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

func (doc *DoctorList) SignByDoc(patInfo string) string {
	//preparing data for signatures
	rng := rand.Reader
	message := []byte(patInfo)
	hashed := sha256.Sum256(message)

	//get doctor's private key
	value, errr := doc.prMap[doc.selfId]
	fmt.Println("HERE IS DOC")
	fmt.Println(errr)
	fmt.Println(doc.selfId)
	fmt.Println(doc.prMap)
	fmt.Println(value)
	fmt.Println("HERE IS DOC")

	key := BytesToPrivateKey(value)

	//get the signature
	signature, err := rsa.SignPKCS1v15(rng, key, crypto.SHA256, hashed[:])
	if err != nil {
		fmt.Println(err)
	}
	/*	fmt.Println("-----------------------")
		fmt.Println(string(signature))
		fmt.Println(patInfo)

		fmt.Println("-----------------------")*/

	fmt.Println(doc.VerifyDocSign(string(signature), patInfo, doc.selfId))
	return string(signature)
}

func (doc *DoctorList) VerifyDocSign(signature string, patInfo string, docID string) bool {
	rng := rand.Reader
	message := []byte(patInfo)
	hashed := sha256.Sum256(message)
	value, _ := doc.prMap[doc.selfId]
	key := BytesToPrivateKey(value)

	sign2, err := rsa.SignPKCS1v15(rng, key, crypto.SHA256, hashed[:])
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%%%55")
	fmt.Println(reflect.DeepEqual(sign2, signature))
	fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%%%55")
	/*if reflect.DeepEqual(sign2, signature) {
		return true
	}
	return false*/

	hashed2 := sha256.Sum256([]byte(message))
	value2, _ := doc.pubMap[docID]
	if value == nil {
		return false
	}
	key2 := BytesToPublicKey(value2)

	//verify the signature, respond true if verified
	err = rsa.VerifyPKCS1v15(key2, crypto.SHA256, hashed2[:], []byte(signature))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from verification: %s\n", err)
		fmt.Println("91aa")
		return false
	}
	return true
}
