package configuration

import (
	"fmt"
	"net/http"

	"github.com/livechat/lc-sdk-go/authorization"
	i "github.com/livechat/lc-sdk-go/internal"
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

// RegisterWebhook allows to register specified webhook
func (a *API) RegisterWebhook(webhook *Webhook) (string, error) {
	var resp registerWebhookResponse
	err := a.Call("register_webhook", webhook, &resp)

	return resp.ID, err
}

// ListRegisteredWebhooks returns configurations of all registered webhooks
func (a *API) ListRegisteredWebhooks() ([]RegisteredWebhook, error) {
	var resp listRegisteredWebhooksResponse
	err := a.Call("list_registered_webhooks", nil, &resp)

	return resp, err
}

// UnregisterWebhook removes webhook with given id from registered webhooks
func (a *API) UnregisterWebhook(id string) error {
	return a.Call("unregister_webhook", unregisterWebhookRequest{
		ID: id,
	}, &emptyResponse{})
}

// CreateBot allows to create bot and returns its ID
func (a *API) CreateBot(name, avatar string, status BotStatus, maxChats uint, defaultPriority GroupPriority, groups []*BotGroupConfig, webhooks *BotWebhooks) (string, error) {
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

// UpdateBot allows to update bot
func (a *API) UpdateBot(id, name, avatar string, status BotStatus, maxChats uint, defaultPriority GroupPriority, groups []*BotGroupConfig, webhooks *BotWebhooks) error {
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

// RemoveBot removes bot with given ID
func (a *API) RemoveBot(id string) error {
	return a.Call("remove_bot", &removeBotRequest{
		BotID: id,
	}, &emptyResponse{})
}

// ListBots returns list of bots (all or caller's only, depending on getAll parameter)
func (a *API) ListBots(getAll bool) ([]*BotAgent, error) {
	var resp listBotsResponse
	err := a.Call("list_bots", &listBotsRequest{
		All: getAll,
	}, &resp)

	return resp.BotAgents, err
}

// GetBot returns bot
func (a *API) GetBot(id string) (*BotAgentDetails, error) {
	var resp getBotResponse
	err := a.Call("get_bot", &getBotRequest{
		BotID: id,
	}, &resp)

	return resp.BotAgent, err
}

// CreateProperties allows to create properties
func (a *API) CreateProperties(properties map[string]*PropertyConfig) error {
	return a.Call("create_properties", properties, &emptyResponse{})
}

// GetPropertyConfigs return list of properties along with their configuration
func (a *API) GetPropertyConfigs(getAll bool) (map[string]*PropertyConfig, error) {
	var resp getPropertyConfigsResponse
	err := a.Call("get_property_configs", &getPropertyConfigsRequest{
		All: getAll,
	}, &resp)

	return resp, err
}

func (a *API) GetGroup(id int) (*Group, error) {
	var resp getGroupResponse
	err := a.Call("get_group", &getGroupRequest{
		ID: id,
	}, &resp)

	return resp, err
}

func validateBotGroupsAssignment(groups []*BotGroupConfig) error {
	for _, group := range groups {
		if group.Priority == DoNotAssign {
			return fmt.Errorf("DoNotAssign priority is allowed only as default group priority")
		}
	}

	return nil
}
