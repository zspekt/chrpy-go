package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	jwtwrappers "github.com/zspekt/chrpy-go/internal/jwtWrappers"
)

func GetFromHeader(r *http.Request, httpHeader string, prefix string) (string, error) {
	header := r.Header.Get(httpHeader)
	if header == "" {
		return "", errors.New("Header is missing")
	}

	prefix += " "

	if !strings.HasPrefix(header, prefix) {
		return "", fmt.Errorf("Invalid Auth header format.\n")
	}

	token := strings.TrimPrefix(header, prefix)

	// fmt.Printf("\n\nRETRIEVED TOKEN: %v\n\n", token)
	return token, nil
}

// returns the ID number associated with the user that's sent the token
// and err if the user doesn't have access
func validateAccess(token string) (int, error) {
	claims, err := jwtwrappers.ValidateAndReturn(token)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	if claims.Issuer != "chirpy-access" {
		log.Println(err)
		return 0, err
	}
	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return userId, nil
}
