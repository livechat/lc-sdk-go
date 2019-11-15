package customer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"

	api_errors "github.com/livechat/lc-sdk-go/errors"
	"github.com/livechat/lc-sdk-go/objects"
)

const apiVersion = "v3.1"

// API provides the API operation methods for making requests to Customer Chat API via Web API.
// See this package's package overview docs for details on the service.
type API struct {
	httpClient *http.Client
	// APIURL defines base of API endpoint (without version, component and action).
	APIURL      string
	clientID    string
	tokenGetter func() *Token
}

// Token represents SSO token from Customer Chat API's perspective.
type Token struct {
	// LicenseID specifies ID of license which owns the token.
	LicenseID int
	// AccessToken is a customer access token returned by LiveChat OAuth Server.
	AccessToken string
	// Region is a datacenter for LicenseID (`dal` or `fra`).
	Region string
}

// TokenGetter is called by each API method to obtain valid Token.
// If TokenGetter returns nil, the method won't be executed on Customer Chat API.
type TokenGetter func() *Token

// NewAPI returns ready to use API.
//
// If provided client is nil, then default http client with 20s timeout is used.
func NewAPI(t TokenGetter, client *http.Client, clientID string) (*API, error) {
	if t == nil {
		return nil, errors.New("cannot initialize api without TokenGetter")
	}

	if client == nil {
		client = &http.Client{
			Timeout: 20 * time.Second,
		}
	}

	return &API{
		tokenGetter: t,
		APIURL:      "https://api.livechatinc.com",
		clientID:    clientID,
		httpClient:  client,
	}, nil
}

// StartChat starts new chat with access, properties and initial thread as defined in initialChat.
// It returns respectively chat ID, thread ID and initial event IDs (except for server-generated events).
func (a *API) StartChat(initialChat *InitialChat, continuous bool) (string, string, []string, error) {
	req := &startChatRequest{
		Chat:       initialChat,
		Continuous: continuous,
	}
	var resp startChatResponse
	err := a.call("start_chat", req, &resp)
	return resp.ChatID, resp.ThreadID, resp.EventIDs, err
}

// SendMessage sends event of type message to given chat.
// It returns event ID.
func (a *API) SendMessage(chatID, text string, recipients Recipients) (string, error) {
	e := objects.Message{
		Event: &objects.Event{
			Type:       "message",
			Recipients: string(recipients),
		},
		Text: text,
	}

	return a.SendEvent(chatID, &e)
}

// SendSystemMessage sends event of type system_message to given chat.
// It returns event ID.
func (a *API) SendSystemMessage(chatID, text, messageType string) (string, error) {
	e := objects.SystemMessage{
		Event: objects.Event{
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
func (a *API) SendEvent(chatID string, e interface{}) (string, error) {
	switch v := e.(type) {
	case *objects.Event:
	case *objects.Message:
	case *objects.SystemMessage:
	default:
		return "", fmt.Errorf("event type %T not supported", v)
	}

	var resp sendEventResponse
	err := a.call("send_event", &sendEventRequest{
		ChatID: chatID,
		Event:  e,
	}, &resp)

	return resp.EventID, err
}

// ActivateChat activates chat initialChat.ID with access, properties and initial thread
// as defined in initialChat.
// It returns respectively thread ID and initial event IDs (except for server-generated events).
func (a *API) ActivateChat(initialChat *InitialChat, continuous bool) (string, []string, error) {
	var resp activateChatResponse

	if initialChat.Thread != nil {
		for _, e := range initialChat.Thread.Events {
			switch v := e.(type) {
			case *objects.Event:
			case *objects.Message:
			case *objects.SystemMessage:
			default:
				return "", nil, fmt.Errorf("event type %T not supported", v)
			}
		}
	}

	err := a.call("activate_chat", &activateChatRequest{
		Chat:       initialChat,
		Continuous: continuous,
	}, &resp)

	return resp.ThreadID, resp.EventIDs, err
}

// GetChatsSummary returns chats summary.
func (a *API) GetChatsSummary(offset, limit uint) ([]objects.Chat, uint, error) {
	var resp getChatsSummaryResponse
	err := a.call("get_chats_summary", &getChatsSummaryRequest{
		Limit:  limit,
		Offset: offset,
	}, &resp)

	return resp.Chats, resp.TotalChats, err
}

// GetChatThreadsSummary returns threads summary for given chat.
func (a *API) GetChatThreadsSummary(chatID string, offset, limit uint) ([]ThreadSummary, uint, error) {
	var resp getChatThreadsSummaryResponse
	err := a.call("get_chat_threads_summary", &getChatThreadsSummaryRequest{
		ChatID: chatID,
		Limit:  limit,
		Offset: offset,
	}, &resp)

	return resp.ThreadsSummary, resp.TotalThreads, err
}

// GetChatThreads returns given threads, or all if no threads are provided, for given chat.
func (a *API) GetChatThreads(chatID string, threadIDs ...string) (objects.Chat, error) {
	var resp getChatThreadsResponse
	err := a.call("get_chat_threads", &getChatThreadsRequest{
		ChatID:    chatID,
		ThreadIDs: threadIDs,
	}, &resp)

	return resp.Chat, err
}

// CloseThread closes active thread for given chat. If no thread is active, then this
// method is a no-op.
func (a *API) CloseThread(chatID string) error {
	return a.call("close_thread", &closeThreadRequest{
		ChatID: chatID,
	}, &emptyResponse{})
}

// UploadFile uploads a file to LiveChat CDN.
// Returned URL shall be used in call to SendFile or SendEvent or it'll become invalid
// in about 24 hours.
func (a *API) UploadFile(filename string, file []byte) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	w, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("couldn't create form file: %v", err)
	}
	if _, err := w.Write(file); err != nil {
		return "", fmt.Errorf("couldn't write file to multipart writer: %v", err)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("couldn't close multipart writer: %v", err)
	}
	token := a.tokenGetter()
	if token == nil {
		return "", fmt.Errorf("couldn't get token")
	}
	url := fmt.Sprintf("%s/%s/customer/action/upload_file?license_id=%v", a.APIURL, apiVersion, token.LicenseID)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("couldn't create new POST request to '%v': %v", url, err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Set("User-agent", fmt.Sprintf("GO SDK Application %s", a.clientID))
	req.Header.Set("X-Region", token.Region)

	var resp uploadFileResponse
	err = a.send(req, &resp)
	return resp.URL, err
}

// SendRichMessagePostback sends postback for given rich message event.
func (a *API) SendRichMessagePostback(chatID, threadID, eventID, postbackID string, toggled bool) error {
	return a.call("send_rich_message_postback", &sendRichMessagePostbackRequest{
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
	return a.call("send_sneak_peek", &sendSneakPeekRequest{
		ChatID:        chatID,
		SneakPeekText: text,
	}, &emptyResponse{})
}

// UpdateChatProperties updates given chat's properties.
func (a *API) UpdateChatProperties(chatID string, properties objects.Properties) error {
	return a.call("update_chat_properties", &updateChatPropertiesRequest{
		ChatID:     chatID,
		Properties: properties,
	}, &emptyResponse{})
}

// DeleteChatProperties deletes given chat's properties.
func (a *API) DeleteChatProperties(chatID string, properties map[string][]string) error {
	return a.call("delete_chat_properties", &deleteChatPropertiesRequest{
		ChatID:     chatID,
		Properties: properties,
	}, &emptyResponse{})
}

// UpdateChatThreadProperties updates given chat thread's properties.
func (a *API) UpdateChatThreadProperties(chatID, threadID string, properties objects.Properties) error {
	return a.call("update_chat_thread_properties", &updateChatThreadPropertiesRequest{
		ChatID:     chatID,
		ThreadID:   threadID,
		Properties: properties,
	}, &emptyResponse{})
}

// DeleteChatThreadProperties deletes given chat thread's properties.
func (a *API) DeleteChatThreadProperties(chatID, threadID string, properties map[string][]string) error {
	return a.call("delete_chat_thread_properties", &deleteChatThreadPropertiesRequest{
		ChatID:     chatID,
		ThreadID:   threadID,
		Properties: properties,
	}, &emptyResponse{})
}

// UpdateEventProperties updates given event's properties.
func (a *API) UpdateEventProperties(chatID, threadID, eventID string, properties objects.Properties) error {
	return a.call("update_event_properties", &updateEventPropertiesRequest{
		ChatID:     chatID,
		ThreadID:   threadID,
		EventID:    eventID,
		Properties: properties,
	}, &emptyResponse{})
}

// DeleteEventProperties deletes given event's properties.
func (a *API) DeleteEventProperties(chatID, threadID, eventID string, properties map[string][]string) error {
	return a.call("delete_event_properties", &deleteEventPropertiesRequest{
		ChatID:     chatID,
		ThreadID:   threadID,
		EventID:    eventID,
		Properties: properties,
	}, &emptyResponse{})
}

// UpdateCustomer updates current customer's info.
func (a *API) UpdateCustomer(name, email, avatarURL string, fields map[string]string) error {
	return a.call("update_customer", &updateCustomerRequest{
		Name:   name,
		Email:  email,
		Avatar: avatarURL,
		Fields: fields,
	}, &emptyResponse{})
}

// SetCustomerFields sets current customer's fields.
func (a *API) SetCustomerFields(fields map[string]string) error {
	return a.call("set_customer_fields", &setCustomerFieldsRequest{
		Fields: fields,
	}, &emptyResponse{})
}

// GetGroupsStatus returns status of provided groups.
//
// Possible values are: GroupStatusOnline, GroupStatusOffline and GroupStatusOnlineForQueue.
// GroupStatusUnknown should never be returned.
func (a *API) GetGroupsStatus(groups []int) (map[int]GroupStatus, error) {
	req := &getGroupsStatusRequest{}
	if len(groups) == 0 {
		req.All = true
	} else {
		req.Groups = groups
	}
	var resp getGroupsStatusResponse
	err := a.call("get_groups_status", req, &resp)

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
	return a.call("check_goals", &checkGoalsRequest{
		PageURL:        pageURL,
		GroupID:        groupID,
		CustomerFields: customerFields,
	}, &emptyResponse{})
}

// GetForm returns an empty prechat, postchat or ticket form and indication whether
// the form is enabled on the license.
func (a *API) GetForm(groupID int, formType FormType) (*Form, bool, error) {
	var resp getFormResponse
	err := a.call("get_form", &getFormRequest{
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
	err := a.call("get_predicted_agent", nil, &resp)
	return &resp, err
}

// GetURLDetails returns info on a given URL.
func (a *API) GetURLDetails(url string) (*URLDetails, error) {
	var resp URLDetails
	err := a.call("get_url_details", &getURLDetailsRequest{
		URL: url,
	}, &resp)
	return &resp, err
}

// MarkEventsAsSeen marks all events up to given date in given chat as seen for current customer.
func (a *API) MarkEventsAsSeen(chatID string, seenUpTo time.Time) error {
	return a.call("mark_events_as_seen", &markEventsAsSeenRequest{
		ChatID:   chatID,
		SeenUpTo: seenUpTo.Format(time.RFC3339Nano),
	}, &emptyResponse{})
}

// GetCustomer returns current Customer.
func (a *API) GetCustomer() (*objects.Customer, error) {
	var resp objects.Customer
	err := a.call("get_customer", nil, &resp)
	return &resp, err
}

func (a *API) send(req *http.Request, respPayload interface{}) error {
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		apiErr := &api_errors.ErrAPI{}
		if err := json.Unmarshal(bodyBytes, apiErr); err != nil {
			return fmt.Errorf("couldn't unmarshal error response: %s (code: %d, raw body: %s)", err.Error(), resp.StatusCode, string(bodyBytes))
		}
		if apiErr.Error() == "" {
			return fmt.Errorf("couldn't unmarshal error response (code: %d, raw body: %s)", resp.StatusCode, string(bodyBytes))
		}
		return apiErr
	}

	if err != nil {
		return err
	}

	return json.Unmarshal(bodyBytes, respPayload)
}

func (a *API) call(action string, reqPayload interface{}, respPayload interface{}) error {
	rawBody, err := json.Marshal(reqPayload)
	if err != nil {
		return err
	}
	token := a.tokenGetter()
	if token == nil {
		return fmt.Errorf("couldn't get token")
	}

	url := fmt.Sprintf("%s/%s/customer/action/%s?license_id=%v", a.APIURL, apiVersion, action, token.LicenseID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(rawBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Set("User-agent", fmt.Sprintf("GO SDK Application %s", a.clientID))
	req.Header.Set("X-Region", token.Region)

	return a.send(req, respPayload)
}
