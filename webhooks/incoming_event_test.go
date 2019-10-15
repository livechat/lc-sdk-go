package webhooks_test

import (
	"testing"

	"github.com/livechat/lc-sdk-go/webhooks"
)

const incomingEvent = `{"webhook_id":"1188f559c4bae6c4b9a87a1b32d78202","secret_key":"1234567890","action":"incoming_event","payload":{"chat_id":"PS0X0L086G","thread_id":"PZ070E0W1B","event":{"type":"message","text":"14","id":"PZ070E0W1B_3","custom_id":"1dnepb4z00t","recipients":"all","created_at":"2019-10-11T09:41:00.877000Z","author_id":"345f8235-d60d-433e-63c5-7f813a6ffe25"}},"additional_data":{}}`

func TestParseIncomingEventPayload(t *testing.T) {
	p, err := webhooks.ParseIncomingEventPayload([]byte(incomingEvent))

	if err != nil {
		t.Error(err)
	}

	eq(p.WebhookID, "1188f559c4bae6c4b9a87a1b32d78202", t)
	eq(p.SecretKey, "1234567890", t)
	eq(p.Action, "incoming_event", t)

	eq(p.Payload.ChatID, "PS0X0L086G", t)
	eq(p.Payload.ThreadID, "PZ070E0W1B", t)

	eq(p.Payload.Event.ID, "PZ070E0W1B_3", t)

	e := p.Payload.Event.Message()
	eq(e.ID, "PZ070E0W1B_3", t)
	eq(e.Type, "message", t)
	eq(e.Text, "14", t)
	eq(e.CustomID, "1dnepb4z00t", t)
	eq(e.Recipients, "all", t)
	eq(e.CreatedAt.String(), "2019-10-11 09:41:00.877 +0000 UTC", t)
	eq(e.AuthorID, "345f8235-d60d-433e-63c5-7f813a6ffe25", t)
}

func BenchmarkParseIncomingEventPayload(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = webhooks.ParseIncomingEventPayload([]byte(incomingEvent))
	}
}
