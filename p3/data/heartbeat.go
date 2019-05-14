package data

type HeartBeatData struct {
	IfNewBlock  bool   `json:"ifNewBlock"`
	Id          int32  `json:"id"`
	BlockJson   string `json:"blockJson"`
	PeerMapJson string `json:"peerMapJson"`
	Addr        string `json:"addr"`
	Hops        int32  `json:"hops"`
}

//var IQ p5.ItemQueue

func NewHeartBeatData(ifNewBlock bool, id int32, blockJson string, peerMapJson string, addr string) HeartBeatData {
	return HeartBeatData{IfNewBlock: ifNewBlock, Id: id, BlockJson: blockJson, PeerMapJson: peerMapJson, Addr: addr, Hops: 3}
}

func PrepareHeartBeatData(selfId int32, peerMapJson string, addr string) HeartBeatData {
	return NewHeartBeatData(false, selfId, "", peerMapJson, addr)
}

/*func GenerateRandomMPT() p1.MerklePatriciaTrie {
	//random words generation library
	babbler := babble.NewBabbler()
	babbler.Count = 1
	mpt := p1.MerklePatriciaTrie{}
	mpt.Insert(babbler.Babble(), babbler.Babble())
	num := rand.Intn(4)
	for i := 0; i < num; i++ {
		mpt.Insert(babbler.Babble(), babbler.Babble())
	}
	return mpt
}*/
