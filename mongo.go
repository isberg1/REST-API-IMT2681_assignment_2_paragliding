package main

import (
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// set up the track collection to have a unique key "id"
func (db *mongoDbStruct) initTrackCollection(dbName, collec, host string) {
	db.DatabaseName = dbName //
	db.Host = host           //"mongodb://127.0.0.1:27017"
	db.collection = collec   //"teststrudentdb"

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

	err2 := session.DB(db.DatabaseName).C(db.collection).EnsureIndex(index)
	if err2 != nil {
		panic(err2)
	}


}




// set up the webhook collection to have a unique key "id"
func (db *mongoDbStruct) initWebHookCollection(dbName, collec, host string) {
	db.DatabaseName = dbName //
	db.Host = host           //"mongodb://127.0.0.1:27017"
	db.collection = collec   //

	index := mgo.Index{
		Key:        []string{"web_hook_url"},
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

	err2 := session.DB(db.DatabaseName).C(db.collection).EnsureIndex(index)
	if err2 != nil {
		panic(err2)
	}
}

// returns the collection count
func (db *mongoDbStruct) count() int {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	count, err := session.DB(db.DatabaseName).C(db.collection).Count()
	if err != nil {
		fmt.Println("error in count", err)
		return -1
	}
	return count
}

// adds a document to the database
func (db *mongoDbStruct) add(s interface{}) error {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err1 := session.DB(db.DatabaseName).C(db.collection).Insert(s)
	if err1 != nil {
		return err1
	}

	return nil
}

// gets a IGC META document, converts and returns it to a struct
func (db *mongoDbStruct) getMetaByID(keyID string) (Meta, bool) {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	allWasGood := true
	igcMeta := Meta{}

	err1 := session.DB(db.DatabaseName).C(db.collection).Find(bson.M{"id": keyID}).One(&igcMeta)
	if err1 != nil {
		allWasGood = false
	}
	return igcMeta, allWasGood
}

// delete a document based on the unique key "id"
func (db *mongoDbStruct) delete(keyID string) bool {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err1 := session.DB(db.DatabaseName).C(db.collection).Remove(bson.M{"id": keyID})
	if err1 != nil {
		fmt.Println("error deleting from database")
		return false
	}

	return true
}

// drops a collection from the DB
func (db *mongoDbStruct) dropCollection() bool {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err1 := session.DB(db.DatabaseName).C(db.collection).DropCollection()
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
		if name == db.collection {
			fmt.Println("unable to drop collection")
			return false
		}
	}

	return true
}

// gets all keys "id" from a collection, and returns an string array of ids
func (db *mongoDbStruct) getAllKeys() ([]string, bool) {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	ok := true
	var ids []string

	err1 := session.DB(db.DatabaseName).C(db.collection).Find(bson.M{}).Distinct("id", &ids)
	if err1 != nil {
		fmt.Println("error retriving from DB")
		ok = false
	}

	return ids, ok
}

// get the timestamp of the latest posted IGC meta track to be posted
func (db *mongoDbStruct) getLatestMetaTimestamp() (int64, bool) {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	ok := true
	var timestamp Meta

	err1 := session.DB(db.DatabaseName).C(db.collection).Find(nil).Sort("-timestamp").Limit(1).One(&timestamp)
	if err1 != nil {
		fmt.Println("error(getLatestMetaTimestamp) retriving from DB", err1)
		ok = false
	}

	return timestamp.TimeStamp, ok
}

// get the timestamp of the oldest posted IGC meta track to be posted
func (db *mongoDbStruct) getOldestMetaByTimeStamp() (int64, bool) {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	ok := true
	var timestamp Meta

	err1 := session.DB(db.DatabaseName).C(db.collection).Find(nil).Sort("timestamp").Limit(1).One(&timestamp)
	if err1 != nil {
		fmt.Println("error(getOldestMetaByTimeStamp) retriving from DB", err1)
		ok = false
	}

	return timestamp.TimeStamp, ok
}

// gets and returns a IGC MEATA struct based on a timestamp
func (db *mongoDbStruct) getMetaByTimstamp(timeStamp int64) (Meta, bool) {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	ok := true
	var igcFile Meta

	err1 := session.DB(db.DatabaseName).C(db.collection).Find(bson.M{"timestamp": timeStamp}).One(&igcFile)
	if err1 != nil {
		fmt.Println("error(getOldestMetaByTimeStamp) retriving from DB", err1)
		ok = false
	}

	return igcFile, ok
}

// gets and returns a webhook struct struct based on a timestamp
func (db *mongoDbStruct) getWebHookByTimstamp(timeStamp int64) (WebHookStruct, bool) {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	ok := true
	var wHook WebHookStruct

	err1 := session.DB(db.DatabaseName).C(db.collection).Find(bson.M{"timestamp": timeStamp}).One(&wHook)
	if err1 != nil {
		fmt.Println("error(getOldestMetaByTimeStamp) retriving from DB", err1)
		ok = false
	}

	return wHook, ok
}

// get Meta struct bigger then timstamp parameter
func (db *mongoDbStruct) getMetaBiggerThen(timeStamp int64) (Meta, error) {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var igcFile Meta

	err1 := session.DB(db.DatabaseName).C(db.collection).Find(bson.M{"timestamp": bson.M{"$gt": timeStamp}}).Sort("timestamp").Limit(1).One(&igcFile) //
	if err1 != nil {
		fmt.Println("error(getMetaBiggerThen) retriving from DB", err1, igcFile)
	}

	return igcFile, err1
}

// count down the counter in all webhook documents
func (db *mongoDbStruct) counterDown() {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	_, err1 := session.DB(db.DatabaseName).C(db.collection).UpdateAll(bson.M{}, bson.M{"$inc": bson.M{"counter": -1}})
	if err1 != nil {
		fmt.Println("(counterDown)", err1)
	}
}

// get all webhook document wher counter is les then 1
func (db *mongoDbStruct) getPostArray() ([]WebHookStruct, error) {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var webHook []WebHookStruct
	err1 := session.DB(db.DatabaseName).C(db.collection).Find(bson.M{"counter": bson.M{"$lt": 1}}).All(&webHook)
	if err1 != nil {
		fmt.Println("(counterDown)", err1)
	}

	return webHook, err1
}

// get the latest "lastNr" nr of IGC document entries
func (db *mongoDbStruct) getLatestMetaIDs(lastNr int) ([]ResponsID, error) {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var ids []ResponsID
	err1 := session.DB(db.DatabaseName).C(db.collection).Find(nil).Sort("-timestamp").Limit(lastNr).All(&ids)
	if err1 != nil {
		fmt.Println("(counterDown)", err1)
	}

	return ids, err1
}

// resets all webhook counters to their default( MinTriggerValue) value
func (db *mongoDbStruct) counterReset(webHookArray []WebHookStruct) {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	for _, val := range webHookArray {
		val.Counter = val.MinTriggerValue
		err1 := session.DB(db.DatabaseName).C(db.collection).Update(bson.M{"id": val.ID}, bson.M{"$set": bson.M{"counter": val.MinTriggerValue}})
		if err1 != nil {
			fmt.Println("(counterDown)", err1)
		}
	}
}

// gets and returns a webhook document as a struct
func (db *mongoDbStruct) getWebHookByID(keyID string) (WebHookStruct, bool) {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	allWasGood := true
	webHook := WebHookStruct{}
	//Find(bson.M{"id": keyID}).One(&igcMeta)
	err1 := session.DB(db.DatabaseName).C(db.collection).Find(bson.M{"id": keyID}).One(&webHook)
	if err1 != nil {
		allWasGood = false
	}
	return webHook, allWasGood
}

// delete a webhook document based on a "id"
func (db *mongoDbStruct) deleteWebHook(keyID string) (WebHookStruct, error) {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	webHook := WebHookStruct{}
	err1 := session.DB(db.DatabaseName).C(db.collection).Find(bson.M{"id": keyID}).One(&webHook)
	if err1 != nil {
		return WebHookStruct{}, err1
	}
	err2 := session.DB(db.DatabaseName).C(db.collection).Remove(bson.M{"id": keyID})
	if err2 != nil {
		return WebHookStruct{}, err2
	}
	return webHook, nil
}

// deletes a collection form the DB
func (db *mongoDbStruct) dropTable() error {
	if db.count() == 0 {
		return nil
	}

	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err1 := session.DB(db.DatabaseName).C(db.collection).DropCollection()
	if err1 != nil {
		return err1
	}

	return nil
}
