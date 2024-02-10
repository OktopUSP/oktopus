package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//TODO: create another package fo structs and interfaces

type Database struct {
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
	log.Println("Trying to ping Mongo database...")
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Couldn't connect to MongoDB --> ", err)
	}

	log.Println("Connected to MongoDB-->", mongoUri)
	devices := client.Database("oktopus").Collection("devices")
	users := client.Database("oktopus").Collection("users")
	db.devices = devices
	db.users = users
	db.ctx = ctx
	return db
}
