package configuration

type registerWebhookResponse struct {
	ID string `json:"webhook_id"`
}

type unregisterWebhookRequest struct {
	ID string `json:"webhook_id"`
}

type listRegisteredWebhooksResponse []RegisteredWebhook

type createBotRequest struct {
	Name                 string            `json:"name"`
	Status               BotStatus         `json:"status"`
	Avatar               string            `json:"string,omitempty"`
	DefaultGroupPriority GroupPriority     `json:"default_group_priority,omitempty"`
	MaxChatsCount        *uint             `json:"max_chats_count,omitempty"`
	Groups               []*BotGroupConfig `json:"groups,omitempty"`
	Webhooks             *BotWebhooks      `json:"webhooks,omitempty"`
}

type createBotResponse struct {
	BotID string `json:"bot_agent_id"`
}

type removeBotRequest struct {
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

type getBotAgentDetailsRequest struct {
	BotID string `json:"bot_agent_id"`
}

type getBotAgentDetailsResponse struct {
	BotAgent *BotAgentDetails `json:"bot_agent"`
}

type createPropertiesRequest map[string]*PropertyConfig

type getPropertyConfigsRequest struct {
	All bool `json:"all"`
}

type getPropertyConfigsResponse map[string]*PropertyConfig

type getGroupRequest struct {
	ID int `json:"id"`
}

type getGroupResponse *Group

type emptyResponse struct{}
