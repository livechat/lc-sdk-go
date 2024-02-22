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
	Name string `json:"name"`
	CreateBotRequestOptions
}

type deleteBotRequest struct {
	BotID string `json:"id"`
}

type updateBotRequest struct {
	BotID string `json:"id"`
	UpdateBotRequestOptions
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

type createBotTemplateRequest struct {
	Name string `json:"name"`
	CreateBotTemplateRequestOptions
}

type updateBotTemplateRequest struct {
	BotTemplateID string `json:"id"`
	UpdateBotTemplateRequestOptions
}

type deleteBotTemplateRequest struct {
	BotTemplateID               string `json:"id"`
	OwnerClientID               string `json:"owner_client_id,omitempty"`
	AffectExistingInstallations bool   `json:"affect_existing_installations"`
}

type listBotTemplatesRequest struct {
	OwnerClientID string `json:"owner_client_id,omitempty"`
}

type ListBotTemplatesResponse []BotTemplate

type resetBotSecretRequest struct {
	BotID         string `json:"id"`
	OwnerClientID string `json:"owner_client_id,omitempty"`
}

type resetBotSecretResponse struct {
	Secret string `json:"secret"`
}

type issueBotTokenResponse struct {
	Token string `json:"token"`
}

type resetBotTemplateSecretRequest struct {
	BotTemplateID               string `json:"id"`
	OwnerClientID               string `json:"owner_client_id,omitempty"`
	AffectExistingInstallations bool   `json:"affect_existing_installations"`
}

type resetBotTemplateSecretResponse struct {
	Secret string `json:"secret"`
}

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
	ListLicensePropertiesRequestOptions
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
	Filters *AgentsFilters `json:"filters,omitempty"`
	Fields  []string       `json:"fields,omitempty"`
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
	AgentPriorities map[string]GroupPriority `json:"agent_priorities"`
	CreateGroupRequestOptions
}

type createGroupResponse struct {
	ID int32 `json:"id"`
}

type updateGroupRequest struct {
	ID int32 `json:"id"`
	UpdateGroupRequestOptions
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
	Access     Access               `json:"access"`
	Conditions AutoAccessConditions `json:"conditions"`
	AddAutoAccessRequestOptions
}

type addAutoAccessResponse struct {
	ID string `json:"id"`
}

type updateAutoAccessRequest struct {
	ID string `json:"id"`
	UpdateAutoAccessRequestOptions
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
	GroupIDs []int `json:"group_ids"`
	ListGroupsPropertiesRequestOptions
}

type reactivateEmailRequest struct {
	AgentID string `json:"agent_id"`
}

type updateCompanyDetailsRequest struct {
	CompanyDetails
	Enrich bool `json:"enrich"`
}
