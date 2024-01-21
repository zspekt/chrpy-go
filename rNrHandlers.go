package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/zspekt/chrpy-go/internal/database"
	jwtwrappers "github.com/zspekt/chrpy-go/internal/jwtWrappers"
)

func refreshPostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n\n\n")
	log.Println("RUNNING refreshPostHandler\n\n")

	resp := tokenResp{}

	token, err := jwtwrappers.GetTokenFromHeader(r)
	if err != nil {
		log.Println(err)
		respondWithError(w, 401, "Unauthorized access")
		return
	}

	fmt.Printf("\nretrieved token: %v\n", token)

	claims, err := jwtwrappers.ValidateAndReturn(token)
	if err != nil {
		fmt.Println("\n\n\n\n\nPUTOOOOOO\n\n\n\n\n\n")
		log.Printf("Error on ValidateAndReturn -> %v\n", err)
		fmt.Printf("Here is the token the error originates from < %v >\n", token)
		respondWithError(w, 401, "Unauthorized access")
		// log.Println()
		return
	}
	if claims.Issuer != "chirpy-refresh" {
		log.Println("Access token being passed as refresh...")
		respondWithError(w, 401, "Unauthorized access")
		return
	}
	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		log.Printf("Error converting subject claim to int -> %v\n", err)
		log.Println(err)
		respondWithError(w, 401, "Unauthorized access")
	}

	cfg := jwtwrappers.JWTRequestConfig{
		UserID:    userId,
		TokenType: "access",
	}

	newAccessToken, err := jwtwrappers.CreateToken(&cfg)
	if err != nil {
		log.Println(err)
		return
	}

	resp = tokenResp{
		Token: newAccessToken,
	}

	respondWithJSON(w, 200, resp)

	bytes, _ := json.Marshal(resp)

	fmt.Printf("\n--> <--\n\n%v\n", string(bytes))
}

func revokePostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n\n\n")
	log.Println("RUNNING revokePostHandler\n\n")

	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Println(err)
		return
	}

	token, err := jwtwrappers.GetTokenFromHeader(r)
	if err != nil {
		log.Println(err)
		respondWithError(w, 401, "Unauthorized access")
		return
	}

	fmt.Printf("\nretrieved token: %v\n", token)

	claims, err := jwtwrappers.ValidateAndReturn(token)
	if err != nil {
		log.Println(err)
		respondWithError(w, 401, "Unauthorized access")
		return
	}
	if claims.Issuer != "chirpy-refresh" {
		log.Println("Access token being passed as refresh...")
		respondWithError(w, 401, "Unauthorized access")
		return
	}

	db.RevokeToken(token)
}
