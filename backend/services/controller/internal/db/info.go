package db

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type VendorsCount struct {
	Vendor string `bson:"_id" json:"vendor"`
	Count  int    `bson:"count" json:"count"`
}

type ProductClassCount struct {
	ProductClass string `bson:"_id" json:"productClass"`
	Count        int    `bson:"count" json:"count"`
}

type StatusCount struct {
	Status int `bson:"_id" json:"status"`
	Count  int `bson:"count" json:"count"`
}

func (d *Database) RetrieveVendorsInfo() ([]VendorsCount, error) {
	var results []VendorsCount
	cursor, err := d.devices.Aggregate(d.ctx, []bson.M{
		{
			"$group": bson.M{
				"_id":   "$vendor",
				"count": bson.M{"$sum": 1},
			},
		},
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cursor.Close(d.ctx)
	if err := cursor.All(d.ctx, &results); err != nil {
		log.Println(err)
		return nil, err
	}
	for _, result := range results {
		log.Println(result)
	}
	return results, nil
}

func (d *Database) RetrieveStatusInfo() ([]StatusCount, error) {
	var results []StatusCount
	cursor, err := d.devices.Aggregate(d.ctx, []bson.M{
		{
			"$group": bson.M{
				"_id":   "$status",
				"count": bson.M{"$sum": 1},
			},
		},
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cursor.Close(d.ctx)
	if err := cursor.All(d.ctx, &results); err != nil {
		log.Println(err)
		return nil, err
	}
	for _, result := range results {
		log.Println(result)
	}
	return results, nil
}

func (d *Database) RetrieveProductsClassInfo() ([]ProductClassCount, error) {
	var results []ProductClassCount
	cursor, err := d.devices.Aggregate(d.ctx, []bson.M{
		{
			"$group": bson.M{
				"_id":   "$productclass",
				"count": bson.M{"$sum": 1},
			},
		},
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cursor.Close(d.ctx)
	if err := cursor.All(d.ctx, &results); err != nil {
		log.Println(err)
		return nil, err
	}
	for _, result := range results {
		log.Println(result)
	}
	return results, nil
}
