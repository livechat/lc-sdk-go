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

type listChatsRequest struct {
	*hashedPaginationRequest
}

type listChatsResponse struct {
	hashedPaginationResponse
	ChatsSummary []objects.ChatSummary `json:"chats_summary"`
	TotalChats   uint                  `json:"total_chats"`
}

type getChatRequest struct {
	ChatID   string `json:"chat_id"`
	ThreadID string `json:"thread_id,omitempty"`
}

type getChatResponse struct {
	Chat objects.Chat `json:"chat"`
}

type listThreadsRequest struct {
	*hashedPaginationRequest
	ChatID         string `json:"chat_id"`
	MinEventsCount uint   `json:"min_events_count,omitempty"`
}

type listThreadsResponse struct {
	hashedPaginationResponse
	Threads      []objects.Thread `json:"threads"`
	FoundThreads uint             `json:"found_threads"`
}

type deactivateChatRequest struct {
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

type updateThreadPropertiesRequest struct {
	ChatID     string             `json:"chat_id"`
	ThreadID   string             `json:"thread_id"`
	Properties objects.Properties `json:"properties"`
}

type deleteThreadPropertiesRequest struct {
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
	Name          string              `json:"name,omitempty"`
	Email         string              `json:"email,omitempty"`
	Avatar        string              `json:"avatar,omitempty"`
	SessionFields []map[string]string `json:"session_fields,omitempty"`
}

type setCustomerSessionFieldsRequest struct {
	SessionFields []map[string]string `json:"session_fields"`
}

type listGroupStatusesRequest struct {
	All      bool  `json:"all,omitempty"`
	GroupIDs []int `json:"group_ids,omitempty"`
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

type listLicensePropertiesRequest struct {
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
}

type listGroupPropertiesRequest struct {
	GroupID   uint   `json:"group_id"`
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
}

type acceptGreetingRequest struct {
	GreetingID int    `json:"greeting_id"`
	UniqueID   string `json:"unique_id"`
}

type cancelGreetingRequest struct {
	UniqueID string `json:"unique_id"`
}

type hashedPaginationRequest struct {
	PageID    string `json:"page_id,omitempty"`
	Limit     uint   `json:"limit,omitempty"`
	SortOrder string `json:"sort_order,omitempty"`
}

type hashedPaginationResponse struct {
	PreviousPageID string `json:"previous_page_id,omitempty"`
	NextPageID     string `json:"next_page_id,omitempty"`
}
