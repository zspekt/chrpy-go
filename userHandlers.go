package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

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

	DBStruct, err := db.LoadDB()
	if err != nil {
		log.Printf("Error loading DB into memory --> %v\n", err)
		return
	}

	requestedUserId, err := database.GetUserID(DBStruct.Users, decdRequest.Email)
	if err != nil {
		log.Fatalf("Error retrieving USERID in AUTH HANDLER FUNC -> %v\n", err)
		return
	}

	// we retrieve the password a
	hashedPass := []byte(DBStruct.Users[requestedUserId].Password)

	err = bcrypt.CompareHashAndPassword(hashedPass, []byte(decdRequest.Password))
	// nil if match
	if err != nil {
		respondWithError(w, 401, "Unauthorized access")
		return
	}

	jwtCfg := &jwtwrappers.JWTRequestConfig{
		UserID:           DBStruct.Users[requestedUserId].Id,
		ExpiresInSeconds: decdRequest.ExpiresInSeconds,
	}

	signedToken, err := jwtwrappers.CreateToken(jwtCfg)
	fmt.Printf("\ncreated   token: %v\n", signedToken)
	if err != nil {
		log.Println("Error creating and signing token -> ", err)
		return
	}

	resp = userLoginResp{
		Id:    requestedUserId,
		Email: DBStruct.Users[requestedUserId].Email,
		Token: signedToken,
	}

	respondWithJSON(w, 200, resp)
}

func usersEditHandler(w http.ResponseWriter, r *http.Request) {
	decdRequest := decodeUserPost{}

	log.Println("\n\nRUNNING UserEditHandler\n\n")

	decodeJson(r.Body, &decdRequest)

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

	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Fatal(err)
	}

	dbStruct, err := db.LoadDB()
	if err != nil {
		log.Fatal(err)
	}

	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		log.Fatalf("Error converting subject claim to int -> %v\n", err)
	}

	if user, ok := dbStruct.Users[userId]; ok {
		user.Email = decdRequest.Email
		user.Password = decdRequest.Password
		dbStruct.Users[user.Id] = user
	} else {
		log.Fatal("User not found when trying to update fields\n")
	}

	fmt.Println(dbStruct.Users[userId])

	fmt.Printf(
		"\nPrinting claims...\n ISSUER -> %v\n SUBJCT -> %v\n",
		claims.Issuer,
		claims.Subject,
	)
}
