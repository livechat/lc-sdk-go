package webhooks_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/livechat/lc-sdk-go/webhooks"
)

var verifiers = map[string]webhooks.Handler{
	"incoming_chat_thread":           incomingChatThread,
	"thread_closed":                  threadClosed,
	"chat_deactivated":               chatDeactivated,
	"access_granted":                 accessGranted,
	"access_revoked":                 accessRevoked,
	"access_set":                     accessSet,
	"chat_user_added":                chatUserAdded,
	"chat_user_removed":              chatUserRemoved,
	"incoming_event":                 incomingEvent,
	"event_updated":                  eventUpdated,
	"incoming_rich_message_postback": incomingRichMessagePostback,
	"chat_properties_updated":        chatPropertiesUpdated,
	"chat_properties_deleted":        chatPropertiesDeleted,
	"chat_thread_properties_updated": chatThreadPropertiesUpdated,
	"chat_thread_properties_deleted": chatThreadPropertiesDeleted,
	"event_properties_updated":       eventPropertiesUpdated,
	"event_properties_deleted":       eventPropertiesDeleted,
	"chat_thread_tagged":             chatThreadTagged,
	"chat_thread_untagged":           chatThreadUntagged,
	"agent_status_changed":           agentStatusChanged,
	"agent_deleted":                  agentDeleted,
	"customer_created":               customerCreated,
	"events_marked_as_seen":          eventsMarkedAsSeen,
	"follow_up_requested":            followUpRequested,
}

func TestRejectWebhooksIfNoHandlersAreConnected(t *testing.T) {
	cfg := webhooks.NewConfiguration()
	h := webhooks.NewWebhookHandler(cfg)
	action := "incoming_chat_thread"
	payload, err := ioutil.ReadFile("./testdata/" + action + ".json")
	if err != nil {
		t.Errorf("Missing test payload for action %v", action)
		return
	}
	req := httptest.NewRequest("POST", "https://example.com", bytes.NewBuffer(payload))
	resp := httptest.NewRecorder()
	h(resp, req)
	if resp.Code != http.StatusBadRequest {
		t.Errorf("invalid code: %v", resp.Code)
		return
	}
}

func TestRejectWebhooksIfFormatIsInvalid(t *testing.T) {
	action := "incoming_chat_thread"
	payload, err := ioutil.ReadFile("./testdata/" + action + ".json")
	if err != nil {
		t.Errorf("Missing test payload for action %v", action)
		return
	}
	payload = append(payload, '}')
	cfg := webhooks.NewConfiguration()
	h := webhooks.NewWebhookHandler(cfg)
	req := httptest.NewRequest("POST", "https://example.com", bytes.NewBuffer(payload))
	resp := httptest.NewRecorder()
	h(resp, req)
	if resp.Code != http.StatusInternalServerError {
		t.Errorf("invalid code: %v", resp.Code)
		return
	}
}

func TestErrorHappensWithCustomErrorHandler(t *testing.T) {
	action := "incoming_chat_thread"
	payload, err := ioutil.ReadFile("./testdata/" + action + ".json")
	if err != nil {
		t.Errorf("Missing test payload for action %v", action)
		return
	}
	payload = append(payload, '}')
	errHandler := func(w http.ResponseWriter, err string, statusCode int) {
		if statusCode != http.StatusInternalServerError {
			t.Errorf("invalid status code in error handler: %v", statusCode)
		}
	}
	cfg := webhooks.NewConfiguration().WithErrorHandler(errHandler)
	h := webhooks.NewWebhookHandler(cfg)
	req := httptest.NewRequest("POST", "https://example.com", bytes.NewBuffer(payload))
	resp := httptest.NewRecorder()
	h(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("invalid code: %v", resp.Code)
		return
	}
}

func TestRejectWebhooksIfSecretKeyDoesntMatch(t *testing.T) {
	verifier := func(int, interface{}) error { return nil }
	action := "incoming_chat_thread"
	cfg := webhooks.NewConfiguration().WithAction(action, verifier, "other_dummy_key")
	h := webhooks.NewWebhookHandler(cfg)
	payload, err := ioutil.ReadFile("./testdata/" + action + ".json")
	if err != nil {
		t.Errorf("Missing test payload for action %v", action)
		return
	}
	req := httptest.NewRequest("POST", "https://example.com", bytes.NewBuffer(payload))
	resp := httptest.NewRecorder()
	h(resp, req)
	if resp.Code != http.StatusBadRequest {
		t.Errorf("invalid code: %v", resp.Code)
		return
	}
}

func TestPayloadParsingOK(t *testing.T) {
	withLicenseCheck := func(verifier webhooks.Handler) webhooks.Handler {
		return func(licenseID int, payload interface{}) error {
			var errors string
			propEq("LicenseID", licenseID, 21377312, &errors)
			if errors != "" {
				return fmt.Errorf(errors)
			}
			return verifier(licenseID, payload)
		}
	}
	testAction := func(action string, verifier webhooks.Handler) error {
		cfg := webhooks.NewConfiguration().WithAction(action, withLicenseCheck(verifier), "dummy_key")
		h := webhooks.NewWebhookHandler(cfg)
		payload, err := ioutil.ReadFile("./testdata/" + action + ".json")
		if err != nil {
			return fmt.Errorf("Missing test payload for action %v", action)
		}
		req := httptest.NewRequest("POST", "https://example.com", bytes.NewBuffer(payload))
		resp := httptest.NewRecorder()
		h(resp, req)
		if resp.Code != http.StatusOK {
			return fmt.Errorf("%v", resp.Body)
		}
		return nil
	}

	for action, verifier := range verifiers {
		stepError := testAction(action, verifier)
		if stepError != nil {
			t.Errorf("Payload incorrectly parsed for %v, error: %v", action, stepError)
			return
		}
	}
}
