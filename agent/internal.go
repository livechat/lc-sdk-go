package agent

import (
	"encoding/json"

	"github.com/livechat/lc-sdk-go/objects"
)

type getChatsSummaryRequest struct {
	Filters    *ChatsFilters     `json:"filters"`
	Pagination paginationRequest `json:"pagination"`
}

type getChatsSummaryResponse struct {
	ChatsSummary []objects.ChatSummary `json:"chats_summary"`
	FoundChats   uint                  `json:"found_chats"`
	*hashedPaginationResponse
}

type getChatThreadsSummaryRequest struct {
	ChatID string `json:"chat_id"`
	*hashedPaginationRequest
}

type getChatThreadsSummaryResponse struct {
	*hashedPaginationResponse
	ThreadsSummary []objects.ThreadSummary `json:"threads_summary"`
	FoundThreads   uint                    `json:"found_threads"`
}

type getChatThreadsRequest struct {
	ChatID    string   `json:"chat_id"`
	ThreadIDs []string `json:"thread_ids"`
}

type getChatThreadsResponse struct {
	Chat objects.Chat `json:"chat"`
}

type getArchivesRequest struct {
	Filters    *ArchivesFilters  `json:"filters"`
	Pagination paginationRequest `json:"pagination"`
}

type getArchivesResponse struct {
	Chats      []objects.Chat     `json:"chats"`
	Pagination paginationResponse `json:"pagination"`
}

type startChatRequest struct {
	Chat       *objects.InitialChat `json:"chat"`
	Continuous bool                 `json:"continuous"`
}

type startChatResponse struct {
	ChatID   string   `json:"chat_id"`
	ThreadID string   `json:"thread_id"`
	EventIDs []string `json:"event_ids,omitempty"`
}

type activateChatRequest struct {
	Chat       *objects.InitialChat `json:"chat"`
	Continuous bool                 `json:"continuous"`
}

type activateChatResponse struct {
	ThreadID string   `json:"thread_id"`
	EventIDs []string `json:"event_ids"`
}

type closeThreadRequest struct {
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
	ChatID string         `json:"chat_id"`
	Target TransferTarget `json:"target"`
	Force  bool           `json:"force"`
}

// used to add and remove user from chat
type changeChatUsersRequest struct {
	ChatID   string `json:"chat_id"`
	UserID   string `json:"user_id"`
	UserType string `json:"user_type"` //todo - should be enum?
}

type sendEventRequest struct {
	ChatID             string        `json:"chat_id"`
	Event              objects.Event `json:"event"`
	AttachToLastThread bool          `json:"attach_to_last_thread"`
}

type sendEventResponse struct {
	EventID string `json:"event_id"`
}

type sendRichMessagePostbackRequest struct {
	ChatID   string   `json:"chat_id"`
	EventID  string   `json:"event_id"`
	ThreadID string   `json:"thread_id"`
	Postback Postback `json:"postback"`
}

type updateChatPropertiesRequest struct {
	ChatID     string             `json:"chat_id"`
	Properties objects.Properties `json:"properties"`
}

type deleteChatPropertiesRequest struct {
	ChatID     string             `json:"chat_id"`
	Properties objects.Properties `json:"properties"`
}

type updateChatThreadPropertiesRequest struct {
	ChatID     string             `json:"chat_id"`
	ThreadID   string             `json:"thread_id"`
	Properties objects.Properties `json:"properties"`
}

type deleteChatThreadPropertiesRequest struct {
	ChatID     string             `json:"chat_id"`
	ThreadID   string             `json:"thread_id"`
	Properties objects.Properties `json:"properties"`
}

type updateEventPropertiesRequest struct {
	ChatID     string             `json:"chat_id"`
	ThreadID   string             `json:"thread_id"`
	EventID    string             `json:"event_id"`
	Properties objects.Properties `json:"properties"`
}

type deleteEventPropertiesRequest struct {
	ChatID     string             `json:"chat_id"`
	ThreadID   string             `json:"thread_id"`
	EventID    string             `json:"event_id"`
	Properties objects.Properties `json:"properties"`
}

// used for both tagging and untagging
type changeChatThreadTagRequest struct {
	ChatID   string `json:"chat_id"`
	ThreadID string `json:"thread_id"`
	Tag      string `json:"tag"`
}

type getCustomersRequest struct {
	PageID  string            `json:"page_id"`
	Limit   uint              `json:"limit"`
	Order   string            `json:"order"`
	Filters *CustomersFilters `json:"filters"`
}

type getCustomersResponse struct {
	*hashedPaginationResponse
	Customers      []objects.Customer `json:"customers"`
	TotalCustomers uint               `json:"total_customers"`
}

type createCustomerRequest struct {
	Name   string            `json:"name"`
	Email  string            `json:"email"`
	Avatar string            `json:"avatar"`
	Fields map[string]string `json:"fields"`
}

type createCustomerResponse struct {
	CustomerID string `json:"customer_id"`
}

type updateCustomerRequest struct {
	CustomerID string            `json:"customer_id"`
	Name       string            `json:"name"`
	Email      string            `json:"email"`
	Avatar     string            `json:"avatar"`
	Fields     map[string]string `json:"fields"`
}

type updateCustomerResponse struct {
	Customer objects.Customer `json:"customer"`
}

type banCustomerRequest struct {
	CustomerID string `json:"customer_id"`
	Ban        Ban    `json:"ban"`
}

type updateAgentRequest struct {
	AgentID       string `json:"agent_id"`
	RoutingStatus string `json:"routing_status"`
}

type markEventsAsSeenRequest struct {
	ChatID   string `json:"chat_id"`
	SeenUpTo string `json:"seen_up_to"`
}

type sendTypingIndicatorRequest struct {
	ChatID     string `json:"chat_id"`
	Recipients string `json:"recipients"`
	IsTyping   bool   `json:"is_typing"`
}

type multicastRequest struct {
	Scopes  MulticastScopes `json:"scopes"`
	Content json.RawMessage `json:"content"`
	Type    string          `json:"type"`
}

type emptyResponse struct{}

type hashedPaginationRequest struct {
	PageID string `json:"page_id"`
	Limit  uint   `json:"limit"`
	Order  string `json:"order"`
}

type hashedPaginationResponse struct {
	PreviousPageID string `json:"previous_page_id,omitempty"`
	NextPageID     string `json:"next_page_id,omitempty"`
}

type paginationRequest struct {
	Page  uint `json:"page"`
	Limit uint `json:"limit"`
}

type paginationResponse struct {
	Page  uint `json:"page"`
	Total uint `json:"total"`
}