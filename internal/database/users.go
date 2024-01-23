package database

import (
	"errors"
	"fmt"
	"log"
)

var UserIDCount int = 0

func (db *DB) UpgradeUser(userID int) (User, error) {
	DBStruct := DBStructure{}
	userReturn := User{}

	db.mutex.Lock()
	defer db.mutex.Unlock()

	DBStruct, err := db.LoadDB()
	if err != nil {
		log.Println(err)
		return User{}, err
	}

	if DBStruct.Users == nil {
		log.Println("Users map is nil")
		return User{}, fmt.Errorf("Users map is nil")
	}

	var ok bool

	// if user doesn't exist or has already upgraded...
	if userReturn, ok = DBStruct.Users[userID]; !ok {
		log.Println("User not in database")
		return userReturn, errors.New("User does not exist in database")
	}
	if userReturn.IsChirpyRed {
		log.Printf("User %v has already upgraded to Chirpy Red\n", userReturn.Id)
		return userReturn, fmt.Errorf(
			"User %v has already upgraded to Chirpy Red\n",
			userReturn.Id,
		)
	}

	userReturn.IsChirpyRed = true
	DBStruct.Users[userID] = userReturn
	err = MarshalAndWrite(DBStruct, db.path)
	if err != nil {
		log.Println(err)
		return User{}, err
	}
	return userReturn, nil
}

func (db *DB) CreateUser(email string, password string) (User, error) {
	DBStruct := DBStructure{}

	DBStruct.Users = map[int]User{}

	// locking access to the file so no one writes to it, or reads before
	// we are done updating it
	db.mutex.Lock()
	defer db.mutex.Unlock()

	err := UnmarshalToStruct[DBStructure](&DBStruct, db.path)
	if err != nil {
		log.Println(err)
		return User{}, err
	}

	UserIDCount++

	newUser := User{
		Id:          UserIDCount,
		Email:       email,
		Password:    password,
		IsChirpyRed: false,
	}

	DBStruct.Users[UserIDCount] = newUser
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

	return newUser, nil
}

/*





















 */

func (db *DB) UpdateUserFields(user User) error {
	dbStruct := DBStructure{}

	db.mutex.Lock()

	dbStruct, err := db.LoadDB()
	if err != nil {
		return err
	}

	dbStruct.Users[user.Id] = user

	MarshalAndWrite(dbStruct, db.path)

	db.mutex.Unlock()
	return nil
}

func GetUserID(userMap map[int]User, email string) (int, error) {
	for k, v := range userMap {
		if v.Email == email {
			return k, nil
		}
	}
	log.Println("GetUserID couldn't find user", email)
	return 0, fmt.Errorf("\nUser %v couldn't be found in our database\n", email)
}

/*





















 */
