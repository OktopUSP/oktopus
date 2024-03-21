package db

import (
	"context"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	client  *mongo.Client
	devices *mongo.Collection
	ctx     context.Context
	m       *sync.Mutex
}

func NewDatabase(ctx context.Context, mongoUri string) Database {
	var db Database

	clientOptions := options.Client().ApplyURI(mongoUri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	db.client = client

	log.Println("Trying to ping Mongo database...")
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Couldn't connect to MongoDB --> ", err)
	}

	log.Println("Connected to MongoDB-->", mongoUri)

	devices := client.Database("adapter").Collection("devices")
	createIndexes(ctx, devices)
	resetDeviceStatus(ctx, devices)

	db.devices = devices
	db.ctx = ctx
	db.m = &sync.Mutex{}

	return db
}

func resetDeviceStatus(ctx context.Context, devices *mongo.Collection) {
	_, err := devices.UpdateMany(ctx, bson.D{{}}, bson.D{
		{
			"$set", bson.D{
				{"mqtt", 0},
				{"stomp", 0},
				{"websockets", 0},
				{"status", 0},
			},
		},
	})
	if err != nil {
		log.Fatalln("ERROR to reset device status in database:", err)
	}
}

func createIndexes(ctx context.Context, devices *mongo.Collection) {
	indexField := bson.M{"sn": 1}
	_, err := devices.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    indexField,
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		log.Println("ERROR to create index in database:", err)
	}
}
