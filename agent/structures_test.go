package agent_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/livechat/lc-sdk-go/v7/agent"
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

func FilledFormFromJSON(t *testing.T, rawEvent json.RawMessage) *agent.FilledForm {
	t.Helper()

	var event agent.Event
	if err := json.Unmarshal(rawEvent, &event); err != nil {
		t.Fatalf("unmarshalling base event failed: %s", err)
	}

	filledForm := event.FilledForm()
	if filledForm == nil {
		t.Fatalf("FilledForm() returned nil")
	}

	return filledForm
}

func FilledFormFieldFromJSON(t *testing.T, rawFields json.RawMessage) *agent.FilledFormField {
	t.Helper()

	rawEvent := json.RawMessage(fmt.Sprintf(`{"type":"filled_form","fields":%s}`, rawFields))
	filledForm := FilledFormFromJSON(t, rawEvent)

	if len(filledForm.Fields) != 1 {
		t.Fatalf("FilledForm() did not return expected number of fields: %d", len(filledForm.Fields))
	}

	return &filledForm.Fields[0]
}

func newTestForm() agent.Form {
	return agent.Form{
		Event:    agent.Event{Type: "form", Visibility: "all"},
		FormID:   "841954e1-49dd-40bd-8f60-019df106c34a",
		FormType: "custom_form",
		Fields: []agent.FormField{
			{
				ID:                 "field-1",
				Type:               "select",
				Label:              "What is your industry?",
				CustomerPropertyID: "14a484e0-591b-4b90-b9b4-6e64d3e7abc1",
				Options:            []agent.FormOption{{ID: "0", Label: "Automotive"}},
			},
			{ID: "field-2", Type: "email", Label: "What is your email?", Required: true},
		},
	}
}

func TestValidateEventAcceptsForm(t *testing.T) {
	form := newTestForm()
	if err := agent.ValidateEvent(form); err != nil {
		t.Errorf("ValidateEvent(Form) failed: %v", err)
	}
	if err := agent.ValidateEvent(&form); err != nil {
		t.Errorf("ValidateEvent(*Form) failed: %v", err)
	}
	if err := agent.ValidateEventPreview(form); err == nil {
		t.Error("ValidateEventPreview(Form) should fail")
	}
}

func TestFormMarshalUsesOuterFields(t *testing.T) {
	raw, err := json.Marshal(newTestForm())
	if err != nil {
		t.Fatalf("marshalling form failed: %s", err)
	}

	var got map[string]interface{}
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshalling marshalled form failed: %s", err)
	}

	var want map[string]interface{}
	wantJSON := `{
		"type": "form",
		"form_id": "841954e1-49dd-40bd-8f60-019df106c34a",
		"form_type": "custom_form",
		"visibility": "all",
		"created_at": "0001-01-01T00:00:00Z",
		"author_id": "",
		"fields": [
			{
				"id": "field-1",
				"type": "select",
				"label": "What is your industry?",
				"customer_property_id": "14a484e0-591b-4b90-b9b4-6e64d3e7abc1",
				"options": [{"id": "0", "label": "Automotive"}]
			},
			{"id": "field-2", "type": "email", "label": "What is your email?", "required": true}
		]
	}`
	if err := json.Unmarshal([]byte(wantJSON), &want); err != nil {
		t.Fatalf("unmarshalling expected JSON failed: %s", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected form JSON.\nExpected: %s\nGot: %s", wantJSON, raw)
	}
}

func TestEventFormRoundTrip(t *testing.T) {
	raw, err := json.Marshal(newTestForm())
	if err != nil {
		t.Fatalf("marshalling form failed: %s", err)
	}

	var event agent.Event
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("unmarshalling base event failed: %s", err)
	}

	form := event.Form()
	if form == nil {
		t.Fatal("Form() returned nil")
	}
	if form.FormID != "841954e1-49dd-40bd-8f60-019df106c34a" || form.FormType != "custom_form" {
		t.Errorf("Form() did not set form_id and form_type correctly: %q, %q", form.FormID, form.FormType)
	}
	if !reflect.DeepEqual(form.Fields, newTestForm().Fields) {
		t.Errorf("Form() did not return expected fields: %+v", form.Fields)
	}

	event.Type = "message"
	if event.Form() != nil {
		t.Error("Form() should return nil for other event type")
	}
}

func TestFilledFormFormID(t *testing.T) {
	rawEvent := json.RawMessage(`{"type":"filled_form","form_id":"abc","fields":[]}`)
	if got := FilledFormFromJSON(t, rawEvent).FormID; got != "abc" {
		t.Errorf("FilledForm() did not set FormID correctly. Expected %q, got %q", "abc", got)
	}
}
