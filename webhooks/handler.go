package webhooks

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Configuration struct {
	Actions map[string]*ActionConfiguration
}

type ActionConfiguration struct {
	SecretKey string
	Handler   Handler
}

type Handler func(licenseID int, webhookPayload interface{}) error
type parser func(json.RawMessage) (interface{}, error)

func NewActionConfiguration(handler Handler) *ActionConfiguration {
	return &ActionConfiguration{
		Handler: handler,
	}
}

func (ac *ActionConfiguration) WithSecretKeyValidation(secretKey string) *ActionConfiguration {
	ac.SecretKey = secretKey
	return ac
}

func NewWebhookHandler(cfg *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var wh WebhookBase
		if err := json.Unmarshal(body, &wh); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		acfg, exists := cfg.Actions[wh.Action]
		if !exists {
			http.Error(w, "Unsupported action", http.StatusBadRequest)
			return
		}
		if acfg.SecretKey != "" && wh.SecretKey != acfg.SecretKey {
			http.Error(w, "Invalid webhook secret key", http.StatusBadRequest)
			return
		}

		var payload interface{}
		switch wh.Action {
		case "chat_user_removed":
			payload = &ChatUserRemoved{}
		case "follow_up_requested":
			payload = &FollowUpRequested{}
		case "incoming_chat_thread":
			payload = &IncomingChatThread{}
		case "incoming_event":
			payload = &IncomingEvent{}
		case "thread_closed":
			payload = &ThreadClosed{}
		}

		if err := json.Unmarshal(wh.Payload, payload); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = acfg.Handler(wh.LicenseID, payload); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
