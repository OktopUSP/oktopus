package db

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// TODO: fix this function to also change device status at different mtp
func (d *Database) UpdateStatus(sn string, status uint8) error {
	var result bson.M
	err := d.devices.FindOneAndUpdate(d.ctx, bson.D{{"sn", sn}}, bson.D{{"$set", bson.D{{"status", status}}}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Device %s is not mapped into database", sn)
			return nil
		}
		log.Println(err)
	}
	log.Printf("%s is now offline.", sn)
	return err
}
