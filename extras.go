package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

func profaneCheck(str *string, curseWords []string, censored string) bool {
	// censored := "****"
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

func chirpsHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		CleanedBody string `json:"cleaned_body"`
	}
	type decodeBody struct {
		Body string `json:"body"`
	}

	decdRequest := decodeBody{}

	err := decodeJson[decodeBody](r.Body, &decdRequest)
	if err != nil {
		log.Println(err)
		respondWithError(w, 500, "\nServer error --> Error decoding parameters\n")
	}

	if len(decdRequest.Body) > 140 {
		log.Println("Exceeds 140 characters.")
		respondWithError(w, 400, "\"error\": \"Exceeds 140 characters.\"")
		return
	}

	curseWords := []string{"kerfuffle", "sharbert", "fornax"}
	cursePresent := profaneCheck(&decdRequest.Body, curseWords, "****")

  chirp := Chirp{
    Body: decdRequest.Body,
    id: 
  }

	resp := response{
		CleanedBody: decdRequest.Body,
	}

	respondWithJSON(w, 201, resp)

	log.Printf("bad words present: %v\n", cursePresent)
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
