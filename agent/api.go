package agent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	i "github.com/livechat/lc-sdk-go/internal"
	"github.com/livechat/lc-sdk-go/objects/events"

	"github.com/livechat/lc-sdk-go/objects"
)

// AgentAPI provides the API operation methods for making requests to Agent Chat API via Web API.
// See this package's package overview docs for details on the service.
type AgentAPI struct {
	*i.API
}

// NewAPI returns ready to use Agent API.
//
// If provided client is nil, then default http client with 20s timeout is used.
func NewAPI(t i.TokenGetter, client *http.Client, clientID string) (*AgentAPI, error) {
	api, err := i.NewAPI(t, client, clientID, "agent")
	if err != nil {
		return nil, err
	}
	return &AgentAPI{api}, nil
}

// GetChatsSummary returns chats summary.
func (a *AgentAPI) GetChatsSummary(filters *chatsFilters, page, limit uint) ([]objects.ChatSummary, uint, string, string, error) {
	var resp getChatsSummaryResponse
	err := a.Call("get_chats_summary", &getChatsSummaryRequest{
		Filters: filters,
		Pagination: paginationRequest{
			Page:  page,
			Limit: limit,
		},
	}, &resp)

	return resp.ChatsSummary, resp.FoundChats, resp.PreviousPageID, resp.NextPageID, err
}

// GetChatThreadsSummary returns threads summary for given chat.
func (a *AgentAPI) GetChatThreadsSummary(chatID, order, pageID string, limit uint) ([]objects.ThreadSummary, uint, string, string, error) {
	var resp getChatThreadsSummaryResponse
	err := a.Call("get_chat_threads_summary", &getChatThreadsSummaryRequest{
		ChatID: chatID,
		hashedPaginationRequest: &hashedPaginationRequest{
			Order:  order,
			Limit:  limit,
			PageID: pageID,
		},
	}, &resp)

	return resp.ThreadsSummary, resp.FoundThreads, resp.PreviousPageID, resp.NextPageID, err
}

// GetChatThreads returns given threads, or all if no threads are provided, for given chat.
func (a *AgentAPI) GetChatThreads(chatID string, threadIDs ...string) (objects.Chat, error) {
	var resp getChatThreadsResponse
	err := a.Call("get_chat_threads", &getChatThreadsRequest{
		ChatID:    chatID,
		ThreadIDs: threadIDs,
	}, &resp)

	return resp.Chat, err
}

// GetArchives returns archived chats.
func (a *AgentAPI) GetArchives(filters *archivesFilters, page, limit uint) ([]objects.Chat, uint, uint, error) {
	var resp getArchivesResponse
	err := a.Call("get_archives", &getArchivesRequest{
		Filters: filters,
		Pagination: paginationRequest{
			Page:  page,
			Limit: limit,
		},
	}, &resp)

	return resp.Chats, resp.Pagination.Page, resp.Pagination.Total, err
}

// StartChat starts new chat with access, properties and initial thread as defined in initialChat.
// It returns respectively chat ID, thread ID and initial event IDs (except for server-generated events).
func (a *AgentAPI) StartChat(initialChat *objects.InitialChat, continuous bool) (string, string, []string, error) {
	var resp startChatResponse

	if err := a.validateInitialChat(initialChat); err != nil {
		return "", "", nil, err
	}

	err := a.Call("start_chat", &startChatRequest{
		Chat:       initialChat,
		Continuous: continuous,
	}, &resp)
	return resp.ChatID, resp.ThreadID, resp.EventIDs, err
}

// ActivateChat activates chat initialChat.ID with access, properties and initial thread
// as defined in initialChat.
// It returns respectively thread ID and initial event IDs (except for server-generated events).
func (a *AgentAPI) ActivateChat(initialChat *objects.InitialChat, continuous bool) (string, []string, error) {
	var resp activateChatResponse

	if err := a.validateInitialChat(initialChat); err != nil {
		return "", nil, err
	}

	err := a.Call("activate_chat", &activateChatRequest{
		Chat:       initialChat,
		Continuous: continuous,
	}, &resp)

	return resp.ThreadID, resp.EventIDs, err
}

// CloseThread closes active thread for given chat. If no thread is active, then this
// method is a no-op.
func (a *AgentAPI) CloseThread(chatID string) error {
	return a.Call("close_thread", &closeThreadRequest{
		ChatID: chatID,
	}, &emptyResponse{})
}

// FollowChat marks given chat as followed by requester.
func (a *AgentAPI) FollowChat(chatID string) error {
	return a.Call("follow_chat", &followChatRequest{
		ChatID: chatID,
	}, &emptyResponse{})
}

// UnfollowChat removes requester from chat followers.
func (a *AgentAPI) UnfollowChat(chatID string) error {
	return a.Call("unfollow_chat", &unfollowChatRequest{
		ChatID: chatID,
	}, &emptyResponse{})
}

// GrantAccess grants access to a new resource without overwriting the existing ones.
func (a *AgentAPI) GrantAccess(resource, id string, access objects.Access) error {
	return a.Call("grant_access", &modifyAccessRequest{
		Resource: resource,
		ID:       id,
		Access:   access,
	}, &emptyResponse{})
}

// RevokeAccess removes access to given resource.
func (a *AgentAPI) RevokeAccess(resource, id string, access objects.Access) error {
	return a.Call("revoke_access", &modifyAccessRequest{
		Resource: resource,
		ID:       id,
		Access:   access,
	}, &emptyResponse{})
}

// SetAccess gives access to a new resource overwriting the existing ones.
func (a *AgentAPI) SetAccess(resource, id string, access objects.Access) error {
	return a.Call("set_access", &modifyAccessRequest{
		Resource: resource,
		ID:       id,
		Access:   access,
	}, &emptyResponse{})
}

// TransferChat transfers chat to agent or group.
func (a *AgentAPI) TransferChat(chatID, targetType string, ids []interface{}, force bool) error {
	return a.Call("transfer_chat", &transferChatRequest{
		ChatID: chatID,
		Target: transferTarget{
			Type: targetType,
			IDs:  ids,
		},
		Force: force,
	}, &emptyResponse{})
}

// AddUserToChat adds user to the chat. You can't add more than one customer type user to the chat.
func (a *AgentAPI) AddUserToChat(chatID, userID, userType string) error {
	return a.Call("add_user_to_chat", &changeChatUsersRequest{
		ChatID:   chatID,
		UserID:   userID,
		UserType: userType,
	}, &emptyResponse{})
}

// RemoveUserFromChat Removes a user from chat. Removing customer user type is not allowed.
// It's always possible to remove the requester from the chat.
func (a *AgentAPI) RemoveUserFromChat(chatID, userID, userType string) error {
	return a.Call("remove_user_from_chat", &changeChatUsersRequest{
		ChatID:   chatID,
		UserID:   userID,
		UserType: userType,
	}, &emptyResponse{})
}

// SendEvent sends event of supported type to given chat.
// It returns event ID.
//
// Supported event types are: event, message and system_message.
func (a *AgentAPI) SendEvent(chatID string, event events.Event, attachToLastThread bool) (string, error) {
	var resp sendEventResponse
	err := a.Call("send_event", &sendEventRequest{
		ChatID:             chatID,
		Event:              event,
		AttachToLastThread: attachToLastThread,
	}, &resp)

	return resp.EventID, err
}

// SendRichMessagePostback sends postback for given rich message event.
func (a *AgentAPI) SendRichMessagePostback(chatID, eventID, threadID, postbackID string, toggled bool) error {
	return a.Call("send_rich_message_postback", &sendRichMessagePostbackRequest{
		ChatID:   chatID,
		EventID:  eventID,
		ThreadID: threadID,
		Postback: postback{
			ID:      postbackID,
			Toggled: toggled,
		},
	}, &emptyResponse{})
}

// UpdateChatProperties updates given chat's properties.
func (a *AgentAPI) UpdateChatProperties(chatID string, properties objects.Properties) error {
	return a.Call("update_chat_properties", &updateChatPropertiesRequest{
		ChatID:     chatID,
		Properties: properties,
	}, &emptyResponse{})
}

// DeleteChatProperties deletes given chat's properties.
func (a *AgentAPI) DeleteChatProperties(chatID string, properties map[string][]string) error {
	return a.Call("delete_chat_properties", &deleteChatPropertiesRequest{
		ChatID:     chatID,
		Properties: properties,
	}, &emptyResponse{})
}

// UpdateChatThreadProperties updates given chat thread's properties.
func (a *AgentAPI) UpdateChatThreadProperties(chatID, threadID string, properties objects.Properties) error {
	return a.Call("update_chat_thread_properties", &updateChatThreadPropertiesRequest{
		ChatID:     chatID,
		ThreadID:   threadID,
		Properties: properties,
	}, &emptyResponse{})
}

// DeleteChatThreadProperties deletes given chat thread's properties.
func (a *AgentAPI) DeleteChatThreadProperties(chatID, threadID string, properties map[string][]string) error {
	return a.Call("delete_chat_thread_properties", &deleteChatThreadPropertiesRequest{
		ChatID:     chatID,
		ThreadID:   threadID,
		Properties: properties,
	}, &emptyResponse{})
}

// UpdateEventProperties updates given event's properties.
func (a *AgentAPI) UpdateEventProperties(chatID, threadID, eventID string, properties objects.Properties) error {
	return a.Call("update_event_properties", &updateEventPropertiesRequest{
		ChatID:     chatID,
		ThreadID:   threadID,
		EventID:    eventID,
		Properties: properties,
	}, &emptyResponse{})
}

// DeleteEventProperties deletes given event's properties.
func (a *AgentAPI) DeleteEventProperties(chatID, threadID, eventID string, properties map[string][]string) error {
	return a.Call("delete_event_properties", &deleteEventPropertiesRequest{
		ChatID:     chatID,
		ThreadID:   threadID,
		EventID:    eventID,
		Properties: properties,
	}, &emptyResponse{})
}

// TagChatThread adds given tag to chat thread.
func (a *AgentAPI) TagChatThread(chatID, threadID, tag string) error {
	return a.Call("tag_chat_thread", &changeChatThreadTagRequest{
		ChatID:   chatID,
		ThreadID: threadID,
		Tag:      tag,
	}, &emptyResponse{})
}

// UntagChatThread removes given tag from chat thread.
func (a *AgentAPI) UntagChatThread(chatID, threadID, tag string) error {
	return a.Call("untag_chat_thread", &changeChatThreadTagRequest{
		ChatID:   chatID,
		ThreadID: threadID,
		Tag:      tag,
	}, &emptyResponse{})
}

// GetCustomers returns the list of Customers.
func (a *AgentAPI) GetCustomers(limit uint, pageID, order string, filters *customersFilters) ([]objects.Customer, uint, string, string, error) {
	var resp getCustomersResponse
	err := a.Call("get_customers", &getCustomersRequest{
		PageID:  pageID,
		Limit:   limit,
		Order:   order,
		Filters: filters,
	}, &resp)

	return resp.Customers, resp.TotalCustomers, resp.PreviousPageID, resp.NextPageID, err
}

// CreateCustomer creates new Customer.
func (a *AgentAPI) CreateCustomer(name, email, avatar string, fields map[string]string) (string, error) {
	var resp createCustomerResponse
	err := a.Call("create_customer", &createCustomerRequest{
		Name:   name,
		Email:  email,
		Avatar: avatar,
		Fields: fields,
	}, &resp)

	return resp.CustomerID, err
}

// UpdateCustomer updates customer's info.
func (a *AgentAPI) UpdateCustomer(customerID, name, email, avatar string, fields map[string]string) (objects.Customer, error) {
	var resp updateCustomerResponse
	err := a.Call("update_customer", &updateCustomerRequest{
		CustomerID: customerID,
		Name:       name,
		Email:      email,
		Avatar:     avatar,
		Fields:     fields,
	}, &resp)

	return resp.Customer, err
}

// BanCustomer bans customer for specific period of time (expressed in days).
func (a *AgentAPI) BanCustomer(customerID string, days uint) error {
	return a.Call("ban_customer", &banCustomerRequest{
		CustomerID: customerID,
		Ban: ban{
			Days: days,
		},
	}, &emptyResponse{})
}

// UpdateAgent updates agent's info.
func (a *AgentAPI) UpdateAgent(agentID, routingStatus string) error {
	return a.Call("update_agent", &updateAgentRequest{
		AgentID:       agentID,
		RoutingStatus: routingStatus,
	}, &emptyResponse{})
}

// MarkEventsAsSeen marks all events up to given date in given chat as seen for current agent.
func (a *AgentAPI) MarkEventsAsSeen(chatID string, seenUpTo time.Time) error {
	return a.Call("mark_events_as_seen", &markEventsAsSeenRequest{
		ChatID:   chatID,
		SeenUpTo: seenUpTo.Format(time.RFC3339Nano),
	}, &emptyResponse{})
}

// SendTypingIndicator sends a notification about typing to defined recipients.
func (a *AgentAPI) SendTypingIndicator(chatID, recipients string, isTyping bool) error {
	return a.Call("send_typing_indicator", &sendTypingIndicatorRequest{
		ChatID:     chatID,
		Recipients: recipients,
		IsTyping:   isTyping,
	}, &emptyResponse{})
}

// Multicast method serves for the chat-unrelated communication. Messages sent using multicast are not being saved.
func (a *AgentAPI) Multicast(scopes MulticastScopes, content json.RawMessage, multicastType string) error {
	return a.Call("multicast", &multicastRequest{
		Scopes:  scopes,
		Content: content,
		Type:    multicastType,
	}, &emptyResponse{})
}

func (a *AgentAPI) validateInitialChat(chat *objects.InitialChat) error {
	if chat.Thread != nil {
		for _, e := range chat.Thread.Events {
			switch v := e.(type) {
			case *events.Event:
			case *events.Message:
			case *events.SystemMessage:
			default:
				return fmt.Errorf("event type %T not supported", v)
			}
		}
	}
	return nil
}
