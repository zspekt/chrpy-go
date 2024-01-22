package main

import (
	"log"
	"strconv"

	jwtwrappers "github.com/zspekt/chrpy-go/internal/jwtWrappers"
)

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
