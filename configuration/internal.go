package configuration

type registerWebhookRequest struct {
	*Webhook
	OwnerClientID string `json:"owner_client_id,omitempty"`
}
type registerWebhookResponse struct {
	ID string `json:"id"`
}

type unregisterWebhookRequest struct {
	ID            string `json:"id"`
	OwnerClientID string `json:"owner_client_id,omitempty"`
}

type listWebhooksResponse []RegisteredWebhook

type createBotRequest struct {
	Name                 string         `json:"name"`
	Avatar               string         `json:"avatar,omitempty"`
	DefaultGroupPriority GroupPriority  `json:"default_group_priority,omitempty"`
	MaxChatsCount        *uint          `json:"max_chats_count,omitempty"`
	Groups               []*GroupConfig `json:"groups,omitempty"`
	Webhooks             *BotWebhooks   `json:"webhooks,omitempty"`
}

type createBotResponse struct {
	BotID string `json:"bot_agent_id"`
}

type deleteBotRequest struct {
	BotID string `json:"bot_agent_id"`
}

type updateBotRequest struct {
	BotID string `json:"id"`
	*createBotRequest
}

type listBotsRequest struct {
	All bool `json:"all"`
}

type listBotsResponse struct {
	BotAgents []*BotAgent `json:"bot_agents"`
}

type getBotRequest struct {
	BotID string `json:"bot_agent_id"`
}

type getBotResponse struct {
	BotAgent *BotAgentDetails `json:"bot_agent"`
}

type unregisterPropertyRequest struct {
	Name          string `json:"name"`
	OwnerClientID string `json:"owner_client_id"`
}

type publishPropertyRequest struct {
	Name          string   `json:"name"`
	OwnerClientID string   `json:"owner_client_id"`
	AccessType    []string `json:"access_type"`
}

type listPropertiesRequest struct {
	OwnerClientID string `json:"owner_client_id"`
}

type listPropertiesResponse map[string]*PropertyConfig

type getGroupRequest struct {
	ID     int      `json:"id"`
	Fields []string `json:"fields,omitempty"`
}

type getGroupResponse *Group

type emptyResponse struct{}

type listLicensePropertiesRequest struct {
	NamespacePrefix string `json:"namespace_prefix,omitempty"`
	NamePrefix      string `json:"name_prefix,omitempty"`
}

type listGroupPropertiesRequest struct {
	GroupID         uint   `json:"group_id"`
	NamespacePrefix string `json:"namespace_prefix,omitempty"`
	NamePrefix      string `json:"name_prefix,omitempty"`
}

type createAgentResponse struct {
	ID string `json:"id"`
}

type getAgentRequest struct {
	ID     string   `json:"id"`
	Fields []string `json:"fields,omitempty"`
}

type getAgentResponse *Agent

type listAgentsRequest struct {
	Filters AgentsFilters `json:"filters,omitempty"`
	Fields  []string      `json:"fields,omitempty"`
}

type listAgentsResponse []*Agent

type deleteAgentRequest struct {
	ID string `json:"id"`
}

type suspendAgentRequest struct {
	ID string `json:"id"`
}

type unsuspendAgentRequest struct {
	ID string `json:"id"`
}

type approveAgentRequest struct {
	ID string `json:"id"`
}

type createGroupRequest struct {
	Name            string                   `json:"name"`
	LanguageCode    string                   `json:"language_code,omitempty"`
	AgentPriorities map[string]GroupPriority `json:"agent_priorities"`
}

type createGroupResponse struct {
	ID int32 `json:"id"`
}

type updateGroupRequest struct {
	ID              int32                    `json:"id"`
	Name            string                   `json:"name,omitempty"`
	LanguageCode    string                   `json:"language_code,omitempty"`
	AgentPriorities map[string]GroupPriority `json:"agent_priorities,omitempty"`
}

type deleteGroupRequest struct {
	ID int32 `json:"id"`
}

type listGroupsRequest struct {
	Fields []string `json:"fields,omitempty"`
}

type listGroupsResponse []*Group

type listWebhookNamesRequest struct {
	Version string `json:"version,omitempty"`
}

type listWebhooksRequest struct {
	OwnerClientID string `json:"owner_client_id,omitempty"`
}

type manageWebhooksStateRequest struct {
	ClientID string `json:"client_id,omitempty"`
}
