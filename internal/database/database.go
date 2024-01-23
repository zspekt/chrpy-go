package database

import (
	"log"
	"os"
	"sync"
)

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(path)
			if err != nil {
				log.Println(err)
				return &DB{}, err
			}
			log.Printf("DB file %v did not exist. Creating it...\n", path)
		} else {
			log.Println(err)
			return &DB{}, err
		}
	}

	return &DB{
		path:  path,
		mutex: &sync.RWMutex{},
	}, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	if _, err := os.Stat(db.path); err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(db.path)
			if err != nil {
				log.Println(err)
				return err
			}
			log.Printf("\tDB file %v did not exist. Creating it...\n", db.path)
		} else {
			log.Println(err)
			return err
		}
	}

	return nil
}

// loadDB reads the database file into memory
func (db *DB) LoadDB() (DBStructure, error) {
	var DBStruct DBStructure

	err := UnmarshalToStruct[DBStructure](&DBStruct, db.path)
	if err != nil {
		log.Println(err)
		return DBStructure{}, err
	}

	return DBStruct, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	err := MarshalAndWrite[DBStructure](dbStructure, db.path)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
