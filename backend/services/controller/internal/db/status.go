package db

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d *Database) UpdateStatus(sn string, status Status, mtp MTP) error {
	var result Device

	err := d.devices.FindOne(d.ctx, bson.D{{"sn", sn}}, nil).Decode(&result)
	if err != nil {
		log.Println(err)
	}

	//TODO: abolish this logic, find another approach, microservices design maybe?
	/*
		In case the device status is online, we must check if the mtp
		changing is going to affect the global status. In case it does,
		we must update the global status accordingly.
	*/

	/*
		mix the existent device status to the updated one
	*/
	switch mtp {
	case MQTT:
		result.Mqtt = status
	case STOMP:
		result.Stomp = status
	case WEBSOCKETS:
		result.Websockets = status
	}

	/*
		check if the global status needs update
	*/
	var globalStatus primitive.E
	if result.Mqtt == Offline && result.Stomp == Offline && result.Websockets == Offline {
		globalStatus = primitive.E{"status", Offline}
	}
	if result.Mqtt == Online || result.Stomp == Online || result.Websockets == Online {
		globalStatus = primitive.E{"status", Online}
	}

	_, err = d.devices.UpdateOne(d.ctx, bson.D{{"sn", sn}}, bson.D{
		{
			"$set", bson.D{
				{mtp.String(), status},
				globalStatus,
			},
		},
	})

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
