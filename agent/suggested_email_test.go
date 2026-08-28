package agent_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/livechat/lc-sdk-go/v6/agent"
)

func captureUpdateCustomerBody(t *testing.T, call func(*agent.API) error) map[string]interface{} {
	t.Helper()
	var body map[string]interface{}
	client := NewTestClient(func(req *http.Request) *http.Response {
		raw, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("reading request body: %v", err)
		}
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Fatalf("decoding request body: %v", err)
		}
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString("{}")),
			Header:     make(http.Header),
		}
	})

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Fatal("API creation failed")
	}
	if err := call(api); err != nil {
		t.Fatalf("UpdateCustomer failed: %v", err)
	}
	return body
}

func TestUpdateCustomerSendsSuggestedEmail(t *testing.T) {
	body := captureUpdateCustomerBody(t, func(api *agent.API) error {
		return api.UpdateCustomer("customer_id", "", "", "", "", "guessed@example.com", nil, nil)
	})

	got, ok := body["suggested_email"]
	if !ok {
		t.Fatalf("suggested_email missing from the request body: %v", body)
	}
	if got != "guessed@example.com" {
		t.Errorf("suggested_email = %v, want guessed@example.com", got)
	}
}

// This version cannot clear a suggestion: the field is a plain string with omitempty, so an empty
// value drops out of the payload entirely. Dismissal is a 3.7 capability.
func TestUpdateCustomerOmitsEmptySuggestedEmail(t *testing.T) {
	body := captureUpdateCustomerBody(t, func(api *agent.API) error {
		return api.UpdateCustomer("customer_id", "Thomas", "", "", "", "", nil, nil)
	})

	if _, ok := body["suggested_email"]; ok {
		t.Errorf("suggested_email should be omitted when empty, got: %v", body)
	}
}
