package database

import (
	"fmt"
	"sort"
)

var ChirpIDCount int = 0

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
