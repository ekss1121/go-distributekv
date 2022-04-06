package db

import (
	"fmt"

	"github.com/boltdb/bolt"
)

type Database struct {
	db *bolt.DB
}

var defaultBucket = "defaultBucket"

func NewDatabase(dbLocation string) (*Database, error, func() error) {
	db, err := bolt.Open(dbLocation, 0600, nil)
	datbase := Database{
		db: db,
	}
	if err != nil {
		return &datbase, err, db.Close
	}
	if err = datbase.createDefaultBucket(); err != nil {
		return &datbase, err, db.Close
	}

	return &datbase, nil, db.Close
}

func (DB *Database) Get(key string) ([]byte, error) {
	var value []byte
	err := DB.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(defaultBucket))
		value = b.Get([]byte(key))
		return nil
	})

	if err == nil {
		return value, nil
	}
	return nil, err
}

func (DB *Database) Set(key string, value []byte) error {
	return DB.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(defaultBucket))
		err := b.Put([]byte(key), []byte(value))
		return err
	})
}

func (DB *Database) createDefaultBucket() error {
	return DB.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(defaultBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}
