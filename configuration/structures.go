package configuration

import "github.com/livechat/lc-sdk-go/configuration/action"

type Webhook struct {
	ID             string          `json:"webhook_id"`
	Action         action.Webhook  `json:"action"`
	SecretKey      string          `json:"secret_key"`
	URL            string          `json:"url"`
	AdditionalData []string        `json:"additional_data,omitempty"`
	Description    string          `json:"description,omitempty"`
	Filters        *WebhookFilters `json:"filters,omitempty"`
}

type WebhookFilters struct {
	AuthorType    string `json:"author_type"`
	OnlyMyChats   bool   `json:"only_my_chats"`
	ChatMemberIds struct {
		AgentsAny     []string `json:"agents_any,omitempty"`
		AgentsExclude []string `json:"agents_exclude,omitempty"`
	} `json:"chat_member_ids"`
}
