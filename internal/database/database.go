package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
)

var (
	ChirpIDCount int = 0
	UserIDCount  int = 0
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

// Reads the file at fPath, assumes json content, and unmarshalls it to
// the provided memory address, expecting a struct of type T.
func UnmarshalToStruct[T any](structure *T, fPath string) error {
	fileBytes, err := os.ReadFile(fPath)
	if err != nil {
		log.Println(err)
		return err
	}

	if !json.Valid(fileBytes) {
		log.Println("\n", fPath, "does not contain valid JSON data.")
		return nil
	}

	// unmarshalling into the provided struct
	err = json.Unmarshal(fileBytes, structure)
	if err != nil {
		panic(err)
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

	DBStruct.Chirps = map[int]Chirp{}

	// locking access to the file so no one writes to it, or reads before
	// we are done updating it
	db.mutex.Lock()

	UnmarshalToStruct[DBStructure](&DBStruct, db.path)

	ChirpIDCount++

	newChirp := Chirp{
		Body: body,
		Id:   ChirpIDCount,
	}

	DBStruct.Chirps[ChirpIDCount] = newChirp
	fmt.Printf(
		"\n\n\tCreated new chirp with ID -> %v\n\t\tBody -> %v\n\n",
		newChirp.Id,
		newChirp.Body,
	)

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

func (db *DB) CreateUser(email string, password string) (User, error) {
	DBStruct := DBStructure{}

	DBStruct.Users = map[string]User{}

	// locking access to the file so no one writes to it, or reads before
	// we are done updating it
	db.mutex.Lock()

	err := UnmarshalToStruct[DBStructure](&DBStruct, db.path)
	if err != nil {
		log.Println(err)
		return User{}, err
	}

	UserIDCount++

	newUser := User{
		Id:       UserIDCount,
		Email:    email,
		Password: password,
	}

	DBStruct.Users[email] = newUser
	fmt.Printf(
		"\n\n\tCreated new User with ID -> %v\n\t\tBody -> %v\n\tPassword -> %v\n\n",
		newUser.Id,
		newUser.Email,
		newUser.Password,
	)

	err = MarshalAndWrite[DBStructure](DBStruct, db.path)
	if err != nil {
		return User{}, err
	}

	db.mutex.Unlock()
	return newUser, nil
}
