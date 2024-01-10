package database

import (
	"encoding/json"
	"log"
	"os"
	"sort"
	"sync"
)

type Chirp struct {
	Body string `json:"body"`
	Id   int    `json:"id"`
}

type DB struct {
	path  string
	mutex *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

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
			log.Printf("\tDB file %v did not exist. Creating it...\n", path)
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

func UnmarshalToStruct[T any](structure *T, fPath string) error {
	fileBytes, err := os.ReadFile(fPath)
	if err != nil {
		return err
	}

	// unmarshalling into the provided struct
	err = json.Unmarshal(fileBytes, structure)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Marshals the struct and writes it to the file.
func MarshalAndWrite[T any](structure T, fPath string) error {
	bytes, err := json.Marshal(structure)
	if err != nil {
		return err
	}

	err = os.WriteFile(fPath, bytes, 0700)
	if err != nil {
		return err
	}

	return nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	DBStruct := DBStructure{}

	DBStruct.Chirps = make(map[int]Chirp)

	// locking access to the file so no one writes to it, or reads before
	// we are done updating it
	db.mutex.Lock()

	UnmarshalToStruct[DBStructure](&DBStruct, db.path)

	newChirpId := len(DBStruct.Chirps) + 1
	newChirp := Chirp{
		Body: body,
		Id:   newChirpId,
	}

	DBStruct.Chirps[newChirpId] = newChirp

	err := MarshalAndWrite[DBStructure](DBStruct, db.path)
	if err != nil {
		return Chirp{}, err
	}

	db.mutex.Unlock()
	return newChirp, nil
}

// func unmarshalDB()

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	var chirpList []Chirp
	DBStruct := DBStructure{}

	// locking access to the file so no one writes to it, or reads before
	// we are done updating it
	db.mutex.Lock()

	err := UnmarshalToStruct[DBStructure](&DBStruct, db.path)
	if err != nil {
		return []Chirp{}, err
	}

	for _, chirp := range DBStruct.Chirps {
		chirpList = append(chirpList, chirp)
	}

	sort.Slice(chirpList, func(i, j int) bool {
		return chirpList[i].Id < chirpList[j].Id
	})

	return chirpList, nil
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
func (db *DB) loadDB() (DBStructure, error) {
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

func (db *DB) GetIdCount() (int, error) {
	d, err := db.loadDB()
	if err != nil {
		return 0, err
	}

	return len(d.Chirps), nil
}
