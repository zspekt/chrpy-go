package main

type decodeChirpPost struct {
	Body string `json:"body"`
}

type decodeUserPost struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type decodeUserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userPostResp struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type userLoginResp struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type tokenResp struct {
	Token string `json:"token"`
}
