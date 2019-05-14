package dataPr5

import (
	"fmt"
	"reflect"
	"sync"
)

type PatientData struct {
	PatInfo string
	PatId   string
	DocId   string
}

type DataPool struct {
	DB map[string]string
	//DocId H< PatID - Info >
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

//add dataPr5 to pool
func AddToPool(docID string, patID string, patInfo string, pat PatientList) DataPool {
	//Encrypt PatientInfo with patient's Private Key
	//Hash <PatID, [PatInfo]PK>
	//Sign H<PatID, [PatInfo]PKpat> with doc's Private key
	//Send to miner <DocID, [H<PatID, [PatInfo]PKpat>]PKdoc>
	//TODO: create [H<PatID, [PatInfo]PKpat>]PKdoc
	//hash := patID + patInfo
	hash := pat.EncryptPatInfo(patID, patInfo)
	fmt.Println("===============================")
	fmt.Println(hash)
	data := make(map[string]string)
	data[docID] = hash
	dataPool := DataPool{data, 3}
	fmt.Println("===============================")
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
