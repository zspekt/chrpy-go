package database

import (
	"encoding/json"
	"log"
	"os"
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

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	DBStruct := DBStructure{}

	// locking access to the file so no one writes to it, or reads before
	// we are done updating it
	db.mutex.Lock()
	fileBytes, err := os.ReadFile(db.path)
	if err != nil {
		return Chirp{}, err
	}

	// unmarshalling into the provided struct
	err = json.Unmarshal(fileBytes, &DBStruct)
	if err != nil {
		log.Println(err)
		return Chirp{}, err
	}

	newChirpId := len(DBStruct.Chirps) + 1
	newChirp := Chirp{
		Body: body,
		Id:   newChirpId,
	}

	DBStruct.Chirps[newChirpId] = newChirp
	db.mutex.Unlock()

	return newChirp, nil
}

func unmarshalDB()

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	DBStruct := DBStructure{}

	// locking access to the file so no one writes to it, or reads before
	// we are done updating it
	db.mutex.Lock()
	fileBytes, err := os.ReadFile(db.path)
	if err != nil {
		return []Chirp{}, err
	}

	// unmarshalling into the provided struct
	err = json.Unmarshal(fileBytes, &DBStruct)
	if err != nil {
		log.Println(err)
		return []Chirp{}, err
	}

	return []Chirp{}, nil
}

/*












 */

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error)

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error
