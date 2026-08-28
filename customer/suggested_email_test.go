package customer_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/livechat/lc-sdk-go/v7/customer"
)

func captureUpdateCustomerBody(t *testing.T, call func(*customer.API) error) map[string]interface{} {
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

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Fatal("API creation failed")
	}
	if err := call(api); err != nil {
		t.Fatalf("UpdateCustomer failed: %v", err)
	}
	return body
}

func TestUpdateCustomerSendsSuggestedEmail(t *testing.T) {
	suggestion := "guessed@example.com"
	body := captureUpdateCustomerBody(t, func(api *customer.API) error {
		return api.UpdateCustomer(nil, nil, nil, nil, &suggestion, nil, nil, nil)
	})

	got, ok := body["suggested_email"]
	if !ok {
		t.Fatalf("suggested_email missing from the request body: %v", body)
	}
	if got != suggestion {
		t.Errorf("suggested_email = %v, want %v", got, suggestion)
	}
}

// The whole point of the pointer in this version: a dismissal is an empty value, and it has to
// reach the wire as "" rather than being dropped by omitempty.
func TestUpdateCustomerSendsDismissalAsEmptyValue(t *testing.T) {
	dismissal := ""
	body := captureUpdateCustomerBody(t, func(api *customer.API) error {
		return api.UpdateCustomer(nil, nil, nil, nil, &dismissal, nil, nil, nil)
	})

	got, ok := body["suggested_email"]
	if !ok {
		t.Fatalf("a dismissal must reach the wire, got: %v", body)
	}
	if got != "" {
		t.Errorf("suggested_email = %v, want an empty value", got)
	}
}

func TestUpdateCustomerOmitsUntouchedSuggestedEmail(t *testing.T) {
	body := captureUpdateCustomerBody(t, func(api *customer.API) error {
		return api.UpdateCustomer(nil, nil, nil, nil, nil, nil, nil, nil)
	})

	if _, ok := body["suggested_email"]; ok {
		t.Errorf("an untouched suggestion must not be sent, got: %v", body)
	}
}
