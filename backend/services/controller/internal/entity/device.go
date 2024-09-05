package entity

type Device struct {
	SN           string
	Model        string
	Customer     string
	Vendor       string
	Version      string
	ProductClass string
	Alias        string
	Status       Status
	Mqtt         Status
	Stomp        Status
	Websockets   Status
	Cwmp         Status
}

type VendorsCount struct {
	Vendor string `bson:"_id" json:"vendor"`
	Count  int    `bson:"count" json:"count"`
}

type ProductClassCount struct {
	ProductClass string `bson:"_id" json:"productClass"`
	Count        int    `bson:"count" json:"count"`
}

type StatusCount struct {
	Status int `bson:"_id" json:"status"`
	Count  int `bson:"count" json:"count"`
}

type DevicesList struct {
	Devices []Device `json:"devices" bson:"documents"`
	Total   int64    `json:"total"`
}

type FilterOptions struct {
	Models         []string `json:"models"`
	ProductClasses []string `json:"productClasses"`
	Vendors        []string `json:"vendors"`
	Versions       []string `json:"versions"`
}
