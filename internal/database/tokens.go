package database

import (
	"log"
	"time"
)

// revokes token. adds it to the DBStruct.Tokens map. key is the token string.
// the value is the point in time where the token was revoked (time.Now())
func (db *DB) RevokeToken(token string) error {
	DBStruct := DBStructure{}

	// locking access to the file so no one writes to it, or reads before
	// we are done updating it
	db.mutex.Lock()
	defer db.mutex.Unlock()

	err := UnmarshalToStruct(&DBStruct, db.path)
	if err != nil {
		log.Fatal(err)
		return err
	}

	if DBStruct.Tokens == nil {
		DBStruct.Tokens = make(map[string]time.Time)
	}

	DBStruct.Tokens[token] = time.Now()

	err = MarshalAndWrite(DBStruct, db.path)
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Println("Token has now been revoked")

	return nil
}

func (db *DB) IsRevoked(token string) (bool, error) {
	DBStruct := DBStructure{}

	db.mutex.Lock()
	defer db.mutex.Unlock()

	DBStruct, err := db.LoadDB()
	if err != nil {
		return true, err
	}

	tokens := DBStruct.Tokens

	if _, ok := tokens[token]; ok {
		log.Println("Token is revoked...")
		return true, nil
	}

	return false, nil
}

/*












 */
