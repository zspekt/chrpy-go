package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)

	_, err := w.Write([]byte(msg))
	if err != nil {
		log.Printf("\nWrite error --> %v\n", err)
		return
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(jsonPayload))
	w.Write(jsonPayload)
}

func decodeJson[T any](r io.ReadCloser, sillycat *T) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(sillycat)
	if err != nil {
		return err
	}
	return nil
}
