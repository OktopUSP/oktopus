package db

import (
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

type UserLevels int32

const (
	NormalUser UserLevels = iota
	AdminUser
)

type User struct {
	Email    string     `json:"email"`
	Name     string     `json:"name"`
	Password string     `json:"password,omitempty"`
	Level    UserLevels `json:"level"`
	Phone    string     `json:"phone"`
}

var ErrorUserExists = errors.New("User already exists")

func (d *Database) RegisterUser(user User) error {
	err := d.users.FindOne(d.ctx, bson.D{{"email", user.Email}}).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, err = d.users.InsertOne(d.ctx, user)
			return err
		}
		log.Println(err)
		return err
	} else {
		return ErrorUserExists
	}
}

func (d *Database) UpdatePassword(user User) error {
	if !validEmail(user.Email) {
		return errors.New("invalid email format")
	}
	_, err := d.users.UpdateOne(d.ctx, bson.D{{"email", user.Email}}, bson.D{{"$set", bson.D{{"password", user.Password}}}})
	return err
}

func (d *Database) FindAllUsers() ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	cursor, err := d.users.Find(d.ctx, bson.D{{}})
	if err != nil {
		return []map[string]interface{}{}, err
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

func (d *Database) DeleteUser(email string) error {
	_, err := d.users.DeleteOne(d.ctx, bson.D{{"email", email}})
	return err
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

func validEmail(email string) bool {
	// Simple regex for email validation
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
