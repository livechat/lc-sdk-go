package customer

import "github.com/livechat/lc-sdk-go/objects"

type startChatRequest struct {
	Chat       *objects.InitialChat `json:"chat,omitempty"`
	Continuous bool                 `json:"continuous,omitempty"`
}

type startChatResponse struct {
	ChatID   string   `json:"chat_id"`
	ThreadID string   `json:"thread_id"`
	EventIDs []string `json:"event_ids"`
}

type sendEventRequest struct {
	ChatID             string      `json:"chat_id"`
	Event              interface{} `json:"event"`
	AttachToLastThread *bool       `json:"attach_to_last_thread,omitempty"`
}

type sendEventResponse struct {
	EventID string `json:"event_id"`
}

type activateChatRequest struct {
	Chat       *objects.InitialChat `json:"chat"`
	Continuous bool                 `json:"continuous,omitempty"`
}

type activateChatResponse struct {
	ThreadID string   `json:"thread_id"`
	EventIDs []string `json:"event_ids"`
}

type getChatsSummaryRequest struct {
	Limit  uint `json:"limit,omitempty"`
	Offset uint `json:"offset,omitempty"`
}

type getChatsSummaryResponse struct {
	Chats      []objects.Chat `json:"chats_summary"`
	TotalChats uint           `json:"total_chats"`
}

type getChatThreadsSummaryRequest struct {
	ChatID string `json:"chat_id"`
	Limit  uint   `json:"limit,omitempty"`
	Offset uint   `json:"offset,omitempty"`
}

type getChatThreadsSummaryResponse struct {
	ThreadsSummary []objects.ThreadSummary `json:"threads_summary"`
	TotalThreads   uint                    `json:"total_threads"`
}

type getChatThreadsRequest struct {
	ChatID    string   `json:"chat_id"`
	ThreadIDs []string `json:"thread_ids,omitempty"`
}

type getChatThreadsResponse struct {
	Chat objects.Chat `json:"chat"`
}

type closeThreadRequest struct {
	ChatID string `json:"chat_id"`
}

type sendRichMessagePostbackRequest struct {
	ChatID   string   `json:"chat_id"`
	ThreadID string   `json:"thread_id"`
	EventID  string   `json:"event_id"`
	Postback postback `json:"postback"`
}

type postback struct {
	ID      string `json:"id"`
	Toggled bool   `json:"toggled"`
}

type sendSneakPeekRequest struct {
	ChatID        string `json:"chat_id"`
	SneakPeekText string `json:"sneak_peek_text"`
}

type updateChatPropertiesRequest struct {
	ChatID     string             `json:"chat_id"`
	Properties objects.Properties `json:"properties"`
}

type deleteChatPropertiesRequest struct {
	ChatID     string              `json:"chat_id"`
	Properties map[string][]string `json:"properties"`
}

type updateChatThreadPropertiesRequest struct {
	ChatID     string             `json:"chat_id"`
	ThreadID   string             `json:"thread_id"`
	Properties objects.Properties `json:"properties"`
}

type deleteChatThreadPropertiesRequest struct {
	ChatID     string              `json:"chat_id"`
	ThreadID   string              `json:"thread_id"`
	Properties map[string][]string `json:"properties"`
}

type updateEventPropertiesRequest struct {
	ChatID     string             `json:"chat_id"`
	ThreadID   string             `json:"thread_id"`
	EventID    string             `json:"event_id"`
	Properties objects.Properties `json:"properties"`
}

type deleteEventPropertiesRequest struct {
	ChatID     string              `json:"chat_id"`
	ThreadID   string              `json:"thread_id"`
	EventID    string              `json:"event_id"`
	Properties map[string][]string `json:"properties"`
}

type updateCustomerRequest struct {
	Name   string            `json:"name,omitempty"`
	Email  string            `json:"email,omitempty"`
	Avatar string            `json:"avatar,omitempty"`
	Fields map[string]string `json:"fields,omitempty"`
}

type setCustomerFieldsRequest struct {
	Fields map[string]string `json:"fields"`
}

type listGroupStatusesRequest struct {
	All    bool  `json:"all,omitempty"`
	Groups []int `json:"groups,omitempty"`
}

type listGroupStatusesResponse struct {
	Status map[int]string `json:"groups_status"`
}

type checkGoalsRequest struct {
	PageURL        string            `json:"page_url"`
	GroupID        int               `json:"group_id"`
	CustomerFields map[string]string `json:"customer_fields"`
}

type getFormRequest struct {
	GroupID int    `json:"group_id"`
	Type    string `json:"type"`
}

type getFormResponse struct {
	Form    *Form `json:"form"`
	Enabled bool  `json:"enabled"`
}

type getURLInfoRequest struct {
	URL string `json:"url"`
}

type markEventsAsSeenRequest struct {
	ChatID   string `json:"chat_id"`
	SeenUpTo string `json:"seen_up_to"`
}

type emptyResponse struct{}
