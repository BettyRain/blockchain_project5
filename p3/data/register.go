package data

import "encoding/json"

type RegisterData struct {
	AssignedId  int32  `json:"assignedId"`
	PeerMapJson string `json:"peerMapJson"`
}

func NewRegisterData(id int32, peerMapJson string) RegisterData {
	return RegisterData{AssignedId: id, PeerMapJson: peerMapJson}
}

func (data *RegisterData) EncodeToJson() (string, error) {
	encodeJson, err := json.Marshal(data)
	return string(encodeJson), err
}
