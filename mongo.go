package main

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type MongoDbStruct struct {
	Host           string
	DatabaseName   string
	collectionName string
}

func (db *MongoDbStruct) Init(dbName, collec, host string) {
	db.DatabaseName = dbName   //
	db.Host = host             //"mongodb://127.0.0.1:27017"
	db.collectionName = collec //"teststrudentdb"

	index := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err2 := session.DB(db.DatabaseName).C(db.collectionName).EnsureIndex(index)
	if err2 != nil {
		panic(err2)
	}
}

func (db *MongoDbStruct) Count() int {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	count, err := session.DB(db.DatabaseName).C(db.collectionName).Count()
	if err != nil {
		fmt.Println("error in Count", err)
		return -1
	}
	return count
}

func (db *MongoDbStruct) Add(s Meta) error {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err1 := session.DB(db.DatabaseName).C(db.collectionName).Insert(s)
	if err1 != nil {
		return err1
	}

	return nil
}

func (db *MongoDbStruct) Get(keyID string) (Meta, bool) {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	allWasGood := true
	igcMeta := Meta{}

	err1 := session.DB(db.DatabaseName).C(db.collectionName).Find(bson.M{"id": keyID}).One(&igcMeta)
	if err1 != nil {
		allWasGood = false
	}
	return igcMeta, allWasGood
}

func (db *MongoDbStruct) Delete(keyID string) bool {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err1 := session.DB(db.DatabaseName).C(db.collectionName).Remove(bson.M{"id": keyID})
	if err1 != nil {
		fmt.Println("error deleting from database")
		return false
	}

	return true
}

func (db *MongoDbStruct) DropCollection() bool {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err1 := session.DB(db.DatabaseName).C(db.collectionName).DropCollection()
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
		if name == db.collectionName {
			fmt.Println("unable to drop collection")
			return false
		}
	}

	return true
}

func (db *MongoDbStruct) GetAllKeys() ([]string, bool) {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	ok := true
	var ids []string

	err1 := session.DB(db.DatabaseName).C(db.collectionName).Find(bson.M{}).Distinct("id", &ids)
	if err1 != nil {
		fmt.Println("error retriving from DB")
		ok = false
	}

	return ids, ok
}

/*
func (db *MongoDbStruct) GetAllKeys() ([]string, bool) {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	ok := true
	var ids []string

	err1 := session.DB(db.DatabaseName).C(db.collectionName).Find().Sort({"time_stamp": -1})
	db.collection.find().sort({_id: -1}).limit(1)
	if err1 != nil {
		fmt.Println("error retriving from DB")
		ok = false
	}

	return ids, ok
}
*/
