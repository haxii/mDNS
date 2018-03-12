package badger

import (
	"time"

	"github.com/dgraph-io/badger"
)

var (
	db *badger.DB
)

//InitDB init badger db
func InitDB(dir, valueDir string) error {
	opts := badger.DefaultOptions
	opts.Dir = dir
	opts.ValueDir = valueDir
	_db, err := badger.Open(opts)
	if err != nil {
		return err
	}
	db = _db
	return nil
}

func CloseDB() error {
	var err error
	if db != nil {
		err = db.Close()
		db = nil
	}
	return err
}

func Set(key, val []byte) error {
	return db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, val)
		return err
	})
}

func Get(key []byte) ([]byte, error) {
	var value []byte
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		val, err := item.Value()
		if err != nil {
			return err
		}
		value = val
		return nil
	})
	return value, err
}

func SetWithTTL(key, val []byte, ttl time.Duration) error {
	return db.Update(func(txn *badger.Txn) error {
		err := txn.SetWithTTL([]byte(key), []byte(val), ttl)
		return err
	})
}
