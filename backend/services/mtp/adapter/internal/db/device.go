package db

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type DevicesList struct {
	Devices []Device `json:"devices" bson:"documents"`
	Total   int64    `json:"total" bson:"totalCount"`
}

type FilterOptions struct {
	Models         []string `json:"models"`
	ProductClasses []string `json:"productClasses"`
	Vendors        []string `json:"vendors"`
	Versions       []string `json:"versions"`
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

	/* ------------------------- Do not overwrite alias ------------------------- */
	if deviceExistent.Alias != "" {
		device.Alias = deviceExistent.Alias
	}
	/* -------------------------------------------------------------------------- */

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
func (d *Database) RetrieveDevices(filter bson.A) (*DevicesList, error) {

	var results []DevicesList

	cursor, err := d.devices.Aggregate(d.ctx, filter)
	if err != nil {
		return nil, err
	}
	if cursor.Err() != nil {
		return nil, cursor.Err()
	}
	defer cursor.Close(d.ctx)
	if err := cursor.All(d.ctx, &results); err != nil {
		log.Println(err)
		return nil, err
	}

	//log.Printf("results: %++v", results)

	return &results[0], err
}

func (d *Database) RetrieveDeviceFilterOptions() (FilterOptions, error) {
	filter := bson.A{
		bson.D{
			{"$group",
				bson.D{
					{"_id", primitive.Null{}},
					{"vendors", bson.D{{"$addToSet", "$vendor"}}},
					{"versions", bson.D{{"$addToSet", "$version"}}},
					{"productClasses", bson.D{{"$addToSet", "$productclass"}}},
					{"models", bson.D{{"$addToSet", "$model"}}},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"_id", 0},
					{"vendors", 1},
					{"versions", 1},
					{"productClasses", 1},
					{"models", 1},
				},
			},
		},
	}

	var results []FilterOptions
	cursor, err := d.devices.Aggregate(d.ctx, filter)
	if err != nil {
		log.Println(err)
		return FilterOptions{}, err
	}
	defer cursor.Close(d.ctx)

	if err := cursor.All(d.ctx, &results); err != nil {
		log.Println(err)
		return FilterOptions{}, err
	}

	if len(results) > 0 {
		return results[0], nil
	} else {
		return FilterOptions{
			Models:         []string{},
			ProductClasses: []string{},
			Vendors:        []string{},
			Versions:       []string{},
		}, nil
	}
}

func (d *Database) DeleteDevices(filter bson.D) (int64, error) {

	result, err := d.devices.DeleteMany(d.ctx, filter)
	if err != nil {
		log.Println(err)
	}
	return result.DeletedCount, err
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
