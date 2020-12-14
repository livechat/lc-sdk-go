package configuration

// Webhook represents webhook to be registered
type Webhook struct {
	Action         WebhookAction   `json:"action"`
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

// AvailableWebhook represents available webhook definition
type AvailableWebhook struct {
	Action         string   `json:"action"`
	AdditionalData []string `json:"additional_data,omitempty"`
	Filters        []string `json:"filters,omitempty"`
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
	MaxChatsCount uint           `json:"max_chats_count"`
	Groups        []*GroupConfig `json:"groups"`
	Webhooks      *BotWebhooks   `json:"webhooks"`
}

// BotWebhooks represents webhooks configuration for bot agent
type BotWebhooks struct {
	URL       string              `json:"url"`
	SecretKey string              `json:"secret_key"`
	Actions   []*BotWebhookAction `json:"actions"`
}

// GroupConfig defines bot's priority and membership in group
type GroupConfig struct {
	ID       uint          `json:"id"`
	Priority GroupPriority `json:"priority"`
}

// BotWebhookAction represents action that should trigger bot's webhook
type BotWebhookAction struct {
	Name           WebhookAction   `json:"name"`
	Filters        *WebhookFilters `json:"filters"`
	AdditionalData []string        `json:"additional_data"`
}

// PropertyConfig defines configuration of a property
type PropertyConfig struct {
	Name          string                     `json:"name"`
	OwnerClientID string                     `json:"owner_client_id"`
	Type          string                     `json:"type"`
	Access        map[string]*PropertyAccess `json:"access"`
	Description   string                     `json:"description,omitempty"`
	Domain        []interface{}              `json:"domain,omitempty"`
	Range         *struct {
		From int `json:"from"`
		To   int `json:"to"`
	} `json:"range,omitempty"`
	PublicAccess []string `json:"public_access,omitempty"`
}

// PropertyAccess defines read/write rights of a property
type PropertyAccess struct {
	Agent    []string `json:"agent"`
	Customer []string `json:"customer"`
}

// Group defines basic group information
type Group struct {
	ID              int                      `json:"id"`
	Name            string                   `json:"name"`
	LanguageCode    string                   `json:"language_code"`
	AgentPriorities map[string]GroupPriority `json:"agent_priorities"`
	RoutingStatus   string                   `json:"routing_status"`
}

// Agent defines basic Agent information
type Agent struct {
	ID string `json:"id"`
	*AgentFields
}

// Agent defines set of configurable Agent fields
type AgentFields struct {
	Name               string        `json:"name,omitempty"`
	Role               string        `json:"role,omitempty"`
	AvatarPath         string        `json:"avatar_path,omitempty"`
	JobTitle           string        `json:"job_title,omitempty"`
	Mobile             string        `json:"mobile,omitempty"`
	MaxChatsCount      uint          `json:"max_chats_count,omitempty"`
	AwaitingApproval   bool          `json:"awaiting_approval,omitempty"`
	Groups             []GroupConfig `json:"groups,omitempty"`
	WorkScheduler      WorkScheduler `json:"work_scheduler,omitempty"`
	Notifications      []string      `json:"notifications,omitempty"`
	EmailSubscriptions []string      `json:"email_subscriptions,omitempty"`
}

// WorkScheduler represents work schedule data
type WorkScheduler map[Weekday]WorkSchedulerDay

// WorkSchedulerDay represents single day work schedule
type WorkSchedulerDay struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// Weekday represents allowed weekday names for work scheduler
type Weekday string

const (
	Monday    Weekday = "monday"
	Tuesday   Weekday = "tuesday"
	Wednesday Weekday = "wednesday"
	Thursday  Weekday = "thursday"
	Friday    Weekday = "friday"
	Saturday  Weekday = "saturday"
	Sunday    Weekday = "sunday"
)

// AgentsFilters defines set of filters for getting agents
type AgentsFilters struct {
	GroupIDs []int32 `json:"group_ids"`
}
