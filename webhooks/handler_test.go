package webhooks_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/livechat/lc-sdk-go/webhooks"
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
	checker := func(int, interface{}) error { return nil }
	cfg := webhooks.NewConfiguration().WithAction("follow_up_requested", checker, "other_dummy_key")
	h := webhooks.NewWebhookHandler(cfg)
	req := httptest.NewRequest("POST", "https://example.com", bytes.NewBufferString(followUpRequested))
	resp := httptest.NewRecorder()
	h(resp, req)
	if resp.Code != http.StatusBadRequest {
		t.Errorf("invalid code: %v", resp.Code)
		return
	}
}

func TestOK(t *testing.T) {
	checker := func(licenseID int, payload interface{}) error {
		if licenseID != 100012582 {
			return fmt.Errorf("Invalid licenseID: %v", licenseID)
		}
		wh, ok := payload.(*webhooks.FollowUpRequested)
		if !ok {
			return fmt.Errorf("invalid payload type: %T", payload)
		}
		if wh.ChatID != "XXXX" {
			return fmt.Errorf("invalid ChatID: %s", wh.ChatID)
		}
		if wh.ThreadID != "YYYY" {
			return fmt.Errorf("invalid ThreadID: %s", wh.ThreadID)
		}
		if wh.CustomerID != "AAA-BBB-CCC" {
			return fmt.Errorf("invalid CustomerID: %s", wh.CustomerID)
		}
		return nil
	}
	cfg := webhooks.NewConfiguration().WithAction("follow_up_requested", checker, "dummy_key")
	h := webhooks.NewWebhookHandler(cfg)
	req := httptest.NewRequest("POST", "https://example.com", bytes.NewBufferString(followUpRequested))
	resp := httptest.NewRecorder()
	h(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("invalid code: %v", resp.Code)
		return
	}
}
