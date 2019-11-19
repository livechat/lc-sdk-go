package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

	url := fmt.Sprintf("%s/%s/agent/action/%s?license_id=%v", a.APIURL, apiVersion, action, token.LicenseID)
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

func (a *API) validateInitialChat(chat *InitialChat) error {
	if chat.Thread != nil {
		for _, e := range chat.Thread.Events {
			switch v := e.(type) {
			case *objects.Event:
			case *objects.Message:
			case *objects.SystemMessage:
			default:
				return fmt.Errorf("event type %T not supported", v)
			}
		}
	}
	return nil
}

func (a *API) GetChatsSummary(filters *ChatsFilters, page, limit uint64) ([]ChatSummary, uint, error) {
	var resp getChatsSummaryResponse
	err := a.call("get_chats_summary", &getChatsSummaryRequest{
		Filters: filters,
		Pagination: paginationRequest{
			Page:  page,
			Limit: limit,
		},
	}, &resp)

	return resp.ChatsSummary, resp.FoundChats, err
}

func (a *API) GetChatThreadsSummary(chatID, order, pageID string, limit uint64) ([]ThreadSummary, int, string, string, error) {
	var resp getChatThreadsSummaryResponse
	err := a.call("get_chat_threads_summary", &getChatThreadsSummaryRequest{
		ChatID: chatID,
		hashedPaginationRequest: &hashedPaginationRequest{
			Order:  order,
			Limit:  limit,
			PageID: pageID,
		},
	}, &resp)

	return resp.ThreadsSummary, resp.FoundThreads, resp.PreviousPageID, resp.NextPageID, err
}

func (a *API) GetChatThreads(chatID string, threadIDs ...string) (objects.Chat, error) {
	var resp getChatThreadsResponse
	err := a.call("get_chat_threads", &getChatThreadsRequest{
		ChatID:    chatID,
		ThreadIDs: threadIDs,
	}, &resp)

	return resp.Chat, err
}

func (a *API) GetArchives(filters *ArchivesFilters, page, limit uint64) ([]objects.Chat, uint64, uint64, error) {
	var resp getArchivesResponse
	err := a.call("get_archives", &getArchivesRequest{
		Filters: filters,
		Pagination: paginationRequest{
			Page:  page,
			Limit: limit,
		},
	}, &resp)

	return resp.Chats, resp.Pagination.Page, resp.Pagination.Total, err
}

func (a *API) StartChat(initialChat *InitialChat, continuous bool) (string, string, []string, error) {
	var resp startChatResponse

	if e := a.validateInitialChat(initialChat); e != nil {
		return "", "", nil, e
	}

	err := a.call("start_chat", &startChatRequest{
		Chat:       initialChat,
		Continuous: continuous,
	}, &resp)
	return resp.ChatID, resp.ThreadID, resp.EventIDs, err
}

func (a *API) ActivateChat(initialChat *InitialChat, continuous bool) (string, []string, error) {
	var resp activateChatResponse

	if e := a.validateInitialChat(initialChat); e != nil {
		return "", nil, e
	}

	err := a.call("activate_chat", &activateChatRequest{
		Chat:       initialChat,
		Continuous: continuous,
	}, &resp)

	return resp.ThreadID, resp.EventIDs, err
}

func (a *API) CloseThread(chatID string) error {
	return a.call("close_thread", &closeThreadRequest{
		ChatID: chatID,
	}, &emptyResponse{})
}

func (a *API) FollowChat(chatID string) error {
	return a.call("follow_chat", &followChatRequest{
		ChatID: chatID,
	}, &emptyResponse{})
}

func (a *API) UnfollowChat(chatID string) error {
	return a.call("unfollow_chat", &unfollowChatRequest{
		ChatID: chatID,
	}, &emptyResponse{})
}

func (a *API) GrantAccess(resource, id string, access objects.Access) error {
	return a.call("grant_access", &modifyAccessRequest{
		Resource: resource,
		ID:       id,
		Access:   access,
	}, &emptyResponse{})
}

func (a *API) RevokeAccess(resource, id string, access objects.Access) error {
	return a.call("revoke_access", &modifyAccessRequest{
		Resource: resource,
		ID:       id,
		Access:   access,
	}, &emptyResponse{})
}

func (a *API) SetAccess(resource, id string, access objects.Access) error {
	return a.call("set_access", &modifyAccessRequest{
		Resource: resource,
		ID:       id,
		Access:   access,
	}, &emptyResponse{})
}

func (a *API) TransferChat(chatID, targetType string, ids []uint, force bool) error {
	return a.call("transfer_chat", &transferChatRequest{
		ChatID: chatID,
		Target: target{
			Type: targetType,
			IDs:  ids,
		},
		Force: force,
	}, &emptyResponse{})
}

func (a *API) AddUserToChat(chatID, userID, userType string) error {
	return a.call("add_user_to_chat", &changeChatUsersRequest{
		ChatID:   chatID,
		UserID:   userID,
		UserType: userType,
	}, &emptyResponse{})
}

func (a *API) RemoveUserFromChat(chatID, userID, userType string) error {
	return a.call("remove_user_from_chat", &changeChatUsersRequest{
		ChatID:   chatID,
		UserID:   userID,
		UserType: userType,
	}, &emptyResponse{})
}

func (a *API) SendEvent(chatID string, event objects.Event, attachToLastThread bool) (string, error) {
	var resp sendEventResponse
	err := a.call("send_event", &sendEventRequest{
		ChatID:             chatID,
		Event:              event,
		AttachToLastThread: attachToLastThread,
	}, &resp)

	return resp.EventID, err
}

// func (a *API) UploadFile() {use common}

func (a *API) SendRichMessagePostback(chatID, eventID, threadID, postbackID string, toggled bool) error {
	return a.call("send_rich_message_postback", &sendRichMessagePostbackRequest{
		ChatID:   chatID,
		EventID:  eventID,
		ThreadID: threadID,
		Postback: postback{
			ID:      postbackID,
			Toggled: toggled,
		},
	}, &emptyResponse{})
}

func (a *API) UpdateChatProperties(chatID string, properties objects.Properties) error {
	return a.call("update_chat_properties", &updateChatPropertiesRequest{
		ChatID:     chatID,
		Properties: properties,
	}, &emptyResponse{})
}

func (a *API) DeleteChatProperties(chatID string, properties objects.Properties) error {
	return a.call("delete_chat_properties", &deleteChatPropertiesRequest{
		ChatID:     chatID,
		Properties: properties,
	}, &emptyResponse{})
}

func (a *API) UpdateChatThreadProperties(chatID, threadID string, properties objects.Properties) error {
	return a.call("update_chat_thread_properties", &updateChatThreadPropertiesRequest{
		ChatID:     chatID,
		ThreadID:   threadID,
		Properties: properties,
	}, &emptyResponse{})
}

func (a *API) DeleteChatThreadProperties(chatID, threadID string, properties objects.Properties) error {
	return a.call("delete_chat_thread_properties", &deleteChatThreadPropertiesRequest{
		ChatID:     chatID,
		ThreadID:   threadID,
		Properties: properties,
	}, &emptyResponse{})
}

func (a *API) UpdateEventProperties(chatID, threadID, eventID string, properties objects.Properties) error {
	return a.call("update_event_properties", &updateEventPropertiesRequest{
		ChatID:     chatID,
		ThreadID:   threadID,
		EventID:    eventID,
		Properties: properties,
	}, &emptyResponse{})
}

func (a *API) DeleteEventProperties(chatID, threadID, eventID string, properties objects.Properties) error {
	return a.call("delete_event_properties", &deleteEventPropertiesRequest{
		ChatID:     chatID,
		ThreadID:   threadID,
		EventID:    eventID,
		Properties: properties,
	}, &emptyResponse{})
}

func (a *API) TagChatThread(chatID, threadID, tag string) error {
	return a.call("tag_chat_thread", &changeChatThreadTagRequest{
		ChatID:   chatID,
		ThreadID: threadID,
		Tag:      tag,
	}, &emptyResponse{})
}

func (a *API) UntagChatThread(chatID, threadID, tag string) error {
	return a.call("untag_chat_thread", &changeChatThreadTagRequest{
		ChatID:   chatID,
		ThreadID: threadID,
		Tag:      tag,
	}, &emptyResponse{})
}

func (a *API) GetCustomers(limit uint, pageID, order string, filters *CustomersFilters) ([]objects.Customer, uint64, string, string, error) {
	var resp getCustomersResponse
	err := a.call("get_customers", &getCustomersRequest{
		PageID:  pageID,
		Limit:   limit,
		Order:   order,
		Filters: filters,
	}, &resp)

	return resp.Customers, resp.TotalCustomers, resp.PreviousPageID, resp.NextPageID, err
}

func (a *API) CreateCustomer(name, email, avatar string, fields map[string]string) (string, error) {
	var resp createCustomerResponse
	err := a.call("create_customer", &createCustomerRequest{
		Name:   name,
		Email:  email,
		Avatar: avatar,
		Fields: fields,
	}, &resp)

	return resp.CustomerID, err
}

func (a *API) UpdateCustomer(customerID, name, email, avatar string, fields map[string]string) (objects.Customer, error) {
	var resp updateCustomerResponse
	err := a.call("update_customer", &updateCustomerRequest{
		CustomerID: customerID,
		Name:       name,
		Email:      email,
		Avatar:     avatar,
		Fields:     fields,
	}, &resp)

	return resp.Customer, err
}

func (a *API) BanCustomer(customerID string, days uint64) error {
	return a.call("ban_customer", &banCustomerRequest{
		CustomerID: customerID,
		Ban: ban{
			Days: days,
		},
	}, &emptyResponse{})
}

func (a *API) UpdateAgent(agentID, routingStatus string) error {
	return a.call("update_agent", &updateAgentRequest{
		AgentID:       agentID,
		RoutingStatus: routingStatus,
	}, &emptyResponse{})
}

func (a *API) MarkEventsAsSeen(chatID string, seenUpTo time.Time) error {
	return a.call("mark_events_as_seen", &markEventsAsSeenRequest{
		ChatID:   chatID,
		SeenUpTo: seenUpTo.Format(time.RFC3339Nano),
	}, &emptyResponse{})
}

func (a *API) SendTypingIndicator(chatID, recipients string, isTyping bool) error {
	return a.call("send_typing_indicator", &sendTypingIndicatorRequest{
		ChatID:     chatID,
		Recipients: recipients,
		IsTyping:   isTyping,
	}, &emptyResponse{})
}

func (a *API) Multicast(scopes Scopes, content json.RawMessage, multicastType string) error {
	return a.call("multicast", &multicastRequest{
		Scopes:  scopes,
		Content: content,
		Type:    multicastType,
	}, &emptyResponse{})
}
