package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func lengthValidationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type response struct {
		Valid bool `json:"valid"`
	}

	type decodeBody struct {
		Body string `json:"body"`
	}

	decdRequest := decodeBody{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&decdRequest)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(decdRequest.Body) > 140 {
		log.Println("Exceeds 140 characters.")
		w.WriteHeader(400)

		// :   bytesResp := []byte{"poop"}
		w.Write([]byte("\"error\": \"Something went wrong\""))
		return
	}

	resp := response{
		Valid: true,
	}

	toWrite, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(200)
	w.Write(toWrite)

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
