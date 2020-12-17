package customer_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/livechat/lc-sdk-go/v3/authorization"
	"github.com/livechat/lc-sdk-go/v3/customer"
	"github.com/livechat/lc-sdk-go/v3/objects"
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
	"resume_chat": `{
		"thread_id": "PGDGHT5G"
	}`,
	"send_event": `{
		"event_id": "K600PKZON8"
	}`,
	"list_chats": `{
		"chats_summary": [{
			"id": "PJ0MRSHTDG",
			"last_thread_id": "K600PKZON8",
			"last_thread_summary": {
				"id": "K600PKZON8",
				"created_at": "2020-05-07T07:11:28.288340Z",
				"user_ids": ["b5657aff34dd32e198160d54666df9d8"],
				"properties": {},
				"tags": ["bug_report"]
			},
			"users": [],
			"properties": {},
			"access": {},
			"last_event_per_type": {
				"message": {
					"thread_id": "K600PKZON8",
					"thread_created_at": "2020-05-07T07:11:28.288340Z",
					"event": {
						"id": "K600PKZON8_1",
						"created_at": "2020-05-07T07:11:28.288340Z",
						"type": "message",
						"properties": {
							"lc2": {
								"welcome_message": true
							}
						},
						"text": "Hello. What can I do for you?",
						"author_id": "b5657aff34dd32e198160d54666df9d8"
					}
				},
				"system_message": {
					"thread_id": "K600PKZON8",
					"thread_created_at": "2020-05-07T07:11:28.288340Z",
					"event": {}
				}
			}
		}],
		"total_chats": 1,
		"previous_page_id": "MTUxNzM5ODEzMTQ5Ng=="
	}`,
	"get_chat": `{
		"id": "PJ0MRSHTDG",
		"thread": {
		  "id": "K600PKZON8",
		  "created_at": "2020-05-07T07:11:28.288340Z",
		  "active": true,
		  "user_ids": [
			"b7eff798-f8df-4364-8059-649c35c9ed0c",
			"b5657aff34dd32e198160d54666df9d8"
		  ],
		  "events": [{
			"id": "Q20N9CKRX2_1",
			"created_at": "2019-12-17T07:57:41.512000Z",
			"type": "message",
			"text": "Hello",
			"author_id": "b5657aff34dd32e198160d54666df9d8"
		  }],
		  "properties": {
			"0805e283233042b37f460ed8fbf22160": {
			  "string_property": "string_value"
			}
		  },
		  "access": {
			"group_ids": [0]
		  },
		  "previous_thread_id": "K600PKZOM8",
		  "next_thread_id": "K600PKZOO8"
		},
		"users": [{
		  "id": "b7eff798-f8df-4364-8059-649c35c9ed0c",
		  "type": "customer",
		  "present": true
		}, {
		  "id": "b5657aff34dd32e198160d54666df9d8",
		  "name": "Agent Smith",
		  "type": "agent",
		  "present": true,
		  "avatar": "https://example.com/avatar.jpg",
		  "job_title": "Support Agent"
		}],
		"access": {
		  "group_ids": [0]
		},
		"properties": {
		  "0805e283233042b37f460ed8fbf22160": {
			"string_property": "string_value"
		  }
		}
	  }`,
	"list_threads": `{
    "threads": [{
      "id": "K600PKZON8",
      "active": true,
      "user_ids": [
        "b7eff798-f8df-4364-8059-649c35c9ed0c",
        "smith@example.com"
      ],
      "events": [{
        "id": "Q20N9CKRX2_1",
        "created_at": "2019-12-17T07:57:41.512000Z",
        "recipients": "all",
        "type": "message",
        "text": "Hello",
        "author_id": "smith@example.com"
      }],
      "properties": {},
      "access": {
        "group_ids": [0]
      },
			"created_at": "2019-12-17T07:57:41.512000Z",
      "previous_thread_id": "K600PKZOM8",
      "next_thread_id": "K600PKZOO8"
    }],
    "found_threads": 1,
    "next_page_id": "MTUxNzM5ODEzMTQ5Ng==",
    "previous_page_id": "MTUxNzM5ODEzMTQ5Nw=="
	}`,
	"deactivate_chat": `{}`,
	"upload_file": `{
		"url": "https://cdn.livechat-static.com/api/file/lc/att/8948324/45a3581b59a7295145c3825c86ec7ab3/image.png"
	}`,
	"send_rich_message_postback":  `{}`,
	"send_sneak_peek":             `{}`,
	"update_chat_properties":      `{}`,
	"delete_chat_properties":      `{}`,
	"update_thread_properties":    `{}`,
	"delete_thread_properties":    `{}`,
	"update_event_properties":     `{}`,
	"delete_event_properties":     `{}`,
	"update_customer":             `{}`,
	"set_customer_session_fields": `{}`,
	"list_group_statuses": `{
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
				},
				{
					"id": "157986144052009331",
					"type": "group_chooser",
					"label": "Choose a department:",
					"required": true,
					"options": [
						{
							"id": "0",
							"group_id": 1,
							"label": "Marketing"
						},
						{
							"id": "1",
							"group_id": 2,
							"label": "Sales"
						},
						{
							"id": "2",
							"group_id": 0,
							"label": "General"
						}
					]
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
	"get_url_info": `{
		"title": "LiveChat | Live Chat Software and Help Desk Software",
		"description": "LiveChat - premium live chat software and help desk software for business. Over 24 000 companies from 150 countries use LiveChat. Try now, chat for free!",
		"image_url": "s3.eu-central-1.amazonaws.com/labs-fraa-livechat-thumbnails/96979c3552cf3fa4ae326086a3048d9354c27324.png",
		"image_width": 200,
		"image_height": 200,
		"url": "https://livechatinc.com"
	}`,
	"mark_events_as_seen": `{}`,
	"list_license_properties": `{
		"0805e283233042b37f460ed8fbf22160": {
				"string_property": "string value"
		}
	}`,
	"list_group_properties": `{
		"0805e283233042b37f460ed8fbf22160": {
				"string_property": "string value"
		}
	}`,
	"get_customer":    `{}`, //TODO - create some real structure here
	"accept_greeting": `{}`,
	"cancel_greeting": `{}`,
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

		if req.URL.String() != "https://api.livechatinc.com/v3.3/customer/action/"+method+"?license_id=12345" {
			t.Errorf("Invalid URL for Customer API request: %s", req.URL.String())
			return createServerError("Invalid URL")
		}

		expectedMethod := "POST"
		if method == "list_license_properties" || method == "list_group_properties" {
			expectedMethod = "GET"
		}
		if expectedMethod != req.Method {
			t.Errorf("Invalid method: %s for Customer API action: %s", req.Method, method)
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
	_, err := customer.NewAPI(nil, nil, "client_id")
	if err == nil {
		t.Errorf("API should not be created without token getter")
	}
}

func TestStartChatShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "start_chat"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	m := &objects.Message{
		Postback: &objects.Postback{
			ID:    "123",
			Value: "abc",
		},
	}
	ic := &objects.InitialChat{
		Thread: &objects.InitialThread{
			Events: []interface{}{m},
		},
	}
	chatID, threadID, _, rErr := api.StartChat(ic, true, true)
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

func TestSendEventShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "send_event"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	eventID, rErr := api.SendEvent("stubChatID", &objects.Event{}, false)
	if rErr != nil {
		t.Errorf("SendEvent failed: %v", rErr)
	}

	if eventID != "K600PKZON8" {
		t.Errorf("Invalid eventID: %v", eventID)
	}
}

func TestSendMessageShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "send_event"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	eventID, rErr := api.SendMessage("stubChatID", "Hello World", customer.All)
	if rErr != nil {
		t.Errorf("SendMessage failed: %v", rErr)
	}

	if eventID != "K600PKZON8" {
		t.Errorf("Invalid eventID: %v", eventID)
	}
}

func TestSendSystemMessageShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "send_event"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	textVars := map[string]string{
		"var1": "val1",
		"var2": "val2",
	}
	eventID, rErr := api.SendSystemMessage("stubChatID", "text", "messagetype", textVars, customer.All, false)
	if rErr != nil {
		t.Errorf("SendSystemMessage failed: %v", rErr)
	}

	if eventID != "K600PKZON8" {
		t.Errorf("Invalid eventID: %v", eventID)
	}
}

func TestResumeChatShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "resume_chat"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	threadID, _, rErr := api.ResumeChat(&objects.InitialChat{}, true, true)
	if rErr != nil {
		t.Errorf("ResumeChat failed: %v", rErr)
	}

	if threadID != "PGDGHT5G" {
		t.Errorf("Invalid threadID: %v", threadID)
	}
}

func TestListChatsShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "list_chats"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	chats, total, prevPage, nextPage, rErr := api.ListChats("", "", 20)
	if rErr != nil {
		t.Errorf("ListChats failed: %v", rErr)
	}

	// TODO add better validation

	if len(chats) != 1 {
		t.Errorf("Invalid chats")
	}
	if chats[0].ID != "PJ0MRSHTDG" {
		t.Errorf("Invalid chat id: %v", chats[0].ID)
	}
	if chats[0].LastThreadSummary.ID != "K600PKZON8" {
		t.Errorf("Invalid last thread summary")
	}
	if chats[0].LastEventPerType["message"].ThreadID != "K600PKZON8" {
		t.Errorf("Invalid last event per type")
	}
	e := chats[0].LastEventPerType["message"].Event
	if e.Message().Text != "Hello. What can I do for you?" {
		t.Errorf("Invalid last message event")
	}
	if total != 1 {
		t.Errorf("Invalid total chats: %v", total)
	}
	if prevPage != "MTUxNzM5ODEzMTQ5Ng==" {
		t.Errorf("Invalid previous page ID: %v", prevPage)
	}
	if nextPage != "" {
		t.Errorf("Invalid next page ID: %v", nextPage)
	}
}

func TestGetChatShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "get_chat"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	chat, rErr := api.GetChat("stubChatID", "stubThreadID")
	if rErr != nil {
		t.Errorf("GetChat failed: %v", rErr)
	}

	if chat.ID != "PJ0MRSHTDG" {
		t.Errorf("Received chat.ID invalid: %v", chat.ID)
	}
}

func TestListThreadsShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "list_threads"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	threads, found, prevPage, nextPage, rErr := api.ListThreads("stubChatID", "", "", 20, 0)
	if rErr != nil {
		t.Errorf("ListThreads failed: %v", rErr)
	}

	if len(threads) != 1 {
		t.Errorf("Received invalid threads length: %v", len(threads))
	}

	if found != 1 {
		t.Errorf("Received invalid total threads: %v", found)
	}

	if prevPage != "MTUxNzM5ODEzMTQ5Nw==" {
		t.Errorf("Invalid previous page ID: %v", prevPage)
	}

	if nextPage != "MTUxNzM5ODEzMTQ5Ng==" {
		t.Errorf("Invalid next page ID: %v", nextPage)
	}
}

func TestDeactivateChatShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "deactivate_chat"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.DeactivateChat("stubChatID")
	if rErr != nil {
		t.Errorf("DeactivateChat failed: %v", rErr)
	}
}

func TestUploadFileShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
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

func TestSendRichMessagePostbackShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
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

func TestSendSneakPeekShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
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

func TestUpdateChatPropertiesShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "update_chat_properties"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.UpdateChatProperties("stubChatID", objects.Properties{})
	if rErr != nil {
		t.Errorf("UpdateChatProperties failed: %v", rErr)
	}
}

func TestDeleteChatPropertiesShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
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

func TestUpdateThreadPropertiesShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "update_thread_properties"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.UpdateThreadProperties("stubChatID", "stubThreadID", objects.Properties{})
	if rErr != nil {
		t.Errorf("UpdateThreadProperties failed: %v", rErr)
	}
}

func TestDeleteThreadPropertiesShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "delete_thread_properties"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.DeleteThreadProperties("stubChatID", "stubThreadID", map[string][]string{})
	if rErr != nil {
		t.Errorf("DeleteThreadProperties failed: %v", rErr)
	}
}

func TestUpdateEventPropertiesShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "update_event_properties"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.UpdateEventProperties("stubChatID", "stubThreadID", "stubEventID", objects.Properties{})
	if rErr != nil {
		t.Errorf("UpdateEventProperties failed: %v", rErr)
	}
}

func TestDeleteEventPropertiesShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
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

func TestUpdateCustomerShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "update_customer"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.UpdateCustomer("stubName", "stub@mail.com", "http://stub.url", []map[string]string{})
	if rErr != nil {
		t.Errorf("UpdateCustomer failed: %v", rErr)
	}
}

func TestSetCustomerSessionFieldsShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "set_customer_session_fields"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.SetCustomerSessionFields([]map[string]string{})
	if rErr != nil {
		t.Errorf("SetCustomerSessionFields failed: %v", rErr)
	}
}

func TestListGroupStatusesShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "list_group_statuses"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	groupStatuses, rErr := api.ListGroupStatuses([]int{1, 2, 3})
	if rErr != nil {
		t.Errorf("ListGroupStatuses failed: %v", rErr)
	}

	expectedStatus := map[int]customer.GroupStatus{
		1: customer.GroupStatusOnline,
		2: customer.GroupStatusOffline,
		3: customer.GroupStatusOnlineForQueue,
	}

	if len(groupStatuses) != 3 {
		t.Errorf("Invalid size of groupStatuses map: %v, expected 3", len(groupStatuses))
	}

	for group, status := range groupStatuses {
		if status != expectedStatus[group] {
			t.Errorf("Incorrect status: %v, for group: %v", status, group)
		}
	}
}

func TestCheckGoalsShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
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

func TestGetFormShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
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

	if len(form.Fields) != 4 {
		t.Errorf("Invalid length of form fields array: %v", len(form.Fields))
	}

	if len(form.Fields[3].Options) != 3 {
		t.Errorf("Invalid length of form group_chooser field options array: %v", len(form.Fields[3].Options))
	}
}

func TestGetPredictedAgentShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
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

func TestGetURLInfoShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "get_url_info"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	info, rErr := api.GetURLInfo("http://totally.unsuspicious.url.com")
	if rErr != nil {
		t.Errorf("GetURLInfo failed: %v", rErr)
	}
	// TODO add better validation

	if info == nil {
		t.Errorf("Incorrect info")
	}
}

func TestMarkEventsAsSeenShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
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

func TestGetCustomerShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
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

// TESTS Error Cases

func TestStartChatShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "start_chat"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	_, _, _, rErr := api.StartChat(&objects.InitialChat{}, true, true)
	verifyErrorResponse("StartChat", rErr, t)
}

func TestSendEventShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "send_event"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	_, rErr := api.SendEvent("stubChatID", &objects.Event{}, false)
	verifyErrorResponse("SendEvent", rErr, t)
}

func TestResumeChatShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "resume_chat"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	_, _, rErr := api.ResumeChat(&objects.InitialChat{}, true, true)
	verifyErrorResponse("ResumeChat", rErr, t)
}

func TestListChatsShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "list_chats"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	_, _, _, _, rErr := api.ListChats("", "", 20)
	verifyErrorResponse("ListChats", rErr, t)
}

func TestGetChatShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "get_chat"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	_, rErr := api.GetChat("stubChatID", "stubThreadID")
	verifyErrorResponse("GetChat", rErr, t)
}

func TestListThreadsShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "list_threads"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	_, _, _, _, rErr := api.ListThreads("stubChatID", "", "", 20, 0)
	verifyErrorResponse("ListThreads", rErr, t)
}

func TestDeactivateChatShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "deactivate_chat"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.DeactivateChat("stubChatID")
	verifyErrorResponse("DeactivateChat", rErr, t)
}

func TestUploadFileShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "upload_file"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	_, rErr := api.UploadFile("filename", []byte{})
	verifyErrorResponse("UploadFile", rErr, t)
}

func TestSendRichMessagePostbackShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "send_rich_message_postback"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.SendRichMessagePostback("stubChatID", "stubThreadID", "stubEventID", "stubPostbackID", false)
	verifyErrorResponse("SendRichMessagePostback", rErr, t)
}

func TestSendSneakPeekShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "send_sneak_peek"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.SendSneakPeek("stubChatID", "sneaky freaky baby")
	verifyErrorResponse("SendSneakPeek", rErr, t)
}

func TestUpdateChatPropertiesShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "update_chat_properties"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.UpdateChatProperties("stubChatID", objects.Properties{})
	verifyErrorResponse("UpdateChatProperties", rErr, t)
}

func TestDeleteChatPropertiesShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "delete_chat_properties"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.DeleteChatProperties("stubChatID", map[string][]string{})
	verifyErrorResponse("DeleteChatProperties", rErr, t)
}

func TestUpdateThreadPropertiesShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "update_thread_properties"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.UpdateThreadProperties("stubChatID", "stubThreadID", objects.Properties{})
	verifyErrorResponse("UpdateThreadProperties", rErr, t)
}

func TestDeleteThreadPropertiesShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "delete_thread_properties"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.DeleteThreadProperties("stubChatID", "stubThreadID", map[string][]string{})
	verifyErrorResponse("DeleteThreadProperties", rErr, t)
}

func TestUpdateEventPropertiesShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "update_event_properties"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.UpdateEventProperties("stubChatID", "stubThreadID", "stubEventID", objects.Properties{})
	verifyErrorResponse("UpdateEventProperties", rErr, t)
}

func TestDeleteEventPropertiesShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "delete_event_properties"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.DeleteEventProperties("stubChatID", "stubThreadID", "stubEventID", map[string][]string{})
	verifyErrorResponse("DeleteEventProperties", rErr, t)
}

func TestUpdateCustomerShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "update_customer"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.UpdateCustomer("stubName", "stub@mail.com", "http://stub.url", []map[string]string{})
	verifyErrorResponse("UpdateCustomer", rErr, t)
}

func TestSetCustomerSessionFieldsShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "set_customer_session_fields"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.SetCustomerSessionFields([]map[string]string{})
	verifyErrorResponse("SetCustomerSessionFields", rErr, t)
}

func TestListGroupStatusesShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "list_group_statuses"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	_, rErr := api.ListGroupStatuses([]int{1, 2, 3})
	verifyErrorResponse("ListGroupStatuses", rErr, t)
}

func TestCheckGoalsShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "check_goals"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.CheckGoals("http://page.url", 0, map[string]string{})
	verifyErrorResponse("CheckGoals", rErr, t)
}

func TestGetFormShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "get_form"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	_, _, rErr := api.GetForm(0, customer.FormTypePrechat)
	verifyErrorResponse("GetForm", rErr, t)
}

func TestGetPredictedAgentShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "get_predicted_agent"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	_, rErr := api.GetPredictedAgent()
	verifyErrorResponse("GetPredictedAgent", rErr, t)
}

func TestGetURLInfoShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "get_url_info"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	_, rErr := api.GetURLInfo("http://totally.unsuspicious.url.com")
	verifyErrorResponse("GetURLInfo", rErr, t)
}

func TestMarkEventsAsSeenShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "mark_events_as_seen"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.MarkEventsAsSeen("stubChatID", time.Time{})
	verifyErrorResponse("MarkEventsAsSeen", rErr, t)
}

func TestGetCustomerShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "get_customer"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	_, rErr := api.GetCustomer()
	verifyErrorResponse("GetCustomer", rErr, t)
}

func TestListLicensePropertiesShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "list_license_properties"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	resp, rErr := api.ListLicenseProperties("", "")
	if rErr != nil {
		t.Errorf("ListLicenseProperties failed: %v", rErr)
	}

	if len(resp) != 1 {
		t.Errorf("Invalid license properties: %v", resp)
	}

	if resp["0805e283233042b37f460ed8fbf22160"]["string_property"] != "string value" {
		t.Errorf("Invalid license property 0805e283233042b37f460ed8fbf22160.string_property: %v", resp["0805e283233042b37f460ed8fbf22160"]["string_property"])
	}
}

func TestListGroupPropertiesShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "list_group_properties"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	resp, rErr := api.ListGroupProperties(0, "", "")
	if rErr != nil {
		t.Errorf("ListGroupProperties failed: %v", rErr)
	}

	if len(resp) != 1 {
		t.Errorf("Invalid group properties: %v", resp)
	}

	if resp["0805e283233042b37f460ed8fbf22160"]["string_property"] != "string value" {
		t.Errorf("Invalid group property 0805e283233042b37f460ed8fbf22160.string_property: %v", resp["0805e283233042b37f460ed8fbf22160"]["string_property"])
	}
}

func TestAcceptGreetingShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "accept_greeting"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.AcceptGreeting(1337, "foo")
	if rErr != nil {
		t.Errorf("AcceptGreeting failed: %v", rErr)
	}
}

func TestAcceptGreetingShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "accept_greeting"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.AcceptGreeting(1337, "foo")
	verifyErrorResponse("AcceptGreeting", rErr, t)
}

func TestCancelGreetingShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "cancel_greeting"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.CancelGreeting("foo")
	if rErr != nil {
		t.Errorf("CancelGreeting failed: %v", rErr)
	}
}
func TestCancelGreetingShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "cancel_greeting"))

	api, err := customer.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.CancelGreeting("foo")
	verifyErrorResponse("CancelGreeting", rErr, t)
}
