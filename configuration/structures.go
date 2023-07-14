package configuration

// Properties represents LiveChat properties in form of property_namespace -> property -> value.
type Properties map[string]map[string]interface{}

// Webhook represents webhook to be registered
type Webhook struct {
	Action         WebhookAction   `json:"action"`
	SecretKey      string          `json:"secret_key"`
	URL            string          `json:"url"`
	Type           string          `json:"type"`
	AdditionalData []string        `json:"additional_data,omitempty"`
	Description    string          `json:"description,omitempty"`
	Filters        *WebhookFilters `json:"filters,omitempty"`
}

// RegisteredWebhook represents data for webhook registered in Configuration API
type RegisteredWebhook struct {
	ID             string          `json:"id"`
	Action         string          `json:"action"`
	SecretKey      string          `json:"secret_key"`
	URL            string          `json:"url"`
	Type           string          `json:"type"`
	AdditionalData []string        `json:"additional_data,omitempty"`
	Description    string          `json:"description,omitempty"`
	Filters        *WebhookFilters `json:"filters,omitempty"`
	OwnerClientID  string          `json:"owner_client_id"`
}

// WebhookData represents available webhook definition
type WebhookData struct {
	Action         string   `json:"action"`
	AdditionalData []string `json:"additional_data,omitempty"`
	Filters        []string `json:"filters,omitempty"`
}

// WebhooksState represents state of webhooks for given clientID on given license
type WebhooksState struct {
	Enabled bool `json:"license_webhooks_enabled"`
}

// ManageWebhooksStateOptions are options for methods responsible for webhooks' state management:
// EnableWebhooks, DisableWebhooks and GetWebhooksState
type ManageWebhooksStateOptions struct {
	ClientID string
}

// WebhookFilters represent set of properties that webhook will use for filtering triggers
type WebhookFilters struct {
	AuthorType   string              `json:"author_type,omitempty"`
	OnlyMyChats  bool                `json:"only_my_chats,omitempty"`
	ChatPresence *chatPresenceFilter `json:"chat_presence,omitempty"`
	SourceType   []string            `json:"source_type,omitempty"`
}

type chatPresenceFilter struct {
	UserIDs *userIDsFilter `json:"user_ids,omitempty"`
	MyBots  bool           `json:"my_bots,omitempty"`
}

type userIDsFilter struct {
	Values        []string `json:"values,omitempty"`
	ExcludeValues []string `json:"exclude_values,omitempty"`
}

// NewChatPresenceFilter creates new filter for triggering webhooks based on chat members
func NewChatPresenceFilter() *chatPresenceFilter {
	return &chatPresenceFilter{}
}

// WithMyBots causes webhooks to be triggered if there's any bot owned by the integration
// in the chat
func (cpf *chatPresenceFilter) WithMyBots() *chatPresenceFilter {
	cpf.MyBots = true
	return cpf
}

// WithUserIDs causes webhooks to be triggered based on chat presence of any provided user_id
// `inclusive` parameter controls if the provided user_ids should match or exclude users present in the chat
func (cpf *chatPresenceFilter) WithUserIDs(user_ids []string, inclusive bool) *chatPresenceFilter {
	if inclusive {
		cpf.UserIDs = &userIDsFilter{
			Values: user_ids,
		}
	} else {
		cpf.UserIDs = &userIDsFilter{
			ExcludeValues: user_ids,
		}
	}
	return cpf
}

// Bot represents basic bot agent information
type Bot struct {
	ID                   string         `json:"id"`
	Name                 string         `json:"name,omitempty"`
	AvatarPath           string         `json:"avatar_path,omitempty"`
	DefaultGroupPriority GroupPriority  `json:"default_group_priority,omitempty"`
	ClientID             string         `json:"owner_client_id,omitempty"`
	MaxChatsCount        uint           `json:"max_chats_count,omitempty"`
	Groups               []*GroupConfig `json:"groups,omitempty"`
	JobTitle             string         `json:"job_title,omitempty"`
	WorkScheduler        *WorkScheduler `json:"work_scheduler,omitempty"`
}

// GroupConfig defines bot's priority and membership in group
type GroupConfig struct {
	ID       uint          `json:"id"`
	Priority GroupPriority `json:"priority"`
}

// PropertyConfig defines configuration of a property
type PropertyConfig struct {
	Name          string                     `json:"name"`
	OwnerClientID string                     `json:"owner_client_id,omitempty"`
	Type          string                     `json:"type"`
	Access        map[string]*PropertyAccess `json:"access"`
	Description   string                     `json:"description,omitempty"`
	Domain        []interface{}              `json:"domain,omitempty"`
	Range         *struct {
		From int `json:"from"`
		To   int `json:"to"`
	} `json:"range,omitempty"`
	PublicAccess []string    `json:"public_access,omitempty"`
	DefaultValue interface{} `json:"default_value,omitempty"`
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
	ID         string `json:"id"`
	AccountID  string `json:"account_id,omitempty"`
	LastLogout string `json:"last_logout,omitempty"`
	*AgentFields
}

// AgentFields defines set of configurable Agent fields
type AgentFields struct {
	Name               string         `json:"name,omitempty"`
	Role               string         `json:"role,omitempty"`
	Avatar             string         `json:"avatar,omitempty"`
	JobTitle           string         `json:"job_title,omitempty"`
	Mobile             string         `json:"mobile,omitempty"`
	MaxChatsCount      uint           `json:"max_chats_count,omitempty"`
	AwaitingApproval   bool           `json:"awaiting_approval,omitempty"`
	Suspended          bool           `json:"suspended,omitempty"`
	Groups             []GroupConfig  `json:"groups,omitempty"`
	WorkScheduler      *WorkScheduler `json:"work_scheduler,omitempty"`
	Notifications      []string       `json:"notifications,omitempty"`
	EmailSubscriptions []string       `json:"email_subscriptions,omitempty"`
}

// WorkScheduler represents work schedule data
type WorkScheduler struct {
	Timezone string     `json:"timezone"`
	Schedule []Schedule `json:"schedule"`
}

// Schedule represent a single day work schedule
type Schedule struct {
	Enabled bool    `json:"enabled"`
	Day     Weekday `json:"day"`
	Start   string  `json:"start"`
	End     string  `json:"end"`
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

// ManageWebhooksDefinitionOptions are options for methods responsible for webhooks' definition management:
// ListWebhooks, RegisterWebhook and UnregisterWebhook
type ManageWebhooksDefinitionOptions struct {
	ClientID string
}

// Condition is option for methods responsible for auto access management:
// AddAutoAccess, UpdateAutoAccess
type Condition struct {
	Values        []Match `json:"values"`
	ExcludeValues []Match `json:"exclude_values"`
}

// Match represents possible match conditions for Condition
type Match struct {
	Value      string `json:"value"`
	ExactMatch bool   `json:"exact_match,omitempty"`
}

// GeolocationCondition is option for methods responsible for auto access management:
// AddAutoAccess, UpdateAutoAccess
type GeolocationCondition struct {
	Values []GeolocationMatch `json:"values"`
}

// GeolocationMatch represents possible match conditions for GeolocationCondition
type GeolocationMatch struct {
	Country     string `json:"country,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
	Region      string `json:"region,omitempty"`
	City        string `json:"city,omitempty"`
}

type AutoAccess struct {
	ID     string `json:"id"`
	Access struct {
		Groups []int `json:"groups"`
	} `json:"access"`
	Conditions struct {
		Url         *Condition            `json:"url,omitempty"`
		Domain      *Condition            `json:"domain,omitempty"`
		Geolocation *GeolocationCondition `json:"geolocation,omitempty"`
	} `json:"conditions"`
	Description string `json:"description,omitempty"`
	NextID      string `json:"next_id,omitempty"`
}

type PlanLimits []struct {
	Resource     string `json:"resource"`
	LimitBalance int32  `json:"limit_balance"`
	Id           string `json:"id,omitempty"`
}

type ChannelActivity []struct {
	ChannelType            string `json:"channel_type"`
	ChannelSubtype         string `json:"channel_subtype"`
	FirstActivityTimestamp string `json:"first_activity_timestamp"`
}

type Tag struct {
	Name      string `json:"name"`
	GroupIDs  []int  `json:"group_ids"`
	CreatedAt string `json:"created_at"`
	AuthorID  string `json:"author_id"`
}
