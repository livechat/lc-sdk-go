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
	Avatar               string         `json:"avatar_path,omitempty"`
	DefaultGroupPriority GroupPriority  `json:"default_group_priority,omitempty"`
	MaxChatsCount        *uint          `json:"max_chats_count,omitempty"`
	Groups               []*GroupConfig `json:"groups,omitempty"`
	OwnerClientID        string         `json:"owner_client_id,omitempty"`
	WorkScheduler        WorkScheduler  `json:"work_scheduler,omitempty"`
	Timezone             string         `json:"timezone,omitempty"`
}

type createBotResponse struct {
	BotID string `json:"id"`
}

type deleteBotRequest struct {
	BotID string `json:"id"`
}

type updateBotRequest struct {
	BotID string `json:"id"`
	*createBotRequest
}

type listBotsRequest struct {
	All    bool     `json:"all"`
	Fields []string `json:"fields,omitempty"`
}

type listBotsResponse []*Bot

type getBotRequest struct {
	BotID  string   `json:"id"`
	Fields []string `json:"fields,omitempty"`
}

type getBotResponse *Bot

type unregisterPropertyRequest struct {
	Name          string `json:"name"`
	OwnerClientID string `json:"owner_client_id,omitempty"`
}

type publishPropertyRequest struct {
	Name          string   `json:"name"`
	OwnerClientID string   `json:"owner_client_id,omitempty"`
	AccessType    []string `json:"access_type"`
}

type listPropertiesRequest struct {
	OwnerClientID string `json:"owner_client_id,omitempty"`
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
	ID              uint   `json:"id"`
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
	OwnerClientID string `json:"owner_client_id,omitempty"`
}

type updateLicensePropertiesRequest struct {
	Properties Properties `json:"properties"`
}

type updateGroupPropertiesRequest struct {
	ID         int        `json:"id"`
	Properties Properties `json:"properties"`
}

type deleteLicensePropertiesRequest struct {
	Properties map[string][]string `json:"properties"`
}

type deleteGroupPropertiesRequest struct {
	ID         int                 `json:"id"`
	Properties map[string][]string `json:"properties"`
}

type addAutoAccessRequest struct {
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

type addAutoAccessResponse struct {
	ID string `json:"id"`
}

type updateAutoAccessRequest struct {
	addAutoAccessRequest
	ID string `json:"id"`
}

type deleteAutoAccessRequest struct {
	ID string `json:"id"`
}

type listAutoAccessesRequest struct {
}

type checkProductLimitsForPlanRequest struct {
	Plan string `json:"plan"`
}

type listChannelsRequest struct {
}

type createTagRequest struct {
	Name     string `json:"name"`
	GroupIDs []int  `json:"group_ids"`
}

type deleteTagRequest struct {
	Name string `json:"name"`
}

type listTagsRequest struct {
	GroupIDs []int `json:"group_ids"`
}

type updateTagRequest struct {
	Name     string `json:"name"`
	GroupIDs []int  `json:"group_ids"`
}

type listGroupsPropertiesRequest struct {
	GroupIDs   []int  `json:"group_ids,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
	NamePrefix string `json:"name_prefix,omitempty"`
}
