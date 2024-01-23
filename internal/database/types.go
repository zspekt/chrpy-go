package database

import (
	"sync"
	"time"
)

type DB struct {
	path  string
	mutex *sync.RWMutex
}

type DBStructure struct {
	Chirps        map[int]Chirp        `json:"chirps"`
	Users         map[int]User         `json:"users"`
	RevokedTokens map[string]time.Time `json:"revoked_tokens"`
}

type User struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type Chirp struct {
	Body     string `json:"body"`
	ChirpId  int    `json:"id"`
	AuthorId int    `json:"author_id"`
}
