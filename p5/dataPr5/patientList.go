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
	pubMap map[string][]byte
	//id -> public key
	prMap map[string][]byte
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
	//pat.patID = id
	if len(pat.pubMap) < 1 {
		pat.pubMap = make(map[string][]byte)
		pat.prMap = make(map[string][]byte)
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
	defer pat.mux.Unlock()
}

func (pat *PatientList) Show() string {
	//	pat.mux.Lock()
	//	defer pat.mux.Unlock()
	res := "This is a Patients Map: \n"
	/*	for key, value := range pat.pubMap {
		//res += "ID = " + key + ", Public Key = " + strconv.Itoa(int(value.E)) + "\n"
	}*/
	return res
}

func GenerateKeys() ([]byte, []byte) {
	reader := rand.Reader
	bitSize := 2048
	privateKey, _ := rsa.GenerateKey(reader, bitSize)
	publicKey := privateKey.PublicKey
	pr := PrivateKeyToBytes(privateKey)
	pub := PublicKeyToBytes(&publicKey)
	return pr, pub
}

func (pat *PatientList) EncryptPatInfo(patID string, patInfo string) string {
	rng := rand.Reader
	message := []byte(patInfo)
	value, _ := pat.pubMap[patID]
	key := BytesToPublicKey(value)
	signature, _ := rsa.EncryptPKCS1v15(rng, key, message[:])
	patInf := make(map[string][]byte)
	patInf[patID] = signature
	jsonData, _ := json.Marshal(patInf)
	pat.DecryptPatInfo(string(jsonData))
	return string(jsonData)
}

func (pat *PatientList) DecryptPatInfo(hash string) map[string]string {
	fmt.Println("DECRYPTED")
	//hash := []byte(h)
	fmt.Println(hash)
	var m map[string][]byte
	err := json.Unmarshal([]byte(hash), &m)
	if err != nil {
		fmt.Println(err)
	}
	rng := rand.Reader
	fmt.Println(m)
	patInf := make(map[string]string)
	for key, val := range m {
		prKey, _ := pat.prMap[key]

		private := BytesToPrivateKey(prKey)
		infoSign, er := rsa.DecryptPKCS1v15(rng, private, val)
		if err != nil {
			fmt.Println(er)
		}
		fmt.Println(string(infoSign))
		fmt.Println("DECRYPTED")
		patInf[key] = string(infoSign)
	}
	fmt.Println(patInf)
	return patInf
}

// PrivateKeyToBytes private key to bytes
func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

// PublicKeyToBytes public key to bytes
func PublicKeyToBytes(pub *rsa.PublicKey) []byte {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		//log.Error(err)
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes
}

// BytesToPrivateKey bytes to private key
func BytesToPrivateKey(priv []byte) *rsa.PrivateKey {
	fmt.Println("BytesToPrivateKey")
	fmt.Println(priv)
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	fmt.Println("BytesToPrivateKey")

	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			//log.Error(err)
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		//log.Error(err)
	}
	return key
}

// BytesToPublicKey bytes to public key
func BytesToPublicKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			//log.Error(err)
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		//log.Error(err)
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		//log.Error("not ok")
	}
	return key
}

func (pat *PatientList) GetPublicKeys() map[string][]byte {
	return pat.pubMap
}
