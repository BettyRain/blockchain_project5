package dataPr5

import (
	"reflect"
	"sync"
)

type DataPool struct {
	DB   map[string]string
	Sign []byte
	Hops int
}

type ItemQueue struct {
	Items []DataPool
	Lock  sync.RWMutex
}

func (iq *ItemQueue) Initialize() {
}

func NewItemQueue() ItemQueue {
	return ItemQueue{}
}

func (iq *ItemQueue) RemoveItem(index int) ItemQueue {
	sliceA := iq.Items
	sliceA = append(sliceA[:index], sliceA[index+1:]...)
	iq.Items = sliceA
	return *iq
}

func (iq *ItemQueue) RemoveAddedData(trie map[string]string) {
	for key, value := range trie {
		for i := 0; i < len(iq.Items); i++ {
			for k, v := range iq.Items[i].DB {
				if (k == key) && (v == value) {
					iq.RemoveItem(i)
				}
			}
		}
	}
}

//add dataPr5 to pool
func AddToPool(docID string, patID string, patInfo string, pat PatientList, doc DoctorList) DataPool {
	//Encrypt PatientInfo with patient's Private Key
	hash := pat.EncryptPatInfo(patID, patInfo)
	docHash := doc.SignByDoc(hash, docID)
	data := make(map[string]string)
	data[docID] = hash
	dataPool := DataPool{data, docHash, 3}
	return dataPool
}

//add dataPr5 to queue
func (iq *ItemQueue) AddToQueue(dp DataPool) {
	//if data in, don't add
	iq.Lock.Lock()

	eq := false
	for i := 0; i < len(iq.Items); i++ {
		eq = reflect.DeepEqual(iq.Items[i].DB, dp.DB)
	}
	if eq == false {
		iq.Items = append(iq.Items, dp)
	}
	iq.Lock.Unlock()
}

//return dataPr5 from pool
func (iq *ItemQueue) GetFromPool() []DataPool {
	iq.Lock.Lock()
	defer iq.Lock.Unlock()
	return iq.Items
}

/*func (iq *ItemQueue) GenerateMPT() p1.MerklePatriciaTrie {
	//random number how many lines to insert in a block (assumption: <=4)
	num := rand.Intn(4)
	count := 0
	iq.GetFromPool()
	mpt := p1.MerklePatriciaTrie{}
	for _, value := range iq.Items {
		if count < num {
			for k, v := range value.DB {
				mpt.Insert(k, v)
				count++
			}
		} else {
			break
		}
	}
	return mpt
}*/
