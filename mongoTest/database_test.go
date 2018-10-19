package main

import (
	"github.com/globalsign/mgo"
	"testing"
)

func setupDB(t *testing.T) *StudentMongoDB {
	db := StudentMongoDB{
		"mongodb://127.0.0.1:27017",
		"teststrudentdb",
		"students",
	}
	session, err := mgo.Dial(db.Host)
	if err != nil {
		t.Error(err)
	}
	defer session.Close()
	return &db
}

func tearDownDB(t *testing.T, db *StudentMongoDB) {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		t.Error(err)
	}
	defer session.Close()

	err1 := session.DB(db.DatabaseName).DropDatabase()
	if err1 != nil {
		t.Error(err1)
	}

}

func TestStudentMongoDB_Add(t *testing.T) {

	db := setupDB(t)
	defer tearDownDB(t, db)

	db.Init()
	if db.Count() != 0 {
		t.Error("database not properly initialized, student count should be 0")
	}
	studnt := Student{"tom", 21, "id1"}
	err := db.Add(studnt)
	if err != nil {
		t.Error("unable to add students")
	}

	if db.Count() != 1 {
		t.Error("count should be 1")
	}
}

func TestStudentMongoDB_Duplicates(t *testing.T) {

	db := setupDB(t)
	defer tearDownDB(t, db)

	db.Init()
	if db.Count() != 0 {
		t.Error("database not properly initialized, student count should be 0")
	}
	studnt := Student{"tom", 21, "id1"}
	err := db.Add(studnt)
	if err != nil {
		t.Error("unable to add students", err)
	}
	if db.Count() != 1 {
		t.Error("count should be 1")
	}

	err2 := db.Add(studnt)
	if err2 == nil {
		t.Error("error duplicate added to database", err2)
	}

	if db.Count() != 1 {
		t.Error("count should be 1")
	}
}

func TestStudentMongoDB_Get(t *testing.T) {
	db := setupDB(t)
	defer tearDownDB(t, db)

	db.Init()
	if db.Count() != 0 {
		t.Error("database not properly initialized, student count should be 0")

	}
	student := Student{"tom", 21, "id1"}
	db.Add(student)

	if db.Count() != 1 {
		t.Error("adding new student failed")
	}

	newStudent, ok := db.Get(student.StudentID)
	if !ok {
		t.Error("couldn't find tom")
	}

	if newStudent.Name != student.Name ||
		newStudent.Age != student.Age ||
		newStudent.StudentID != student.StudentID {

		t.Error("student do not mach", newStudent.Name, "'", student.Name)

	}
}

func TestStudentMongoDB_Delete(t *testing.T) {
	db := setupDB(t)
	defer tearDownDB(t, db)

	db.Init()
	if db.Count() != 0 {
		t.Error("database not properly initialized, student count should be 0")
	}
	studnt := Student{"tom", 21, "id1"}
	db.Add(studnt)

	if db.Count() != 1 {
		t.Error("count should be 1")
	}

	ok := db.Delete(studnt.StudentID)
	if !ok {
		t.Error("unable to delete student")
	}

	if db.Count() != 0 {
		t.Error("database not properly initialized, student count should be 0")
	}
}

func TestStudentMongoDB_DropCollection(t *testing.T) {
	db := setupDB(t)
	defer tearDownDB(t, db)

	db.Init()
	if db.Count() != 0 {
		t.Error("database not properly initialized, student count should be 0")
	}
	studnt := Student{"tom", 21, "id1"}
	db.Add(studnt)

	if db.Count() != 1 {
		t.Error("count should be 1")
	}

	ok := db.DropCollection()
	if !ok {
		t.Error("unable to delete collection")
	}

}
