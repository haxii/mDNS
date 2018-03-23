package badger

import (
	"time"

	"github.com/dgraph-io/badger"
)

var (
	//ErrKeyNotFound means Key not found
	ErrKeyNotFound = badger.ErrKeyNotFound
)

//DB badger db
type DB struct {
	db *badger.DB
}

// OpenBadger opens badger db
func OpenBadger(dir, valueDir string) (*DB, error) {
	opts := badger.DefaultOptions
	opts.Dir = dir
	opts.ValueDir = valueDir
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

// Close closes db
func (db *DB) Close() error {
	if db.db != nil {
		err := db.db.Close()
		if err != nil {
			return err
		}
		db.db = nil
	}
	return nil
}

// Set sets key, val to db
func (db *DB) Set(key, val []byte) error {
	return db.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, val)
		return err
	})
}

// Get returns val and error
func (db *DB) Get(key []byte) ([]byte, error) {
	var value []byte
	err := db.db.View(func(txn *badger.Txn) error {
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

// SetWithTTL sets key, val to db with ttl
func (db *DB) SetWithTTL(key, val []byte, ttl time.Duration) error {
	return db.db.Update(func(txn *badger.Txn) error {
		err := txn.SetWithTTL([]byte(key), []byte(val), ttl)
		return err
	})
}
