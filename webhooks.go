package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/zspekt/chrpy-go/internal/database"
)

var polkaApiKey string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("\nError loading .env --> %v\n", err)
		return
	}

	polkaApiKey = os.Getenv("polkaApiKey")
	log.Println("polkaApiKey has been set...")
}

func isAuthdRequest(r *http.Request, keyToMatch string) (bool, error) {
	keyFromRequest, err := GetFromHeader(r, "Authorization", "ApiKey")
	if err != nil {
		return false, err
	}

	if keyFromRequest != keyToMatch {
		log.Println("Invalid ApiKey being passed to webhook...")
		return false, nil
	}

	return true, nil
}

func polkaPostHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("RUNNING polkaPostHandler...")

	boolean, err := isAuthdRequest(r, polkaApiKey)
	if err != nil {
		log.Println(err)
		respondWithError(w, 401, "Unauthorized access")
		return
	}

	if !boolean {
		log.Println("Invalid ApiKey...")
		respondWithError(w, 401, "Unauthorized access")
		return
	}

	decdRequest := webhookRequest{}

	err = decodeJson(r.Body, &decdRequest)
	if err != nil {
		log.Println(err)
		respondWithError(w, 500, "Server error")
		return
	}

	if decdRequest.Event != "user.upgraded" {
		w.WriteHeader(200)
		log.Println("Incorrect event...")
		return
	}

	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Println(err)
		respondWithError(w, 500, "Server error")
		return
	}

	user, err := db.UpgradeUser(decdRequest.Data["user_id"])
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("\nUpgraded user %v to chirpy red: %v\n", user.Email, user.IsChirpyRed)

	w.WriteHeader(200)
}
