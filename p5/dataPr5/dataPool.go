package dataPr5

import (
	"../../p3"
	"fmt"
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
}

type ItemQueue struct {
	Items []DataPool
	Lock  sync.RWMutex
}

func (iq *ItemQueue) Initialize() {

}

//add dataPr5 to pool
func AddToPool(docID string, patID string, patInfo string) DataPool {
	//Encrypt PatientInfo with patient's Private Key
	//Hash <PatID, [PatInfo]PK>
	//Sign H<PatID, [PatInfo]PKpat> with doc's Private key
	//Send to miner <DocID, [H<PatID, [PatInfo]PKpat>]PKdoc>
	//TODO: create [H<PatID, [PatInfo]PKpat>]PKdoc
	hash := patID + patInfo
	data := make(map[string]string)
	data[docID] = hash
	dataPool := DataPool{data}
	fmt.Println("===============================")
	fmt.Println(dataPool)
	//TODO: send data to first node
	//TODO: hops?
	//	p3.ForwardNewData(dataPool)
	fmt.Println("===============================")

	return dataPool
}

//add dataPr5 to queue
func (iq *ItemQueue) AddToQueue(dp DataPool) {
	//and send to miners??
	iq.Lock.Lock()
	iq.Items = append(iq.Items, dp)
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
	//TODO: delete info which is in canonical chain
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
