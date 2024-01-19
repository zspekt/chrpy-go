package main

import (
	"fmt"
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
		log.Fatal(err)
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

	resp := userLoginResp{}

	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Fatal(err)
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
	fmt.Printf("\ncreated   token: %v\n", signedToken)
	if err != nil {
		log.Println("Error creating and signing token -> ", err)
		return
	}

	resp = userLoginResp{
		Id:    DBStruct.Users[requestedUser].Id,
		Email: requestedUser,
		Token: signedToken,
	}

	respondWithJSON(w, 200, resp)
}

func usersEditHandler(w http.ResponseWriter, r *http.Request) {
	token, err := jwtwrappers.GetTokenFromHeader(r)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("\nretrieved token: %v\n", token)

	claims, err := jwtwrappers.ValidateAndReturn(token)
	if err != nil {
		// error is here
		log.Println(err)
		return
	}

	// unfinished code, obviously

	fmt.Printf("\nPrinting claims...\n%v\t%v\t%v", claims.Issuer, claims.IssuedAt, claims.ExpiresAt)
}
