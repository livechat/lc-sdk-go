package customer

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/livechat/lc-sdk-go/internal"

	"github.com/livechat/lc-sdk-go/objects/authorization"

	"github.com/livechat/lc-sdk-go/objects/events"
)

const apiVersion = "v3.2"

type API struct {
	base internal.APIBase
}

func NewAPI(t authorization.TokenGetter, client *http.Client, clientID string) (*API, error) {
	if t == nil {
		return nil, errors.New("cannot initialize api without TokenGetter")
	}

	if client == nil {
		client = &http.Client{
			Timeout: 20 * time.Second,
		}
	}

	return &API{
		internal.APIBase{
			ApiVersion:  apiVersion,
			ApiName:     "customer",
			TokenGetter: t,
			ApiURL:      "https://api.livechatinc.com",
			ClientID:    clientID,
			HttpClient:  client,
		},
	}, nil
}

func (a *API) StartChat(initialChat *InitialChat, continuous bool) (string, string, []string, error) {
	req := &startChatRequest{
		Chat:       initialChat,
		Continuous: continuous,
	}
	var resp startChatResponse
	err := a.base.Call("start_chat", req, &resp)
	return resp.ChatID, resp.ThreadID, resp.EventIDs, err
}

func (a *API) SendMessage(chatID, text string, recipients Recipients) (string, error) {
	e := events.Message{
		Event: &events.Event{
			Type:       "message",
			Recipients: string(recipients),
		},
		Text: text,
	}

	return a.SendEvent(chatID, e)
}

func (a *API) SendSystemMessage(chatID, text, messageType string) (string, error) {
	e := events.SystemMessage{
		Event: events.Event{
			Type: "system_message",
		},
		Text: text,
		Type: messageType,
	}

	return a.SendEvent(chatID, e)
}

func (a *API) SendEvent(chatID string, e interface{}) (string, error) {
	switch v := e.(type) {
	case *events.Event:
	case *events.Message:
	case *events.SystemMessage:
	default:
		return "", fmt.Errorf("event type %T not supported", v)
	}

	var resp sendEventResponse
	err := a.base.Call("send_event", &sendEventRequest{
		ChatID: chatID,
		Event:  e,
	}, &resp)

	return resp.EventID, err
}

func (a *API) ActivateChat(initialChat *InitialChat, continuous bool) (string, []string, error) {
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

	err := a.base.Call("activate_chat", &activateChatRequest{
		Chat:       initialChat,
		Continuous: continuous,
	}, &resp)

	return resp.ThreadID, resp.EventIDs, err
}

func (a *API) GetChatsSummary(offset, limit uint) ([]Chat, uint, error) {
	var resp getChatsSummaryResponse
	err := a.base.Call("get_chats_summary", &getChatsSummaryRequest{
		Limit:  limit,
		Offset: offset,
	}, &resp)

	return resp.Chats, resp.TotalChats, err
}

func (a *API) GetChatThreadsSummary(chatID string, offset, limit uint) ([]ThreadSummary, uint, error) {
	var resp getChatThreadsSummaryResponse
	err := a.base.Call("get_chat_threads_summary", &getChatThreadsSummaryRequest{
		ChatID: chatID,
		Limit:  limit,
		Offset: offset,
	}, &resp)

	return resp.ThreadsSummary, resp.TotalThreads, err
}

func (a *API) GetChatThreads(chatID string, threadIDs ...string) (Chat, error) {
	var resp getChatThreadsResponse
	err := a.base.Call("get_chat_threads", &getChatThreadsRequest{
		ChatID:    chatID,
		ThreadIDs: threadIDs,
	}, &resp)

	return resp.Chat, err
}

func (a *API) CloseThread(chatID string) error {
	return a.base.Call("close_thread", &closeThreadRequest{
		ChatID: chatID,
	}, &emptyResponse{})
}

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
	token := a.base.TokenGetter()
	if token == nil {
		return "", fmt.Errorf("couldn't get token")
	}
	url := fmt.Sprintf("%s/%s/customer/action/upload_file?license_id=%v", a.base.ApiURL, apiVersion, token.LicenseID)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("couldn't create new POST request to '%v': %v", url, err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Set("User-agent", fmt.Sprintf("GO SDK Application %s", a.base.ClientID))
	req.Header.Set("X-Region", token.Region)

	var resp uploadFileResponse
	err = a.base.Send(req, &resp)
	return resp.URL, err
}

func (a *API) SendRichMessagePostback(chatID, threadID, eventID, postbackID string, toggled bool) error {
	return a.base.Call("send_rich_message_postback", &sendRichMessagePostbackRequest{
		ChatID:   chatID,
		ThreadID: threadID,
		EventID:  eventID,
		Postback: postback{
			ID:      postbackID,
			Toggled: toggled,
		},
	}, &emptyResponse{})
}

func (a *API) SendSneakPeek(chatID, text string) error {
	return a.base.Call("send_sneak_peek", &sendSneakPeekRequest{
		ChatID:        chatID,
		SneakPeekText: text,
	}, &emptyResponse{})
}

func (a *API) UpdateChatProperties(chatID string, properties Properties) error {
	return a.base.Call("update_chat_properties", &updateChatPropertiesRequest{
		ChatID:     chatID,
		Properties: properties,
	}, &emptyResponse{})
}

func (a *API) DeleteChatProperties(chatID string, properties map[string][]string) error {
	return a.base.Call("delete_chat_properties", &deleteChatPropertiesRequest{
		ChatID:     chatID,
		Properties: properties,
	}, &emptyResponse{})
}

func (a *API) UpdateChatThreadProperties(chatID, threadID string, properties Properties) error {
	return a.base.Call("update_chat_thread_properties", &updateChatThreadPropertiesRequest{
		ChatID:     chatID,
		ThreadID:   threadID,
		Properties: properties,
	}, &emptyResponse{})
}

func (a *API) DeleteChatThreadProperties(chatID, threadID string, properties map[string][]string) error {
	return a.base.Call("delete_chat_thread_properties", &deleteChatThreadPropertiesRequest{
		ChatID:     chatID,
		ThreadID:   threadID,
		Properties: properties,
	}, &emptyResponse{})
}

func (a *API) UpdateEventProperties(chatID, threadID, eventID string, properties Properties) error {
	return a.base.Call("update_event_properties", &updateEventPropertiesRequest{
		ChatID:     chatID,
		ThreadID:   threadID,
		EventID:    eventID,
		Properties: properties,
	}, &emptyResponse{})
}

func (a *API) DeleteEventProperties(chatID, threadID, eventID string, properties map[string][]string) error {
	return a.base.Call("delete_event_properties", &deleteEventPropertiesRequest{
		ChatID:     chatID,
		ThreadID:   threadID,
		EventID:    eventID,
		Properties: properties,
	}, &emptyResponse{})
}

func (a *API) UpdateCustomer(name, email, avatarURL string, fields map[string]string) error {
	return a.base.Call("update_customer", &updateCustomerRequest{
		Name:   name,
		Email:  email,
		Avatar: avatarURL,
		Fields: fields,
	}, &emptyResponse{})
}

func (a *API) SetCustomerFields(fields map[string]string) error {
	return a.base.Call("set_customer_fields", &setCustomerFieldsRequest{
		Fields: fields,
	}, &emptyResponse{})
}

func (a *API) GetGroupsStatus(groups []int) (map[int]GroupStatus, error) {
	req := &getGroupsStatusRequest{}
	if len(groups) == 0 {
		req.All = true
	} else {
		req.Groups = groups
	}
	var resp getGroupsStatusResponse
	err := a.base.Call("get_groups_status", req, &resp)

	r := map[int]GroupStatus{}

	if err == nil {
		for g, s := range resp.Status {
			r[g] = toGroupStatus(s)
		}
	}

	return r, err
}

func (a *API) CheckGoals(pageURL string, groupID int, customerFields map[string]string) error {
	return a.base.Call("check_goals", &checkGoalsRequest{
		PageURL:        pageURL,
		GroupID:        groupID,
		CustomerFields: customerFields,
	}, &emptyResponse{})
}

func (a *API) GetForm(groupID int, formType FormType) (*Form, bool, error) {
	var resp getFormResponse
	err := a.base.Call("get_form", &getFormRequest{
		GroupID: groupID,
		Type:    string(formType),
	}, &resp)

	return resp.Form, resp.Enabled, err
}

func (a *API) GetPredictedAgent() (*PredictedAgent, error) {
	var resp PredictedAgent
	err := a.base.Call("get_predicted_agent", nil, &resp)
	return &resp, err
}

func (a *API) GetURLDetails(url string) (*URLDetails, error) {
	var resp URLDetails
	err := a.base.Call("get_url_details", &getURLDetailsRequest{
		URL: url,
	}, &resp)
	return &resp, err
}

func (a *API) MarkEventsAsSeen(chatID string, seenUpTo time.Time) error {
	return a.base.Call("mark_events_as_seen", &markEventsAsSeenRequest{
		ChatID:   chatID,
		SeenUpTo: seenUpTo.Format(time.RFC3339Nano),
	}, &emptyResponse{})
}

func (a *API) GetCustomer() (*Customer, error) {
	var resp Customer
	err := a.base.Call("get_customer", nil, &resp)
	return &resp, err
}
