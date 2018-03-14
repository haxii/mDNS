package badger

import (
	"bytes"
	"os"
	"testing"
	"time"
)

var (
	badgerDir = "./badger"

	testKey   = []byte("testKey")
	testValue = []byte("testValue")
)

func TestBadger(t *testing.T) {
	t.Run("testInitDB", testInitDB)
	t.Run("testSet", testSet)
	t.Run("testGet", testGet)
	t.Run("testSetWithTTL", testSetWithTTL)
	t.Run("testCloseDB", testCloseDB)

	os.RemoveAll(badgerDir)
}

func testInitDB(t *testing.T) {
	err := InitDB(badgerDir, badgerDir)
	if err != nil {
		t.Fatal(err)
	}
	if db == nil {
		t.Fatal("init db error")
	}
}

func testCloseDB(t *testing.T) {
	err := CloseDB()
	if err != nil {
		t.Fatal(err)
	}
	if db != nil {
		t.Fatal("close db error")
	}
}

func testSet(t *testing.T) {
	err := Set(testKey, testValue)
	if err != nil {
		t.Error(err)
	}
}

func testGet(t *testing.T) {
	val, err := Get(testKey)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(val, testValue) {
		t.Fail()
	}
}

func testSetWithTTL(t *testing.T) {
	err := SetWithTTL(testKey, testValue, time.Second)
	if err != nil {
		t.Error(err)
	}

	val, err := Get(testKey)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(val, testValue) {
		t.Fail()
	}

	time.Sleep(time.Second)
	val, _ = Get(testKey)
	if len(val) > 0 {
		t.Fail()
	}
}
