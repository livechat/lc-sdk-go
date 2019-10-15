package webhooks_test

import (
	"testing"

	"github.com/livechat/lc-sdk-go/webhooks"
)

const threadClosed = `{"webhook_id":"37901609f767dea5bc3a8acf5458bd80","secret_key":"1234567890","action":"thread_closed","payload":{"chat_id":"PS0X0L086G","thread_id":"PZ070E0W1B","user_id":"l.wojciechowski@livechatinc.com"},"additional_data":{}}`

func TestParseThreadClosedPayload(t *testing.T) {
	p, err := webhooks.ParseThreadClosedPayload([]byte(threadClosed))

	if err != nil {
		t.Error(err)
	}

	eq(p.WebhookID, "37901609f767dea5bc3a8acf5458bd80", t)
	eq(p.SecretKey, "1234567890", t)
	eq(p.Action, "thread_closed", t)

	eq(p.Payload.ChatID, "PS0X0L086G", t)
	eq(p.Payload.ThreadID, "PZ070E0W1B", t)
	eq(p.Payload.UserID, "l.wojciechowski@livechatinc.com", t)
}
