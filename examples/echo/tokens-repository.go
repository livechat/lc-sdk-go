package main

import "time"

type Token struct {
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token"`
	AccountID      string `json:"account_id"`
	OrganizationID string `json:"organization_id"`
	ClientID       string `json:"client_id"`
	ExpiresIn      int    `json:"expires_in"`
	ExpirationDate time.Time
	Region         string
}

type tokensRepository interface {
	Set(string, *Token)
	Get(string) *Token
}

type TokensRepository struct {
	storage map[string]*Token
}

func NewTokenRepository() *TokensRepository {
	return &TokensRepository{make(map[string]*Token)}
}

func (r *TokensRepository) Set(key string, t *Token) {
	r.storage[key] = t
}

func (r *TokensRepository) Get(key string) *Token {
	return r.storage[key]
}
