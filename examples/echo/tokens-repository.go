package main

import (
	"sync"
	"time"
)

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
	mu      *sync.RWMutex
}

func NewTokenRepository() *TokensRepository {
	return &TokensRepository{make(map[string]*Token), &sync.RWMutex{}}
}

func (r *TokensRepository) Set(key string, t *Token) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage[key] = t
}

func (r *TokensRepository) Get(key string) *Token {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.storage[key]
}
