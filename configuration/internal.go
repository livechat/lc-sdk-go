package configuration

type registerWebhookResponse struct {
	ID string `json:"webhook_id"`
}

type unregisterWebhookRequest struct {
	ID string `json:"webhook_id"`
}

type getWebhookConfigResponse []RegisteredWebhook

type createBotAgentRequest struct {
	Name                 string           `json:"name"`
	Status               BotStatus        `json:"status"`
	Avatar               string           `json:"string"`
	DefaultGroupPriority GroupPriority    `json:"default_group_priority"`
	MaxChatsCount        uint             `json:"max_chats_count"`
	Groups               []BotGroupConfig `json:"groups"`
	Webhooks             BotWebhooks      `json:"webhooks"`
}

type createBotAgentResponse struct {
	BotID string `json:"bot_agent_id"`
}

type removeBotAgentRequest struct {
	BotID string `json:"bot_agent_id"`
}

type emptyResponse struct{}

type updateBotAgentRequest struct {
	BotID string `json:"id"`
	*createBotAgentRequest
}

type getBotAgentsRequest struct {
	All bool `json:"all"`
}

type getBotAgentsResponse struct {
	BotAgents []BotAgent `json:"bot_agents"`
}

type getBotAgentDetailsRequest struct {
	BotID string `json:"bot_agent_id"`
}

type getBotAgentDetailsResponse struct {
	BotAgent BotAgentDetails `json:"bot_agent"`
}

type createPropertiesRequest map[string]PropertyConfig

type getPropertyConfigsRequest struct {
	All bool `json:"all"`
}

type getPropertyConfigsResponse map[string]PropertyConfig
