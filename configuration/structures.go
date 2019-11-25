package configuration

// Webhook represents webhook to be registered
type Webhook struct {
	Action         *WebhookAction  `json:"action"`
	SecretKey      string          `json:"secret_key"`
	URL            string          `json:"url"`
	AdditionalData []string        `json:"additional_data,omitempty"`
	Description    string          `json:"description,omitempty"`
	Filters        *WebhookFilters `json:"filters,omitempty"`
}

// RegisteredWebhook represents data for webhook registered in Configuration API
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

// WebhookFilters represent set of properties that webhook will use for filtering triggers
type WebhookFilters struct {
	AuthorType    string               `json:"author_type,omitempty"`
	OnlyMyChats   bool                 `json:"only_my_chats,omitempty"`
	ChatMemberIDs *chatMemberIDsFilter `json:"chat_member_ids,omitempty"`
}

type chatMemberIDsFilter struct {
	AgentsAny     []string `json:"agents_any,omitempty"`
	AgentsExclude []string `json:"agents_exclude,omitempty"`
}

// NewChatMemberIDsFilter creates new filter for triggering webhooks based on agents in chat
// `inclusive` parameter controls if the filtered agents should match or exclude given agents
func NewChatMemberIDsFilter(agents []string, inclusive bool) *chatMemberIDsFilter {
	cmf := &chatMemberIDsFilter{}
	switch {
	case inclusive:
		cmf.AgentsAny = agents
	default:
		cmf.AgentsExclude = agents
	}
	return cmf
}

// BotAgent represents basic bot agent information
type BotAgent struct {
	ID     string    `json:"id"`
	Name   string    `json:"name"`
	Avatar string    `json:"avatar"`
	Status BotStatus `json:"status"`
}

// BotAgentDetails represents detailed bot agent information
type BotAgentDetails struct {
	BotAgent
	DefaultGroupPriority GroupPriority `json:"default_group_priority"`
	Application          struct {
		ClientID string `json:"client_id"`
	} `json:"application"`
	MaxChatsCount uint              `json:"max_chats_count"`
	Groups        []*BotGroupConfig `json:"groups"`
	Webhooks      *BotWebhooks      `json:"webhooks"`
}

// BotWebhooks represents webhooks configuration for bot agent
type BotWebhooks struct {
	URL       string              `json:"url"`
	SecretKey string              `json:"secret_key"`
	Actions   []*BotWebhookAction `json:"actions"`
}

// BotGroupConfig defines bot's priority and membership in group
type BotGroupConfig struct {
	ID       uint          `json:"id"`
	Priority GroupPriority `json:"priority"`
}

// BotWebhookAction represents action that should trigger bot's webhook
type BotWebhookAction struct {
	Name           *WebhookAction  `json:"name"`
	Filters        *WebhookFilters `json:"filters"`
	AdditionalData []string        `json:"additional_data"`
}

// PropertyConfig defines configuration of a property
type PropertyConfig struct {
	Type        string               `json:"type"`
	Locations   map[string]*Location `json:"locations"`
	Description string               `json:"description"`
	Domain      []interface{}        `json:"domain"`
	Range       struct {
		From int `json:"from"`
		To   int `json:"to"`
	} `json:"range"`
}

// Location represents property location
type Location struct {
	Access map[string]*PropertyAccess `json:"access"`
}

// PropertyAccess defines read/write rights of a property
type PropertyAccess struct {
	Read  bool `json:"read"`
	Write bool `json:"write"`
}
