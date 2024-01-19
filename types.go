package main

type decodeChirpPost struct {
	Body string `json:"body"`
}

type decodeUserPost struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type decodeUserLogin struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type userPostResp struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type userLoginResp struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}
