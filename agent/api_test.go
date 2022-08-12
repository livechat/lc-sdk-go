package agent_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/livechat/lc-sdk-go/v5/agent"
	"github.com/livechat/lc-sdk-go/v5/authorization"
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

func stubBearerTokenGetter() *authorization.Token {
	return stubTokenGetter(authorization.BearerToken)()
}

func stubTokenGetter(tokenType authorization.TokenType) func() *authorization.Token {
	return func() *authorization.Token {
		return &authorization.Token{
			AccessToken: "access_token",
			Region:      "region",
			Type:        tokenType,
		}
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
			"last_thread_summary": {
				"id": "K600PKZON8",
				"created_at": "2020-05-07T07:11:28.288340Z",
				"user_ids": [
					"smith@example.com",
					"b7eff798-f8df-4364-8059-649c35c9ed0c"
				],
				"properties": {
					"property_namespace": {
						"property_name": "property_value"
					}
				},
				"active": false,
				"access": {
					"group_ids": [
						0
					]
				},
				"tags": ["bug_report"],
				"queue": {
					"position": 42,
					"wait_time": 1337,
					"queued_at": "2020-05-12T11:42:47.383000Z"
				}
			},
			"users": [
				{
					"id": "b7eff798-f8df-4364-8059-649c35c9ed0c",
					"name": "Thomas Anderson",
					"events_seen_up_to": "2020-05-12T12:31:46.463000Z",
					"type": "customer",
					"present": true,
					"created_at": "2019-11-02T19:19:50.625101Z",
					"last_visit": {
						"started_at": "2020-05-12T11:32:03.497479Z",
						"ended_at": "2020-05-12T11:33:33.497000Z",
						"ip": "<customer_ip>",
						"user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36",
						"geolocation": {
							"country": "Poland",
							"country_code": "PL",
							"region": "Dolnoslaskie",
							"city": "Wroclaw",
							"timezone": "Europe/Warsaw",
							"latitude": "51.1043015",
							"longitude": "17.0335007"
						},
						"last_pages": [
							{
								"opened_at": "2020-05-12T11:32:03.497479Z",
								"url": "https://cdn.livechatinc.com/preview/11442778",
								"title": "Sample Page | Preview your chat window"
							}
						]
					},
					"statistics": {
						"chats_count": 1,
						"threads_count": 3,
						"visits_count": 6,
						"page_views_count": 2,
						"greetings_shown_count": 2,
						"greetings_accepted_count": 1
					},
					"agent_last_event_created_at": "2020-05-12T11:42:47.393002Z",
					"customer_last_event_created_at": "2020-05-12T12:31:46.463000Z"
				},
				{
					"id": "smith@example.com",
					"name": "Agent Smith",
					"email": "smith@example.com",
					"events_seen_up_to": "2020-05-12T12:31:46.999999Z",
					"type": "agent",
					"present": false,
					"avatar": "https://cdn.livechatinc.com/cloud/?uri=https://livechat.s3.amazonaws.com/default/avatars/a7.png"
				}
			],
			"properties": {
				"property_namespace": {
					"property_name": "property_value"
				}
			},
			"access": {
				"group_ids": [
					0
				]
			},
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
			},
			"is_followed": false
		}],
		"found_chats": 1,
		"previous_page_id": "MTUxNzM5ODEzMTQ5Ng=="
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
        "visibility": "all",
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
	"get_chat": `{
		"id": "PJ0MRSHTDG",
		"thread": {
		  "id": "K600PKZON8",
		  "created_at": "2020-05-07T07:11:28.288340Z",
		  "active": true,
		  "user_ids": [
			"b7eff798-f8df-4364-8059-649c35c9ed0c",
			"smith@example.com"
		  ],
		  "events": [{
			"id": "Q20N9CKRX2_1",
			"created_at": "2019-12-17T07:57:41.512000Z",
			"visibility": "all",
			"type": "message",
			"text": "Hello",
			"author_id": "smith@example.com"
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
		  "present": true,
		  "created_at": "2019-12-17T08:53:20.693553+01:00",
		  "statistics": {
			"chats_count": 1
		  },
		  "agent_last_event_created_at": "2019-12-17T09:04:05.239000+01:00"
		}, {
		  "id": "smith@example.com",
		  "name": "Agent Smith",
		  "email": "smith@example.com",
		  "type": "agent",
		  "present": true,
		  "avatar": "https://example.com/avatar.jpg"
		}],
		"properties": {
		  "0805e283233042b37f460ed8fbf22160": {
			"string_property": "string_value"
		  }
		},
		"access": {
		  "group_ids": [0]
		},
		"is_followed": true
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
		"found_chats": 1,
		"next_page_id": "nextpagehash"
	}`,
	"deactivate_chat":       `{}`,
	"follow_chat":           `{}`,
	"unfollow_chat":         `{}`,
	"grant_chat_access":     `{}`,
	"revoke_chat_access":    `{}`,
	"set_chat_access":       `{}`,
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
		"limited_customers": 1,
		"previous_page_id": "prevpagehash"
	}`,
	"create_customer": `{
		"customer_id": "mister_customer"
	}`,
	"update_customer":       `{}`,
	"ban_customer":          `{}`,
	"mark_events_as_seen":   `{}`,
	"set_routing_status":    `{}`,
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
	"follow_customer":   `{}`,
	"unfollow_customer": `{}`,
	"list_routing_statuses": `[{
		"agent_id": "smith@example.com",
		"status": "accepting_chats"
	}, {
		"agent_id": "agent@example.com",
		"status": "not_accepting_chats"
	}]`,
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
				Body:       io.NopCloser(bytes.NewBufferString(responseError)),
				Header:     make(http.Header),
			}
		}

		if req.URL.String() != "https://api.livechatinc.com/v3.5/agent/action/"+method {
			t.Errorf("Invalid URL for Agent API request: %s", req.URL.String())
			return createServerError("Invalid URL")
		}

		if req.Method != "POST" {
			t.Errorf("Invalid method: %s for Agent API action: %s", req.Method, method)
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
			Body:       io.NopCloser(bytes.NewBufferString(mockedResponses[method])),
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
			Body:       io.NopCloser(bytes.NewBufferString(responseError)),
			Header:     make(http.Header),
		}
	}
}

func createMockedMultipleAuthErrorsResponder(t *testing.T, fails int) func(req *http.Request) *http.Response {
	var n int

	responseError := `{
		"error": {
			"type": "authentication",
			"message": "Invalid access token"
		}
	}`

	return func(req *http.Request) *http.Response {
		n++
		if n > fails {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
				Header:     make(http.Header),
			}
		}
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       io.NopCloser(bytes.NewBufferString(responseError)),
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
		t.Error("API should not be created without token getter")
	}
}

func TestAuthorIDHeader(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		if xAuthorID := req.Header.Get("X-Author-Id"); xAuthorID != "my_bot" {
			t.Errorf("Invalid X-Author-Id header: %s", xAuthorID)
		}
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(nil)),
			Header:     make(http.Header),
		}
	})
	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}
	api.SetAuthorID("my_bot")
	api.Call("", nil, nil)
}

func TestStartChatShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "start_chat"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	chatID, threadID, _, rErr := api.StartChat(&agent.InitialChat{}, true, true)
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

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	eventID, rErr := api.SendEvent("stubChatID", agent.Event{}, false)
	if rErr != nil {
		t.Errorf("SendEvent failed: %v", rErr)
	}

	if eventID != "K600PKZON8" {
		t.Errorf("Invalid eventID: %v", eventID)
	}
}

func TestResumeChatShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "resume_chat"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	threadID, _, rErr := api.ResumeChat(&agent.InitialChat{}, true, true)
	if rErr != nil {
		t.Errorf("ResumeChat failed: %v", rErr)
	}

	if threadID != "PGDGHT5G" {
		t.Errorf("Invalid threadID: %v", threadID)
	}
}

func TestListChatsShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "list_chats"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	chats, found, prevPage, nextPage, rErr := api.ListChats(agent.NewChatsFilters(), "", "", 20)
	if rErr != nil {
		t.Errorf("ListChats failed: %v", rErr)
	}

	// TODO add better validation

	if len(chats) != 1 {
		t.Errorf("Invalid chats count. Got: %v, expected: %v", len(chats), 1)
	}
	if chats[0].ID != "PJ0MRSHTDG" {
		t.Errorf("Invalid chat id. Got: %v, expected: %v", chats[0].ID, "PJ0MRSHTDG")
	}
	if len(chats[0].Users) != 2 {
		t.Errorf("Invalid chat users count. Got: %v, expected: %v", len(chats[0].Users), 2)
	}
	if chats[0].IsFollowed {
		t.Error("Chat is followed should be false")
	}
	if chats[0].LastThreadSummary.ID != "K600PKZON8" {
		t.Errorf("Invalid last thread id. Got: %v, expected: %v", chats[0].LastThreadSummary.ID, "K600PKZON8")
	}
	lastThreadCreationDate := chats[0].LastThreadSummary.CreatedAt.Format(time.RFC3339Nano)
	if lastThreadCreationDate != "2020-05-07T07:11:28.28834Z" {
		t.Errorf("Invalid last thread creation date. Got: %v, expected: %v", lastThreadCreationDate, "2020-05-07T07:11:28.28834Z")
	}
	if len(chats[0].LastThreadSummary.Properties) != 1 {
		t.Errorf("Invalid last thread properties count. Got: %v, expected: %v", len(chats[0].LastThreadSummary.Properties), 1)
	}
	if chats[0].LastThreadSummary.Queue.Position != 42 {
		t.Errorf("Invalid last thread queue position. Got: %v, expected: %v", chats[0].LastThreadSummary.Queue.Position, 42)
	}
	if chats[0].LastEventPerType["message"].ThreadID != "K600PKZON8" {
		t.Errorf("Invalid last event per type thread id. Got: %v, expected: %v", chats[0].LastEventPerType["message"].ThreadID, "K600PKZON8")
	}
	e := chats[0].LastEventPerType["message"].Event
	if e.Message().Text != "Hello. What can I do for you?" {
		t.Errorf("Invalid last message event text. Got: %v, expected: %v", e.Message().Text, "Hello. What can I do for you?")
	}
	if found != 1 {
		t.Errorf("Invalid found. Got: %v, expected: %v", found, 1)
	}
	if prevPage != "MTUxNzM5ODEzMTQ5Ng==" {
		t.Errorf("Invalid previous page id. Got: %v, expected: %v", prevPage, "MTUxNzM5ODEzMTQ5Ng==")
	}
	if nextPage != "" {
		t.Errorf("Invalid next page id. Got: %v, expected: %v", nextPage, "")
	}
}

func TestGetChatShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "get_chat"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
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

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	threads, found, prevPage, nextPage, rErr := api.ListThreads("stubChatID", "", "", 20, 0, agent.NewThreadsFilters())
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

func TestListArchivesShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "list_archives"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	chats, found, prevPageID, nextPageID, rErr := api.ListArchives(agent.NewArchivesFilters(), "MTUxNzM5ODEzMTQ5Ng==", 20)
	if rErr != nil {
		t.Errorf("ListArchives failed: %v", rErr)
	}

	if chats[0].ID != "PJ0MRSHTDG" {
		t.Errorf("Received chat.ID invalid: %v", chats[0].ID)
	}
	if found != 1 {
		t.Errorf("Received found invalid: %v", found)
	}
	if prevPageID != "" {
		t.Errorf("Received prevPageID invalid: %v", prevPageID)
	}
	if nextPageID == "" {
		t.Errorf("Received nextPageID invalid: %v", nextPageID)
	}
}

func TestDeactivateChatShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "deactivate_chat"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.DeactivateChat("stubChatID", false)
	if rErr != nil {
		t.Errorf("DeactivateChat failed: %v", rErr)
	}
}

func TestFollowChatShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "follow_chat"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.FollowChat("stubChatID")
	if rErr != nil {
		t.Errorf("FollowChat failed: %v", rErr)
	}
}

func TestUnfollowChatShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "unfollow_chat"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.UnfollowChat("stubChatID")
	if rErr != nil {
		t.Errorf("UnfollowChat failed: %v", rErr)
	}
}

func TestUploadFileShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "upload_file"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	fileUrl, rErr := api.UploadFile("filename", []byte{})
	if rErr != nil {
		t.Errorf("UploadFile failed: %v", rErr)
	}

	if fileUrl != "https://cdn.livechat-static.com/api/file/lc/att/8948324/45a3581b59a7295145c3825c86ec7ab3/image.png" {
		t.Errorf("Invalid file URL: %v", fileUrl)
	}
}

func TestAddUserToChatShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "add_user_to_chat"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.AddUserToChat("chat", "user", "agent", "all", false)
	if rErr != nil {
		t.Errorf("AddUserToChat failed: %v", rErr)
	}
}

func TestRemoveUserFromChatShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "remove_user_from_chat"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.RemoveUserFromChat("chat", "user", "agent", false)
	if rErr != nil {
		t.Errorf("RemoveUserFromChat failed: %v", rErr)
	}
}

func TestSendRichMessagePostbackShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "send_rich_message_postback"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.SendRichMessagePostback("stubChatID", "stubThreadID", "stubEventID", "stubPostbackID", false)
	if rErr != nil {
		t.Errorf("SendRichMessagePostback failed: %v", rErr)
	}
}

func TestUpdateChatPropertiesShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "update_chat_properties"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.UpdateChatProperties("stubChatID", agent.Properties{})
	if rErr != nil {
		t.Errorf("UpdateChatProperties failed: %v", rErr)
	}
}

func TestDeleteChatPropertiesShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "delete_chat_properties"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.DeleteChatProperties("stubChatID", map[string][]string{})
	if rErr != nil {
		t.Errorf("DeleteChatProperties failed: %v", rErr)
	}
}

func TestUpdateThreadPropertiesShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "update_thread_properties"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.UpdateThreadProperties("stubChatID", "stubThreadID", agent.Properties{})
	if rErr != nil {
		t.Errorf("UpdateThreadProperties failed: %v", rErr)
	}
}

func TestDeleteThreadPropertiesShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "delete_thread_properties"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.DeleteThreadProperties("stubChatID", "stubThreadID", map[string][]string{})
	if rErr != nil {
		t.Errorf("DeleteThreadProperties failed: %v", rErr)
	}
}

func TestUpdateEventPropertiesShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "update_event_properties"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.UpdateEventProperties("stubChatID", "stubThreadID", "stubEventID", agent.Properties{})
	if rErr != nil {
		t.Errorf("UpdateEventProperties failed: %v", rErr)
	}
}

func TestDeleteEventPropertiesShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "delete_event_properties"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.DeleteEventProperties("stubChatID", "stubThreadID", "stubEventID", map[string][]string{})
	if rErr != nil {
		t.Errorf("DeleteEventProperties failed: %v", rErr)
	}
}

func TestTagThreadShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "tag_thread"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.TagThread("stubChatID", "stubThreadID", "tag")
	if rErr != nil {
		t.Errorf("TagThread failed: %v", rErr)
	}
}

func TestUntagThreadShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "untag_thread"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.UntagThread("stubChatID", "stubThreadID", "tag")
	if rErr != nil {
		t.Errorf("UntagThread failed: %v", rErr)
	}
}

func TestGetCustomerShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "get_customer"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
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

func TestListCustomersShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "list_customers"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	customers, total, limited, prevPage, nextPage, rErr := api.ListCustomers(100, "page", "asc", "", agent.NewCustomersFilters())
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
	if limited != 1 {
		t.Errorf("Invalid limited customers amount: %v", limited)
	}
}

func TestCreateCustomerShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "create_customer"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
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

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.UpdateCustomer("mister_customer", "stubName", "stub@mail.com", "http://stub.url", []map[string]string{})
	if rErr != nil {
		t.Errorf("UpdateCustomer failed: %v", rErr)
	}
}

func TestBanCustomerShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "ban_customer"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.BanCustomer("mister_customer", 20)
	if rErr != nil {
		t.Errorf("BanCustomer failed: %v", rErr)
	}
}

func TestSetRoutingStatusShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "set_routing_status"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.SetRoutingStatus("some_agent", "accepting chats")
	if rErr != nil {
		t.Errorf("SetRoutingStatus failed: %v", rErr)
	}
}

func TestMarkEventsAsSeenShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "mark_events_as_seen"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.MarkEventsAsSeen("stubChatID", time.Time{})
	if rErr != nil {
		t.Errorf("MarkEventsAsSeen failed: %v", rErr)
	}
}

func TestSendTypingIndicatorShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "send_typing_indicator"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.SendTypingIndicator("stubChatID", "all", true)
	if rErr != nil {
		t.Errorf("SendTypingIndicator failed: %v", rErr)
	}
}

func TestMulticastShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "multicast"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.Multicast(agent.MulticastRecipients{}, []byte("{}"), "type")
	if rErr != nil {
		t.Errorf("Multicast failed: %v", rErr)
	}
}

func TestTransferChatShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "transfer_chat"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	ids := make([]interface{}, 1)
	ids[0] = "1"
	rErr := api.TransferChat("stubChatID", "agents", ids, agent.TransferChatOptions{})
	if rErr != nil {
		t.Errorf("TransferChat failed: %v", rErr)
	}
}

// TESTS Error Cases

func TestStartChatShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "start_chat"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	_, _, _, rErr := api.StartChat(&agent.InitialChat{}, true, true)
	verifyErrorResponse("StartChat", rErr, t)
}

func TestSendEventShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "send_event"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	_, rErr := api.SendEvent("stubChatID", &agent.Event{}, false)
	verifyErrorResponse("SendEvent", rErr, t)
}

func TestResumeChatShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "resume_chat"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	_, _, rErr := api.ResumeChat(&agent.InitialChat{}, true, true)
	verifyErrorResponse("ResumeChat", rErr, t)
}

func TestListChatsShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "list_chats"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	_, _, _, _, rErr := api.ListChats(agent.NewChatsFilters(), "", "", 20)
	verifyErrorResponse("ListChats", rErr, t)
}

func TestGetChatShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "get_chat"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	_, rErr := api.GetChat("stubChatID", "stubThreadID")
	verifyErrorResponse("GetChat", rErr, t)
}

func TestListThreadsShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "list_threads"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	_, _, _, _, rErr := api.ListThreads("stubChatID", "", "", 20, 0, agent.NewThreadsFilters())
	verifyErrorResponse("ListThreads", rErr, t)
}

func TestListArchivesShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "list_archives"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	_, _, _, _, rErr := api.ListArchives(agent.NewArchivesFilters(), "", 20)
	verifyErrorResponse("ListArchives", rErr, t)
}

func TestDeactivateChatShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "deactivate_chat"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.DeactivateChat("stubChatID", false)
	verifyErrorResponse("DeactivateChat", rErr, t)
}

func TestFollowChatShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "follow_chat"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.FollowChat("stubChatID")
	verifyErrorResponse("FollowChat", rErr, t)
}

func TestUnfollowChatShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "unfollow_chat"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.UnfollowChat("stubChatID")
	verifyErrorResponse("UnfollowChat", rErr, t)
}

func TestUploadFileShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "upload_file"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	_, rErr := api.UploadFile("filename", []byte{})
	verifyErrorResponse("UploadFile", rErr, t)

}

func TestAddUserToChatShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "add_user_to_chat"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.AddUserToChat("chat", "user", "agent", "all", false)
	verifyErrorResponse("AddUserToChat", rErr, t)
}

func TestRemoveUserFromChatShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "remove_user_from_chat"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.RemoveUserFromChat("chat", "user", "agent", false)
	verifyErrorResponse("RemoveUserFromChat", rErr, t)

}

func TestSendRichMessagePostbackShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "send_rich_message_postback"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.SendRichMessagePostback("stubChatID", "stubThreadID", "stubEventID", "stubPostbackID", false)
	verifyErrorResponse("SendRichMessagePostback", rErr, t)
}

func TestUpdateChatPropertiesShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "update_chat_properties"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.UpdateChatProperties("stubChatID", agent.Properties{})
	verifyErrorResponse("UpdateChatProperties", rErr, t)
}

func TestDeleteChatPropertiesShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "delete_chat_properties"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.DeleteChatProperties("stubChatID", map[string][]string{})
	verifyErrorResponse("DeleteChatProperties", rErr, t)
}

func TestUpdateThreadPropertiesShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "update_thread_properties"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.UpdateThreadProperties("stubChatID", "stubThreadID", agent.Properties{})
	verifyErrorResponse("UpdateThreadProperties", rErr, t)
}

func TestDeleteThreadPropertiesShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "delete_thread_properties"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.DeleteThreadProperties("stubChatID", "stubThreadID", map[string][]string{})
	verifyErrorResponse("DeleteThreadProperties", rErr, t)
}

func TestUpdateEventPropertiesShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "update_event_properties"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.UpdateEventProperties("stubChatID", "stubThreadID", "stubEventID", agent.Properties{})
	verifyErrorResponse("UpdateEventProperties", rErr, t)
}

func TestDeleteEventPropertiesShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "delete_event_properties"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.DeleteEventProperties("stubChatID", "stubThreadID", "stubEventID", map[string][]string{})
	verifyErrorResponse("DeleteEventProperties", rErr, t)
}

func TestTagThreadShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "tag_thread"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.TagThread("stubChatID", "stubThreadID", "tag")
	verifyErrorResponse("TagThread", rErr, t)
}

func TestUntagThreadShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "untag_thread"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.UntagThread("stubChatID", "stubThreadID", "tag")
	verifyErrorResponse("UntagThread", rErr, t)
}

func TestListCustomersShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "list_customers"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	_, _, _, _, _, rErr := api.ListCustomers(100, "page", "asc", "", agent.NewCustomersFilters())
	verifyErrorResponse("ListCustomers", rErr, t)
}

func TestCreateCustomerShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "create_customer"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	_, rErr := api.CreateCustomer("stubName", "stub@mail.com", "http://stub.url", []map[string]string{})
	verifyErrorResponse("CreateCustomer", rErr, t)
}
func TestUpdateCustomerShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "update_customer"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.UpdateCustomer("mister_customer", "stubName", "stub@mail.com", "http://stub.url", []map[string]string{})
	verifyErrorResponse("UpdateCustomer", rErr, t)
}

func TestBanCustomerShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "ban_customer"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.BanCustomer("mister_customer", 20)
	verifyErrorResponse("BanCustomer", rErr, t)
}

func TestSetRoutingStatusShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "set_routing_status"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.SetRoutingStatus("some_agent", "accepting chats")
	verifyErrorResponse("SetRoutingStatus", rErr, t)
}

func TestMarkEventsAsSeenShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "mark_events_as_seen"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.MarkEventsAsSeen("stubChatID", time.Time{})
	verifyErrorResponse("MarkEventsAsSeen", rErr, t)
}

func TestSendTypingIndicatorShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "send_typing_indicator"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.SendTypingIndicator("stubChatID", "all", true)
	verifyErrorResponse("SendTypingIndicator", rErr, t)
}

func TestMulticastShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "multicast"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.Multicast(agent.MulticastRecipients{}, []byte("{}"), "type")
	verifyErrorResponse("Multicast", rErr, t)
}

func TestTransferChatShouldNotCrashOnErrorResponse(t *testing.T) {
	client := NewTestClient(createMockedErrorResponder(t, "transfer_chat"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}
	ids := make([]interface{}, 1)
	ids[0] = 1
	rErr := api.TransferChat("stubChatID", "group", ids, agent.TransferChatOptions{})
	verifyErrorResponse("SendTypingIndicator", rErr, t)
}

func TestListAgentsForTransferShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "list_agents_for_transfer"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	resp, rErr := api.ListAgentsForTransfer("PJ0MRSHTDG")
	if rErr != nil {
		t.Errorf("ListAgentsForTransfer failed: %v", rErr)
	}

	if len(resp) != 2 {
		t.Errorf("Invalid ListAgentsForTransfer response: %v", resp)
	}
}

func TestBasicAuthorizationScheme(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		if authHeader := req.Header.Get("Authorization"); authHeader != "Basic access_token" {
			t.Errorf("Invalid Authorization header: %s", authHeader)
		}
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(nil)),
			Header:     make(http.Header),
		}
	})

	api, err := agent.NewAPI(stubTokenGetter(authorization.BasicToken), client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}
	api.Call("", nil, nil)
}

func TestBearerAuthorizationScheme(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		if authHeader := req.Header.Get("Authorization"); authHeader !=
			"Bearer access_token" {
			t.Errorf("Invalid Authorization header: %s", authHeader)
		}
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(nil)),
			Header:     make(http.Header),
		}
	})

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}
	api.Call("", nil, nil)
}

func TestInvalidAuthorizationScheme(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response { return nil })

	api, err := agent.NewAPI(stubTokenGetter(authorization.TokenType(2020)), client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}
	err = api.Call("", nil, nil)
	if err == nil {
		t.Error("Err should not be nil")
	}
}

func TestRetryStrategyAllFails(t *testing.T) {
	client := NewTestClient(createMockedMultipleAuthErrorsResponder(t, 10))

	api, err := agent.NewAPI(stubTokenGetter(authorization.BearerToken), client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	var retries uint
	api.SetRetryStrategy(func(attempts uint, err error) bool {
		if attempts < 3 {
			retries++
			return true
		}

		return false
	})

	err = api.Call("", nil, nil)
	if err == nil {
		t.Error("Err should not be nil")
	}

	if retries != 3 {
		t.Error("Retries should be done 3 times")
	}

}

func TestRetryStrategyLastSuccess(t *testing.T) {
	client := NewTestClient(createMockedMultipleAuthErrorsResponder(t, 2))

	api, err := agent.NewAPI(stubTokenGetter(authorization.BearerToken), client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	var retries uint
	api.SetRetryStrategy(func(attempts uint, err error) bool {
		if attempts < 3 {
			retries++
			return true
		}

		return false
	})

	err = api.Call("", nil, &struct{}{})
	if err != nil {
		t.Error("Err should be nil after 2 retries")
	}

	if retries != 2 {
		t.Error("Retries should be done 2 times")
	}

}

func TestFollowCustomerShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "follow_customer"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	err = api.FollowCustomer("foo")
	if err != nil {
		t.Errorf("FollowCustomer failed: %v", err)
	}
}

func TestUnfollowCustomerShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "unfollow_customer"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	err = api.UnfollowCustomer("foo")
	if err != nil {
		t.Errorf("UnfollowCustomer failed: %v", err)
	}
}

func TestListRoutingStatusesShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "list_routing_statuses"))

	api, err := agent.NewAPI(stubBearerTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	resp, err := api.ListRoutingStatuses([]int{})
	if err != nil {
		t.Errorf("ListRoutingStatuses failed: %v", err)
	}

	if len(resp) != 2 {
		t.Errorf("Invalid response length: %v", len(resp))
	}

	if resp[0].AgentID != "smith@example.com" {
		t.Errorf("Invalid agent_id response: %v", resp[0].AgentID)
	}

	if resp[0].Status != "accepting_chats" {
		t.Errorf("Invalid status response: %v", resp[0].Status)
	}
}
