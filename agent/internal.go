package agent

import (
	"encoding/json"

	"github.com/livechat/lc-sdk-go/objects"
)

type listChatsRequest struct {
	Filters    *chatsFilters      `json:"filters,omitempty"`
	Pagination *paginationRequest `json:"pagination,omitempty"`
}

type listChatsResponse struct {
	hashedPaginationResponse
	ChatsSummary []objects.ChatSummary `json:"chats_summary"`
	FoundChats   uint                  `json:"found_chats"`
}

type getChatThreadsSummaryRequest struct {
	ChatID string `json:"chat_id"`
	*hashedPaginationRequest
}

type getChatThreadsSummaryResponse struct {
	hashedPaginationResponse
	ThreadsSummary []objects.ThreadSummary `json:"threads_summary"`
	FoundThreads   uint                    `json:"found_threads"`
}

type getChatThreadsRequest struct {
	ChatID    string   `json:"chat_id"`
	ThreadIDs []string `json:"thread_ids,omitempty"`
}

type getChatThreadsResponse struct {
	Chat objects.Chat `json:"chat"`
}

type listArchivesRequest struct {
	Filters    *archivesFilters   `json:"filters,omitempty"`
	Pagination *paginationRequest `json:"pagination,omitempty"`
}

type listArchivesResponse struct {
	Chats      []objects.Chat     `json:"chats"`
	Pagination paginationResponse `json:"pagination"`
}

type startChatRequest struct {
	Chat       *InitialChat `json:"chat,omitempty"`
	Continuous bool         `json:"continuous,omitempty"`
}

type startChatResponse struct {
	ChatID   string   `json:"chat_id"`
	ThreadID string   `json:"thread_id"`
	EventIDs []string `json:"event_ids,omitempty"`
}

type activateChatRequest struct {
	Chat       *InitialChat `json:"chat"`
	Continuous bool         `json:"continuous,omitempty"`
}

type activateChatResponse struct {
	ThreadID string   `json:"thread_id"`
	EventIDs []string `json:"event_ids"`
}

type deactivateChatRequest struct {
	ChatID string `json:"chat_id"`
}

type followChatRequest struct {
	ChatID string `json:"chat_id"`
}

type unfollowChatRequest struct {
	ChatID string `json:"chat_id"`
}

// used to grant, revoke and set access
type modifyAccessRequest struct {
	Resource string         `json:"resource"`
	ID       string         `json:"id"`
	Access   objects.Access `json:"access"`
}

type transferChatRequest struct {
	ChatID string          `json:"chat_id"`
	Target *transferTarget `json:"target,omitempty"`
	Force  bool            `json:"force"`
}

// used to add and remove user from chat
type changeChatUsersRequest struct {
	ChatID   string `json:"chat_id"`
	UserID   string `json:"user_id"`
	UserType string `json:"user_type"` //todo - should be enum?
}

type sendEventRequest struct {
	ChatID             string      `json:"chat_id"`
	Event              interface{} `json:"event"`
	AttachToLastThread *bool       `json:"attach_to_last_thread,omitempty"`
}

type sendEventResponse struct {
	EventID string `json:"event_id"`
}

type sendRichMessagePostbackRequest struct {
	ChatID   string   `json:"chat_id"`
	EventID  string   `json:"event_id"`
	ThreadID string   `json:"thread_id"`
	Postback postback `json:"postback"`
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

// used for both tagging and untagging
type changeChatThreadTagRequest struct {
	ChatID   string `json:"chat_id"`
	ThreadID string `json:"thread_id"`
	Tag      string `json:"tag"`
}

type getCustomersRequest struct {
	CustomerID string `json:"customer_id"`
}

type getCustomersResponse struct {
}

type listCustomersRequest struct {
	PageID    string            `json:"page_id,omitempty"`
	Limit     uint              `json:"limit,omitempty"`
	SortOrder string            `json:"sort_order,omitempty"`
	Filters   *customersFilters `json:"filters,omitempty"`
}

type listCustomersResponse struct {
	hashedPaginationResponse
	Customers      []objects.Customer `json:"customers"`
	TotalCustomers uint               `json:"total_customers"`
}

type createCustomerRequest struct {
	Name          string              `json:"name,omitempty"`
	Email         string              `json:"email,omitempty"`
	Avatar        string              `json:"avatar,omitempty"`
	SessionFields []map[string]string `json:"session_fields,omitempty"`
}

type createCustomerResponse struct {
	CustomerID string `json:"customer_id"`
}

type updateCustomerRequest struct {
	CustomerID    string              `json:"customer_id"`
	Name          string              `json:"name,omitempty"`
	Email         string              `json:"email,omitempty"`
	Avatar        string              `json:"avatar,omitempty"`
	SessionFields []map[string]string `json:"session_fields,omitempty"`
}

type updateCustomerResponse struct {
	Customer objects.Customer `json:"customer"`
}

type banCustomerRequest struct {
	CustomerID string `json:"customer_id"`
	Ban        ban    `json:"ban"`
}

type updateAgentRequest struct {
	AgentID       string `json:"agent_id,omitempty"`
	RoutingStatus string `json:"routing_status,omitempty"`
}

type markEventsAsSeenRequest struct {
	ChatID   string `json:"chat_id"`
	SeenUpTo string `json:"seen_up_to"`
}

type sendTypingIndicatorRequest struct {
	ChatID     string `json:"chat_id"`
	Recipients string `json:"recipients,omitempty"`
	IsTyping   bool   `json:"is_typing"`
}

type multicastRequest struct {
	Recipients MulticastRecipients `json:"recipients"`
	Content    json.RawMessage     `json:"content"`
	Type       string              `json:"type,omitempty"`
}

type emptyResponse struct{}

type hashedPaginationRequest struct {
	PageID    string `json:"page_id,omitempty"`
	Limit     uint   `json:"limit,omitempty"`
	SortOrder string `json:"sort_order,omitempty"`
}

type hashedPaginationResponse struct {
	PreviousPageID string `json:"previous_page_id,omitempty"`
	NextPageID     string `json:"next_page_id,omitempty"`
}

type paginationRequest struct {
	Page  uint `json:"page,omitempty"`
	Limit uint `json:"limit,omitempty"`
}

type paginationResponse struct {
	Page  uint `json:"page,omitempty"`
	Total uint `json:"total,omitempty"`
}

type listAgentsForTransferRequest struct {
	ChatID string `json:"chat_id"`
}
