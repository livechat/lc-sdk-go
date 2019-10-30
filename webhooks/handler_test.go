package webhooks_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"io/ioutil"

	"github.com/livechat/lc-sdk-go/webhooks"
	wv "github.com/livechat/lc-sdk-go/webhooks/tests/validators"
)

var followUpRequested = `{
	"webhook_id": "dummy_id",
	"secret_key": "dummy_key",
	"action": "follow_up_requested",
	"license_id": 100012582,
	"payload": {
		"chat_id": "XXXX",
		"thread_id": "YYYY",
		"customer_id": "AAA-BBB-CCC"
	}
}`

var verifiers = map[string]webhooks.Handler {
	"incoming_chat_thread": wv.IncomingChatThread,
	"thread_closed": wv.ThreadClosed,
	"access_set": wv.AccessSet,
	"chat_user_added": wv.ChatUserAdded,
	"chat_user_removed": wv.ChatUserRemoved,
	"incoming_event": wv.IncomingEvent,
	"event_updated": wv.EventUpdated,
	"incoming_rich_message_postback": wv.IncomingRichMessagePostback,
	"chat_properties_updated": wv.ChatPropertiesUpdated,
	"chat_properties_deleted": wv.ChatPropertiesDeleted,
	"chat_thread_properties_updated": wv.ChatThreadPropertiesUpdated,
	"chat_thread_properties_deleted": wv.ChatThreadPropertiesDeleted,
	"event_properties_updated": wv.EventPropertiesUpdated,
	"event_properties_deleted": wv.EventPropertiesDeleted,
	"follow_up_requested": wv.FollowUpRequested,
	"chat_thread_tagged": wv.ChatThreadTagged,
	"chat_thread_untagged": wv.ChatThreadUntagged,
	"agent_status_changed": wv.AgentStatusChanged,
	"agent_deleted": wv.AgentDeleted,
	"events_marked_as_seen": wv.EventsMarkedAsSeen,
}

func TestRejectWebhooksIfNoHandlersAreConnected(t *testing.T) {
	cfg := webhooks.NewConfiguration()
	h := webhooks.NewWebhookHandler(cfg)
	req := httptest.NewRequest("POST", "https://example.com", bytes.NewBufferString(followUpRequested))
	resp := httptest.NewRecorder()
	h(resp, req)
	if resp.Code != http.StatusBadRequest {
		t.Errorf("invalid code: %v", resp.Code)
		return
	}
}

func TestRejectWebhooksIfFormatIsInvalid(t *testing.T) {
	hook := followUpRequested + "}"
	cfg := webhooks.NewConfiguration()
	h := webhooks.NewWebhookHandler(cfg)
	req := httptest.NewRequest("POST", "https://example.com", bytes.NewBufferString(hook))
	resp := httptest.NewRecorder()
	h(resp, req)
	if resp.Code != http.StatusInternalServerError {
		t.Errorf("invalid code: %v", resp.Code)
		return
	}
}

func TestErrorHappensWithCustomErrorHandler(t *testing.T) {
	hook := followUpRequested + "}"
	errHandler := func(w http.ResponseWriter, err string, statusCode int) {
		if statusCode != http.StatusInternalServerError {
			t.Errorf("invalid status code in error handler: %v", statusCode)
		}
	}
	cfg := webhooks.NewConfiguration().WithErrorHandler(errHandler)
	h := webhooks.NewWebhookHandler(cfg)
	req := httptest.NewRequest("POST", "https://example.com", bytes.NewBufferString(hook))
	resp := httptest.NewRecorder()
	h(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("invalid code: %v", resp.Code)
		return
	}
}

func TestRejectWebhooksIfSecretKeyDoesntMatch(t *testing.T) {
	verifier := func(int, interface{}) error { return nil }
	cfg := webhooks.NewConfiguration().WithAction("follow_up_requested", verifier, "other_dummy_key")
	h := webhooks.NewWebhookHandler(cfg)
	req := httptest.NewRequest("POST", "https://example.com", bytes.NewBufferString(followUpRequested))
	resp := httptest.NewRecorder()
	h(resp, req)
	if resp.Code != http.StatusBadRequest {
		t.Errorf("invalid code: %v", resp.Code)
		return
	}
}

func TestPayloadParsingOK(t *testing.T) {
	testAction := func (action string, verifier webhooks.Handler) error{
		cfg := webhooks.NewConfiguration().WithAction(action, verifier, "dummy_key")
		h := webhooks.NewWebhookHandler(cfg)
		payload, err := ioutil.ReadFile("./tests/payloads/" + action + ".json")
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
