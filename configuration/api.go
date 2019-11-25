package configuration

import (
	"fmt"
	"net/http"

	i "github.com/livechat/lc-sdk-go/internal"
)

type ConfigurationAPI struct {
	*i.API
}

func NewAPI(t i.TokenGetter, client *http.Client, clientID string) (*ConfigurationAPI, error) {
	api, err := i.NewAPI(t, client, clientID, "configuration")
	if err != nil {
		return nil, err
	}
	return &ConfigurationAPI{api}, nil
}

func (a *ConfigurationAPI) RegisterWebhook(webhook *Webhook) (string, error) {
	var resp registerWebhookResponse
	err := a.Call("register_webhook", webhook, &resp)

	return resp.ID, err
}

func (a *ConfigurationAPI) GetWebhooksConfig() ([]RegisteredWebhook, error) {
	var resp getWebhookConfigResponse
	err := a.Call("get_webhooks_config", nil, &resp)

	return resp, err
}

func (a *ConfigurationAPI) UnregisterWebhook(id string) error {
	return a.Call("unregister_webhook", unregisterWebhookRequest{
		ID: id,
	}, &emptyResponse{})
}

func (a *ConfigurationAPI) CreateBotAgent(name, avatar string, status BotStatus, maxChats uint, defaultPriority GroupPriority, groups []BotGroupConfig, webhooks BotWebhooks) (string, error) {
	var resp createBotAgentResponse
	if err := validateBotGroupsAssignment(groups); err != nil {
		return "", err
	}
	err := a.Call("create_bot_agent", &createBotAgentRequest{
		Name:                 name,
		Avatar:               avatar,
		Status:               status,
		MaxChatsCount:        maxChats,
		DefaultGroupPriority: defaultPriority,
		Groups:               groups,
		Webhooks:             webhooks,
	}, &resp)

	return resp.BotID, err
}

func (a *ConfigurationAPI) UpdateBotAgent(id, name, avatar string, status BotStatus, maxChats uint, defaultPriority GroupPriority, groups []BotGroupConfig, webhooks BotWebhooks) error {
	if err := validateBotGroupsAssignment(groups); err != nil {
		return err
	}
	return a.Call("update_bot_agent", &updateBotAgentRequest{
		BotID: id,
		createBotAgentRequest: &createBotAgentRequest{
			Name:                 name,
			Avatar:               avatar,
			Status:               status,
			MaxChatsCount:        maxChats,
			DefaultGroupPriority: defaultPriority,
			Groups:               groups,
			Webhooks:             webhooks,
		},
	}, &emptyResponse{})
}

func (a *ConfigurationAPI) RemoveBotAgent(id string) error {
	return a.Call("remove_bot_agent", &removeBotAgentRequest{
		BotID: id,
	}, &emptyResponse{})
}

func (a *ConfigurationAPI) GetBotAgents(getAll bool) ([]BotAgent, error) {
	var resp getBotAgentsResponse
	err := a.Call("get_bot_agents", &getBotAgentsRequest{
		All: getAll,
	}, &resp)

	return resp.BotAgents, err
}

func (a *ConfigurationAPI) GetBotAgentDetails(id string) (BotAgentDetails, error) {
	var resp getBotAgentDetailsResponse
	err := a.Call("get_bot_agent_details", &getBotAgentDetailsRequest{
		BotID: id,
	}, &resp)

	return resp.BotAgent, err
}

func (a *ConfigurationAPI) CreateProperties(properties *map[string]PropertyConfig) error {
	return a.Call("create_properties", properties, &emptyResponse{})
}

func (a *ConfigurationAPI) GetPropertyConfigs(getAll bool) (map[string]PropertyConfig, error) {
	var resp getPropertyConfigsResponse
	err := a.Call("get_property_configs", &getPropertyConfigsRequest{
		All: getAll,
	}, &resp)

	return resp, err
}

func (a *ConfigurationAPI) SetAPIURL(url string) {
	a.APIURL = url
}

func validateBotGroupsAssignment(groups []BotGroupConfig) error {
	for _, group := range groups {
		if group.Priority == DoNotAssign {
			return fmt.Errorf("DoNotAssing priority is allowed only as default group priority")
		}
	}

	return nil
}
