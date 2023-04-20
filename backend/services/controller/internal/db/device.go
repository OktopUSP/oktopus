package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Device struct {
	SN       string
	Model    string
	Customer string
	Vendor   string
	Version  string
}

func (d *Database) CreateDevice(device Device) error {
	var result bson.M
	opts := options.FindOneAndReplace().SetUpsert(true)
	err := d.devices.FindOneAndReplace(d.ctx, bson.D{{"sn", device.SN}}, device, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("New device %s added to database", device.SN)
			return nil
		}
		log.Fatal(err)
	}
	log.Printf("Device %s already existed, and got replaced for new info", device.SN)
	return err
}

func (d *Database) RetrieveDevices() ([]Device, error) {
	var results []Device
	//TODO: filter devices by user ownership
	cursor, err := d.devices.Find(d.ctx, bson.D{}, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if err = cursor.All(d.ctx, &results); err != nil {
		log.Println(err)
		return nil, err
	}
	return results, nil
}

func (d *Database) DeleteDevice() {

}
