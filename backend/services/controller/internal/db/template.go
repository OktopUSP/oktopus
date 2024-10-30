package db

import (
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Template struct {
	Name  string `json:"name" bson:"name"`
	Type  string `json:"type" bson:"type"`
	Value string `json:"value" bson:"value"`
}

var ErrorTemplateExists = errors.New("message already exists")
var ErrorTemplateNotExists = errors.New("message don't exist")

func (d *Database) FindTemplate(filter interface{}) (Template, error) {
	var result Template
	err := d.template.FindOne(d.ctx, filter).Decode(&result)
	return result, err
}

func (d *Database) AllTemplates(filter interface{}) ([]Template, error) {
	var results []Template

	cursor, err := d.template.Find(d.ctx, filter)
	if err != nil {
		return results, err
	}
	if err = cursor.All(d.ctx, &results); err != nil {
		log.Println(err)
	}
	return results, err
}

func (d *Database) AddTemplate(name, tr string, t string) error {
	opts := options.FindOneAndReplace().SetUpsert(true)
	err := d.template.FindOneAndReplace(d.ctx, bson.D{{"name", name}}, Template{Name: name, Type: tr, Value: t}, opts).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("New message %s added to database", name)
			return nil
		}
		return err
	}
	log.Printf("Message %s already existed, and got replaced for new payload", name)
	return err
}

func (d *Database) UpdateTemplate(name, t string) error {
	result, err := d.template.UpdateOne(d.ctx, bson.D{{"name", name}}, bson.D{{"$set", bson.D{{"value", t}}}})
	if err == nil {
		if result.MatchedCount == 0 {
			return ErrorTemplateNotExists
		}
	}
	return err
}

func (d *Database) DeleteTemplate(name string) error {
	result, err := d.template.DeleteOne(d.ctx, bson.D{{"name", name}})
	if err == nil {
		if result.DeletedCount == 0 {
			return ErrorTemplateNotExists
		}
	}
	return err
}
