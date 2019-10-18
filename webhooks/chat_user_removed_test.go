package webhooks_test

import (
	"testing"

	"github.com/livechat/lc-sdk-go/webhooks"
)

const chatUserRemoved = `{"webhook_id":"2c0b13904e79b2aca271e5f84b898f35","secret_key":"1234567890","action":"chat_user_removed","payload":{"chat_id":"PS0X0L086G","user_type":"agent","user_id":"l.wojciechowski@livechatinc.com"},"additional_data":{}}`

func TestParseChatUserRemovedPayload(t *testing.T) {
	p, err := webhooks.ParseChatUserRemovedPayload([]byte(chatUserRemoved))

	if err != nil {
		t.Error(err)
	}

	eq(p.WebhookID, "2c0b13904e79b2aca271e5f84b898f35", t)
	eq(p.SecretKey, "1234567890", t)
	eq(p.Action, "chat_user_removed", t)

	eq(p.Payload.ChatID, "PS0X0L086G", t)
	eq(p.Payload.UserType, "agent", t)
	eq(p.Payload.UserID, "l.wojciechowski@livechatinc.com", t)
}

func BenchmarkParseChatUserRemovedPayload(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = webhooks.ParseChatUserRemovedPayload([]byte(chatUserRemoved))
	}
}
