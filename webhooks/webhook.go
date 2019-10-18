package webhooks

import (
	"io/ioutil"
	"net/http"
)

type WebhookDetails struct {
	WebhookID string `json:"webhook_id"`
	SecretKey string `json:"secret_key"`
	Action    string `json:"action"`
}

type handler func(payload interface{}) error
type parser func(body []byte) (interface{}, error)

func webhookHandler(h handler, p parser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		payload, err := p(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h(payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
