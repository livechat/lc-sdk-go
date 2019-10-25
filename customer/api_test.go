package customer_test

import (
	"time"
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/livechat/lc-sdk-go/customer"
	"github.com/livechat/lc-sdk-go/objects/events"
)

// TEST HELPERS

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func stubTokenGetter() *customer.Token {
	return &customer.Token{
		LicenseID: 12345,
		AccessToken: "access_token",
		Region: "region",
	}
}

var mockedResponses = map[string]string {
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
	"get_chats_summary": `{
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
		"total_chats": 1
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
		"total_threads": 2
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
	"close_thread": `{}`,
	"upload_file": `{
		"url": "https://cdn.livechat-static.com/api/file/lc/att/8948324/45a3581b59a7295145c3825c86ec7ab3/image.png"
	}`,
	"send_rich_message_postback": `{}`,
	"send_sneak_peek": `{}`,
	"update_chat_properties": `{}`,
	"delete_chat_properties": `{}`,
	"update_chat_thread_properties": `{}`,
	"delete_chat_thread_properties": `{}`,
	"update_event_properties": `{}`,
	"delete_event_properties": `{}`,
	"update_customer": `{}`,
	"set_customer_fields": `{}`,
	"get_groups_status": `{
		"groups_status": {
			"1": "online",
			"2": "offline",
			"3": "online_for_queue"
		}
	}`,
	"check_goals": `{}`,
	"get_form": `{
		"form": {
			"id": "156630109416307809",
			"fields": [
			  {
				"id": "15663010941630615",
				"type": "header",
				"label": "Welcome to our LiveChat! Please fill in the form below before starting the chat."
			  },
			  {
				"id": "156630109416307759",
				"type": "name",
				"label": "Name:",
				"required": false
			  },
			  {
				"id": "15663010941630515",
				"type": "email",
				"label": "E-mail:",
				"required": false
			  }
			]
		},
		"enabled": true
	}`,
	"get_predicted_agent": `{
		"agent": {
			"id": "agent1@example.com",
			"name": "Name",
			"avatar": "https://example.avatar/example.com",
			"is_bot": false,
			"job_title": "support hero",
			"type": "agent"
		}
	}`,
	"get_url_details": `{
		"title": "LiveChat | Live Chat Software and Help Desk Software",
		"description": "LiveChat - premium live chat software and help desk software for business. Over 24 000 companies from 150 countries use LiveChat. Try now, chat for free!",
		"image_url": "s3.eu-central-1.amazonaws.com/labs-fraa-livechat-thumbnails/96979c3552cf3fa4ae326086a3048d9354c27324.png",
		"image_width": 200,
		"image_height": 200,
		"url": "https://livechatinc.com"
	}`,
	"mark_events_as_seen": `{}`,
	"get_customer": `{}`, //TODO - create some real structure here
}

func createMockedResponder(t *testing.T, method string) func(req *http.Request) *http.Response {
	return func(req *http.Request) *http.Response {
		createServerError := func (message string) *http.Response {
			responseError := `{
				"error": {
					"type": "MOCK_SERVER_ERROR",
					"message": ` + message + `
				}
			}`

			return &http.Response{
				StatusCode: 400,
				Body: ioutil.NopCloser(bytes.NewBufferString(responseError)),
				Header: make(http.Header),
			}
		}

		if req.URL.String() != "https://api.livechatinc.com/v3.2/customer/action/" + method + "?license_id=12345" {
			t.Errorf("Invalid URL for Customer API request: %s", req.URL.String())
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

// TESTS

func TestRejectAPICreationWithoutTokenGetter(t *testing.T) {
	_, err := customer.NewAPI(nil, nil, "client_id")
	if err == nil {
		t.Errorf("API should not be created without token getter")
	}
}

func TestStartChatOK(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "start_chat"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	chatID, threadID, _, rErr := api.StartChat(&customer.InitialChat{}, true)
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

func TestSendEventOK(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "send_event"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	eventID, rErr := api.SendEvent("stubChatID", &events.Event{})
	if rErr != nil {
		t.Errorf("SendEvent failed: %v", rErr)
	}

	if eventID != "K600PKZON8" {
		t.Errorf("Invalid eventID: %v", eventID)
	}
}

func TestActivateChatOK(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "activate_chat"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	threadID, _, rErr := api.ActivateChat(&customer.InitialChat{}, true)
	if rErr != nil {
		t.Errorf("ActivateChat failed: %v", rErr)
	}

	if threadID != "PGDGHT5G" {
		t.Errorf("Invalid threadID: %v", threadID)
	}
}

func TestGetChatsSummaryOK(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "get_chats_summary"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	chatsSummary, total, rErr := api.GetChatsSummary(0, 20)
	if rErr != nil {
		t.Errorf("GetChatsSummary failed: %v", rErr)
	}

	// TODO add better validation

	if chatsSummary == nil {
		t.Errorf("Invalid chats summary")
	}
	if total != 1 {
		t.Errorf("Invalid total chats: %v", total)
	}
}

func TestGetChatThreadsSummaryOK(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "get_chat_threads_summary"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	threadsSummary, total, rErr := api.GetChatThreadsSummary("stubChatID", 0, 20)
	if rErr != nil {
		t.Errorf("GetChatThreadsSummary failed: %v", rErr)
	}

	// TODO add better validation

	if threadsSummary == nil {
		t.Errorf("Invalid threads summary")
	}
	if total != 2 {
		t.Errorf("Invalid total chats: %v", total)
	}
}

func TestGetChatThreadsOK(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "get_chat_threads"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
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

func TestCloseThreadOK(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "close_thread"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.CloseThread("stubChatID")
	if rErr != nil {
		t.Errorf("CloseThread failed: %v", rErr)
	}
}

func TestUploadFileOK(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "upload_file"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
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

func TestSendRichMessagePostbackOK(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "send_rich_message_postback"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.SendRichMessagePostback("stubChatID", "stubThreadID", "stubEventID", "stubPostbackID", false)
	if rErr != nil {
		t.Errorf("SendRichMessagePostback failed: %v", rErr)
	}
}

func TestSendSneakPeekOK(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "send_sneak_peek"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.SendSneakPeek("stubChatID", "sneaky freaky baby")
	if rErr != nil {
		t.Errorf("SendSneakPeek failed: %v", rErr)
	}
}

func TestUpdateChatPropertiesOK(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "update_chat_properties"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.UpdateChatProperties("stubChatID", customer.Properties{})
	if rErr != nil {
		t.Errorf("UpdateChatProperties failed: %v", rErr)
	}
}

func TestDeleteChatPropertiesOK(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "delete_chat_properties"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.DeleteChatProperties("stubChatID", map[string][]string{})
	if rErr != nil {
		t.Errorf("DeleteChatProperties failed: %v", rErr)
	}
}

func TestUpdateChatThreadPropertiesOK(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "update_chat_thread_properties"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.UpdateChatThreadProperties("stubChatID", "stubThreadID", customer.Properties{})
	if rErr != nil {
		t.Errorf("UpdateChatThreadProperties failed: %v", rErr)
	}
}

func TestDeleteChatThreadPropertiesOK(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "delete_chat_thread_properties"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.DeleteChatThreadProperties("stubChatID", "stubThreadID", map[string][]string{})
	if rErr != nil {
		t.Errorf("DeleteChatThreadProperties failed: %v", rErr)
	}
}

func TestUpdateEventPropertiesOK(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "update_event_properties"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.UpdateEventProperties("stubChatID", "stubThreadID", "stubEventID", customer.Properties{})
	if rErr != nil {
		t.Errorf("UpdateEventProperties failed: %v", rErr)
	}
}

func TestDeleteEventPropertiesOK(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "delete_event_properties"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.DeleteEventProperties("stubChatID", "stubThreadID", "stubEventID", map[string][]string{})
	if rErr != nil {
		t.Errorf("DeleteEventProperties failed: %v", rErr)
	}
}

func TestUpdateCustomerOK(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "update_customer"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.UpdateCustomer("stubName", "stub@mail.com", "http://stub.url", map[string]string{})
	if rErr != nil {
		t.Errorf("UpdateCustomer failed: %v", rErr)
	}
}

func TestSetCustomerFieldsOK(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "set_customer_fields"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.SetCustomerFields(map[string]string{})
	if rErr != nil {
		t.Errorf("SetCustomerFields failed: %v", rErr)
	}
}

func TestGetGroupsStatusOK(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "get_groups_status"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	groupsStatus, rErr := api.GetGroupsStatus([]int{1,2,3})
	if rErr != nil {
		t.Errorf("GetGroupsStatus failed: %v", rErr)
	}

	expectedStatus := map[int]customer.GroupStatus{
		1: customer.GroupStatusOnline,
		2: customer.GroupStatusOffline,
		3: customer.GroupStatusOnlineForQueue,
	}

	if len(groupsStatus) != 3 {
		t.Errorf("Invalid size of groupsStatus map: %v, expected 3", len(groupsStatus))
	}

	for group, status := range groupsStatus {
		if status != expectedStatus[group] {
			t.Errorf("Incorrect status: %v, for group: %v", status, group)
		}
	}
}

func TestCheckGoalsOK(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "check_goals"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.CheckGoals("http://page.url", 0, map[string]string{})
	if rErr != nil {
		t.Errorf("CheckGoals failed: %v", rErr)
	}
}

func TestGetFormOK(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "get_form"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	form, enabled, rErr := api.GetForm(0, customer.FormTypePrechat)
	if rErr != nil {
		t.Errorf("GetForm failed: %v", rErr)
	}

	// TODO add better validation
	if !enabled {
		t.Errorf("Invalid enabled state: %v", enabled)
	}

	if form.ID != "156630109416307809" {
		t.Errorf("Invalid form id: %v", form.ID)
	}

	if len(form.Fields) != 3 {
		t.Errorf("Invalid length of form fields array: %v", len(form.Fields))
	}
}

func TestGetPredictedAgentOK(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "get_predicted_agent"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	agent, rErr := api.GetPredictedAgent()
	if rErr != nil {
		t.Errorf("GetPredictedAgent failed: %v", rErr)
	}

	// TODO add better validation

	if agent == nil {
		t.Errorf("Invalid Agent")
	}
}

func TestGetURLDetailsOK(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "get_url_details"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	details, rErr := api.GetURLDetails("http://totally.unsuspicious.url.com")
	if rErr != nil {
		t.Errorf("GetURLDetails failed: %v", rErr)
	}
	// TODO add better validation

	if details == nil {
		t.Errorf("Incorrect details")
	}
}

func TestMarkEventsAsSeenOK(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "mark_events_as_seen"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.MarkEventsAsSeen("stubChatID", time.Time{})
	if rErr != nil {
		t.Errorf("MarkEventsAsSeen failed: %v", rErr)
	}
}

func TestGetCustomerOK(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "get_customer"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	customer, rErr := api.GetCustomer()
	if rErr != nil {
		t.Errorf("GetCustomer failed: %v", rErr)
	}

	// TODO add better validation

	if customer == nil {
		t.Errorf("Invalid Customer")
	}
}
