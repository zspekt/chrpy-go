package database

import (
	"encoding/json"
	"fmt"
	"log"
)

var UserIDCount int = 0

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

/*





















 */

func (db *DB) UpdateUserFields(replacement string, userId int) (struct{}, error) {
	strct := User{}
	dbStruct := DBStructure{}

	err := json.Unmarshal([]byte(replacement), &strct)
	if err != nil {
		log.Println("ERROR IN UpdateUserFields. RETURNING ERR FOR HANDLING...")
		return struct{}{}, err
	}

	dbStruct, err = db.LoadDB()
	if err != nil {
		log.Println(
			"ERROR IN UpdateUserFields WHEN CALLING LOADDB. RETURNING ERROR FOR HANDLING...",
		)
		return struct{}{}, err
	}

	fmt.Println(dbStruct)
	// dbStruct.Users[string(userId)].Password = strct.Password

	return struct{}{}, nil
}

/*





















 */
