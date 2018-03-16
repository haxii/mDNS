package badger

import (
	"bytes"
	"os"
	"testing"
	"time"
)

var (
	db        *BadgerDB
	badgerDir = "./TestBadger"

	testKey   = []byte("testKey")
	testValue = []byte("testValue")
)

func TestBadger(t *testing.T) {
	t.Run("testOpenDB", testOpenDB)
	t.Run("testSet", testSet)
	t.Run("testGet", testGet)
	t.Run("testSetWithTTL", testSetWithTTL)
	t.Run("testCloseDB", testCloseDB)

	os.RemoveAll(badgerDir)
}

func testOpenDB(t *testing.T) {
	var err error
	db, err = OpenBadger(badgerDir, badgerDir)
	if err != nil {
		t.Fatal(err)
	}
	if db == nil {
		t.Fatal("init db error")
	}
}

func testCloseDB(t *testing.T) {
	err := db.Close()
	if err != nil {
		t.Fatal(err)
	}
	if db.db != nil {
		t.Fatal("close db error")
	}
}

func testSet(t *testing.T) {
	err := db.Set(testKey, testValue)
	if err != nil {
		t.Error(err)
	}
}

func testGet(t *testing.T) {
	val, err := db.Get(testKey)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(val, testValue) {
		t.Fail()
	}
}

func testSetWithTTL(t *testing.T) {
	err := db.SetWithTTL(testKey, testValue, time.Second)
	if err != nil {
		t.Error(err)
	}

	val, err := db.Get(testKey)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(val, testValue) {
		t.Fail()
	}

	time.Sleep(time.Second)
	val, _ = db.Get(testKey)
	if len(val) > 0 {
		t.Fail()
	}
}
