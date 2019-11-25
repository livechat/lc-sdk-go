package configuration

type Webhook struct {
	Action         WebhookAction   `json:"action"`
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
	OwnerClientID  string          `json:"owner_client_id"`
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

type BotAgent struct {
	ID     string    `json:"id"`
	Name   string    `json:"name"`
	Avatar string    `json:"avatar"`
	Status BotStatus `json:"status"`
}

type BotAgentDetails struct {
	BotAgent
	DefaultGroupPriority GroupPriority `json:"default_group_priority"`
	Application          struct {
		ClientID string `json:"client_id"`
	} `json:"application"`
	MaxChatsCount uint             `json:"max_chats_count"`
	Groups        []BotGroupConfig `json:"groups"`
	Webhooks      BotWebhooks      `json:"webhooks"`
}

type BotWebhooks struct {
	URL       string             `json:"url"`
	SecretKey string             `json:"secret_key"`
	Actions   []BotWebhookAction `json:"actions"`
}

type BotGroupConfig struct {
	ID       uint          `json:"id"`
	Priority GroupPriority `json:"priority"`
}

type BotWebhookAction struct {
	Name           WebhookAction  `json:"name"`
	Filters        WebhookFilters `json:"filters"`
	AdditionalData []string       `json:"additional_data"`
}

type PropertyConfig struct {
	Type        string              `json:"type"`
	Locations   map[string]Location `json:"locations"`
	Description string              `json:"description"`
	Domain      []interface{}       `json:"domain"`
	Range       struct {
		From int `json:"from"`
		To   int `json:"to"`
	} `json:"range"`
}

type Location struct {
	Access map[string]PropertyAccess `json:"access"`
}

type PropertyAccess struct {
	Read  bool `json:"read"`
	Write bool `json:"write"`
}
