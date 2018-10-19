package main

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Student struct {
	Name      string `json:"name"`
	Age       int    `json:"age"`
	StudentID string `json:"studentID"`
}

type StudentMongoDB struct {
	Host                  string
	DatabaseName          string
	StudentCollectionName string
}

func (db *StudentMongoDB) Init() {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	index := mgo.Index{
		Key:        []string{"studentid"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err2 := session.DB(db.DatabaseName).C(db.StudentCollectionName).EnsureIndex(index)
	if err2 != nil {
		panic(err2)
	}
}

func (db *StudentMongoDB) Count() int {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	count, err := session.DB(db.DatabaseName).C(db.StudentCollectionName).Count()
	if err != nil {
		fmt.Println("error in Count", err)
		return -1
	}
	return count
}

func (db *StudentMongoDB) Add(s Student) error {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err1 := session.DB(db.DatabaseName).C(db.StudentCollectionName).Insert(s)
	if err1 != nil {
		return err1
	}

	return nil
}

func (db *StudentMongoDB) Get(keyID string) (Student, bool) {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	allWasGood := true
	student := Student{}

	c := session.DB(db.DatabaseName).C(db.StudentCollectionName)
	err1 := c.Find(bson.M{"studentid": keyID}).One(&student)

	if err1 != nil {
		//fmt.Println("71 error in Insert", err1)
		allWasGood = false
	}
	return student, allWasGood
}

func (db *StudentMongoDB) Delete(keyID string) bool {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err1 := session.DB(db.DatabaseName).C(db.StudentCollectionName).Remove(bson.M{"studentid": keyID})
	if err1 != nil {
		fmt.Println("error deleting from database")
		return false
	}

	return true
}

func (db *StudentMongoDB) DropCollection() bool {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err1 := session.DB(db.DatabaseName).C(db.StudentCollectionName).DropCollection()
	if err1 != nil {
		fmt.Println("unable to drop collection")
		return false
	}

	collections, err2 := session.DB(db.DatabaseName).CollectionNames()
	if err2 != nil {
		fmt.Println("unable to drop collection")
		return false
	}

	for _, name := range collections {
		if name == db.StudentCollectionName {
			fmt.Println("unable to drop collection")
			return false
		}
	}

	return true
}
