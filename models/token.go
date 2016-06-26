package models

import (
	"time"

	"golang.org/x/oauth2"
)

type Token struct {
	Id           int
	AccessToken  string
	TokenType    string
	RefreshToken string
	Expiry       time.Time
}

func (u Token) Token() *oauth2.Token {
	return &oauth2.Token{
		AccessToken:  u.AccessToken,
		TokenType:    u.TokenType,
		RefreshToken: u.RefreshToken,
		Expiry:       u.Expiry,
	}
}
