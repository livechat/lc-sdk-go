package webhooks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// The ErrorHandler type is used to define custom error handlers for WebhookHandler.
//
// It allows to customize behaviour of WebhookHandler when webhook processing errors,
// eg. to always respond with 200OK.
type ErrorHandler func(w http.ResponseWriter, err string, statusCode int)

// A Configuration structure is used to configure WebhookHandler
type Configuration struct {
	actions     map[string]*actionConfiguration
	handleError ErrorHandler
}

type actionConfiguration struct {
	secretKey string
	handle    Handler
}

// The Handler type is used to define webhook processors.
//
// It can be used with WebhookHandler, in which case WebhookHandler will
// pass decoded webhook payload (ie. one of webhooks structures).
type Handler func(licenseID int, webhookPayload interface{}) error

// NewConfiguration creates basic WebhookHandler configuration that
// processes no webhooks and uses http.Error to handle webhook processing
// errors.
func NewConfiguration() *Configuration {
	return &Configuration{
		actions:     make(map[string]*actionConfiguration),
		handleError: http.Error,
	}
}

// WithAction allows to attach custom webhook Handler for given webhook action.
//
// If secretKey is an empty string, then no validation of webhook's secret is performed.
// Otherwise, webhook's secret is strictly validated. In case of any mismatch between expected and actual secret key,
// webhook processing is stopped and error is returned.
func (cfg *Configuration) WithAction(action string, handler Handler, secretKey string) *Configuration {
	cfg.actions[action] = &actionConfiguration{
		handle:    handler,
		secretKey: secretKey,
	}
	return cfg
}

// WithErrorHandler allows to attach custom ErrorHandler, which acts as sink for all WebhookHandler errors.
//
// Custom ErrorHandler might be used to eg. always return 200OK for incoming webhooks.
func (cfg *Configuration) WithErrorHandler(h ErrorHandler) *Configuration {
	cfg.handleError = h
	return cfg
}

// NewWebhookHandler creates WebhookHandler that can be used with golang HTTP server.
//
// WebhookHandler decodes raw webhook JSON into dedicated webhook structures and, if provided, passes
// those structures into webhook Handlers attached to given webhook type.
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
		acfg, exists := cfg.actions[wh.Action]
		if !exists {
			cfg.handleError(w, fmt.Sprintf("Unsupported action: %v", wh.Action), http.StatusBadRequest)
			return
		}
		if acfg.secretKey != "" && wh.SecretKey != acfg.secretKey {
			cfg.handleError(w, "Invalid webhook secret key", http.StatusBadRequest)
			return
		}

		var payload interface{}
		switch wh.Action {
		case "incoming_chat_thread":
			payload = &IncomingChatThread{}
		case "thread_closed":
			payload = &ThreadClosed{}
		case "access_granted":
			payload = &AccessGranted{}
		case "access_revoked":
			payload = &AccessRevoked{}
		case "access_set":
			payload = &AccessSet{}
		case "chat_user_added":
			payload = &ChatUserAdded{}
		case "chat_user_removed":
			payload = &ChatUserRemoved{}
		case "incoming_event":
			payload = &IncomingEvent{}
		case "event_updated":
			payload = &EventUpdated{}
		case "incoming_rich_message_postback":
			payload = &IncomingRichMessagePostback{}
		case "chat_properties_updated":
			payload = &ChatPropertiesUpdated{}
		case "chat_properties_deleted":
			payload = &ChatPropertiesDeleted{}
		case "chat_thread_properties_deleted":
			payload = &ChatThreadPropertiesDeleted{}
		case "chat_thread_properties_updated":
			payload = &ChatThreadPropertiesUpdated{}
		case "event_properties_updated":
			payload = &EventPropertiesUpdated{}
		case "event_properties_deleted":
			payload = &EventPropertiesDeleted{}
		case "chat_thread_tagged":
			payload = &ChatThreadTagged{}
		case "chat_thread_untagged":
			payload = &ChatThreadUntagged{}
		case "agent_status_changed":
			payload = &AgentStatusChanged{}
		case "agent_deleted":
			payload = &AgentDeleted{}
		case "customer_created":
			payload = &CustomerCreated{}
		case "events_marked_as_seen":
			payload = &EventsMarkedAsSeen{}
		case "follow_up_requested":
			payload = &FollowUpRequested{}
		default:
			cfg.handleError(w, fmt.Sprintf("unknown webhook: %v", wh.Action), http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(wh.Payload, payload); err != nil {
			cfg.handleError(w, fmt.Sprintf("couldn't unmarshal webhook payload: %v", err), http.StatusInternalServerError)
			return
		}

		if err = acfg.handle(wh.LicenseID, payload); err != nil {
			cfg.handleError(w, fmt.Sprintf("webhook handler error: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
