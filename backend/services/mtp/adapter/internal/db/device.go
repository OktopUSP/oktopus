package db

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MTP int32

const (
	UNDEFINED MTP = iota
	MQTT
	STOMP
	WEBSOCKETS
	CWMP
)

type Status uint8

const (
	Offline Status = iota
	Associating
	Online
)

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

func (d *Database) CreateDevice(device Device) error {
	var result bson.M
	var deviceExistent Device

	d.m.Lock()
	defer d.m.Unlock()

	/* ------------------ Do not overwrite status of other mtp ------------------ */
	err := d.devices.FindOne(d.ctx, bson.D{{"sn", device.SN}}, nil).Decode(&deviceExistent)
	if err == nil {
		if deviceExistent.Mqtt == Online {
			device.Mqtt = Online
		}
		if deviceExistent.Stomp == Online {
			device.Stomp = Online
		}
		if deviceExistent.Websockets == Online {
			device.Websockets = Online
		}
		if deviceExistent.Cwmp == Online {
			device.Cwmp = Online
		}
	} else {
		if err != mongo.ErrNoDocuments {
			log.Println(err)
			return err
		}
	}
	/* -------------------------------------------------------------------------- */

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Important: You must pass sessCtx as the Context parameter to the operations for them to be executed in the
		// transaction.
		opts := options.FindOneAndReplace().SetUpsert(true)

		err := d.devices.FindOneAndReplace(d.ctx, bson.D{{"sn", device.SN}}, device, opts).Decode(&result)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				log.Printf("New device %s added to database", device.SN)
				return nil, nil
			}
			return nil, err
		}
		log.Printf("Device %s already existed, and got replaced for new info", device.SN)
		return nil, nil
	}

	session, err := d.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(d.ctx)

	_, err = session.WithTransaction(d.ctx, callback)
	if err != nil {
		return err
	}
	return err
}
func (d *Database) RetrieveDevices(filter bson.A) ([]Device, error) {

	var results []Device

	cursor, err := d.devices.Aggregate(d.ctx, filter)
	if err != nil {
		return results, err
	}

	if cursor.Err() != nil {
		return results, cursor.Err()
	}

	for cursor.Next(d.ctx) {
		var device Device

		err := cursor.Decode(&device)
		if err != nil {
			log.Println("Error to decode device info fields")
			continue
		}

		results = append(results, device)
	}

	return results, err
}

func (d *Database) RetrieveDevice(sn string) (Device, error) {
	var result Device
	err := d.devices.FindOne(d.ctx, bson.D{{"sn", sn}}, nil).Decode(&result)
	if err != nil {
		log.Println(err)
	}
	return result, err
}

func (d *Database) RetrieveDevicesCount(filter bson.M) (int64, error) {
	count, err := d.devices.CountDocuments(d.ctx, filter)
	return count, err
}

func (d *Database) DeleteDevice() {

}

func (d *Database) SetDeviceAlias(sn string, newAlias string) error {
	err := d.devices.FindOneAndUpdate(d.ctx, bson.D{{"sn", sn}}, bson.D{{"$set", bson.D{{"alias", newAlias}}}}).Err()
	return err
}

func (d *Database) DeviceExists(sn string) (bool, error) {
	_, err := d.RetrieveDevice(sn)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (m MTP) String() string {
	switch m {
	case UNDEFINED:
		return "unknown"
	case MQTT:
		return "mqtt"
	case STOMP:
		return "stomp"
	case WEBSOCKETS:
		return "websockets"
	case CWMP:
		return "cwmp"
	}
	return "unknown"
}
