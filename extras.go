package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/zspekt/chrpy-go/internal/database"
)

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
	type response struct {
		CleanedBody string `json:"cleaned_body"`
	}
	type decodeBody struct {
		Body string `json:"body"`
	}

	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Println(err)
		return
	}

	decdRequest := decodeBody{}

	// fmt.Println("\n\t\tRIGHT BEFORE DECODE JSON\n")
	err = decodeJson[decodeBody](r.Body, &decdRequest)
	// fmt.Println("\n\t\tAFFTEEEEER DECODE JSON\n")
	if err != nil {
		log.Fatal(err)
		respondWithError(w, 500, "\nServer error --> Error decoding parameters\n")
	}

	if len(decdRequest.Body) > 140 {
		log.Println("Exceeds 140 characters.")
		respondWithError(w, 400, "\"error\": \"Exceeds 140 characters.\"")
		return
	}

	curseWords := []string{"kerfuffle", "sharbert", "fornax"}
	profaneCheck(&decdRequest.Body, curseWords, "****")

	// id, err := db.GetIdCount()
	// if err != nil {
	//   log.Println(err)
	//   return
	// }

	// fmt.Println("\n\t\tRIGHT BEFORE CREATECHIRP CALL\n\n")
	chirp, err := db.CreateChirp(decdRequest.Body)
	if err != nil {
		log.Println(err)
	}
	// fmt.Println("\n\t\tAFTERRRR CREATECHIRP CALL\n\n")

	fmt.Println(chirp)

	respondWithJSON(w, 201, chirp)

	// log.Printf("bad words present: %v\n", cursePresent)
}

func chirpsGetHandler(w http.ResponseWriter, r *http.Request) {
	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Println(err)
		return
	}

	// DBStruct, err := db.LoadDB()
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	chirps, err := db.GetChirps()
	if err != nil {
		return
	}

	respondWithJSON(w, 200, chirps)
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		valid bool
	}

	resp := response{
		valid: true,
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	var body string
	body = string(bytes)

	if len(body) > 140 {
		w.WriteHeader(400)
		return
	}

	toWrite, err := json.Marshal(resp)
	w.Write(toWrite)
}

/*
















 */

// func lengthValidationHandler(w http.ResponseWriter, r *http.Request) {
// 	type response struct {
// 		Valid bool `json:"valid"`
// 	}
// 	type decodeBody struct {
// 		Body string `json:"body"`
// 	}
//
// 	decdRequest := decodeBody{}
//
// 	err := decodeJson[decodeBody](r.Body, &decdRequest)
// 	if err != nil {
// 		log.Println(err)
// 		respondWithError(w, 500, "\nServer error --> Error decoding parameters\n")
// 	}
//
// 	if len(decdRequest.Body) > 140 {
// 		log.Println("Exceeds 140 characters.")
// 		respondWithError(w, 400, "\"error\": \"Exceeds 140 characters.\"")
// 		return
// 	}
//
// 	resp := response{
// 		Valid: true,
// 	}
//
// 	respondWithJSON(w, 200, resp)
//
// 	log.Println("Chirp does not exceed 140 characters. Well done.")
// }
