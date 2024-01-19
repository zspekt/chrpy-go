package database

import (
	"encoding/json"
	"log"
	"os"
)

// Reads the file at fPath, assumes json content, and unmarshalls it to
// the provided memory address, expecting a struct of type T.
func UnmarshalToStruct[T any](structure *T, fPath string) error {
	fileBytes, err := os.ReadFile(fPath)
	if err != nil {
		log.Println(err)
		return err
	}

	if !json.Valid(fileBytes) {
		log.Println("\n", fPath, "does not contain valid JSON data.")
		return nil
	}

	// unmarshalling into the provided struct
	err = json.Unmarshal(fileBytes, structure)
	if err != nil {
		panic(err)
		return err
	}
	return nil
}

// Marshals the struct and writes it to the file.
func MarshalAndWrite[T any](structure T, fPath string) error {
	bytes, err := json.Marshal(structure)
	if err != nil {
		return err
	}

	err = os.WriteFile(fPath, bytes, 0700)
	if err != nil {
		return err
	}

	return nil
}
