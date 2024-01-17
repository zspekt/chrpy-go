package jwtwrappers

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var jwtSecret string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("\nError loading .env --> %v\n", err)
		return
	}

	jwtSecret = os.Getenv("jwtSecret")
	log.Println("jwtSecret has been set...")
}

func CreateToken(cfg *JWTRequestConfig) (string, error) {
	// placeholders
	var (
		expires_in_seconds int    = cfg.ExpiresInSeconds
		userId             int    = cfg.UserID
		issuer             string = "chirpy"
	)

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(time.Duration(expires_in_seconds))

	claims := jwt.RegisteredClaims{
		Issuer:    issuer,
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		Subject:   strconv.Itoa(userId),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return signedToken, nil
}
