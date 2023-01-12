package agent

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/livechat/lc-sdk-go/v5/authorization"
	i "github.com/livechat/lc-sdk-go/v5/internal"
)

type agentAPI interface {
	Call(string, interface{}, interface{}, ...*i.CallOptions) error
	UploadFile(string, []byte) (string, error)
	SetCustomHost(string)
	SetCustomHeader(string, string)
	SetRetryStrategy(i.RetryStrategyFunc)
	SetStatsSink(i.StatsSinkFunc)
	SetLogger(*log.Logger)
}

// API provides the API operation methods for making requests to Agent Chat API via Web API.
// See this package's package overview docs for details on the service.
type API struct {
	agentAPI
}

// NewAPI returns ready to use Agent API.
//
// If provided client is nil, then default http client with 20s timeout is used.
func NewAPI(t authorization.TokenGetter, client *http.Client, clientID string) (*API, error) {
	api, err := i.NewAPIWithFileUpload(t, client, clientID, i.DefaultHTTPRequestGenerator("agent"))
	if err != nil {
		return nil, err
	}
	return &API{api}, nil
}

// SetAuthorID provides a way to point the actual author of the action (e.g. send an event as a bot)
func (a *API) SetAuthorID(authorID string) {
	a.agentAPI.SetCustomHeader("X-Author-Id", authorID)
}

// ListChats returns chat summaries list.
func (a *API) ListChats(filters *chatsFilters, sortOrder, pageID string, limit uint) (summary []ChatSummary, found uint, previousPage, nextPage string, err error) {
	var resp listChatsResponse
	err = a.Call("list_chats", &listChatsRequest{
		Filters: filters,
		hashedPaginationRequest: &hashedPaginationRequest{
			SortOrder: sortOrder,
			PageID:    pageID,
			Limit:     limit,
		},
	}, &resp)

	return resp.ChatsSummary, resp.FoundChats, resp.PreviousPageID, resp.NextPageID, err
}

// GetChat returns given thread for given chat.
func (a *API) GetChat(chatID string, threadID string) (Chat, error) {
	var resp Chat
	err := a.Call("get_chat", &getChatRequest{
		ChatID:   chatID,
		ThreadID: threadID,
	}, &resp)

	return resp, err
}

// ListChats returns threads list.
func (a *API) ListThreads(chatID, sortOrder, pageID string, limit, minEventsCount uint, filters *threadsFilters) (threads []Thread, found uint, previousPage, nextPage string, err error) {
	var resp listThreadsResponse
	err = a.Call("list_threads", &listThreadsRequest{
		ChatID: chatID,
		hashedPaginationRequest: &hashedPaginationRequest{
			SortOrder: sortOrder,
			PageID:    pageID,
			Limit:     limit,
		},
		MinEventsCount: minEventsCount,
		Filters:        filters,
	}, &resp)

	return resp.Threads, resp.FoundThreads, resp.PreviousPageID, resp.NextPageID, err
}

// ListArchives returns archived chats.
func (a *API) ListArchives(filters *archivesFilters, pageID string, limit uint) (chats []Chat, found uint, previousPage, nextPage string, err error) {
	var resp listArchivesResponse
	err = a.Call("list_archives", &listArchivesRequest{
		Filters: filters,
		hashedPaginationRequest: &hashedPaginationRequest{
			PageID: pageID,
			Limit:  limit,
		},
	}, &resp)

	return resp.Chats, resp.FoundChats, resp.PreviousPageID, resp.NextPageID, err
}

// StartChat starts new chat with access, properties and initial thread as defined in initialChat.
// It returns respectively chat ID, thread ID and initial event IDs (except for server-generated events).
func (a *API) StartChat(initialChat *InitialChat, continuous, active bool) (chatID, threadID string, eventIDs []string, err error) {
	var resp startChatResponse

	if err := initialChat.Validate(); err != nil {
		return "", "", nil, err
	}

	err = a.Call("start_chat", &startChatRequest{
		Chat:       initialChat,
		Continuous: continuous,
		Active:     active,
	}, &resp)
	return resp.ChatID, resp.ThreadID, resp.EventIDs, err
}

// ResumeChat resumes chat initialChat.ID with access, properties and initial thread
// as defined in initialChat.
// It returns respectively thread ID and initial event IDs (except for server-generated events).
func (a *API) ResumeChat(initialChat *InitialChat, continuous, active bool) (threadID string, eventIDs []string, err error) {
	var resp resumeChatResponse

	if err := initialChat.Validate(); err != nil {
		return "", nil, err
	}

	err = a.Call("resume_chat", &resumeChatRequest{
		Chat:       initialChat,
		Continuous: continuous,
		Active:     active,
	}, &resp)

	return resp.ThreadID, resp.EventIDs, err
}

// DeactivateChat deactivates active thread for given chat. If no thread is active, then this
// method is a no-op.
func (a *API) DeactivateChat(chatID string, ignoreRequesterPresence bool) error {
	return a.Call("deactivate_chat", &deactivateChatRequest{
		ID:                      chatID,
		IgnoreRequesterPresence: ignoreRequesterPresence,
	}, &emptyResponse{})
}

// FollowChat marks given chat as followed by requester.
func (a *API) FollowChat(chatID string) error {
	return a.Call("follow_chat", &followChatRequest{
		ID: chatID,
	}, &emptyResponse{})
}

// UnfollowChat removes requester from chat followers.
func (a *API) UnfollowChat(chatID string) error {
	return a.Call("unfollow_chat", &unfollowChatRequest{
		ID: chatID,
	}, &emptyResponse{})
}

// TransferChat transfers chat to agent or group.
func (a *API) TransferChat(chatID, targetType string, ids []interface{}, opts TransferChatOptions) error {
	var target *transferTarget
	if targetType != "" || len(ids) > 0 {
		target = &transferTarget{
			Type: targetType,
			IDs:  ids,
		}
	}
	return a.Call("transfer_chat", &transferChatRequest{
		ID:                       chatID,
		Target:                   target,
		IgnoreRequesterPresence:  opts.IgnoreRequesterPresence,
		IgnoreAgentsAvailability: opts.IgnoreAgentsAvailability,
	}, &emptyResponse{})
}

// AddUserToChat adds user to the chat. You can't add more than one customer type user to the chat.
func (a *API) AddUserToChat(chatID, userID, userType, visibility string, ignoreRequesterPresence bool) error {
	return a.Call("add_user_to_chat", &addUserToChatRequest{
		ChatID:                  chatID,
		UserID:                  userID,
		UserType:                userType,
		Visibility:              visibility,
		IgnoreRequesterPresence: ignoreRequesterPresence,
	}, &emptyResponse{})
}

// RemoveUserFromChat Removes a user from chat. Removing customer user type is not allowed.
// It's always possible to remove the requester from the chat.
func (a *API) RemoveUserFromChat(chatID, userID, userType string, ignoreRequesterPresence bool) error {
	return a.Call("remove_user_from_chat", &removeUserFromChatRequest{
		ChatID:                  chatID,
		UserID:                  userID,
		UserType:                userType,
		IgnoreRequesterPresence: ignoreRequesterPresence,
	}, &emptyResponse{})
}

// SendEvent sends event of supported type to given chat.
// It returns event ID.
//
// Supported event types are: event, message, system_message and file.
func (a *API) SendEvent(chatID string, event interface{}, attachToLastThread bool) (string, error) {
	if err := ValidateEvent(event); err != nil {
		return "", err
	}

	var resp sendEventResponse
	err := a.Call("send_event", &sendEventRequest{
		ChatID:             chatID,
		Event:              event,
		AttachToLastThread: &attachToLastThread,
	}, &resp)

	return resp.EventID, err
}

// SendRichMessagePostback sends postback for given rich message event.
func (a *API) SendRichMessagePostback(chatID, eventID, threadID, postbackID string, toggled bool) error {
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
func (a *API) UpdateChatProperties(chatID string, properties Properties) error {
	return a.Call("update_chat_properties", &updateChatPropertiesRequest{
		ID:         chatID,
		Properties: properties,
	}, &emptyResponse{})
}

// DeleteChatProperties deletes given chat's properties.
func (a *API) DeleteChatProperties(chatID string, properties map[string][]string) error {
	return a.Call("delete_chat_properties", &deleteChatPropertiesRequest{
		ID:         chatID,
		Properties: properties,
	}, &emptyResponse{})
}

// UpdateThreadProperties updates given thread's properties.
func (a *API) UpdateThreadProperties(chatID, threadID string, properties Properties) error {
	return a.Call("update_thread_properties", &updateThreadPropertiesRequest{
		ChatID:     chatID,
		ThreadID:   threadID,
		Properties: properties,
	}, &emptyResponse{})
}

// DeleteThreadProperties deletes given thread's properties.
func (a *API) DeleteThreadProperties(chatID, threadID string, properties map[string][]string) error {
	return a.Call("delete_thread_properties", &deleteThreadPropertiesRequest{
		ChatID:     chatID,
		ThreadID:   threadID,
		Properties: properties,
	}, &emptyResponse{})
}

// UpdateEventProperties updates given event's properties.
func (a *API) UpdateEventProperties(chatID, threadID, eventID string, properties Properties) error {
	return a.Call("update_event_properties", &updateEventPropertiesRequest{
		ChatID:     chatID,
		ThreadID:   threadID,
		EventID:    eventID,
		Properties: properties,
	}, &emptyResponse{})
}

// DeleteEventProperties deletes given event's properties.
func (a *API) DeleteEventProperties(chatID, threadID, eventID string, properties map[string][]string) error {
	return a.Call("delete_event_properties", &deleteEventPropertiesRequest{
		ChatID:     chatID,
		ThreadID:   threadID,
		EventID:    eventID,
		Properties: properties,
	}, &emptyResponse{})
}

// TagThread adds given tag to thread.
func (a *API) TagThread(chatID, threadID, tag string) error {
	return a.Call("tag_thread", &changeThreadTagRequest{
		ChatID:   chatID,
		ThreadID: threadID,
		Tag:      tag,
	}, &emptyResponse{})
}

// UntagThread removes given tag from thread.
func (a *API) UntagThread(chatID, threadID, tag string) error {
	return a.Call("untag_thread", &changeThreadTagRequest{
		ChatID:   chatID,
		ThreadID: threadID,
		Tag:      tag,
	}, &emptyResponse{})
}

// GetCustomer returns Customer.
func (a *API) GetCustomer(customerID string) (customer Customer, err error) {
	var resp Customer
	err = a.Call("get_customer", &getCustomersRequest{
		ID: customerID,
	}, &resp)

	return resp, err
}

// ListCustomers returns the list of Customers.
func (a *API) ListCustomers(limit uint, pageID, sortOrder, sortBy string, filters *customersFilters) (customers []Customer, total uint, limited uint, previousPage, nextPage string, err error) {
	var resp listCustomersResponse
	err = a.Call("list_customers", &listCustomersRequest{
		PageID:    pageID,
		Limit:     limit,
		SortOrder: sortOrder,
		Filters:   filters,
		SortBy:    sortBy,
	}, &resp)

	return resp.Customers, resp.TotalCustomers, resp.LimitedCustomers, resp.PreviousPageID, resp.NextPageID, err
}

// CreateCustomer creates new Customer.
func (a *API) CreateCustomer(name, email, avatar string, sessionFields []map[string]string) (string, error) {
	var resp createCustomerResponse
	err := a.Call("create_customer", &createCustomerRequest{
		Name:          name,
		Email:         email,
		Avatar:        avatar,
		SessionFields: sessionFields,
	}, &resp)

	return resp.CustomerID, err
}

// UpdateCustomer updates customer's info.
func (a *API) UpdateCustomer(customerID, name, email, avatar string, sessionFields []map[string]string) error {
	return a.Call("update_customer", &updateCustomerRequest{
		ID:            customerID,
		Name:          name,
		Email:         email,
		Avatar:        avatar,
		SessionFields: sessionFields,
	}, &emptyResponse{})
}

// BanCustomer bans customer for specific period of time (expressed in days).
func (a *API) BanCustomer(customerID string, days uint) error {
	return a.Call("ban_customer", &banCustomerRequest{
		ID: customerID,
		Ban: ban{
			Days: days,
		},
	}, &emptyResponse{})
}

// SetRoutingStatus changes status of an agent or a bot.
func (a *API) SetRoutingStatus(agentID, status string) error {
	return a.Call("set_routing_status", &setRoutingStatusRequest{
		AgentID: agentID,
		Status:  status,
	}, &emptyResponse{})
}

// MarkEventsAsSeen marks all events up to given date in given chat as seen for current agent.
func (a *API) MarkEventsAsSeen(chatID string, seenUpTo time.Time) error {
	return a.Call("mark_events_as_seen", &markEventsAsSeenRequest{
		ChatID:   chatID,
		SeenUpTo: seenUpTo.Format(time.RFC3339Nano),
	}, &emptyResponse{})
}

// SendTypingIndicator sends a notification about typing to defined recipients.
func (a *API) SendTypingIndicator(chatID, visibility string, isTyping bool) error {
	return a.Call("send_typing_indicator", &sendTypingIndicatorRequest{
		ChatID:     chatID,
		Visibility: visibility,
		IsTyping:   isTyping,
	}, &emptyResponse{})
}

// Multicast method serves for the chat-unrelated communication. Messages sent using multicast are not being saved.
func (a *API) Multicast(recipients MulticastRecipients, content json.RawMessage, multicastType string) error {
	return a.Call("multicast", &multicastRequest{
		Recipients: recipients,
		Content:    content,
		Type:       multicastType,
	}, &emptyResponse{})
}

// ListAgentsForTransfer returns the Agents you can transfer a given chat to.
func (a *API) ListAgentsForTransfer(chatID string) (AgentsForTransfer, error) {
	var resp AgentsForTransfer
	err := a.Call("list_agents_for_transfer", &listAgentsForTransferRequest{
		ChatID: chatID,
	}, &resp)
	return resp, err
}

// FollowCustomer marks a customer as followed. As a result, the requester (an agent) will receive the info about all the changes related to that customer via pushes.
func (a *API) FollowCustomer(customerID string) error {
	return a.Call("follow_customer", &followCustomerRequest{
		ID: customerID,
	}, &emptyResponse{})
}

// UnfollowCustomer removes the agent from the list of customer's followers.
func (a *API) UnfollowCustomer(customerID string) error {
	return a.Call("unfollow_customer", &followCustomerRequest{
		ID: customerID,
	}, &emptyResponse{})
}

func (a *API) ListRoutingStatuses(groupIDs []int) ([]AgentStatus, error) {
	var resp []AgentStatus
	err := a.Call("list_routing_statuses", &listRoutingStatusesRequest{
		Filters: &routingStatusesFilter{
			GroupIDs: groupIDs,
		},
	}, &resp)

	return resp, err
}
