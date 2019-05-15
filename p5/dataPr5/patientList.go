package dataPr5

import (
	//"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"log"
	//"crypto/x509"
	//"encoding/pem"
	"fmt"
	"sync"
)

//Public Key, Private Key, PatID, lock
type PatientList struct {
	//patID  string
	PubMap map[string][]byte
	//id -> public key
	PrMap map[string][]byte
	//id -> private key
	mux  sync.Mutex
	Hops int
}

type PatientMessage struct {
	JsonPubMap []byte
	JsonPrMap  []byte
	Hops       int
}

//New list creation
func NewPatientList(id string) PatientList {
	patList := PatientList{}
	patList.Register(id)
	patList.mux = sync.Mutex{}
	patList.Hops = 3
	return patList
}

//Add new patient data to hashmaps
func (pat *PatientList) AddNewPatient(PubMap map[string][]byte, PrMap map[string][]byte) {
	for key, value := range PubMap {
		pat.PubMap[key] = value
	}
	for key, value := range PrMap {
		pat.PrMap[key] = value
	}
}

//Show the patient list
//Used to check if Doctor has an access to patients data
func (pat *PatientList) Show() string {
	pat.mux.Lock()
	defer pat.mux.Unlock()
	res := "This is a Patient List: \n"
	for key, _ := range pat.PubMap {
		res += "id = " + key + "\n"
	}
	return res
}

//New patient registration
//Public-private keys creation
func (pat *PatientList) Register(id string) {
	pat.mux.Lock()
	if len(pat.PubMap) < 1 {
		pat.PubMap = make(map[string][]byte)
		pat.PrMap = make(map[string][]byte)
	}
	privateKey, publicKey := GenerateKeys()
	pat.PrMap[id] = privateKey
	pat.PubMap[id] = publicKey

	fmt.Println("---- NEW PATIENT ----")
	fmt.Printf("privateKey=%v\n", privateKey)
	fmt.Printf("publicKey=%v\n", publicKey)
	fmt.Printf("SelfIdPat=%v\n", id)
	fmt.Println("---- NEW PATIENT ----")
	defer pat.mux.Unlock()
}

//Keys generation by RSA library
func GenerateKeys() ([]byte, []byte) {
	reader := rand.Reader
	bitSize := 2048
	privateKey, _ := rsa.GenerateKey(reader, bitSize)
	publicKey := privateKey.PublicKey
	pr := PrivateKeyToBytes(privateKey)
	pub := PublicKeyToBytes(&publicKey)
	return pr, pub
}

//Encrypting Patient Info via RSA library
func (pat *PatientList) EncryptPatInfo(patID string, patInfo string) string {
	rng := rand.Reader
	message := []byte(patInfo)
	value, _ := pat.PubMap[patID]
	key := BytesToPublicKey(value)
	signature, _ := rsa.EncryptPKCS1v15(rng, key, message[:])
	patInf := make(map[string][]byte)
	patInf[patID] = signature
	jsonData, _ := json.Marshal(patInf)
	pat.DecryptPatInfo(string(jsonData))
	return string(jsonData)
}

//Decrypting Patient Info via RSA library
func (pat *PatientList) DecryptPatInfo(hash string) map[string]string {
	var m map[string][]byte
	err := json.Unmarshal([]byte(hash), &m)
	if err != nil {
		fmt.Println(err)
	}
	rng := rand.Reader
	patInf := make(map[string]string)
	for key, val := range m {
		prKey, exist := pat.PrMap[key]
		if exist {
			private := BytesToPrivateKey(prKey)
			infoSign, er := rsa.DecryptPKCS1v15(rng, private, val)
			if err != nil {
				fmt.Println(er)
			}
			patInf[key] = string(infoSign)
		}
	}
	return patInf
}

// PrivateKeyToBytes private key to bytes via RSA library
func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

// PublicKeyToBytes public key to bytes via RSA library
func PublicKeyToBytes(pub *rsa.PublicKey) []byte {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		fmt.Println(err)
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes
}

// BytesToPrivateKey bytes to private key via RSA library
func BytesToPrivateKey(priv []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			fmt.Println(err)
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		fmt.Println(err)
	}
	return key
}

// BytesToPublicKey bytes to public key via RSA library
func BytesToPublicKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			fmt.Println(err)
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		fmt.Println(err)
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		fmt.Println("ok error")
	}
	return key
}
