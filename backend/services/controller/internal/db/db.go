package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Database struct {
	devices *mongo.Collection
	ctx     context.Context
}

func NewDatabase(ctx context.Context, mongoUri string) Database {
	var db Database
	clientOptions := options.Client().ApplyURI(mongoUri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB-->", mongoUri)
	devices := client.Database("oktopus").Collection("devices")
	db.devices = devices
	db.ctx = ctx
	return db
}
