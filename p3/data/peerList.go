package data

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	//	"strings"
	"strconv"
	"sync"
)

type PeerList struct {
	selfId    int32
	peerMap   map[string]int32
	maxLength int32
	mux       sync.Mutex
}

func NewPeerList(id int32, maxLength int32) PeerList {
	peerList := PeerList{}
	peerList.Register(id)
	peerList.maxLength = maxLength
	peerList.peerMap = make(map[string]int32)
	peerList.mux = sync.Mutex{}
	return peerList
}

func (peers *PeerList) Add(addr string, id int32) {
	peers.mux.Lock()
	//If exists -> Overwrite the latest ID.
	if id != peers.selfId {
		peers.peerMap[addr] = id
	}
	peers.mux.Unlock()
}

func (peers *PeerList) Delete(addr string) {
	peers.mux.Lock()
	_, ok := peers.peerMap[addr]
	if ok {
		delete(peers.peerMap, addr)
	}
	peers.mux.Unlock()
}

func (peers *PeerList) Rebalance() {
	//MY LOGIC
	//1. Get all ids from map to array
	//2. Insert SelfId into array
	//3. Sort array
	//4. Find the index of selfId
	//5. Create new array
	//5.1. If len of sorted > index + max/2
	//add to array all elements in the range [index-max/2; index+max/2]
	//5.2 if index < max/2 (don't have max/2 in the beginning)
	//add all elements after [index; index+max/2]
	//add all elements before index (if exist), get num = index - length
	//add num elements from the the end
	//5.3. if len < index+max/2+1 (don't have in the end)
	//add all elements before [index-max/2; index]
	//add all elements after index (if exist), get num = index - length
	//add num elements from the beginning
	//6. Fill a new hashmap with ids from a new array.
	peers.mux.Lock()

	if len(peers.peerMap) > (int(peers.maxLength)) {
		sortedArr := []int32{}
		peerRebalance := map[string]int32{}
		for _, value := range peers.peerMap {
			sortedArr = append(sortedArr, value)
		}
		sortedArr = append(sortedArr, peers.selfId)
		sort.Slice(sortedArr, func(i, j int) bool { return sortedArr[i] < sortedArr[j] })
		index := 0
		for k, v := range sortedArr {
			if v == peers.selfId {
				index = k
			}
		}
		newArr := []int32{}
		if len(sortedArr)-1 < index+(int(peers.maxLength/2)) {
			num := index + (int(peers.maxLength) / 2) + 1 - len(sortedArr)
			for i := index - (int(peers.maxLength) / 2); i < index; i++ {
				newArr = append(newArr, sortedArr[i])
				peerRebalance[peers.KeyByValue(sortedArr[i])] = sortedArr[i]
			}
			for j := index + 1; j <= index+num; j++ {
				newArr = append(newArr, sortedArr[j])
				peerRebalance[peers.KeyByValue(sortedArr[j])] = sortedArr[j]
			}
			for k := 0; k < (int(peers.maxLength)/2)-num; k++ {
				newArr = append(newArr, sortedArr[k])
				peerRebalance[peers.KeyByValue(sortedArr[k])] = sortedArr[k]
			}

		} else if index < (int(peers.maxLength / 2)) {
			num := (int(peers.maxLength / 2)) - index
			for i := index - num; i < index; i++ {
				newArr = append(newArr, sortedArr[i])
				peerRebalance[peers.KeyByValue(sortedArr[i])] = sortedArr[i]
			}
			for j := index + 1; j <= index+(int(peers.maxLength/2)); j++ {
				newArr = append(newArr, sortedArr[j])
				peerRebalance[peers.KeyByValue(sortedArr[j])] = sortedArr[j]
			}
			for k := len(sortedArr) - 1; k >= (len(sortedArr) - num); k-- {
				newArr = append(newArr, sortedArr[k])
				peerRebalance[peers.KeyByValue(sortedArr[k])] = sortedArr[k]
			}
		} else {
			for i := index - 1; i >= index-(int(peers.maxLength/2)); i-- {
				newArr = append(newArr, sortedArr[i])
				peerRebalance[peers.KeyByValue(sortedArr[i])] = sortedArr[i]
			}
			for j := index + 1; j <= index+(int(peers.maxLength/2)); j++ {
				newArr = append(newArr, sortedArr[j])
				peerRebalance[peers.KeyByValue(sortedArr[j])] = sortedArr[j]
			}
		}
		peers.peerMap = peerRebalance
	}
	peers.mux.Unlock()
}

//Show() shows all addresses and their corresponding IDs.
//For example, it returns "This is PeerMap: \n addr=127.0.0.1, id=1".
func (peers *PeerList) Show() string {
	peers.mux.Lock()
	defer peers.mux.Unlock()
	res := "This is a PeerMap: \n"
	for key, value := range peers.peerMap {
		res += "addr = " + key + ", id = " + strconv.Itoa(int(value)) + "\n"
	}
	return res
}

func (peers *PeerList) Register(id int32) {
	peers.mux.Lock()
	defer peers.mux.Unlock()
	peers.selfId = id
	fmt.Printf("SelfId=%v\n", id)
}

func (peers *PeerList) Copy() map[string]int32 {
	peers.mux.Lock()
	defer peers.mux.Unlock()
	newMap := make(map[string]int32)
	newMap = peers.peerMap
	return newMap
}

func (peers *PeerList) GetSelfId() int32 {
	peers.mux.Lock()
	defer peers.mux.Unlock()
	return peers.selfId
}

func (peers *PeerList) PeerMapToJson() (string, error) {
	peers.mux.Lock()
	defer peers.mux.Unlock()
	peerMapJson, err := json.Marshal(peers.peerMap)
	return string(peerMapJson), err
}

//InjectPeerMapJson() inserts every entries(every <addr, id> pair) of the parameter
//"peerMapJsonStr" into your own PeerMap, except the entry whose addres is your own local address.
func (peers *PeerList) InjectPeerMapJson(peerMapJsonStr string, selfAddr string) {
	peers.mux.Lock()
	peer := map[string]int32{}
	err := json.Unmarshal([]byte(peerMapJsonStr), &peer)
	if err != nil {
		println(err)
	}
	for key, value := range peer {
		peers.peerMap[key] = value
	}
	_, ok := peers.peerMap[selfAddr]
	if ok {
		delete(peers.peerMap, selfAddr)
	}
	peers.mux.Unlock()
}

func TestPeerListRebalance() {
	peers := NewPeerList(5, 4)
	peers.Add("1111", 1)
	peers.Add("4444", 4)
	peers.Add("-1-1", -1)
	peers.Add("0000", 0)
	peers.Add("2121", 21)
	peers.Rebalance()
	expected := NewPeerList(5, 4)
	expected.Add("1111", 1)
	expected.Add("4444", 4)
	expected.Add("2121", 21)
	expected.Add("-1-1", -1)
	fmt.Println(reflect.DeepEqual(peers, expected))

	peers = NewPeerList(5, 2)
	peers.Add("1111", 1)
	peers.Add("4444", 4)
	peers.Add("-1-1", -1)
	peers.Add("0000", 0)
	peers.Add("2121", 21)
	peers.Rebalance()
	expected = NewPeerList(5, 2)
	expected.Add("4444", 4)
	expected.Add("2121", 21)
	fmt.Println(reflect.DeepEqual(peers, expected))

	peers = NewPeerList(5, 4)
	peers.Add("1111", 1)
	peers.Add("7777", 7)
	peers.Add("9999", 9)
	peers.Add("11111111", 11)
	peers.Add("2020", 20)
	peers.Rebalance()
	expected = NewPeerList(5, 4)
	expected.Add("1111", 1)
	expected.Add("7777", 7)
	expected.Add("9999", 9)
	expected.Add("2020", 20)
	fmt.Println(reflect.DeepEqual(peers, expected))
}

//returns key by giving value from hashmap peerMap
func (peers *PeerList) KeyByValue(id int32) string {

	for key, value := range peers.peerMap {
		if reflect.DeepEqual(id, value) {
			return key
		}
	}
	return ""
}
