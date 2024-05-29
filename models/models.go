package models

type RequestLongURL struct {
	LongURL        string `json:"longURL"`
	ExpirationTime int    `json:"expirationTime"`
}
