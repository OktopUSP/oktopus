package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type User struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Level    int    `json:"level"`
}

func (d *Database) RegisterUser(user User) error {
	err := d.users.FindOne(d.ctx, bson.D{{"email", user.Email}}).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, err = d.users.InsertOne(d.ctx, user)
			return err
		}
		log.Println(err)
	}
	return err
}

func (d *Database) FindAllUsers() ([]User, error) {
	var result []User
	cursor, err := d.users.Find(d.ctx, bson.D{{}})
	if err != nil {
		return []User{}, err
	}
	if err = cursor.All(d.ctx, &result); err != nil {
		log.Fatal(err)
	}
	return result, err
}

func (d *Database) FindUser(email string) (User, error) {
	var result User
	err := d.users.FindOne(d.ctx, bson.D{{"email", email}}).Decode(&result)
	return result, err
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
