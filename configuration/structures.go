package configuration

// Properties represents LiveChat properties in form of property_namespace -> property -> value.
type Properties map[string]map[string]interface{}

type ListGroupsPropertiesRequestOptions struct {
	Namespace  string `json:"namespace,omitempty"`
	NamePrefix string `json:"name_prefix,omitempty"`
}

type ListLicensePropertiesRequestOptions struct {
	Namespace  string `json:"namespace,omitempty"`
	NamePrefix string `json:"name_prefix,omitempty"`
}

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
	Avatar               string         `json:"avatar,omitempty"`
	DefaultGroupPriority GroupPriority  `json:"default_group_priority,omitempty"`
	ClientID             string         `json:"owner_client_id,omitempty"`
	MaxChatsCount        uint           `json:"max_chats_count,omitempty"`
	Groups               []GroupConfig  `json:"groups,omitempty"`
	JobTitle             string         `json:"job_title,omitempty"`
	WorkScheduler        *WorkScheduler `json:"work_scheduler,omitempty"`
}

// BotTemplate represents basic bot template information
type BotTemplate struct {
	ID                   string        `json:"id"`
	Name                 string        `json:"name,omitempty"`
	Avatar               string        `json:"avatar,omitempty"`
	MaxChatsCount        uint          `json:"max_chats_count,omitempty"`
	DefaultGroupPriority GroupPriority `json:"default_group_priority,omitempty"`
	JobTitle             string        `json:"job_title,omitempty"`
}

type CreateBotRequestOptions struct {
	Avatar               string         `json:"avatar,omitempty"`
	DefaultGroupPriority GroupPriority  `json:"default_group_priority,omitempty"`
	JobTitle             string         `json:"job_title,omitempty"`
	MaxChatsCount        *uint          `json:"max_chats_count,omitempty"`
	Groups               []GroupConfig  `json:"groups,omitempty"`
	OwnerClientID        string         `json:"owner_client_id,omitempty"`
	WorkScheduler        *WorkScheduler `json:"work_scheduler,omitempty"`
}

type CreateBotResponse struct {
	BotID  string `json:"id"`
	Secret string `json:"secret"`
}

type UpdateBotRequestOptions struct {
	Name                 string         `json:"name,omitempty"`
	Avatar               string         `json:"avatar,omitempty"`
	DefaultGroupPriority GroupPriority  `json:"default_group_priority,omitempty"`
	JobTitle             string         `json:"job_title,omitempty"`
	MaxChatsCount        *uint          `json:"max_chats_count,omitempty"`
	Groups               []GroupConfig  `json:"groups,omitempty"`
	OwnerClientID        string         `json:"owner_client_id,omitempty"`
	WorkScheduler        *WorkScheduler `json:"work_scheduler,omitempty"`
}

type ResetBotSecretRequestOptions struct {
	OwnerClientID string `json:"owner_client_id,omitempty"`
}

type IssueBotTokenRequest struct {
	BotID          string `json:"bot_id"`
	BotSecret      string `json:"bot_secret"`
	OrganizationID string `json:"organization_id"`
	ClientID       string `json:"client_id"`
}

type CreateBotTemplateRequestOptions struct {
	Avatar                      string        `json:"avatar,omitempty"`
	MaxChatsCount               *uint         `json:"max_chats_count,omitempty"`
	DefaultGroupPriority        GroupPriority `json:"default_group_priority,omitempty"`
	JobTitle                    string        `json:"job_title,omitempty"`
	OwnerClientID               string        `json:"owner_client_id,omitempty"`
	AffectExistingInstallations bool          `json:"affect_existing_installations,omitempty"`
}

type CreateBotTemplateResponse struct {
	BotTemplateID string `json:"id"`
	Secret        string `json:"secret"`
}

type UpdateBotTemplateRequestOptions struct {
	Name                        string        `json:"name,omitempty"`
	Avatar                      string        `json:"avatar,omitempty"`
	MaxChatsCount               *uint         `json:"max_chats_count,omitempty"`
	DefaultGroupPriority        GroupPriority `json:"default_group_priority,omitempty"`
	JobTitle                    string        `json:"job_title,omitempty"`
	OwnerClientID               string        `json:"owner_client_id,omitempty"`
	AffectExistingInstallations bool          `json:"affect_existing_installations,omitempty"`
}

type DeleteBotTemplateRequestOptions struct {
	OwnerClientID               string `json:"owner_client_id,omitempty"`
	AffectExistingInstallations bool   `json:"affect_existing_installations,omitempty"`
}

type ListBotTemplatesRequestOptions struct {
	OwnerClientID string `json:"owner_client_id,omitempty"`
}

type ResetBotTemplateSecretRequestOptions struct {
	OwnerClientID               string `json:"owner_client_id,omitempty"`
	AffectExistingInstallations bool   `json:"affect_existing_installations,omitempty"`
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

type CreateGroupRequestOptions struct {
	LanguageCode string `json:"language_code,omitempty"`
}

type UpdateGroupRequestOptions struct {
	Name            string                   `json:"name,omitempty"`
	LanguageCode    string                   `json:"language_code,omitempty"`
	AgentPriorities map[string]GroupPriority `json:"agent_priorities,omitempty"`
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

// AutoAccessConditions must have at least one of Url, Domain or Geolocation set
type AutoAccessConditions struct {
	Url         *Condition            `json:"url,omitempty"`
	Domain      *Condition            `json:"domain,omitempty"`
	Geolocation *GeolocationCondition `json:"geolocation,omitempty"`
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

type Access struct {
	Groups []int `json:"groups"`
}

type AutoAccess struct {
	ID          string               `json:"id"`
	Access      Access               `json:"access"`
	Conditions  AutoAccessConditions `json:"conditions"`
	Description string               `json:"description,omitempty"`
	NextID      string               `json:"next_id,omitempty"`
}

type AddAutoAccessRequestOptions struct {
	Description string `json:"description,omitempty"`
	NextID      string `json:"next_id,omitempty"`
}

type UpdateAutoAccessRequestOptions struct {
	Access      *Access               `json:"access,omitempty"`
	Conditions  *AutoAccessConditions `json:"conditions,omitempty"`
	Description *string               `json:"description,omitempty"`
	NextID      *string               `json:"next_id,omitempty"`
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

type CompanyDetails struct {
	InvoiceName  *string `json:"invoice_name,omitempty"`
	Company      *string `json:"company,omitempty"`
	Street       *string `json:"street,omitempty"`
	PostalCode   *string `json:"postal_code,omitempty"`
	City         *string `json:"city,omitempty"`
	Country      *string `json:"country,omitempty"`
	NIP          *string `json:"nip,omitempty"`
	State        *string `json:"state,omitempty"`
	Province     *string `json:"province,omitempty"`
	Phone        *string `json:"phone,omitempty"`
	URL          *string `json:"url,omitempty"`
	InvoiceEmail *string `json:"invoice_email,omitempty"`
	CompanySize  *string `json:"company_size,omitempty"`
	ChatPurpose  *string `json:"chat_purpose,omitempty"`
	Audience     *string `json:"audience,omitempty"`
	Industry     *string `json:"industry,omitempty"`
}

// RichMessageElementImage represents an image in a rich message element
type RichMessageElementImage struct {
	URL             string `json:"url"`
	AlternativeText string `json:"alternative_text"`
}

// RichMessageElementButton represents a button in a rich message element
type RichMessageElementButton struct {
	ButtonID   string `json:"button_id,omitempty"`
	PostbackID string `json:"postback_id,omitempty"`
	Text       string `json:"text"`
	Type       string `json:"type"`
	Value      string `json:"value"`
	Role       string `json:"role,omitempty"`
}

// RichMessageElement represents an element within a rich message
type RichMessageElement struct {
	Title    string                      `json:"title,omitempty"`
	Subtitle string                      `json:"subtitle,omitempty"`
	Image    *RichMessageElementImage    `json:"image,omitempty"`
	Buttons  []*RichMessageElementButton `json:"buttons,omitempty"`
}

// RichMessage represents a rich greeting message
type RichMessage struct {
	TemplateID string                `json:"template_id"`
	Elements   []*RichMessageElement `json:"elements"`
}

// ActionURL represents a URL condition within an action rule
type ActionURL struct {
	URL      string `json:"url,omitempty"`
	Operator string `json:"operator,omitempty"`
}

// ActionRule represents a rule for triggering a greeting
type ActionRule struct {
	ID           int64             `json:"id,omitempty"`
	Value        string            `json:"value,omitempty"`
	Type         string            `json:"type,omitempty"`
	Operator     string            `json:"operator,omitempty"`
	Condition    string            `json:"condition"`
	SessionField map[string]string `json:"session_field,omitempty"`
	URLs         []ActionURL       `json:"urls,omitempty"`
}

// Greeting represents a greeting in LiveChat
type Greeting struct {
	ID          int32             `json:"id"`
	Type        string            `json:"type,omitempty"`
	Active      bool              `json:"active"`
	Name        string            `json:"name"`
	Group       int32             `json:"group"`
	ActiveFrom  *string           `json:"active_from,omitempty"`
	ActiveUntil *string           `json:"active_until,omitempty"`
	Rules       []*ActionRule     `json:"rules,omitempty"`
	Properties  map[string]string `json:"properties,omitempty"`
	RichMessage *RichMessage      `json:"rich_message,omitempty"`
}

// CreateGreetingRequestOptions defines optional fields for CreateGreeting
type CreateGreetingRequestOptions struct {
	ActiveFrom  *string           `json:"active_from,omitempty"`
	ActiveUntil *string           `json:"active_until,omitempty"`
	Properties  map[string]string `json:"properties,omitempty"`
	RichMessage *RichMessage      `json:"rich_message,omitempty"`
}

// UpdateGreetingRequestOptions defines optional fields for UpdateGreeting
type UpdateGreetingRequestOptions struct {
	Type        *string           `json:"type,omitempty"`
	Active      *bool             `json:"active,omitempty"`
	ActiveFrom  *string           `json:"active_from,omitempty"`
	ActiveUntil *string           `json:"active_until,omitempty"`
	Name        *string           `json:"name,omitempty"`
	Rules       []*ActionRule     `json:"rules,omitempty"`
	Properties  map[string]string `json:"properties,omitempty"`
	RichMessage *RichMessage      `json:"rich_message,omitempty"`
}

// ListGreetingsRequestOptions defines optional fields for ListGreetings
type ListGreetingsRequestOptions struct {
	Groups []int32 `json:"groups,omitempty"`
}
