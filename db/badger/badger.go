package badger

import (
	"time"

	"github.com/dgraph-io/badger"
)

var (
	ErrKeyNotFound = badger.ErrKeyNotFound
)

type BadgerDB struct {
	db *badger.DB
}

// OpenDB opens badger db
func OpenBadger(dir, valueDir string) (*BadgerDB, error) {
	opts := badger.DefaultOptions
	opts.Dir = dir
	opts.ValueDir = valueDir
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &BadgerDB{db: db}, nil
}

// Close closes db
func (b *BadgerDB) Close() error {
	if b.db != nil {
		err := b.db.Close()
		if err != nil {
			return err
		}
		b.db = nil
	}
	return nil
}

// Set sets key, val to db
func (b *BadgerDB) Set(key, val []byte) error {
	return b.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, val)
		return err
	})
}

// Get returns val and error
func (b *BadgerDB) Get(key []byte) ([]byte, error) {
	var value []byte
	err := b.db.View(func(txn *badger.Txn) error {
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
func (b *BadgerDB) SetWithTTL(key, val []byte, ttl time.Duration) error {
	return b.db.Update(func(txn *badger.Txn) error {
		err := txn.SetWithTTL([]byte(key), []byte(val), ttl)
		return err
	})
}
