package agent

import (
	"encoding/json"

	"github.com/livechat/lc-sdk-go/objects"
)

type getChatsSummaryRequest struct {
	// Filters    ChatsSummaryFilters                  `json:"filters"`
	Pagination *hashedPaginationRequest `json:"pagination"`
}

type getChatThreadsSummaryRequest struct {
	ChatID string `json:"chat_id"`
	*hashedPaginationRequest
}

type getChatThreadsSummaryResponse struct {
	*hashedPaginationResponse
	ThreadsSummary []ThreadSummary `json:"threads_summary`
	FoundThreads   int             `json:"found_threads"`
}

type getChatThreadsRequest struct {
	ChatID    string   `json:"chat_id"`
	ThreadIDs []string `json:"thread_ids"`
}

type getChatThreadsResponse struct {
	Chat objects.Chat `json:"chat"`
}

type getArchivesRequest struct {
	// Filters    ArchivesFilters           `json:"filters"`
	Pagination paginationRequest `json:"pagination"`
}

type getArchivesResponse struct {
	Chats      []objects.Chat     `json:"chats"`
	Pagination paginationResponse `json:"pagination"`
}

type startChatRequest struct {
	Chat       *InitialChat `json:"chat"`
	Continuous bool         `json:"continuous"`
}

type startChatResponse struct {
	ChatID   string   `json:"chat_id"`
	ThreadID string   `json:"thread_id"`
	EventIDs []string `json:"event_ids,omitempty"`
}

type activateChatRequest struct {
	Chat       *InitialChat `json:"chat"`
	Continuous bool         `json:"continuous"`
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
	ChatID string `json:"chat_id"`
	Target target `json:"target"`
	Force  bool   `json:"force"`
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
	Postback postback `json:"postback"`
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
	PageID  string  `json:"page_id"`
	Limit   uint    `json:"limit"`
	Order   string  `json:"order"`
	Filters Filters `json:"filters"`
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
	Ban        ban    `json:"ban"`
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
	Scopes  multicastScopes `json:"scopes"`
	Content json.RawMessage `json:"content"`
	Type    string          `json:"type"`
}

type emptyResponse struct{}

type postback struct {
	ID      string `json:"id"`
	Toggled bool   `json:"toggled"`
}

type ban struct {
	Days uint64 `json:"days"`
}

type multicastScopes struct {
	Agents    *multicastScopesAgents    `json:"agents,omitempty"`
	Customers *multicastScopesCustomers `json:"customers,omitempty"`
}

type multicastScopesAgents struct {
	Groups *[]uint64 `json:"groups,omitempty"`
	IDs    *[]string `json:"ids,omitempty"`
	All    *bool     `json:"all,omitempty"`
}

type multicastScopesCustomers struct {
	IDs *[]string `json:"ids,omitempty"`
}

type hashedPaginationRequest struct {
	PageID string `json:"page_id"`
	Limit  uint64 `json:"limit"`
	Order  string `json:"order"`
}

type hashedPaginationResponse struct {
	PreviousPageID string `json:"previous_page_id,omitempty"`
	NextPageID     string `json:"next_page_id,omitempty"`
}

type paginationRequest struct {
	Page  uint64 `json:"page"`
	Limit uint64 `json:"limit"`
}

type paginationResponse struct {
	Page  uint64 `json:"page"`
	Total uint64 `json:"total"`
}

type target struct {
	Type string `json:"type"`
	IDs  []uint `json:"ids"`
}
