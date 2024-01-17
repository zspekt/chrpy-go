package main

import (
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/zspekt/chrpy-go/internal/database"
	jwtwrappers "github.com/zspekt/chrpy-go/internal/jwtWrappers"
)

func usersPostHandler(w http.ResponseWriter, r *http.Request) {
	decdRequest := decodeUserLogin{}

	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Println(err)
		return
	}

	decodeJson(r.Body, &decdRequest)

	password, err := hashPassword(decdRequest.Password)
	if err != nil {
		return
	}

	user, err := db.CreateUser(decdRequest.Email, password)
	if err != nil {
		log.Fatal(err)
		return
	}

	resp := userPostResp{
		Id:    user.Id,
		Email: user.Email,
	}

	respondWithJSON(w, 201, resp)
}

func usersAuthHandler(w http.ResponseWriter, r *http.Request) {
	decdRequest := decodeUserLogin{}

	resp := userPostResp{}

	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Println(err)
		return
	}

	decodeJson(r.Body, &decdRequest)

	requestedUser := decdRequest.Email

	DBStruct, err := db.LoadDB()
	if err != nil {
		log.Printf("Error loading DB into memory --> %v\n", err)
		return
	}

	// we retrieve the password a
	hashedPass := []byte(DBStruct.Users[requestedUser].Password)

	err = bcrypt.CompareHashAndPassword(hashedPass, []byte(decdRequest.Password))
	// nil if match
	if err != nil {
		respondWithError(w, 401, "Unauthorized access")
		return
	}

	jwtCfg := &jwtwrappers.JWTRequestConfig{
		UserID:           DBStruct.Users[requestedUser].Id,
		ExpiresInSeconds: decdRequest.ExpiresInSeconds,
	}

	signedToken, err := jwtwrappers.CreateToken(jwtCfg)
	if err != nil {
		log.Println("Error creating and signing token -> ", err)
		return
	}

	resp = userPostResp{
		Id:    DBStruct.Users[requestedUser].Id,
		Email: requestedUser,
		Token: signedToken,
	}

	respondWithJSON(w, 200, resp)
}

func usersEditHandler(w http.ResponseWriter, r *http.Request) {
}
