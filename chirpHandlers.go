package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/zspekt/chrpy-go/internal/database"
)

// replaces provided curseWords with provided censored parameter of type string.
// Returns a boolean indicating if there were any curseWords present
func profaneCheck(str *string, curseWords []string, censored string) bool {
	slice := strings.Split(*str, " ")
	var cursedWordPresent bool

	curseWordsMap := make(map[string]bool, len(curseWords))
	for _, curse := range curseWords {
		curseWordsMap[strings.ToLower(curse)] = true
	}

	for i, v := range slice {
		if curseWordsMap[strings.ToLower(v)] {
			slice[i] = censored
			cursedWordPresent = true
		}
	}

	*str = strings.Join(slice, " ")
	return cursedWordPresent
}

func chirpsPostHandler(w http.ResponseWriter, r *http.Request) {
	decdRequest := decodeChirpPost{}

	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Println(err)
		return
	}

	err = decodeJson[decodeChirpPost](r.Body, &decdRequest)
	if err != nil {
		log.Fatal(err)
		respondWithError(w, 500, "\nServer error --> Error decoding body into JSON\n")
	}

	if len(decdRequest.Body) > 140 {
		log.Println("Exceeds 140 characters.")
		respondWithError(w, 400, "\"error\": \"Exceeds 140 characters.\"")
		return
	}
	curseWords := []string{"kerfuffle", "sharbert", "fornax"}
	profaneCheck(&decdRequest.Body, curseWords, "****")

	chirp, err := db.CreateChirp(decdRequest.Body)
	if err != nil {
		log.Println(err)
	}
	respondWithJSON(w, 201, chirp)
}

func chirpsGetHandler(w http.ResponseWriter, r *http.Request) {
	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Println(err)
		return
	}

	chirps, err := db.GetChirps()
	if err != nil {
		return
	}

	respondWithJSON(w, 200, chirps)
}

func chirpsGetByIDHandler(w http.ResponseWriter, r *http.Request) {
	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Println(err)
		return
	}

	chirpsMap, err := db.LoadDB()
	if err != nil {
		log.Println(err)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "*"))
	if err != nil {
		log.Println(err)
		return
	}

	if chirp, ok := chirpsMap.Chirps[id]; ok {
		respondWithJSON(w, 200, chirp)
		return
	}

	respondWithError(w, 404, "404 Not found")
}
