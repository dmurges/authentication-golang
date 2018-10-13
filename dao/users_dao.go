package dao

import (
	. "authentication-service/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type UserDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "user"
)

func (m *UserDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

func (m *UserDAO) Insert(user User) error {
	err := db.C(COLLECTION).Insert(&user)
	return err
}

func (m *UserDAO) GetUsers() ([]User, error) {
	var users []User
	err := db.C(COLLECTION).Find(bson.M{}).All(&users)
	return users, err
}

func (m *UserDAO) GetUser(id string) (User, error) {
	var user User
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&user)
	return user, err
}

func (m *UserDAO) GetUserByEmail(email string) (User, error) {
	var user User
	err := db.C(COLLECTION).Find(bson.M{"email": email}).One(&user)
	return user, err
}