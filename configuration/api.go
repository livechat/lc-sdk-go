package configuration

import (
	"fmt"
	"net/http"

	"github.com/livechat/lc-sdk-go/authorization"
	i "github.com/livechat/lc-sdk-go/internal"
	"github.com/livechat/lc-sdk-go/objects"
)

// API provides the API operation methods for making requests to Livechat Configuration API via Web API.
// See this package's package overview docs for details on the service.
type API struct {
	*i.API
}

// NewAPI returns ready to use Configuration API.
//
// If provided client is nil, then default http client with 20s timeout is used.
func NewAPI(t authorization.TokenGetter, client *http.Client, clientID string) (*API, error) {
	api, err := i.NewAPI(t, client, clientID, "configuration")
	if err != nil {
		return nil, err
	}
	return &API{api}, nil
}

// RegisterWebhook allows to register specified webhook.
func (a *API) RegisterWebhook(webhook *Webhook) (string, error) {
	var resp registerWebhookResponse
	err := a.Call("register_webhook", webhook, &resp)

	return resp.ID, err
}

// ListRegisteredWebhooks returns configurations of all registered webhooks.
func (a *API) ListRegisteredWebhooks() ([]RegisteredWebhook, error) {
	var resp listRegisteredWebhooksResponse
	err := a.Call("list_registered_webhooks", nil, &resp)

	return resp, err
}

// UnregisterWebhook removes webhook with given id from registered webhooks.
func (a *API) UnregisterWebhook(id string) error {
	return a.Call("unregister_webhook", unregisterWebhookRequest{
		ID: id,
	}, &emptyResponse{})
}

// CreateBot allows to create bot and returns its ID.
func (a *API) CreateBot(name, avatar string, status BotStatus, maxChats uint, defaultPriority GroupPriority, groups []*GroupConfig, webhooks *BotWebhooks) (string, error) {
	var resp createBotResponse
	if err := validateBotGroupsAssignment(groups); err != nil {
		return "", err
	}
	err := a.Call("create_bot", &createBotRequest{
		Name:                 name,
		Avatar:               avatar,
		Status:               status,
		MaxChatsCount:        &maxChats,
		DefaultGroupPriority: defaultPriority,
		Groups:               groups,
		Webhooks:             webhooks,
	}, &resp)

	return resp.BotID, err
}

// UpdateBot allows to update bot.
func (a *API) UpdateBot(id, name, avatar string, status BotStatus, maxChats uint, defaultPriority GroupPriority, groups []*GroupConfig, webhooks *BotWebhooks) error {
	if err := validateBotGroupsAssignment(groups); err != nil {
		return err
	}
	return a.Call("update_bot", &updateBotRequest{
		BotID: id,
		createBotRequest: &createBotRequest{
			Name:                 name,
			Avatar:               avatar,
			Status:               status,
			MaxChatsCount:        &maxChats,
			DefaultGroupPriority: defaultPriority,
			Groups:               groups,
			Webhooks:             webhooks,
		},
	}, &emptyResponse{})
}

// DeleteBot deletes bot with given ID.
func (a *API) DeleteBot(id string) error {
	return a.Call("delete_bot", &deleteBotRequest{
		BotID: id,
	}, &emptyResponse{})
}

// ListBots returns list of bots (all or caller's only, depending on getAll parameter).
func (a *API) ListBots(getAll bool) ([]*BotAgent, error) {
	var resp listBotsResponse
	err := a.Call("list_bots", &listBotsRequest{
		All: getAll,
	}, &resp)

	return resp.BotAgents, err
}

// GetBot returns bot.
func (a *API) GetBot(id string) (*BotAgentDetails, error) {
	var resp getBotResponse
	err := a.Call("get_bot", &getBotRequest{
		BotID: id,
	}, &resp)

	return resp.BotAgent, err
}

// CreateAgent creates a new Agent with specified parameters within a license.
func (a *API) CreateAgent(id string, fields *AgentFields) (string, error) {
	var resp createAgentResponse
	request := &Agent{
		ID:          id,
		AgentFields: fields,
	}
	err := a.Call("create_agent", request, &resp)

	return resp.ID, err
}

// GetAgent returns the info about an Agent specified by id (i.e. login).
func (a *API) GetAgent(id string, fields []string) (*Agent, error) {
	var resp getAgentResponse
	err := a.Call("get_agent", &getAgentRequest{
		ID:     id,
		Fields: fields,
	}, &resp)

	return resp, err
}

// ListAgents returns all Agents within a license.
func (a *API) ListAgents(groupIDs []int32, fields []string) ([]*Agent, error) {
	var resp listAgentsResponse
	request := &listAgentsRequest{
		Fields: fields,
	}

	if len(groupIDs) > 0 {
		request.Filters = AgentsFilters{
			GroupIDs: groupIDs,
		}
	}

	err := a.Call("list_agents", request, &resp)
	return resp, err
}

// UpdateAgent updates the properties of an Agent specified by id.
func (a *API) UpdateAgent(id string, fields *AgentFields) error {
	request := &Agent{
		ID:          id,
		AgentFields: fields,
	}
	return a.Call("update_agent", request, &emptyResponse{})
}

// DeleteAgent deletes an Agent specified by id.
func (a *API) DeleteAgent(id string) error {
	return a.Call("delete_agent", &deleteAgentRequest{
		ID: id,
	}, &emptyResponse{})
}

// SuspendAgent suspends an Agent specified by id.
func (a *API) SuspendAgent(id string) error {
	return a.Call("suspend_agent", &suspendAgentRequest{
		ID: id,
	}, &emptyResponse{})
}

// UnsuspendAgent unsuspends an Agent specified by id.
func (a *API) UnsuspendAgent(id string) error {
	return a.Call("unsuspend_agent", &unsuspendAgentRequest{
		ID: id,
	}, &emptyResponse{})
}

// RequestAgentUnsuspension sends a request to license owners and vice owners with an unsuspension request
func (a *API) RequestAgentUnsuspension() error {
	return a.Call("request_agent_unsuspension", nil, &emptyResponse{})
}

// ApproveAgent approves an Agent thus allowing the Agent to use the application.
func (a *API) ApproveAgent(id string) error {
	return a.Call("approve_agent", &approveAgentRequest{
		ID: id,
	}, &emptyResponse{})
}

// RegisterProperties allows to create properties
func (a *API) RegisterProperties(properties map[string]*PropertyConfig) error {
	return a.Call("register_properties", properties, &emptyResponse{})
}

// ListRegisteredProperties return list of properties along with their configuration
func (a *API) ListRegisteredProperties(getAll bool) (map[string]*PropertyConfig, error) {
	var resp listRegisteredPropertiesResponse
	err := a.Call("list_registered_properties", &listRegisteredPropertiesRequest{
		All: getAll,
	}, &resp)

	return resp, err
}

// CreateGroup creates new group
func (a *API) CreateGroup(name, language string, agentPriorities map[string]GroupPriority) (int32, error) {
	var resp createGroupResponse
	err := a.Call("create_group", &createGroupRequest{
		Name:            name,
		LanguageCode:    language,
		AgentPriorities: agentPriorities,
	}, &resp)

	return resp.ID, err
}

// UpdateGroup updates existing group
func (a *API) UpdateGroup(id int32, name, language string, agentPriorities map[string]GroupPriority) error {
	return a.Call("update_group", &updateGroupRequest{
		ID:              id,
		Name:            name,
		LanguageCode:    language,
		AgentPriorities: agentPriorities,
	}, &emptyResponse{})
}

// DeleteGroup deletes existing group
func (a *API) DeleteGroup(id int32) error {
	return a.Call("delete_group", &deleteGroupRequest{
		ID: id,
	}, &emptyResponse{})
}

// ListGroups lists all existing groups
func (a *API) ListGroups(fields []string) ([]*Group, error) {
	var resp listGroupsResponse
	err := a.Call("list_groups", &listGroupsRequest{
		Fields: fields,
	}, &resp)

	return resp, err
}

// GetGroup returns details about a group specified by its id
func (a *API) GetGroup(id int) (*Group, error) {
	var resp getGroupResponse
	err := a.Call("get_group", &getGroupRequest{
		ID: id,
	}, &resp)

	return resp, err
}

func validateBotGroupsAssignment(groups []*GroupConfig) error {
	for _, group := range groups {
		if group.Priority == DoNotAssign {
			return fmt.Errorf("DoNotAssign priority is allowed only as default group priority")
		}
	}

	return nil
}

// ListLicenseProperties returns the properties set within a license.
func (a *API) ListLicenseProperties(namespacePrefix, namePrefix string) (objects.Properties, error) {
	var resp objects.Properties
	err := a.Call("list_license_properties", &listLicensePropertiesRequest{
		NamespacePrefix: namespacePrefix,
		NamePrefix:      namePrefix,
	}, &resp)
	return resp, err
}

// ListGroupProperties returns the properties set within a group.
func (a *API) ListGroupProperties(groupID uint, namespacePrefix, namePrefix string) (objects.Properties, error) {
	var resp objects.Properties
	err := a.Call("list_group_properties", &listGroupPropertiesRequest{
		GroupID:         groupID,
		NamespacePrefix: namespacePrefix,
		NamePrefix:      namePrefix,
	}, &resp)
	return resp, err
}
