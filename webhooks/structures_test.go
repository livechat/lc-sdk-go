package webhooks_test

import (
	"encoding/json"
	"testing"

	"github.com/livechat/lc-sdk-go/v5/webhooks"
)

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

func FilledFormFieldFromJSON(t *testing.T, rawFields json.RawMessage) *webhooks.FilledFormField {
	t.Helper()

	event := webhooks.Event{Type: "filled_form"}
	event.Fields = rawFields

	filledForm := event.FilledForm()
	if filledForm == nil {
		t.Fatalf("FilledForm() returned nil")
	}
	if len(filledForm.Fields) != 1 {
		t.Fatalf("FilledForm() did not return expected number of fields: %d", len(filledForm.Fields))
	}

	return &filledForm.Fields[0]
}
