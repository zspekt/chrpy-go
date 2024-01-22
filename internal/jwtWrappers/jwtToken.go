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

	"github.com/zspekt/chrpy-go/internal/database"
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

func CreateAccessNdRefresh(cfg *JWTRequestConfig) (string, string, error) {
	var (
		accessToken  string
		refreshToken string
	)

	cfg.TokenType = "access"

	accessToken, err := CreateToken(cfg)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	cfg.TokenType = "refresh"

	refreshToken, err = CreateToken(cfg)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func CreateToken(cfg *JWTRequestConfig) (string, error) {
	var expiresInSeconds int
	var issuer string

	switch cfg.TokenType {
	case "access":
		issuer = "chirpy-access"
		expiresInSeconds = 3600
		log.Println("Token is type access and expires in 3600 seconds")
	case "refresh":
		issuer = "chirpy-refresh"
		expiresInSeconds = 5184000
		log.Println("Token is type refresh and expires in 5184000 seconds")
	}

	// jwt token claims
	var (
		userId    int       = cfg.UserID
		issuedAt  time.Time = time.Now().UTC()
		expiresAt time.Time = issuedAt.Add(time.Duration(expiresInSeconds) * time.Second)
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

// validates token and checks if it has been revoked. also returns token type.
// // no, it doesn't return token type. but one can easily get that from the claims.
// what was i on?
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
		log.Println(token)
		// log.Fatal(err)
		return jwt.RegisteredClaims{}, err
	}

	if !jwtToken.Valid {
		log.Fatalf("Token is invalid.\n")
		return jwt.RegisteredClaims{}, fmt.Errorf("Token is invalid")
	}

	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Println(err)
		return jwt.RegisteredClaims{}, err
	}

	isRevoked, err := db.IsRevoked(token)
	if err != nil {
		return jwt.RegisteredClaims{}, err
	}
	if isRevoked {
		log.Println("Token is revoked...")
		return jwt.RegisteredClaims{}, fmt.Errorf("Token is revoked")
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
