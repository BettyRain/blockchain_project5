package dataPr5

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	//"crypto/x509"
	//"encoding/pem"
	"fmt"
	"os"
	"strconv"
	"sync"
)

//Public Key, Private Key, PatID, lock
type PatientList struct {
	//patID  string
	pubMap map[string]*rsa.PublicKey
	//id -> public key
	prMap map[string]*rsa.PrivateKey
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
	if os.Args[1] == "9913" || os.Args[1] == "8813" {
		pat.pubMap = make(map[string]*rsa.PublicKey)
		pat.prMap = make(map[string]*rsa.PrivateKey)
	}

	privateKey, publicKey := GenerateKeys()
	pat.prMap[id] = privateKey
	pat.pubMap[id] = publicKey

	fmt.Println("---- NEW PATIENT ----")
	fmt.Printf("privateKey=%v\n", privateKey)
	fmt.Printf("publicKey=%v\n", publicKey)
	fmt.Printf("SelfIdPat=%v\n", id)
	fmt.Printf("PAT LIST=%v\n", pat)
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

func GenerateKeys() (*rsa.PrivateKey, *rsa.PublicKey) {
	reader := rand.Reader
	bitSize := 2048
	privateKey, _ := rsa.GenerateKey(reader, bitSize)
	publicKey := privateKey.PublicKey
	return privateKey, &publicKey
}

func (pat *PatientList) EncryptPatInfo(patID string, patInfo string) string {
	fmt.Println("******************************************")
	//rng := rand.Reader
	message := []byte(patInfo)
	//hashed := sha256.Sum256(message)

	//signature := []byte("")

	value, isExist := pat.prMap[patID]
	if isExist == false {
		fmt.Println("DOESNT EXIST")
		NewPatientList(patID)
		value, _ = pat.prMap[patID]
	}
	//TODO: signature doesn't give anything back
	//	fmt.Println(hashed)
	//slice := hashed[:]
	//fmt.Println(slice)

	//signature, err := SignPKCS1v15(rng, rsaPrivateKey, crypto.SHA256, hashed[:])
	//signature, ok := rsa.SignPKCS1v15(rng, value, crypto.SHA256, slice)
	//func SignPSS(rand io.Reader, priv *PrivateKey, hash crypto.Hash, hashed []byte, opts *PSSOptions) ([]byte, error)
	//signature, ok := rsa.SignPSS(rng, value, slice )
	//fmt.Println(ok)
	//fmt.Println(signature)

	EncryptWithPublicKey(message, value)

	patInf := make(map[string]string)
	//patInf[patID] = string(signature[:])

	jsonData, _ := json.Marshal(patInf)
	hash := sha256.Sum256(jsonData)
	fmt.Println(patInf)
	fmt.Println(jsonData)
	//fmt.Println(hashed)
	fmt.Println("******************************************")
	return string(hash[:])

}
func EncryptWithPublicKey(message []byte, pub *rsa.PrivateKey) {
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example
	PSSmessage := message
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPSS(
		rand.Reader,
		pub,
		newhash,
		hashed,
		&opts,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("PSS Signature : %x\n", signature)
}
