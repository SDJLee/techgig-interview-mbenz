package model

type Request struct {
	Vin         string `json:"vin"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

type ReqTravelDistance struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

type ReqChargeLevel struct {
	Vin string `json:"vin"`
}

type ReqChargeStations struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}
