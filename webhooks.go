package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zspekt/chrpy-go/internal/database"
)

func polkaPostHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("RUNNING polkaPostHandler...")

	decdRequest := webhookRequest{}

	err := decodeJson(r.Body, &decdRequest)
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
