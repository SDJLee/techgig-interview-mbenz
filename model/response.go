package model

import "gopkg.in/guregu/null.v3"

//{ "transactionId": "043020211 //A unique numerical value", "vin": "W1K2062161F0014 //vehicle identification number", "source": "source name", "destination": "destination name", "distance": "100 //distance between the source and destination in miles", "currentChargeLevel": "1 //current charge level in percentage , 0<=charge<=100", "isChargingRequired": "true/false //whether the vehicle has to stop for charging?.If true populate charging stations", "chargingStations": [ "s1", "s2" ], "errors": [ { "Id": 8888, "description": "Unable to reach the destination with the current charge level" }, { "id": 9999, "description": "Technical Exception" } ] }

type ResError struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type Station struct {
	Name     string `json:"name"`
	Limit    int64  `json:"limit"`
	Distance int64  `json:"distance"`
}

type Response struct {
	TransactionID      int64       `json:"transactionId"`
	Vin                null.String `json:"vin"`
	Source             null.String `json:"source"`
	Destination        null.String `json:"destination"`
	Distance           null.Int    `json:"distance,omitempty"`
	CurrentChargeLevel null.Int    `json:"currentChargeLevel,omitempty"`
	IsChargingRequired null.Bool   `json:"isChargingRequired,omitempty"`
	ChargingStations   []string    `json:"chargingStations,omitempty"`
	Errors             []*ResError `json:"errors,omitempty"`
}

// { "source": "source name", "destination": "destination name""distance": "100 //distance between the source and destination in miles", "error": "It will be null if No Error" }
type ResTravelDistance struct {
	Source      string      `json:"source"`
	Destination string      `json:"destination"`
	Distance    int64       `json:"distance"`
	Error       null.String `json:"error"`
}

//{ "vin": "vehicle identification number", "currentChargeLevel": "current battery charge level in percentage, 0<=charge<=100 eg: 1", "error": "It will be null if No Error" }
type ResChargeLevel struct {
	Vin                string      `json:"vin"`
	CurrentChargeLevel int64       `json:"currentChargeLevel"`
	Error              null.String `json:"error"`
}

// { "source": "source name", "destination": "destination name", "chargingStations": [ { "name": "s1 //Name of the charging station", "distance": "10 //Distance from the source in miles", "limit": "60 //Available charge level. If the current battery level is 1%, vehicle can charge upto 61% from this station " }, { "name": "s2", "distance": "20", "limit": "30" }, { "name": "s3", "distance": "30", "limit": "30" }, { "name": "s4", "distance": "60", "limit": "40" } ], "error": "It will be null if No Error" }
type ResChargeStations struct {
	Source           string      `json:"source"`
	Destination      string      `json:"destination"`
	Error            null.String `json:"error"`
	ChargingStations []*Station  `json:"chargingStations"`
}
