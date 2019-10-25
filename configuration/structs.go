package configuration

type Webhook struct {
	Action         string          `json:"action"`
	SecretKey      string          `json:"secret_key"`
	URL            string          `json:"url"`
	AdditionalData []string        `json:"additional_data,omitempty"`
	Description    string          `json:"description,omitempty"`
	Filters        *WebhookFilters `json:"filters,omitempty"`
}

type RegisteredWebhook struct {
	ID             string          `json:"webhook_id"`
	Action         string          `json:"action"`
	SecretKey      string          `json:"secret_key"`
	URL            string          `json:"url"`
	AdditionalData []string        `json:"additional_data,omitempty"`
	Description    string          `json:"description,omitempty"`
	Filters        *WebhookFilters `json:"filters,omitempty"`
	Owner          string          `json:"owner_client_id"`
}

type WebhookFilters struct {
	AuthorType    string               `json:"author_type,omitempty"`
	OnlyMyChats   bool                 `json:"only_my_chats,omitempty"`
	ChatMemberIDs *ChatMemberIDsFilter `json:"chat_member_ids,omitempty"`
}

type ChatMemberIDsFilter struct {
	AgentsAny     []string `json:"agents_any,omitempty"`
	AgentsExclude []string `json:"agents_exclude,omitempty"`
}
