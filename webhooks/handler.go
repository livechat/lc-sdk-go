package webhooks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ErrorHandler func(w http.ResponseWriter, err string, statusCode int)

type Configuration struct {
	Actions     map[string]*ActionConfiguration
	handleError ErrorHandler
}

type ActionConfiguration struct {
	SecretKey string
	Handler   Handler
}

type Handler func(licenseID int, webhookPayload interface{}) error
type parser func(json.RawMessage) (interface{}, error)

func NewConfiguration() *Configuration {
	return &Configuration{
		Actions:     make(map[string]*ActionConfiguration),
		handleError: http.Error,
	}
}

func (cfg *Configuration) WithAction(action string, handler Handler, secretKey string) *Configuration {
	cfg.Actions[action] = &ActionConfiguration{
		Handler:   handler,
		SecretKey: secretKey,
	}
	return cfg
}

func (cfg *Configuration) WithErrorHandler(h func(w http.ResponseWriter, err string, statusCode int)) *Configuration {
	cfg.handleError = h
	return cfg
}

func NewWebhookHandler(cfg *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			cfg.handleError(w, fmt.Sprintf("couldn't read request body: %v", err), http.StatusInternalServerError)
			return
		}

		var wh WebhookBase
		if err := json.Unmarshal(body, &wh); err != nil {
			cfg.handleError(w, fmt.Sprintf("couldn't unmarshal webhook base: %v", err), http.StatusInternalServerError)
			return
		}
		acfg, exists := cfg.Actions[wh.Action]
		if !exists {
			cfg.handleError(w, fmt.Sprintf("Unsupported action: %v", wh.Action), http.StatusBadRequest)
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
			http.Error(w, fmt.Sprintf("couldn't unmarshal webhook payload: %v", err), http.StatusInternalServerError)
			return
		}

		if err = acfg.Handler(wh.LicenseID, payload); err != nil {
			http.Error(w, fmt.Sprintf("webhook handler error: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
