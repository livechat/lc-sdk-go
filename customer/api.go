package customer

import (
	"fmt"
	"net/http"
	"time"

	"github.com/livechat/lc-sdk-go/authorization"
	i "github.com/livechat/lc-sdk-go/internal"
	"github.com/livechat/lc-sdk-go/objects"
)

// API provides the API operation methods for making requests to Customer Chat API via Web API.
// See this package's package overview docs for details on the service.
type API struct {
	*i.FileUploadAPI
}

func CustomerRequestGetter(r i.RequestGetter) i.RequestGetter {
	return func(t *authorization.Token, a string) (*http.Request, error) {
		req, err := r(t, a)
		if err != nil {
			return nil, err
		}
		if t.LicenseID != nil {
			qs := req.URL.Query()
			qs.Add("license_id", fmt.Sprintf("%v", *t.LicenseID))
			req.URL.RawQuery = qs.Encode()
		}
		if a == "list_license_properties" || a == "list_group_properties" {
			req.Method = "GET"
		}
		return req, nil
	}
}

// NewAPI returns ready to use Customer API.
//
// If provided client is nil, then default http client with 20s timeout is used.
func NewAPI(t authorization.TokenGetter, client *http.Client, clientID string) (*API, error) {
	api, err := i.NewFileUploadAPI(t, client, clientID, CustomerRequestGetter(i.DefaultRequestGetter("customer")))
	if err != nil {
		return nil, err
	}
	return &API{api}, nil
}

// StartChat starts new chat with access, properties and initial thread as defined in initialChat.
// It returns respectively chat ID, thread ID and initial event IDs (except for server-generated events).
func (a *API) StartChat(initialChat *objects.InitialChat, continuous bool) (chatID, threadID string, eventIDs []string, err error) {
	req := &startChatRequest{
		Chat:       initialChat,
		Continuous: continuous,
	}

	if err := initialChat.Validate(); err != nil {
		return "", "", nil, err
	}
	var resp startChatResponse
	err = a.Call("start_chat", req, &resp)
	return resp.ChatID, resp.ThreadID, resp.EventIDs, err
}

// SendMessage sends event of type message to given chat.
// It returns event ID.
func (a *API) SendMessage(chatID, text string, recipients Recipients) (string, error) {
	e := objects.Message{
		Event: objects.Event{
			Type:       "message",
			Recipients: string(recipients),
		},
		Text: text,
	}

	return a.SendEvent(chatID, &e, false)
}

// SendSystemMessage sends event of type system_message to given chat.
// It returns event ID.
func (a *API) SendSystemMessage(chatID, text, messageType string, textVars map[string]string, recipients Recipients, attachToLastThread bool) (string, error) {
	e := objects.SystemMessage{
		Event: objects.Event{
			Type:       "system_message",
			Recipients: string(recipients),
		},
		Text:     text,
		Type:     messageType,
		TextVars: textVars,
	}

	return a.SendEvent(chatID, &e, attachToLastThread)
}

// SendEvent sends event of supported type to given chat.
// It returns event ID.
//
// Supported event types are: event, file, message, rich_message and system_message.
func (a *API) SendEvent(chatID string, e interface{}, attachToLastThread bool) (string, error) {
	if err := objects.ValidateEvent(e); err != nil {
		return "", err
	}

	var resp sendEventResponse
	err := a.Call("send_event", &sendEventRequest{
		ChatID:             chatID,
		Event:              e,
		AttachToLastThread: &attachToLastThread,
	}, &resp)

	return resp.EventID, err
}

// ActivateChat activates chat initialChat.ID with access, properties and initial thread
// as defined in initialChat.
// It returns respectively thread ID and initial event IDs (except for server-generated events).
func (a *API) ActivateChat(initialChat *objects.InitialChat, continuous bool) (threadID string, eventIDs []string, err error) {
	var resp activateChatResponse

	if err := initialChat.Validate(); err != nil {
		return "", nil, err
	}

	err = a.Call("activate_chat", &activateChatRequest{
		Chat:       initialChat,
		Continuous: continuous,
	}, &resp)

	return resp.ThreadID, resp.EventIDs, err
}

// ListChats returns chat summaries list.
func (a *API) ListChats(sortOrder, pageID string, limit uint) (summary []objects.ChatSummary, total uint, previousPage, nextPage string, err error) {
	var resp listChatsResponse
	err = a.Call("list_chats", &listChatsRequest{
		hashedPaginationRequest: &hashedPaginationRequest{
			SortOrder: sortOrder,
			PageID:    pageID,
			Limit:     limit,
		},
	}, &resp)

	return resp.ChatsSummary, resp.TotalChats, resp.PreviousPageID, resp.NextPageID, err
}

// GetChat returns given thread for given chat.
func (a *API) GetChat(chatID string, threadID string) (objects.Chat, error) {
	var resp getChatResponse
	err := a.Call("get_chat", &getChatRequest{
		ChatID:   chatID,
		ThreadID: threadID,
	}, &resp)

	return resp.Chat, err
}

// ListThreads returns threads list.
func (a *API) ListThreads(chatID, sortOrder, pageID string, limit, minEventsCount uint) (threads []objects.Thread, found uint, previousPage, nextPage string, err error) {
	var resp listThreadsResponse
	err = a.Call("list_threads", &listThreadsRequest{
		ChatID: chatID,
		hashedPaginationRequest: &hashedPaginationRequest{
			SortOrder: sortOrder,
			PageID:    pageID,
			Limit:     limit,
		},
	}, &resp)

	return resp.Threads, resp.FoundThreads, resp.PreviousPageID, resp.NextPageID, err
}

// DeactivateChat deactivates active thread for given chat. If no thread is active, then this
// method is a no-op.
func (a *API) DeactivateChat(chatID string) error {
	return a.Call("deactivate_chat", &deactivateChatRequest{
		ChatID: chatID,
	}, &emptyResponse{})
}

// SendRichMessagePostback sends postback for given rich message event.
func (a *API) SendRichMessagePostback(chatID, threadID, eventID, postbackID string, toggled bool) error {
	return a.Call("send_rich_message_postback", &sendRichMessagePostbackRequest{
		ChatID:   chatID,
		ThreadID: threadID,
		EventID:  eventID,
		Postback: postback{
			ID:      postbackID,
			Toggled: toggled,
		},
	}, &emptyResponse{})
}

// SendSneakPeek sends sneak peek of message for given chat.
func (a *API) SendSneakPeek(chatID, text string) error {
	return a.Call("send_sneak_peek", &sendSneakPeekRequest{
		ChatID:        chatID,
		SneakPeekText: text,
	}, &emptyResponse{})
}

// UpdateChatProperties updates given chat's properties.
func (a *API) UpdateChatProperties(chatID string, properties objects.Properties) error {
	return a.Call("update_chat_properties", &updateChatPropertiesRequest{
		ChatID:     chatID,
		Properties: properties,
	}, &emptyResponse{})
}

// DeleteChatProperties deletes given chat's properties.
func (a *API) DeleteChatProperties(chatID string, properties map[string][]string) error {
	return a.Call("delete_chat_properties", &deleteChatPropertiesRequest{
		ChatID:     chatID,
		Properties: properties,
	}, &emptyResponse{})
}

// UpdateThreadProperties updates given thread's properties.
func (a *API) UpdateThreadProperties(chatID, threadID string, properties objects.Properties) error {
	return a.Call("update_thread_properties", &updateThreadPropertiesRequest{
		ChatID:     chatID,
		ThreadID:   threadID,
		Properties: properties,
	}, &emptyResponse{})
}

// DeleteThreadProperties deletes given chat thread's properties.
func (a *API) DeleteThreadProperties(chatID, threadID string, properties map[string][]string) error {
	return a.Call("delete_thread_properties", &deleteThreadPropertiesRequest{
		ChatID:     chatID,
		ThreadID:   threadID,
		Properties: properties,
	}, &emptyResponse{})
}

// UpdateEventProperties updates given event's properties.
func (a *API) UpdateEventProperties(chatID, threadID, eventID string, properties objects.Properties) error {
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

// UpdateCustomer updates current customer's info.
func (a *API) UpdateCustomer(name, email, avatarURL string, sessionFields []map[string]string) error {
	return a.Call("update_customer", &updateCustomerRequest{
		Name:          name,
		Email:         email,
		Avatar:        avatarURL,
		SessionFields: sessionFields,
	}, &emptyResponse{})
}

// SetCustomerSessionFields sets current customer's fields.
func (a *API) SetCustomerSessionFields(sessionFields []map[string]string) error {
	return a.Call("set_customer_session_fields", &setCustomerSessionFieldsRequest{
		SessionFields: sessionFields,
	}, &emptyResponse{})
}

// ListGroupStatuses returns status of provided groups.
//
// Possible values are: GroupStatusOnline, GroupStatusOffline and GroupStatusOnlineForQueue.
// GroupStatusUnknown should never be returned.
func (a *API) ListGroupStatuses(groupIDs []int) (map[int]GroupStatus, error) {
	req := &listGroupStatusesRequest{}
	if len(groupIDs) == 0 {
		req.All = true
	} else {
		req.GroupIDs = groupIDs
	}
	var resp listGroupStatusesResponse
	err := a.Call("list_group_statuses", req, &resp)

	r := map[int]GroupStatus{}

	if err == nil {
		for g, s := range resp.Status {
			r[g] = toGroupStatus(s)
		}
	}

	return r, err
}

// CheckGoals triggers checking if goals were achieved. Then, Agents receive the information.
// You should call this method to provide goals parameters for the server when the customers limit is reached.
// Works only for offline Customers.
func (a *API) CheckGoals(pageURL string, groupID int, customerFields map[string]string) error {
	return a.Call("check_goals", &checkGoalsRequest{
		PageURL:        pageURL,
		GroupID:        groupID,
		CustomerFields: customerFields,
	}, &emptyResponse{})
}

// GetForm returns an empty prechat, postchat or ticket form and indication whether
// the form is enabled on the license.
func (a *API) GetForm(groupID int, formType FormType) (form *Form, enabled bool, err error) {
	var resp getFormResponse
	err = a.Call("get_form", &getFormRequest{
		GroupID: groupID,
		Type:    string(formType),
	}, &resp)

	return resp.Form, resp.Enabled, err
}

// GetPredictedAgent returns the predicted Agent - the one the Customer will chat with
// when the chat starts. To use this method, the Customer needs to be logged in,
// which can be done via Customer Chat RTM Api's login method.
func (a *API) GetPredictedAgent() (*PredictedAgent, error) {
	var resp PredictedAgent
	err := a.Call("get_predicted_agent", nil, &resp)
	return &resp, err
}

// GetURLInfo returns info on a given URL.
func (a *API) GetURLInfo(url string) (*URLInfo, error) {
	var resp URLInfo
	err := a.Call("get_url_info", &getURLInfoRequest{
		URL: url,
	}, &resp)
	return &resp, err
}

// MarkEventsAsSeen marks all events up to given date in given chat as seen for current customer.
func (a *API) MarkEventsAsSeen(chatID string, seenUpTo time.Time) error {
	return a.Call("mark_events_as_seen", &markEventsAsSeenRequest{
		ChatID:   chatID,
		SeenUpTo: seenUpTo.Format(time.RFC3339Nano),
	}, &emptyResponse{})
}

// GetCustomer returns current Customer.
func (a *API) GetCustomer() (*objects.Customer, error) {
	var resp objects.Customer
	err := a.Call("get_customer", nil, &resp)
	return &resp, err
}

// ListLicenseProperties returns the properties of a given license.
func (a *API) ListLicenseProperties(namespace, name string) (objects.Properties, error) {
	var resp objects.Properties
	err := a.Call("list_license_properties", &listLicensePropertiesRequest{
		Namespace: namespace,
		Name:      name,
	}, &resp)
	return resp, err
}

// ListGroupProperties returns the properties of a given group.
func (a *API) ListGroupProperties(groupID uint, namespace, name string) (objects.Properties, error) {
	var resp objects.Properties
	err := a.Call("list_group_properties", &listGroupPropertiesRequest{
		GroupID:   groupID,
		Namespace: namespace,
		Name:      name,
	}, &resp)
	return resp, err
}

// AcceptGreeting marks an incoming greeting as seen.
func (a *API) AcceptGreeting(greetingID int, uniqueID string) error {
	return a.Call("accept_greeting", &acceptGreetingRequest{
		GreetingID: greetingID,
		UniqueID:   uniqueID,
	}, &emptyResponse{})
}

// CancelGreeting cancels a greeting (an invitation to the chat).
func (a *API) CancelGreeting(uniqueID string) error {
	return a.Call("cancel_greeting", &cancelGreetingRequest{
		UniqueID: uniqueID,
	}, &emptyResponse{})
}
