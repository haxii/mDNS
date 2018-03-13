package badger

import (
	"os"
	"testing"
	"time"
)

var (
	badgerDir = "./badger"
)

func TestBadger(t *testing.T) {
	t.Run("testInitDB", testInitDB)
	t.Run("testSet", testSet)
	t.Run("testGet", testGet)
	t.Run("testSetWithTTL", testSetWithTTL)

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
	err := Set([]byte("test_"), []byte("value_"))
	if err != nil {
		t.Error(err)
	}
}

func testGet(t *testing.T) {
	val, err := Get([]byte("test_"))
	if err != nil {
		t.Error(err)
	}
	if string(val) != "value_" {
		t.Fail()
	}
}

func testSetWithTTL(t *testing.T) {
	err := SetWithTTL([]byte("test_"), []byte("value_"), time.Second)
	if err != nil {
		t.Error(err)
	}

	val, err := Get([]byte("test_"))
	if err != nil {
		t.Error(err)
	}
	if string(val) != "value_" {
		t.Fail()
	}

	time.Sleep(time.Second)
	val, _ = Get([]byte("test_"))
	if len(val) > 0 {
		t.Fail()
	}
}
