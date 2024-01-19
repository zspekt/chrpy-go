package jwtwrappers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	if cfg.ExpiresInSeconds <= 0 {
		cfg.ExpiresInSeconds = 86400
	}

	var (
		expires_in_seconds int       = cfg.ExpiresInSeconds
		userId             int       = cfg.UserID
		issuer             string    = "chirpy"
		issuedAt           time.Time = time.Now().UTC()
		expiresAt          time.Time = issuedAt.Add(time.Duration(expires_in_seconds) * time.Second)
	)

	fmt.Printf("\nissuedAt -> %v\nexpiresAt -> %v", issuedAt, expiresAt)

	claims := jwt.RegisteredClaims{
		Issuer:    issuer,
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		Subject:   strconv.Itoa(userId),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Println("jwtToken LINE 48")
		return "", err
	}
	// fmt.Printf("\n\nCREATED TOKEN: %v\n\n", signedToken)
	return signedToken, nil
}

func ValidateAndReturn(token string) (jwt.RegisteredClaims, error) {
	// fmt.Println("jwtSecret right here -> ", jwtSecret)
	claims := &jwt.RegisteredClaims{}

	jwtToken, err := jwt.ParseWithClaims(
		token,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		},
	)
	if err != nil {
		log.Println("jwtToken LINE 67")
		// log.Fatal(err)
		return jwt.RegisteredClaims{}, err
	}

	if !jwtToken.Valid {
		log.Fatalf("Token is not valid.\n")
		// return jwt.RegisteredClaims{}, fmt.Errorf("Token is invalid")
	}

	return *claims, nil
}

func GetTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Authorization header is missing\n")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("Invalid Auth header format.\n")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	// fmt.Printf("\n\nRETRIEVED TOKEN: %v\n\n", token)
	return token, nil
}
