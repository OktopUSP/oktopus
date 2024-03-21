package entity

type Device struct {
	SN           string
	Model        string
	Customer     string
	Vendor       string
	Version      string
	ProductClass string
	Status       Status
	Mqtt         Status
	Stomp        Status
	Websockets   Status
}
