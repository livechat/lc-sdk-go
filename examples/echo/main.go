package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/livechat/lc-sdk-go/webhooks"
)

func fillConfig(cfg *Configuration) {
	f, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	config, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(config, cfg)
}

func main() {
	cfg := &Configuration{}
	fillConfig(cfg)

	tr := NewTokenRepository()
	as := NewAccountsService(cfg)
	installationHandler := NewInstallationHandler(cfg, tr, as)
	incominEventHandler := NewIncomingEventHandler(cfg, tr)
	whConfig := webhooks.NewConfiguration().
		WithAction("incoming_event", incominEventHandler.Handle, cfg.WebhookSecret).
		WithErrorHandler(func(w http.ResponseWriter, err string, statusCode int) {
			fmt.Printf("Error when handling webhook: %v\n", err)
			w.WriteHeader(http.StatusOK)
		})

	http.HandleFunc("/oauth", installationHandler.Handle)
	http.HandleFunc("/webhook", webhooks.NewWebhookHandler(whConfig))
	http.ListenAndServe(":8080", nil)
}
