package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type accountsService interface {
	ExchangeCode(string) (*Token, error)
}

type AccountsService struct {
	clientID         string
	clientSecret     string
	redirectURI      string
	tokenExchangeURL string
}

func NewAccountsService(cfg *Configuration) *AccountsService {
	return &AccountsService{cfg.ClientID, cfg.ClientSecret, cfg.RedirectURI, cfg.AccountsURL + "/v2/token"}
}

func (s *AccountsService) ExchangeCode(code string) (*Token, error) {
	reqBody, err := json.Marshal(map[string]string{
		"grant_type":    "authorization_code",
		"code":          code,
		"client_id":     s.clientID,
		"client_secret": s.clientSecret,
		"redirect_uri":  s.redirectURI,
	})

	if err != nil {
		return nil, err
	}

	resp, err := http.Post(s.tokenExchangeURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid_request")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	t := &Token{}
	if err := json.Unmarshal(body, t); err != nil {
		return nil, err
	}

	if t.ClientID == "" {
		t.ClientID = s.clientID
	}

	parts := strings.Split(t.AccessToken, ":")
	if len(parts) != 2 {
		return nil, errors.New(("invalid_region"))
	}
	t.Region = parts[0]
	t.ExpirationDate = time.Now().Add(time.Duration(t.ExpiresIn) * time.Second)

	return t, nil
}
