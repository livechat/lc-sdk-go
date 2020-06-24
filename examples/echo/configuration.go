package main

type Configuration struct {
	ClientID      string `json:"client_id"`
	ClientSecret  string `json:"client_secret"`
	RedirectURI   string `json:"redirect_uri"`
	AccountsURL   string `json:"accounts_url"`
	WebhookURL    string `json:"webhook_url"`
	WebhookSecret string `json:"webhook_secret"`
}
