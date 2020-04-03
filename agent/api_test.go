package agent_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/livechat/lc-sdk-go/agent"
	"github.com/livechat/lc-sdk-go/authorization"
	"github.com/livechat/lc-sdk-go/objects"
)

// TEST HELPERS

type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn roundTripFunc) *http.Client {
	return &http.Client{
		Transport: roundTripFunc(fn),
	}
}

func stubTokenGetter() *authorization.Token {
	licenseID := 12345
	return &authorization.Token{
		LicenseID:   &licenseID,
		AccessToken: "access_token",
		Region:      "region",
	}
}

var mockedResponses = map[string]string{
	"start_chat": `{
		"chat_id": "PJ0MRSHTDG",
		"thread_id": "PGDGHT5G"
	}`,
	"activate_chat": `{
		"thread_id": "PGDGHT5G"
	}`,
	"send_event": `{
		"event_id": "K600PKZON8"
	}`,
	"list_chats": `{
		"chats_summary": [
			{
			  "id": "123",
			  "order": 343544565,
			  "last_thread_id": "xyz",
			  "users": [],
			  "properties": {},
			  "access": {},
			  "last_event_per_type": {
				"message": {
				  "thread_id": "K600PKZON8",
				  "thread_order": 3,
				  "event": {}
				},
				"system_message": {
				  "thread_id": "K600PKZON8",
				  "thread_order": 3,
				  "event": {}
				}
			  }
			}
		],
		"found_chats": 1,
		"previous_page_id": "prevpagehash"
	}`,
	`get_chat_threads_summary`: `{
		"threads_summary": [
			{
				"id": "a0c22fdd-fb71-40b5-bfc6-a8a0bc3117f5",
				"order": 2,
				"total_events": 1
			},
			{
				"id": "b0c22fdd-fb71-40b5-bfc6-a8a0bc3117f6",
				"order": 1,
				"total_events": 0
			}
		],
		"found_threads": 2,
		"previous_page_id": "prevpagehash"

	}`,
	"get_chat_threads": `{
		"chat": {
			"id": "PJ0MRSHTDG",
			"users": [],
			"properties": {},
			"access": {},
			"threads": []
		}
	}`,
	"list_archives": `{
		"chats": [
			{
				"id": "PJ0MRSHTDG",
				"users": [],
				"properties": {},
				"access": {},
				"threads": []
			}
		],
		"pagination": {
			"page": 1,
			"total": 3
		}
	}`,
	"deactivate_chat":       `{}`,
	"follow_chat":           `{}`,
	"unfollow_chat":         `{}`,
	"grant_access":          `{}`,
	"revoke_access":         `{}`,
	"set_access":            `{}`,
	"add_user_to_chat":      `{}`,
	"remove_user_from_chat": `{}`,
	"tag_thread":            `{}`,
	"untag_thread":          `{}`,
	"upload_file": `{
		"url": "https://cdn.livechat-static.com/api/file/lc/att/8948324/45a3581b59a7295145c3825c86ec7ab3/image.png"
	}`,
	"send_rich_message_postback": `{}`,
	"update_chat_properties":     `{}`,
	"delete_chat_properties":     `{}`,
	"update_thread_properties":   `{}`,
	"delete_thread_properties":   `{}`,
	"update_event_properties":    `{}`,
	"delete_event_properties":    `{}`,
	"get_customer": `{
		"id": "b7eff798-f8df-4364-8059-649c35c9ed0c",
		"type": "customer",
		"created_at": "2017-10-11T15:19:21.010200Z",
		"name": "John Smith",
		"email": "customer1@example.com",
		"avatar": "example.com/avatars/1.jpg",
		"session_fields": [{
			"custom_key": "custom_value"
		}, {
			"another_custom_key": "another_custom_value"
		}],
		"last_visit": {
			"started_at": "2017-10-12T15:19:21.010200Z",
			"referrer": "http://www.google.com/",
			"ip": "194.181.146.130",
			"user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.116 Safari/537.36",
			"geolocation": {
				"latitude": "-14.6973803",
				"longitude": "-75.1266898",
				"country": "Poland",
				"country_code": "PL",
				"region": "Dolnoslaskie",
				"city": "Wroclaw",
				"timezone": "Europe/Warsaw"
			},
			"last_pages": [
				{
					"opened_at": "2017-10-12T15:19:21.010200Z",
					"url": "https://www.livechatinc.com/",
					"title": "LiveChat - Homepage"
				},
				{
					"opened_at": "2017-10-12T15:19:21.010200Z",
					"url": "https://www.livechatinc.com/tour",
					"title": "LiveChat - Tour"
				}
			]
		},
		"statistics": {
			"chats_count": 3,
			"threads_count": 9,
			"visits_count": 5,
			"page_views_count": 1337,
			"greetings_shown_count": 69,
			"greetings_accepted_count": 42
		},
		"__priv_lc2_customer_id": "S1525771305.dafea66e5c",
		"agent_last_event_created_at": "2017-10-12T15:19:21.010200Z",
		"customer_last_event_created_at": "2017-10-12T15:19:21.010200Z",
		"chat_ids": [
				"PWJ8Y4THAV"
		]
	}`,
	"list_customers": `{
		"customers": [],
		"total_customers": 0,
		"previous_page_id": "prevpagehash"
	}`,
	"create_customer": `{
		"customer_id": "mister_customer"
	}`,
	"update_customer": `{
		"customer": {}
	}`,
	"ban_customer":          `{}`,
	"mark_events_as_seen":   `{}`,
	"update_agent":          `{}`,
	"send_typing_indicator": `{}`,
	"multicast":             `{}`,
	"transfer_chat":         `{}`,
	"list_agents_for_transfer": `[
		{
			"agent_id": "agent1@example.com",
			"total_active_chats": 2
		},
		{
			"agent_id": "agent2@example.com",
			"total_active_chats": 5
		}
	]`,
}

func createMockedResponder(t *testing.T, method string) roundTripFunc {
	return func(req *http.Request) *http.Response {
		createServerError := func(message string) *http.Response {
			responseError := `{
				"error": {
					"type": "MOCK_SERVER_ERROR",
					"message": "` + message + `"
				}
			}`

			return &http.Response{
				StatusCode: 400,
				Body:       ioutil.NopCloser(bytes.NewBufferString(responseError)),
				Header:     make(http.Header),
			}
		}

		if req.URL.String() != "https://api.livechatinc.com/v3.2/agent/action/"+method+"?license_id=12345" {
			t.Errorf("Invalid URL for Agent API request: %s", req.URL.String())
			return createServerError("Invalid URL")
		}

		if authHeader := req.Header.Get("Authorization"); authHeader != "Bearer access_token" {
			t.Errorf("Invalid Authorization header: %s", authHeader)
			return createServerError("Invalid Authorization")
		}

		if regionHeader := req.Header.Get("X-Region"); regionHeader != "region" {
			t.Errorf("Invalid X-Region header: %s", regionHeader)
			return createServerError("Invalid X-Region")
		}

		// TODO: validate also req body

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(mockedResponses[method])),
			Header:     make(http.Header),
		}
	}
}

func createMockedErrorResponder(t *testing.T, method string) func(req *http.Request) *http.Response {
	return func(req *http.Request) *http.Response {
		responseError := `{
			"error": {
				"type": "Validation",
				"message": "Wrong format of request"
			}
		}`

		return &http.Response{
			StatusCode: 400,
			Body:       ioutil.NopCloser(bytes.NewBufferString(responseError)),
			Header:     make(http.Header),
		}
	}
}

func verifyErrorResponse(method string, resp error, t *testing.T) {
	if resp == nil {
		t.Errorf("%v should fail", method)
		return
	}

	if resp.Error() != "API error: Validation - Wrong format of request" {
		t.Errorf("%v failed with wrong error: %v", method, resp)
	}
}

// TESTS OK Cases

func TestRejectAPICreationWithoutTokenGetter(t *testing.T) {
	_, err := agent.NewAPI(nil, nil, "client_id")
	if err == nil {
		t.Errorf("API should not be created without token getter")
	}
}

func TestStartChatShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "start_chat"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	chatID, threadID, _, rErr := api.StartChat(&agent.InitialChat{}, true)
	if rErr != nil {
		t.Errorf("StartChat failed: %v", rErr)
	}
	if chatID != "PJ0MRSHTDG" {
		t.Errorf("Invalid chatID: %v", chatID)
	}

	if threadID != "PGDGHT5G" {
		t.Errorf("Invalid threadID: %v", threadID)
	}
}

func TestSendEventShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "send_event"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	eventID, rErr := api.SendEvent("stubChatID", objects.Event{}, false)
	if rErr != nil {
		t.Errorf("SendEvent failed: %v", rErr)
	}

	if eventID != "K600PKZON8" {
		t.Errorf("Invalid eventID: %v", eventID)
	}
}

func TestActivateChatShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "activate_chat"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	threadID, _, rErr := api.ActivateChat(&agent.InitialChat{}, true)
	if rErr != nil {
		t.Errorf("ActivateChat failed: %v", rErr)
	}

	if threadID != "PGDGHT5G" {
		t.Errorf("Invalid threadID: %v", threadID)
	}
}

func TestGetChatsSummaryShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "list_chats"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	chatsSummary, found, prevPage, nextPage, rErr := api.ListChats(agent.NewChatsFilters(), 0, 20)
	if rErr != nil {
		t.Errorf("ListChats failed: %v", rErr)
	}

	// TODO add better validation

	if chatsSummary == nil {
		t.Errorf("Invalid chats summary")
	}
	if found != 1 {
		t.Errorf("Invalid total chats: %v", found)
	}
	if prevPage != "prevpagehash" {
		t.Errorf("Invalid previous page ID: %v", prevPage)
	}
	if nextPage != "" {
		t.Errorf("Invalid next page ID: %v", nextPage)
	}
}

func TestGetChatThreadsSummaryShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "get_chat_threads_summary"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	threadsSummary, found, prevPage, nextPage, rErr := api.GetChatThreadsSummary("stubChatID", "asc", "pageid", 20)
	if rErr != nil {
		t.Errorf("GetChatThreadsSummary failed: %v", rErr)
	}

	// TODO add better validation

	if threadsSummary == nil {
		t.Errorf("Invalid threads summary")
	}
	if found != 2 {
		t.Errorf("Invalid found threads: %v", found)
	}
	if prevPage != "prevpagehash" {
		t.Errorf("Invalid previous page ID: %v", prevPage)
	}
	if nextPage != "" {
		t.Errorf("Invalid next page ID: %v", nextPage)
	}
}

func TestGetChatThreadsShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "get_chat_threads"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	chat, rErr := api.GetChatThreads("stubChatID", "stubThreadID")
	if rErr != nil {
		t.Errorf("GetChatThreads failed: %v", rErr)
	}

	if chat.ID != "PJ0MRSHTDG" {
		t.Errorf("Received chat.ID invalid: %v", chat.ID)
	}
}

func TestGetArchivesShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "list_archives"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	chats, page, total, rErr := api.ListArchives(agent.NewArchivesFilters(), 1, 20)
	if rErr != nil {
		t.Errorf("GetChatThreads failed: %v", rErr)
	}

	if chats[0].ID != "PJ0MRSHTDG" {
		t.Errorf("Received chat.ID invalid: %v", chats[0].ID)
	}
	if page != 1 {
		t.Errorf("Received pagination.page invalid: %v", page)
	}
	if total != 3 {
		t.Errorf("Received pagination.total invalid: %v", total)
	}
}

func TestCloseThreadShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "deactivate_chat"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.DeactivateChat("stubChatID")
	if rErr != nil {
		t.Errorf("DeactivateChat failed: %v", rErr)
	}
}

func TestFollowChatShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "follow_chat"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.FollowChat("stubChatID")
	if rErr != nil {
		t.Errorf("FollowChat failed: %v", rErr)
	}
}

func TestUnfollowChatShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "unfollow_chat"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.UnfollowChat("stubChatID")
	if rErr != nil {
		t.Errorf("UnfollowChat failed: %v", rErr)
	}
}

func TestUploadFileShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "upload_file"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	fileUrl, rErr := api.UploadFile("filename", []byte{})
	if rErr != nil {
		t.Errorf("UploadFile failed: %v", rErr)
	}

	if fileUrl != "https://cdn.livechat-static.com/api/file/lc/att/8948324/45a3581b59a7295145c3825c86ec7ab3/image.png" {
		t.Errorf("Invalid file URL: %v", fileUrl)
	}
}

func TestGrantAccessShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "grant_access"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.GrantAccess("resource", "id", objects.Access{})
	if rErr != nil {
		t.Errorf("GrantAccess failed: %v", rErr)
	}
}

func TestRevokeAccessShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "revoke_access"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.RevokeAccess("resource", "id", objects.Access{})
	if rErr != nil {
		t.Errorf("RevokeAccess failed: %v", rErr)
	}
}

func TestSetAccessShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "set_access"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.SetAccess("resource", "id", objects.Access{})
	if rErr != nil {
		t.Errorf("SetAccess failed: %v", rErr)
	}
}

func TestAddUserToChatShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "add_user_to_chat"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.AddUserToChat("chat", "user", "agent")
	if rErr != nil {
		t.Errorf("AddUserToChat failed: %v", rErr)
	}
}

func TestRemoveUserFromChatShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "remove_user_from_chat"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.RemoveUserFromChat("chat", "user", "agent")
	if rErr != nil {
		t.Errorf("RemoveUserFromChat failed: %v", rErr)
	}
}

func TestSendRichMessagePostbackShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "send_rich_message_postback"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.SendRichMessagePostback("stubChatID", "stubThreadID", "stubEventID", "stubPostbackID", false)
	if rErr != nil {
		t.Errorf("SendRichMessagePostback failed: %v", rErr)
	}
}

func TestUpdateChatPropertiesShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "update_chat_properties"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.UpdateChatProperties("stubChatID", objects.Properties{})
	if rErr != nil {
		t.Errorf("UpdateChatProperties failed: %v", rErr)
	}
}

func TestDeleteChatPropertiesShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "delete_chat_properties"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.DeleteChatProperties("stubChatID", map[string][]string{})
	if rErr != nil {
		t.Errorf("DeleteChatProperties failed: %v", rErr)
	}
}

func TestUpdateChatThreadPropertiesShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "update_thread_properties"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.UpdateThreadProperties("stubChatID", "stubThreadID", objects.Properties{})
	if rErr != nil {
		t.Errorf("UpdateThreadProperties failed: %v", rErr)
	}
}

func TestDeleteChatThreadPropertiesShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "delete_thread_properties"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.DeleteThreadProperties("stubChatID", "stubThreadID", map[string][]string{})
	if rErr != nil {
		t.Errorf("DeleteThreadProperties failed: %v", rErr)
	}
}

func TestUpdateEventPropertiesShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "update_event_properties"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.UpdateEventProperties("stubChatID", "stubThreadID", "stubEventID", objects.Properties{})
	if rErr != nil {
		t.Errorf("UpdateEventProperties failed: %v", rErr)
	}
}

func TestDeleteEventPropertiesShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "delete_event_properties"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.DeleteEventProperties("stubChatID", "stubThreadID", "stubEventID", map[string][]string{})
	if rErr != nil {
		t.Errorf("DeleteEventProperties failed: %v", rErr)
	}
}

func TestTagChatThreadPropertiesShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "tag_thread"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.TagThread("stubChatID", "stubThreadID", "tag")
	if rErr != nil {
		t.Errorf("TagThread failed: %v", rErr)
	}
}

func TestUntagChatThreadPropertiesShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "untag_thread"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.UntagThread("stubChatID", "stubThreadID", "tag")
	if rErr != nil {
		t.Errorf("UntagThread failed: %v", rErr)
	}
}

func TestGetCustomerShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "get_customer"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	customer, rErr := api.GetCustomer("b7eff798-f8df-4364-8059-649c35c9ed0c")
	if rErr != nil {
		t.Errorf("GetCustomer failed: %v", rErr)
	}

	if customer.ID != "b7eff798-f8df-4364-8059-649c35c9ed0c" {
		t.Errorf("Invalid customer ID: %v", customer.ID)
	}

	if customer.Type != "customer" {
		t.Errorf("Invalid customer type: %v", customer.Type)
	}

	if customer.Name != "John Smith" {
		t.Errorf("Invalid customer name: %v", customer.Name)
	}

	if customer.Email != "customer1@example.com" {
		t.Errorf("Invalid customer email: %v", customer.Email)
	}

	if customer.Avatar != "example.com/avatars/1.jpg" {
		t.Errorf("Invalid customer avatar: %v", customer.Avatar)
	}

	if len(customer.SessionFields) != 2 {
		t.Errorf("Invalid customer session fields: %+v", customer.SessionFields)
	}
}

func TestGetCustomersShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "list_customers"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	customers, total, prevPage, nextPage, rErr := api.ListCustomers(100, "page", "asc", agent.NewCustomersFilters())
	if rErr != nil {
		t.Errorf("ListCustomers failed: %v", rErr)
	}

	if len(customers) != 0 {
		t.Errorf("Invalid customers len: %v", len(customers))
	}
	if total != 0 {
		t.Errorf("Invalid total: %v", total)
	}
	if prevPage != "prevpagehash" {
		t.Errorf("Invalid previous page ID: %v", prevPage)
	}
	if nextPage != "" {
		t.Errorf("Invalid next page ID: %v", nextPage)
	}
}

func TestCreateCustomerShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "create_customer"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	customerID, rErr := api.CreateCustomer("stubName", "stub@mail.com", "http://stub.url", []map[string]string{})
	if rErr != nil {
		t.Errorf("CreateCustomer failed: %v", rErr)
	}

	if customerID != "mister_customer" {
		t.Errorf("Invalid customer ID: %v", customerID)
	}
}
func TestUpdateCustomerShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "update_customer"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	_, rErr := api.UpdateCustomer("mister_customer", "stubName", "stub@mail.com", "http://stub.url", []map[string]string{})
	if rErr != nil {
		t.Errorf("UpdateCustomer failed: %v", rErr)
	}
}

func TestBanCustomerShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "ban_customer"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.BanCustomer("mister_customer", 20)
	if rErr != nil {
		t.Errorf("BanCustomer failed: %v", rErr)
	}
}

func TestUpdateAgentShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "update_agent"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.UpdateAgent("some_agent", "accepting chats")
	if rErr != nil {
		t.Errorf("UpdateAgent failed: %v", rErr)
	}
}

func TestMarkEventsAsSeenShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "mark_events_as_seen"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.MarkEventsAsSeen("stubChatID", time.Time{})
	if rErr != nil {
		t.Errorf("MarkEventsAsSeen failed: %v", rErr)
	}
}

func TestSendTypingIndicatorShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "send_typing_indicator"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.SendTypingIndicator("stubChatID", "all", true)
	if rErr != nil {
		t.Errorf("SendTypingIndicator failed: %v", rErr)
	}
}

func TestMulticastShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "multicast"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.Multicast(agent.MulticastRecipients{}, []byte("{}"), "type")
	if rErr != nil {
		t.Errorf("Multicast failed: %v", rErr)
	}
}

func TestTransferChatShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "transfer_chat"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	ids := make([]interface{}, 1)
	ids[0] = "1"
	rErr := api.TransferChat("stubChatID", "agents", ids, true)
	if rErr != nil {
		t.Errorf("TransferChat failed: %v", rErr)
	}
}

// TESTS Error Cases

func TestStartChatShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "start_chat"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	_, _, _, rErr := api.StartChat(&agent.InitialChat{}, true)
	verifyErrorResponse("StartChat", rErr, t)
}

func TestSendEventShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "send_event"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	_, rErr := api.SendEvent("stubChatID", &objects.Event{}, false)
	verifyErrorResponse("SendEvent", rErr, t)
}

func TestActivateChatShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "activate_chat"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	_, _, rErr := api.ActivateChat(&agent.InitialChat{}, true)
	verifyErrorResponse("ActivateChat", rErr, t)
}

func TestGetChatsSummaryShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "list_chats"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	_, _, _, _, rErr := api.ListChats(agent.NewChatsFilters(), 0, 20)
	verifyErrorResponse("ListChats", rErr, t)
}

func TestGetChatThreadsSummaryShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "get_chat_threads_summary"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	_, _, _, _, rErr := api.GetChatThreadsSummary("stubChatID", "asc", "pageid", 20)
	verifyErrorResponse("GetChatThreadsSummary", rErr, t)
}

func TestGetChatThreadsShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "get_chat_threads"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	_, rErr := api.GetChatThreads("stubChatID", "stubThreadID")
	verifyErrorResponse("GetChatThreads", rErr, t)
}

func TestGetArchivesShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "list_archives"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	_, _, _, rErr := api.ListArchives(agent.NewArchivesFilters(), 1, 20)
	verifyErrorResponse("ListArchives", rErr, t)
}

func TestCloseThreadShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "deactivate_chat"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.DeactivateChat("stubChatID")
	verifyErrorResponse("DeactivateChat", rErr, t)
}

func TestFollowChatShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "follow_chat"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.FollowChat("stubChatID")
	verifyErrorResponse("FollowChat", rErr, t)
}

func TestUnfollowChatShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "unfollow_chat"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.UnfollowChat("stubChatID")
	verifyErrorResponse("UnfollowChat", rErr, t)
}

func TestUploadFileShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "upload_file"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	_, rErr := api.UploadFile("filename", []byte{})
	verifyErrorResponse("UploadFile", rErr, t)

}

func TestGrantAccessShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "grant_access"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.GrantAccess("resource", "id", objects.Access{})
	verifyErrorResponse("GrantAccess", rErr, t)
}

func TestRevokeAccessShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "revoke_access"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.RevokeAccess("resource", "id", objects.Access{})
	verifyErrorResponse("RevokeAccess", rErr, t)
}

func TestSetAccessShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "set_access"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.SetAccess("resource", "id", objects.Access{})
	verifyErrorResponse("SetAccess", rErr, t)
}

func TestAddUserToChatShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "add_user_to_chat"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.AddUserToChat("chat", "user", "agent")
	verifyErrorResponse("AddUserToChat", rErr, t)
}

func TestRemoveUserFromChatShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "remove_user_from_chat"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.RemoveUserFromChat("chat", "user", "agent")
	verifyErrorResponse("RemoveUserFromChat", rErr, t)

}

func TestSendRichMessagePostbackShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "send_rich_message_postback"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.SendRichMessagePostback("stubChatID", "stubThreadID", "stubEventID", "stubPostbackID", false)
	verifyErrorResponse("SendRichMessagePostback", rErr, t)
}

func TestUpdateChatPropertiesShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "update_chat_properties"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.UpdateChatProperties("stubChatID", objects.Properties{})
	verifyErrorResponse("UpdateChatProperties", rErr, t)
}

func TestDeleteChatPropertiesShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "delete_chat_properties"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.DeleteChatProperties("stubChatID", map[string][]string{})
	verifyErrorResponse("DeleteChatProperties", rErr, t)
}

func TestUpdateChatThreadPropertiesShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "update_thread_properties"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.UpdateThreadProperties("stubChatID", "stubThreadID", objects.Properties{})
	verifyErrorResponse("UpdateThreadProperties", rErr, t)
}

func TestDeleteChatThreadPropertiesShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "delete_thread_properties"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.DeleteThreadProperties("stubChatID", "stubThreadID", map[string][]string{})
	verifyErrorResponse("DeleteThreadProperties", rErr, t)
}

func TestUpdateEventPropertiesShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "update_event_properties"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.UpdateEventProperties("stubChatID", "stubThreadID", "stubEventID", objects.Properties{})
	verifyErrorResponse("UpdateEventProperties", rErr, t)
}

func TestDeleteEventPropertiesShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "delete_event_properties"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.DeleteEventProperties("stubChatID", "stubThreadID", "stubEventID", map[string][]string{})
	verifyErrorResponse("DeleteEventProperties", rErr, t)
}

func TestTagChatThreadPropertiesShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "tag_thread"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.TagThread("stubChatID", "stubThreadID", "tag")
	verifyErrorResponse("TagThread", rErr, t)
}

func TestUntagChatThreadPropertiesShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "untag_thread"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.UntagThread("stubChatID", "stubThreadID", "tag")
	verifyErrorResponse("UntagThread", rErr, t)
}

func TesGetCustomersShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "list_customers"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	_, _, _, _, rErr := api.ListCustomers(100, "page", "asc", agent.NewCustomersFilters())
	verifyErrorResponse("ListCustomers", rErr, t)
}

func TestCreateCustomerShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "create_customer"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	_, rErr := api.CreateCustomer("stubName", "stub@mail.com", "http://stub.url", []map[string]string{})
	verifyErrorResponse("CreateCustomer", rErr, t)
}
func TestUpdateCustomerShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "update_customer"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	_, rErr := api.UpdateCustomer("mister_customer", "stubName", "stub@mail.com", "http://stub.url", []map[string]string{})
	verifyErrorResponse("UpdateCustomer", rErr, t)
}

func TestBanCustomerShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "ban_customer"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.BanCustomer("mister_customer", 20)
	verifyErrorResponse("BanCustomer", rErr, t)
}

func TestUpdateAgentShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "update_agent"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.UpdateAgent("some_agent", "accepting chats")
	verifyErrorResponse("UpdateAgent", rErr, t)
}

func TestMarkEventsAsSeenShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "mark_events_as_seen"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.MarkEventsAsSeen("stubChatID", time.Time{})
	verifyErrorResponse("MarkEventsAsSeen", rErr, t)
}

func TestSendTypingIndicatorShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "send_typing_indicator"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.SendTypingIndicator("stubChatID", "all", true)
	verifyErrorResponse("SendTypingIndicator", rErr, t)
}

func TestMulticastShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "multicast"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.Multicast(agent.MulticastRecipients{}, []byte("{}"), "type")
	verifyErrorResponse("Multicast", rErr, t)
}

func TestTransferChatShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "transfer_chat"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}
	ids := make([]interface{}, 1)
	ids[0] = 1
	rErr := api.TransferChat("stubChatID", "group", ids, false)
	verifyErrorResponse("SendTypingIndicator", rErr, t)
}

func TestListAgentsForTransferShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "list_agents_for_transfer"))

	api, err := agent.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	resp, rErr := api.ListAgentsForTransfer("PJ0MRSHTDG")
	if rErr != nil {
		t.Errorf("ListAgentsForTransfer failed: %v", rErr)
	}

	if len(resp) != 2 {
		t.Errorf("Invalid ListAgentsForTransfer response: %v", resp)
	}
}
