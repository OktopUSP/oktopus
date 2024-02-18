package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//TODO: create another package fo structs and interfaces

type Database struct {
	client  *mongo.Client
	devices *mongo.Collection
	users   *mongo.Collection
	ctx     context.Context
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

	devices := client.Database("oktopus").Collection("devices")
	createIndexes(ctx, devices)

	users := client.Database("oktopus").Collection("users")
	db.devices = devices
	db.users = users
	db.ctx = ctx
	return db
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
