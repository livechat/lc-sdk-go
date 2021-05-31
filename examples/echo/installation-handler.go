package main

import (
	"fmt"
	"net/http"

	"github.com/livechat/lc-sdk-go/v4/authorization"
	"github.com/livechat/lc-sdk-go/v4/configuration"
)

type InstallationHandler struct {
	cfg *Configuration
	tr  tokensRepository
	as  accountsService
}

func NewInstallationHandler(cfg *Configuration, tr tokensRepository, as accountsService) *InstallationHandler {
	return &InstallationHandler{cfg, tr, as}
}

func (h *InstallationHandler) Handle(w http.ResponseWriter, r *http.Request) {
	code, exists := r.URL.Query()["code"]

	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	t, err := h.as.ExchangeCode(code[0])
	if err != nil {
		fmt.Println("Error when handling installation: code exchange failed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tg := func() *authorization.Token {
		return &authorization.Token{
			AccessToken: t.AccessToken,
			Region:      t.Region,
		}
	}

	api, err := configuration.NewAPI(tg, nil, t.ClientID)
	if err != nil {
		fmt.Println("Error when handling installation: configuration-api initilization failed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	wh := &configuration.Webhook{
		Action:      "incoming_event",
		SecretKey:   h.cfg.WebhookSecret,
		URL:         h.cfg.WebhookURL,
		Description: "echo integration",
		Filters:     &configuration.WebhookFilters{AuthorType: "customer"},
	}

	whID, err := api.RegisterWebhook(wh, nil)
	if err != nil {
		fmt.Println("Error when handling installation: webhook registration failed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.tr.Set(whID, t)

	w.WriteHeader(http.StatusOK)
}
