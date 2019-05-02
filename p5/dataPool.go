package p5

import "sync"

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

//add data to pool
func (dp *DataPool) AddToPool(docID string, patID string, patInfo string) {
	//Encrypt PatientInfo with patient's Private Key
	//Hash <PatID, [PatInfo]PK>
	//Sign H<PatID, [PatInfo]PKpat> with doc's Private key
	//Send to miner <DocID, [H<PatID, [PatInfo]PKpat>]PKdoc>
	hash := patID + patInfo
	dp.DB = make(map[string]string)
	dp.DB[docID] = hash
}

//add data to queue
func (iq *ItemQueue) AddToQueue(dp DataPool) {
	iq.Lock.Lock()
	iq.Items = append(iq.Items, dp)
	iq.Lock.Unlock()
}

//return data from pool
func (iq *ItemQueue) GetFromPool() []DataPool {
	iq.Lock.Lock()
	defer iq.Lock.Unlock()
	return iq.Items
}
