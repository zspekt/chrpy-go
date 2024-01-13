package main

type decodeChirpPost struct {
	Body string `json:"body"`
}

type decodeUserPost struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userPostResp struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}
