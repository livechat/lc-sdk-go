package customer_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/livechat/lc-sdk-go/v7/customer"
)

func TestFilledFormTypesOK(t *testing.T) {
	rawFields := json.RawMessage(`[{"id": "42", "type": "email", "answer":"foo@bar.eu"}]`)
	testCases := []struct {
		name             string
		expectedFormType string
		rawFormType      json.RawMessage
	}{
		{name: "Missing form_type", expectedFormType: "", rawFormType: nil},
		{name: "Null form_type", expectedFormType: "", rawFormType: json.RawMessage(`null`)},
		{name: "Empty string form_type", expectedFormType: "", rawFormType: json.RawMessage(`""`)},
		{name: "Existing form_type", expectedFormType: "prechat", rawFormType: json.RawMessage(`"prechat"`)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var suffixFormType string
			if tc.rawFormType != nil {
				suffixFormType = fmt.Sprintf(`,"form_type":%s`, tc.rawFormType)
			}

			rawEvent := json.RawMessage(fmt.Sprintf(`{"type":"filled_form","fields":%s%s}`, rawFields, suffixFormType))
			filledForm := FilledFormFromJSON(t, rawEvent)

			if filledForm.FormType != tc.expectedFormType {
				t.Fatalf("FilledForm() did not set FormType correctly. Expected %q, got %q", tc.expectedFormType, filledForm.FormType)
			}

			if len(filledForm.Fields) != 1 {
				t.Fatalf("FilledForm() did not return expected number of fields: %d", len(filledForm.Fields))
			}
		})
	}
}

func TestFilledFormFieldTypesOK(t *testing.T) {
	t.Run("Single-answer filled_form field", func(t *testing.T) {
		rawFormFields := json.RawMessage(`[{"id": "42", "type": "email", "answer":"foo@bar.eu"}]`)
		singleField := FilledFormFieldFromJSON(t, rawFormFields).Single()

		if singleField.Answer != "foo@bar.eu" {
			t.Errorf("FilledFormField() did not return expected answer: %s", singleField.Answer)
		}
	})

	t.Run("Single-choice filled_form field", func(t *testing.T) {
		rawFormFields := json.RawMessage(`[{"id": "42", "type": "radio", "answer":{"id": "1", "label": "foo"}}]`)
		singleChoiceField := FilledFormFieldFromJSON(t, rawFormFields).SingleChoice()

		if singleChoiceField.Answer.ID != "1" || singleChoiceField.Answer.Label != "foo" {
			t.Errorf("FilledFormField() did not return expected answer: %s", singleChoiceField.Answer)
		}
	})

	t.Run("Multiple-choice filled_form field", func(t *testing.T) {
		rawFormFields := json.RawMessage(`[{"id": "42", "type": "checkbox", "answers":[{"id": "1", "label": "foo"}, {"id": "2", "label": "bar"}]}]`)
		multiChoiceField := FilledFormFieldFromJSON(t, rawFormFields).MultiChoice()

		if len(multiChoiceField.Answers) != 2 {
			t.Errorf("FilledFormField() did not return expected number of answers: %d", len(multiChoiceField.Answer))
		}
		if multiChoiceField.Answers[0].ID != "1" || multiChoiceField.Answers[0].Label != "foo" {
			t.Errorf("FilledFormField() did not return expected answer: %s", multiChoiceField.Answers[0])
		}
		if multiChoiceField.Answers[1].ID != "2" || multiChoiceField.Answers[1].Label != "bar" {
			t.Errorf("FilledFormField() did not return expected answer: %s", multiChoiceField.Answers[1])
		}
	})

	t.Run("Email checkbox filled_form field", func(t *testing.T) {
		rawFormFields := json.RawMessage(`[{"id": "42", "type": "checkbox_for_email", "answer": true}]`)
		emailCheckboxField := FilledFormFieldFromJSON(t, rawFormFields).EmailCheckbox()

		if !emailCheckboxField.Answer {
			t.Errorf("FilledFormField() did not return expected answer: %t", emailCheckboxField.Answer)
		}
	})

	t.Run("Group chooser filled_form field", func(t *testing.T) {
		rawFormFields := json.RawMessage(`[{"id": "42", "type": "group_chooser", "answer":{"id": "1", "label": "foo", "group_id": 42}}]`)
		groupChooserField := FilledFormFieldFromJSON(t, rawFormFields).GroupChooser()

		if groupChooserField.Answer.ID != "1" || groupChooserField.Answer.Label != "foo" || groupChooserField.Answer.GroupID != 42 {
			t.Errorf("FilledFormField() did not return expected answer: %v", groupChooserField.Answer)
		}
	})
}

func FilledFormFromJSON(t *testing.T, rawEvent json.RawMessage) *customer.FilledForm {
	t.Helper()

	var event customer.Event
	if err := json.Unmarshal(rawEvent, &event); err != nil {
		t.Fatalf("unmarshalling base event failed: %s", err)
	}

	filledForm := event.FilledForm()
	if filledForm == nil {
		t.Fatalf("FilledForm() returned nil")
	}

	return filledForm
}

func FilledFormFieldFromJSON(t *testing.T, rawFields json.RawMessage) *customer.FilledFormField {
	t.Helper()

	rawEvent := json.RawMessage(fmt.Sprintf(`{"type":"filled_form","fields":%s}`, rawFields))
	filledForm := FilledFormFromJSON(t, rawEvent)

	if len(filledForm.Fields) != 1 {
		t.Fatalf("FilledForm() did not return expected number of fields: %d", len(filledForm.Fields))
	}

	return &filledForm.Fields[0]
}

// maxEventKeys is the maximum number of keys the Customer API accepts in a single event object.
const maxEventKeys = 25

func fullBaseEvent(eventType string) customer.Event {
	return customer.Event{
		ID:         "event-id",
		CustomID:   "custom-id",
		CreatedAt:  time.Date(2026, 11, 20, 12, 0, 0, 0, time.UTC),
		AuthorID:   "author-id",
		Properties: customer.Properties{"routing": {"keyword": "help"}, "source": {"channel": "sms"}},
		Recipients: "all",
		Type:       eventType,
	}
}

func serializedEventKeys(t *testing.T, event interface{}) map[string]json.RawMessage {
	t.Helper()

	raw, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("marshalling %T failed: %s", event, err)
	}

	var keys map[string]json.RawMessage
	if err := json.Unmarshal(raw, &keys); err != nil {
		t.Fatalf("unmarshalling %T into map failed: %s", event, err)
	}

	return keys
}

func assertEventSerializationWithinLimits(t *testing.T, event interface{}) {
	t.Helper()

	keys := serializedEventKeys(t, event)
	for key, value := range keys {
		if string(value) == "null" {
			t.Errorf("%T serialized unset key %q as null", event, key)
		}
	}
	if len(keys) > maxEventKeys {
		t.Errorf("%T serialized to %d keys, Customer API allows at most %d", event, len(keys), maxEventKeys)
	}
}

func allTypedEvents() map[string]interface{} {
	return map[string]interface{}{
		"Event": fullBaseEvent("custom"),
		"Message": customer.Message{
			Event:    fullBaseEvent("message"),
			Text:     "hello",
			Postback: &customer.Postback{ID: "postback-id", ThreadID: "thread-id", EventID: "event-id", Type: "message", Value: "value"},
		},
		"File": customer.File{
			Event:           fullBaseEvent("file"),
			ContentType:     "image/png",
			Name:            "image.png",
			URL:             "https://example.com/image.png",
			ThumbnailURL:    "https://example.com/thumb.png",
			Thumbnail2xURL:  "https://example.com/thumb2x.png",
			Width:           640,
			Height:          480,
			Size:            1024,
			AlternativeText: "an image",
		},
		"RichMessage": customer.RichMessage{
			Event:      fullBaseEvent("rich_message"),
			TemplateID: "cards",
			Layout:     "list",
			Elements:   []customer.RichMessageElement{{Title: "title", Subtitle: "subtitle"}},
		},
		"SystemMessage": customer.SystemMessage{
			Event:             fullBaseEvent("system_message"),
			SystemMessageType: "followup_requested",
			Text:              "Follow-up requested by %name%",
			TextVars:          map[string]string{"name": "John"},
		},
		"System": customer.System{
			Event:   fullBaseEvent("system"),
			Source:  "source",
			Subtype: "subtype",
			Details: "details",
			Version: 1,
		},
	}
}

func TestTypedEventsSerializationWithinCustomerAPILimits(t *testing.T) {
	for name, event := range allTypedEvents() {
		t.Run(name, func(t *testing.T) {
			if err := customer.ValidateEvent(event); err != nil {
				t.Fatalf("ValidateEvent() rejected test event: %s", err)
			}
			assertEventSerializationWithinLimits(t, event)
		})
	}
}

func TestMessageSerialization(t *testing.T) {
	message := customer.Message{
		Event: customer.Event{
			CustomID:   "custom-id",
			Properties: customer.Properties{"routing": {"keyword": "help"}},
			Recipients: "all",
			Type:       "message",
		},
		Text: "hello",
	}

	assertEventSerializationWithinLimits(t, message)

	keys := serializedEventKeys(t, message)
	for _, key := range []string{"custom_id", "properties", "recipients", "type", "text"} {
		if _, ok := keys[key]; !ok {
			t.Errorf("Message serialization is missing %q key: %v", key, keys)
		}
	}
}

func TestSystemMessageWithTextVarsSerialization(t *testing.T) {
	systemMessage := customer.SystemMessage{
		Event:             customer.Event{Type: "system_message", Recipients: "all"},
		SystemMessageType: "followup_requested",
		TextVars:          map[string]string{"email": "john@example.com"},
	}

	assertEventSerializationWithinLimits(t, systemMessage)

	keys := serializedEventKeys(t, systemMessage)
	if string(keys["text_vars"]) != `{"email":"john@example.com"}` {
		t.Errorf("SystemMessage serialized unexpected text_vars: %s", keys["text_vars"])
	}
}

func TestRichMessageLayoutSerialization(t *testing.T) {
	testCases := []struct {
		name           string
		layout         string
		expectedLayout string
	}{
		{name: "Empty layout", layout: "", expectedLayout: ""},
		{name: "Layout set", layout: "list", expectedLayout: `"list"`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			richMessage := customer.RichMessage{
				Event:      customer.Event{Type: "rich_message"},
				TemplateID: "cards",
				Layout:     tc.layout,
				Elements:   []customer.RichMessageElement{{Title: "title"}},
			}

			assertEventSerializationWithinLimits(t, richMessage)

			layout, ok := serializedEventKeys(t, richMessage)["layout"]
			if tc.expectedLayout == "" && ok {
				t.Errorf("RichMessage serialized empty layout: %s", layout)
			}
			if tc.expectedLayout != "" && string(layout) != tc.expectedLayout {
				t.Errorf("RichMessage serialized unexpected layout: expected %s, got %s", tc.expectedLayout, layout)
			}
		})
	}
}

func TestTypedEventsRoundTrip(t *testing.T) {
	converters := map[string]func(e *customer.Event) interface{}{
		"Event":         func(e *customer.Event) interface{} { return e },
		"Message":       func(e *customer.Event) interface{} { return e.Message() },
		"File":          func(e *customer.Event) interface{} { return e.File() },
		"RichMessage":   func(e *customer.Event) interface{} { return e.RichMessage() },
		"SystemMessage": func(e *customer.Event) interface{} { return e.SystemMessage() },
		"System":        func(e *customer.Event) interface{} { return e.System() },
	}

	for name, event := range allTypedEvents() {
		t.Run(name, func(t *testing.T) {
			original, err := json.Marshal(event)
			if err != nil {
				t.Fatalf("marshalling %T failed: %s", event, err)
			}

			var decoded customer.Event
			if err := json.Unmarshal(original, &decoded); err != nil {
				t.Fatalf("unmarshalling %T into Event failed: %s", event, err)
			}

			converted := converters[name](&decoded)
			if string(mustMarshal(t, converted)) == "null" {
				t.Fatalf("converting Event to %s returned nil", name)
			}
			if roundTripped := mustMarshal(t, converted); !bytes.Equal(original, roundTripped) {
				t.Errorf("%s changed after round trip:\noriginal: %s\ngot:      %s", name, original, roundTripped)
			}
		})
	}
}

func TestEventPreservesExplicitNullAfterRoundTrip(t *testing.T) {
	raw := []byte(`{"created_at":"0001-01-01T00:00:00Z","author_id":"","type":"custom","text":null}`)

	var event customer.Event
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("unmarshalling Event failed: %s", err)
	}

	if got := mustMarshal(t, event); !bytes.Equal(raw, got) {
		t.Errorf("Event changed after round trip:\noriginal: %s\ngot:      %s", raw, got)
	}
}

func mustMarshal(t *testing.T, v interface{}) []byte {
	t.Helper()

	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshalling %T failed: %s", v, err)
	}

	return raw
}
