package customer

import (
	"fmt"
	"net/http"
	"time"

	i "github.com/livechat/lc-sdk-go/internal"
	"github.com/livechat/lc-sdk-go/objects"
	"github.com/livechat/lc-sdk-go/objects/events"
)

// CustomerAPI provides the API operation methods for making requests to Customer Chat API via Web API.
// See this package's package overview docs for details on the service.
type CustomerAPI struct {
	*i.API
}

// NewAPI returns ready to use Customer API.
//
// If provided client is nil, then default http client with 20s timeout is used.
func NewAPI(t i.TokenGetter, client *http.Client, clientID string) (*CustomerAPI, error) {
	api, err := i.NewAPI(t, client, clientID, "customer")
	if err != nil {
		return nil, err
	}
	return &CustomerAPI{api}, nil
}

// StartChat starts new chat with access, properties and initial thread as defined in initialChat.
// It returns respectively chat ID, thread ID and initial event IDs (except for server-generated events).
func (a *CustomerAPI) StartChat(initialChat *objects.InitialChat, continuous bool) (string, string, []string, error) {
	req := &startChatRequest{
		Chat:       initialChat,
		Continuous: continuous,
	}
	var resp startChatResponse
	err := a.Call("start_chat", req, &resp)
	return resp.ChatID, resp.ThreadID, resp.EventIDs, err
}

// SendMessage sends event of type message to given chat.
// It returns event ID.
func (a *CustomerAPI) SendMessage(chatID, text string, recipients Recipients) (string, error) {
	e := events.Message{
		Event: &events.Event{
			Type:       "message",
			Recipients: string(recipients),
		},
		Text: text,
	}

	return a.SendEvent(chatID, &e)
}

// SendSystemMessage sends event of type system_message to given chat.
// It returns event ID.
func (a *CustomerAPI) SendSystemMessage(chatID, text, messageType string) (string, error) {
	e := events.SystemMessage{
		Event: events.Event{
			Type: "system_message",
		},
		Text: text,
		Type: messageType,
	}

	return a.SendEvent(chatID, &e)
}

// SendEvent sends event of supported type to given chat.
// It returns event ID.
//
// Supported event types are: event, message and system_message.
func (a *CustomerAPI) SendEvent(chatID string, e interface{}) (string, error) {
	switch v := e.(type) {
	case *events.Event:
	case *events.Message:
	case *events.SystemMessage:
	default:
		return "", fmt.Errorf("event type %T not supported", v)
	}

	var resp sendEventResponse
	err := a.Call("send_event", &sendEventRequest{
		ChatID: chatID,
		Event:  e,
	}, &resp)

	return resp.EventID, err
}

// ActivateChat activates chat initialChat.ID with access, properties and initial thread
// as defined in initialChat.
// It returns respectively thread ID and initial event IDs (except for server-generated events).
func (a *CustomerAPI) ActivateChat(initialChat *objects.InitialChat, continuous bool) (string, []string, error) {
	var resp activateChatResponse

	if initialChat.Thread != nil {
		for _, e := range initialChat.Thread.Events {
			switch v := e.(type) {
			case *events.Event:
			case *events.Message:
			case *events.SystemMessage:
			default:
				return "", nil, fmt.Errorf("event type %T not supported", v)
			}
		}
	}

	err := a.Call("activate_chat", &activateChatRequest{
		Chat:       initialChat,
		Continuous: continuous,
	}, &resp)

	return resp.ThreadID, resp.EventIDs, err
}

// GetChatsSummary returns chats summary.
func (a *CustomerAPI) GetChatsSummary(offset, limit uint) ([]objects.Chat, uint, error) {
	var resp getChatsSummaryResponse
	err := a.Call("get_chats_summary", &getChatsSummaryRequest{
		Limit:  limit,
		Offset: offset,
	}, &resp)

	return resp.Chats, resp.TotalChats, err
}

// GetChatThreadsSummary returns threads summary for given chat.
func (a *CustomerAPI) GetChatThreadsSummary(chatID string, offset, limit uint) ([]objects.ThreadSummary, uint, error) {
	var resp getChatThreadsSummaryResponse
	err := a.Call("get_chat_threads_summary", &getChatThreadsSummaryRequest{
		ChatID: chatID,
		Limit:  limit,
		Offset: offset,
	}, &resp)

	return resp.ThreadsSummary, resp.TotalThreads, err
}

// GetChatThreads returns given threads, or all if no threads are provided, for given chat.
func (a *CustomerAPI) GetChatThreads(chatID string, threadIDs ...string) (objects.Chat, error) {
	var resp getChatThreadsResponse
	err := a.Call("get_chat_threads", &getChatThreadsRequest{
		ChatID:    chatID,
		ThreadIDs: threadIDs,
	}, &resp)

	return resp.Chat, err
}

// CloseThread closes active thread for given chat. If no thread is active, then this
// method is a no-op.
func (a *CustomerAPI) CloseThread(chatID string) error {
	return a.Call("close_thread", &closeThreadRequest{
		ChatID: chatID,
	}, &emptyResponse{})
}

// SendRichMessagePostback sends postback for given rich message event.
func (a *CustomerAPI) SendRichMessagePostback(chatID, threadID, eventID, postbackID string, toggled bool) error {
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
func (a *CustomerAPI) SendSneakPeek(chatID, text string) error {
	return a.Call("send_sneak_peek", &sendSneakPeekRequest{
		ChatID:        chatID,
		SneakPeekText: text,
	}, &emptyResponse{})
}

// UpdateChatProperties updates given chat's properties.
func (a *CustomerAPI) UpdateChatProperties(chatID string, properties objects.Properties) error {
	return a.Call("update_chat_properties", &updateChatPropertiesRequest{
		ChatID:     chatID,
		Properties: properties,
	}, &emptyResponse{})
}

// DeleteChatProperties deletes given chat's properties.
func (a *CustomerAPI) DeleteChatProperties(chatID string, properties map[string][]string) error {
	return a.Call("delete_chat_properties", &deleteChatPropertiesRequest{
		ChatID:     chatID,
		Properties: properties,
	}, &emptyResponse{})
}

// UpdateChatThreadProperties updates given chat thread's properties.
func (a *CustomerAPI) UpdateChatThreadProperties(chatID, threadID string, properties objects.Properties) error {
	return a.Call("update_chat_thread_properties", &updateChatThreadPropertiesRequest{
		ChatID:     chatID,
		ThreadID:   threadID,
		Properties: properties,
	}, &emptyResponse{})
}

// DeleteChatThreadProperties deletes given chat thread's properties.
func (a *CustomerAPI) DeleteChatThreadProperties(chatID, threadID string, properties map[string][]string) error {
	return a.Call("delete_chat_thread_properties", &deleteChatThreadPropertiesRequest{
		ChatID:     chatID,
		ThreadID:   threadID,
		Properties: properties,
	}, &emptyResponse{})
}

// UpdateEventProperties updates given event's properties.
func (a *CustomerAPI) UpdateEventProperties(chatID, threadID, eventID string, properties objects.Properties) error {
	return a.Call("update_event_properties", &updateEventPropertiesRequest{
		ChatID:     chatID,
		ThreadID:   threadID,
		EventID:    eventID,
		Properties: properties,
	}, &emptyResponse{})
}

// DeleteEventProperties deletes given event's properties.
func (a *CustomerAPI) DeleteEventProperties(chatID, threadID, eventID string, properties map[string][]string) error {
	return a.Call("delete_event_properties", &deleteEventPropertiesRequest{
		ChatID:     chatID,
		ThreadID:   threadID,
		EventID:    eventID,
		Properties: properties,
	}, &emptyResponse{})
}

// UpdateCustomer updates current customer's info.
func (a *CustomerAPI) UpdateCustomer(name, email, avatarURL string, fields map[string]string) error {
	return a.Call("update_customer", &updateCustomerRequest{
		Name:   name,
		Email:  email,
		Avatar: avatarURL,
		Fields: fields,
	}, &emptyResponse{})
}

// SetCustomerFields sets current customer's fields.
func (a *CustomerAPI) SetCustomerFields(fields map[string]string) error {
	return a.Call("set_customer_fields", &setCustomerFieldsRequest{
		Fields: fields,
	}, &emptyResponse{})
}

// GetGroupsStatus returns status of provided groups.
//
// Possible values are: GroupStatusOnline, GroupStatusOffline and GroupStatusOnlineForQueue.
// GroupStatusUnknown should never be returned.
func (a *CustomerAPI) GetGroupsStatus(groups []int) (map[int]GroupStatus, error) {
	req := &getGroupsStatusRequest{}
	if len(groups) == 0 {
		req.All = true
	} else {
		req.Groups = groups
	}
	var resp getGroupsStatusResponse
	err := a.Call("get_groups_status", req, &resp)

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
func (a *CustomerAPI) CheckGoals(pageURL string, groupID int, customerFields map[string]string) error {
	return a.Call("check_goals", &checkGoalsRequest{
		PageURL:        pageURL,
		GroupID:        groupID,
		CustomerFields: customerFields,
	}, &emptyResponse{})
}

// GetForm returns an empty prechat, postchat or ticket form and indication whether
// the form is enabled on the license.
func (a *CustomerAPI) GetForm(groupID int, formType FormType) (*Form, bool, error) {
	var resp getFormResponse
	err := a.Call("get_form", &getFormRequest{
		GroupID: groupID,
		Type:    string(formType),
	}, &resp)

	return resp.Form, resp.Enabled, err
}

// GetPredictedAgent returns the predicted Agent - the one the Customer will chat with
// when the chat starts. To use this method, the Customer needs to be logged in,
// which can be done via Customer Chat RTM Api's login method.
func (a *CustomerAPI) GetPredictedAgent() (*PredictedAgent, error) {
	var resp PredictedAgent
	err := a.Call("get_predicted_agent", nil, &resp)
	return &resp, err
}

// GetURLDetails returns info on a given URL.
func (a *CustomerAPI) GetURLDetails(url string) (*URLDetails, error) {
	var resp URLDetails
	err := a.Call("get_url_details", &getURLDetailsRequest{
		URL: url,
	}, &resp)
	return &resp, err
}

// MarkEventsAsSeen marks all events up to given date in given chat as seen for current customer.
func (a *CustomerAPI) MarkEventsAsSeen(chatID string, seenUpTo time.Time) error {
	return a.Call("mark_events_as_seen", &markEventsAsSeenRequest{
		ChatID:   chatID,
		SeenUpTo: seenUpTo.Format(time.RFC3339Nano),
	}, &emptyResponse{})
}

// GetCustomer returns current Customer.
func (a *CustomerAPI) GetCustomer() (*objects.Customer, error) {
	var resp objects.Customer
	err := a.Call("get_customer", nil, &resp)
	return &resp, err
}
