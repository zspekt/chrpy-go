package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func lengthValidationHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Valid bool `json:"valid"`
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

	resp := response{
		Valid: true,
	}

	respondWithJSON(w, 200, resp)

	log.Println("Chirp does not exceed 140 characters. Well done.")
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
