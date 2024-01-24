package database

import (
	"fmt"
	"log"
	"sort"
)

var ChirpIDCount int = 0

// GetChirps returns all chirps belonging to the same author, if provided with
// an authorID. Otherwise it simply returns all chirps. Sorts them in both cases.
func (db *DB) GetChirps(authorID ...int) ([]Chirp, error) {
	var chirpList []Chirp
	DBStruct := DBStructure{}

	log.Printf("Provided authorID -> %v\n", authorID)

	db.mutex.Lock()
	defer db.mutex.Unlock()

	err := UnmarshalToStruct[DBStructure](&DBStruct, db.path)
	if err != nil {
		return []Chirp{}, err
	}

	// if provided an authorID...
	if len(authorID) != 0 {
		for _, chirp := range DBStruct.Chirps {
			if chirp.AuthorId == authorID[0] {
				chirpList = append(chirpList, chirp)
			}
			sort.Slice(chirpList, func(i, j int) bool {
				return chirpList[i].ChirpId < chirpList[j].ChirpId
			})
		}
		return chirpList, nil
	}

	for _, chirp := range DBStruct.Chirps {
		chirpList = append(chirpList, chirp)
	}

	sort.Slice(chirpList, func(i, j int) bool {
		return chirpList[i].ChirpId < chirpList[j].ChirpId
	})

	return chirpList, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string, authorId int) (Chirp, error) {
	DBStruct := DBStructure{}

	// locking access to the file so no one writes to it, or reads before
	// we are done updating it
	db.mutex.Lock()

	UnmarshalToStruct[DBStructure](&DBStruct, db.path)

	if DBStruct.Chirps == nil {
		DBStruct.Chirps = make(map[int]Chirp)
	}

	ChirpIDCount++

	newChirp := Chirp{
		Body:     body,
		ChirpId:  ChirpIDCount,
		AuthorId: authorId,
	}

	DBStruct.Chirps[ChirpIDCount] = newChirp
	fmt.Printf(
		"\n\n\tCreated new chirp with ID -> %v\n\t\tBody -> %v\n\n",
		newChirp.ChirpId,
		newChirp.Body,
	)

	err := MarshalAndWrite[DBStructure](DBStruct, db.path)
	if err != nil {
		return Chirp{}, err
	}

	db.mutex.Unlock()
	return newChirp, nil
}

func (db *DB) DeleteChirp(chirpId int, userId int) error {
	DBStruct := DBStructure{}

	// locking access to the file so no one writes to it, or reads before
	// we are done updating it
	db.mutex.Lock()

	UnmarshalToStruct[DBStructure](&DBStruct, db.path)

	if DBStruct.Chirps == nil {
		return fmt.Errorf("Chirps map is nil")
	}

	if DBStruct.Chirps[chirpId].AuthorId != userId {
		log.Println(
			"User trying to delete chirp belonging to an account they haven't authenticated as",
		)
		return fmt.Errorf(
			"User trying to delete chirp from a different account",
			// userId,
			// DBStruct.Chirps[chirpId].AuthorId,
			// chirpId,
		)
	}

	delete(DBStruct.Chirps, chirpId)

	err := MarshalAndWrite(DBStruct, "./database.json")
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
