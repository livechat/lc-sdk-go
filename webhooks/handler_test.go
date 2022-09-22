package webhooks_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/livechat/lc-sdk-go/v5/webhooks"
)

var verifiers = map[string]webhooks.Handler{
	"incoming_chat":                   incomingChat,
	"incoming_event":                  incomingEvent,
	"event_updated":                   eventUpdated,
	"incoming_rich_message_postback":  incomingRichMessagePostback,
	"chat_deactivated":                chatDeactivated,
	"chat_properties_updated":         chatPropertiesUpdated,
	"thread_properties_updated":       threadPropertiesUpdated,
	"chat_properties_deleted":         chatPropertiesDeleted,
	"thread_properties_deleted":       threadPropertiesDeleted,
	"user_added_to_chat":              userAddedToChat,
	"user_removed_from_chat":          userRemovedFromChat,
	"thread_tagged":                   threadTagged,
	"thread_untagged":                 threadUntagged,
	"agent_created":                   agentCreated,
	"agent_updated":                   agentUpdated,
	"agent_deleted":                   agentDeleted,
	"agent_suspended":                 agentSuspended,
	"agent_unsuspended":               agentUnsuspended,
	"agent_approved":                  agentApproved,
	"events_marked_as_seen":           eventsMarkedAsSeen,
	"chat_access_updated":             chatAccessUpdated,
	"event_properties_updated":        eventPropertiesUpdated,
	"event_properties_deleted":        eventPropertiesDeleted,
	"routing_status_set":              routingStatusSet,
	"chat_transferred":                chatTransferred,
	"incoming_customer":               incomingCustomer,
	"customer_session_fields_updated": customerSessionFieldsUpdated,
	"group_created":                   groupCreated,
	"group_updated":                   groupUpdated,
	"group_deleted":                   groupDeleted,
	"auto_access_added":               autoAccessAdded,
	"auto_access_updated":             autoAccessUpdated,
	"auto_access_deleted":             autoAccessDeleted,
	"bot_created":                     botCreated,
	"bot_updated":                     botUpdated,
	"bot_deleted":                     botDeleted,
}

func TestRejectWebhooksIfNoHandlersAreConnected(t *testing.T) {
	cfg := webhooks.NewConfiguration()
	h := webhooks.NewWebhookHandler(cfg)
	action := "incoming_chat"
	payload, err := os.ReadFile("./testdata/" + action + ".json")
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
	action := "incoming_chat"
	payload, err := os.ReadFile("./testdata/" + action + ".json")
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
	action := "incoming_chat"
	payload, err := os.ReadFile("./testdata/" + action + ".json")
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
	verifier := func(context.Context, *webhooks.Webhook) error { return nil }
	action := "incoming_chat"
	cfg := webhooks.NewConfiguration().WithAction(action, verifier, "other_dummy_key")
	h := webhooks.NewWebhookHandler(cfg)
	payload, err := os.ReadFile("./testdata/" + action + ".json")
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
		return func(ctx context.Context, wh *webhooks.Webhook) error {
			var propEqErrors string
			propEq("OrganizationID", wh.OrganizationID, "390e44e6-f1e6-0368c-z6ddb-74g14508c2ex", &propEqErrors)
			if propEqErrors != "" {
				return errors.New(propEqErrors)
			}
			return verifier(ctx, wh)
		}
	}
	testAction := func(action string, verifier webhooks.Handler) error {
		cfg := webhooks.NewConfiguration().WithAction(action, withLicenseCheck(verifier), "dummy_key")
		h := webhooks.NewWebhookHandler(cfg)
		payload, err := os.ReadFile("./testdata/" + action + ".json")
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
		t.Run(action, func(t *testing.T) {
			stepError := testAction(action, verifier)
			if stepError != nil {
				t.Errorf("Payload incorrectly parsed for %v, error: %v", action, stepError)
			}
		})
	}
}

func TestHandlerContextForwardsRequestContext(t *testing.T) {
	verifier := func(ctx context.Context, wh *webhooks.Webhook) error {
		rawVal := ctx.Value("dummy-key")
		val, ok := rawVal.(string)
		if !ok {
			t.Errorf("invalid type of 'dummy-key' in wh ctx: %T", rawVal)
			return nil
		}
		if val != "dummy-value" {
			t.Errorf("invalid value of 'dummy-key' in wh ctx: %v", val)
			return nil
		}
		return nil
	}
	action := "incoming_chat"
	cfg := webhooks.NewConfiguration().WithAction(action, verifier, "")
	h := webhooks.NewWebhookHandler(cfg)
	payload, err := os.ReadFile("./testdata/" + action + ".json")
	if err != nil {
		t.Errorf("Missing test payload for action %v", action)
		return
	}
	req := httptest.NewRequest("POST", "https://example.com", bytes.NewBuffer(payload))
	req = req.WithContext(context.WithValue(context.Background(), "dummy-key", "dummy-value"))
	resp := httptest.NewRecorder()
	h(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("invalid code: %v", resp.Code)
		return
	}
}
