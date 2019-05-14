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

func (doc *DoctorList) VerifyDocSign(signature string, message string, docID string) bool {
	fmt.Println("-----------------------")
	fmt.Println(signature)
	fmt.Println(message)
	fmt.Println(docID)
	fmt.Println("-----------------------")
	hashed := sha256.Sum256([]byte(message))
	value, _ := doc.pubMap[docID]
	if value == nil {
		fmt.Println("82aa")
		return false
	}
	key := BytesToPublicKey(value)
	fmt.Println("KEY")
	fmt.Println(key)
	//verify the signature, respond true if verified
	err := rsa.VerifyPKCS1v15(key, crypto.SHA256, hashed[:], []byte(signature))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from verification: %s\n", err)
		fmt.Println("91aa")
		return false
	}
	return true
}
