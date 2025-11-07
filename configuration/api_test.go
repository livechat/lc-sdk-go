package configuration_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/livechat/lc-sdk-go/v7/authorization"
	"github.com/livechat/lc-sdk-go/v7/configuration"
)

const (
	ExpectedNewGroupID          = 19
	ExpectedNewAutoAccessID     = "pqi8oasdjahuakndw9nsad9na"
	ExpectedPropertiesNamespace = "0805e283233042b37f460ed8fbf22160"
)

type serverMock struct {
	LastRequest *http.Request
	Method      string
}

func (s *serverMock) RoundTrip(req *http.Request) (*http.Response, error) {
	s.LastRequest = req
	if responseNOK := s.validateCommon(req); responseNOK != nil {
		return responseNOK, nil
	}
	return getMockResponseOK(s.Method), nil
}

// Performs validation of common parameters. Returns nil if everything is OK.
// Does not validate body of the request which is specific for each method.
func (s *serverMock) validateCommon(req *http.Request) *http.Response {
	if req.URL.String() != "https://api.livechatinc.com/v3.7/configuration/action/"+s.Method {
		return getMockResponseNOK("Invalid URL")
	}

	if req.Method != "POST" {
		return getMockResponseNOK("Invalid URL")
	}

	if authHeader := req.Header.Get("Authorization"); authHeader != "Bearer access_token" {
		return getMockResponseNOK("Invalid Authorization")
	}

	if regionHeader := req.Header.Get("X-Region"); regionHeader != "region" {
		return getMockResponseNOK("Invalid X-Region")
	}
	return nil
}

func NewTestClient(s *serverMock) *http.Client {
	return &http.Client{
		Transport: s,
	}
}

func stubTokenGetter() *authorization.Token {
	return &authorization.Token{
		AccessToken: "access_token",
		Region:      "region",
	}
}

var mockedResponses = map[string]string{
	"register_webhook": `{
		"id": "pqi8oasdjahuakndw9nsad9na"
	}`,
	"list_webhooks": `[
		{
			"id": "pqi8oasdjahuakndw9nsad9na",
			"url": "http://myservice.com/webhooks",
			"description": "Test webhook",
			"action": "thread_closed",
			"secret_key": "laudla991lamda0pnoaa0",
			"type": "license",
			"filters": {
				"chat_presence": {
					"my_bots": true,
					"user_ids": {
						"value": ["johndoe@mail.com"]
					}
				}
			},
			"owner_client_id": "asXdesldiAJSq9padj"
		}
	]`,
	"unregister_webhook": `{}`,
	"create_bot": `{
		"id": "5c9871d5372c824cbf22d860a707a578",
		"secret": "laudla991lamda0pnoaa0"
	}`,
	"update_bot": `{}`,
	"delete_bot": `{}`,
	"list_bots": `[
		{
			"id": "5c9871d5372c824cbf22d860a707a578",
			"name": "John Doe",
			"avatar": "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg"
		},
		{
			"id": "8g1231ss112c013cbf11d530b595h987",
			"name": "Jason Brown",
			"avatar": "livechat.s3.amazonaws.com/1011121/all/avatars/wff9482gkdjanzjgdsf88a184jsskaz1.jpg"
		}
	]`,
	"get_bot": `{
		"id": "5c9871d5372c824cbf22d860a707a578",
		"name": "John Doe",
		"avatar": "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg",
		"default_group_priority": "first",
		"owner_client_id": "asXdesldiAJSq9padj",
		"max_chats_count": 6,
		"groups": [{
			"id": 0,
			"priority": "normal"
		}, {
			"id": 1,
			"priority": "normal"
		}, {
			"id": 2,
			"priority": "first"
		}]
	}`,
	"reset_bot_secret": `{
		"secret": "laudla991lamda0pnoaa0"
	}`,
	"issue_bot_token": `{
		"token": "bbd8924fcbcdbddbeaf60c19b238b0b0"
	}`,
	"create_bot_template": `{
		"id": "5c9871d5372c824cbf22d860a707a578",
		"secret": "laudla991lamda0pnoaa0"
	}`,
	"update_bot_template": `{}`,
	"delete_bot_template": `{}`,
	"list_bot_templates": `[
		{
			"id": "5c9871d5372c824cbf22d860a707a578",
			"name": "John Doe",
			"avatar": "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg",
			"max_chats_count": 2137,
			"default_group_priority": "first",
			"job_title": "b0t"
		},
		{
			"id": "8g1231ss112c013cbf11d530b595h987",
			"name": "Jason Brown",
			"avatar": "livechat.s3.amazonaws.com/1011121/all/avatars/wff9482gkdjanzjgdsf88a184jsskaz1.jpg",
			"max_chats_count": 69,
			"default_group_priority": "first",
			"job_title": "b0t"
		}
	]`,
	"register_property":   `{}`,
	"unregister_property": `{}`,
	"publish_property":    `{}`,
	"list_properties": `{
		"dummy_property": {
			"type": "string",
			"description": "This is a dummy property",
			"access": {
				"chat": {
					"agent": ["read", "write"],
					"customer": ["read"]
				},
				"group": {
					"agent": ["write"]
				}
			},
			"domain": [
				"hello",
				"hi"
			],
			"default_value": "hi"
		}
	}`,
	"list_license_properties": fmt.Sprintf(`{
		"%s": {
				"string_property": "string value"
		}
	}`, ExpectedPropertiesNamespace),
	"create_agent": `{
		"id": "smith@example.com"
	}`,
	"get_agent": `{
		"id": "smith@example.com",
		"account_id": "d24fa41e-bc16-41b8-a15b-9ca45ff7e0cf",
		"name": "Agent Smith",
		"avatar": "https://domain.com/avatar.image.jpg",
		"role": "administrator",
		"login_status": "accepting chats"
	}`,
	"list_agents": `[
		{
			"id": "smith@example.com",
			"account_id": "d24fa41e-bc16-41b8-a15b-9ca45ff7e0cf",
			"job_title": "Support Hero",
			"max_chats_count": 5,
			"last_logout": "2022-08-23T14:31:21.000000Z",
			"summaries": [
				"daily_summary",
				"weekly_summary"
			]
		},
		{
			"id": "adam@example.com",
			"account_id": "2fca3315-f6b2-422f-8550-c84b649cef1a",
			"job_title": "Support Hero (Newbie)",
			"max_chats_count": 2,
			"summaries": [
				"weekly_summary"
			]
		}
	]`,
	"update_agent":               `{}`,
	"delete_agent":               `{}`,
	"suspend_agent":              `{}`,
	"unsuspend_agent":            `{}`,
	"request_agent_unsuspension": `{}`,
	"approve_agent":              `{}`,
	"create_group": fmt.Sprintf(`{
		"id": %d
	}`, ExpectedNewGroupID),
	"update_group": `{}`,
	"delete_group": `{}`,
	"list_groups": `[
		{
			"id": 0,
			"name": "General",
			"language_code": "en",
			"routing_status": "offline"
		},
		{
			"id": 19,
			"name": "Sport shoes",
			"language_code": "en",
			"routing_status": "offline"
		}
	]`,
	"get_group": `{
		"id": 1,
		"name": "Sports shoes",
		"language_code": "en",
		"agent_priorities": {
		  "agent1@example.com": "normal",
		  "agent2@example.com": "normal",
		  "agent3@example.com": "last"
		},
		"routing_status": "offline"
	}`,
	"list_webhook_names": `[
		{
			"action": "chat_access_granted",
			"filters": [
				"chat_member_ids",
				"only_my_chats"
			],
			"additional_data": [
				"chat_properties",
				"chat_presence_user_ids"
			]
		},
		{
			"action": "event_properties_deleted",
			"filters": [
				"chat_member_ids",
				"only_my_chats"
			],
			"additional_data": [
				"chat_properties",
				"chat_presence_user_ids"
			]
		}
	]`,
	"enable_license_webhooks":  `{}`,
	"disable_license_webhooks": `{}`,
	"get_license_webhooks_state": `{
		"license_webhooks_enabled": true
	}`,
	"delete_license_properties": `{}`,
	"delete_group_properties":   `{}`,
	"update_license_properties": `{}`,
	"update_group_properties":   `{}`,
	"add_auto_access":           fmt.Sprintf(`{ "id": "%s" }`, ExpectedNewAutoAccessID),
	"delete_auto_access":        `{}`,
	"update_auto_access":        `{}`,
	"list_auto_accesses": `[
		{
			"id": "1faad6f5f1d6e8fdf27e8af9839783b7",
			"description": "Chats on livechat.com from United States",
			"access": {
				"groups": [
					0
				]
			},
			"conditions": {
				"geolocation": {
					"values": [
						{
							"country": "United States",
							"country_code": "US"
						}
					]
				},
				"domain": {
					"values": [
						{
							"value": "livechat.com",
							"exact_match": true
						}
					]
				}
			},
			"next_id": "pqi8oasdjahuakndw9nsad9na"
		}
	]`,
	"add_allowed_domain":    `{}`,
	"delete_allowed_domain": `{}`,
	"list_allowed_domains": `{
		"domains": ["example.com", "test.com"]
	}`,
	"check_product_limits_for_plan": `[
		{
			"resource": "groups",
			"limit_balance": 1
		},
		{
			"resource": "groups_per_agent",
			"limit_balance": 2,
			"id": "user@example.com"
		},
		{
			"resource": "group_chooser_groups",
			"limit_balance": 3,
			"id": "0"
		}
	]`,
	"list_channels": `[
		{
			"channel_type": "code",
			"channel_subtype": "",
			"first_activity_timestamp": "2017-10-12T13:56:16Z"
		},
		{
			"channel_type": "direct_link",
			"channel_subtype": "",
			"first_activity_timestamp": "2017-10-12T15:20:00Z"
		},
		{
			"channel_type": "integration",
			"channel_subtype": "c6e4f62e2a2dab12531235b12c5a2a6b",
			"first_activity_timestamp": "2019-08-16T16:55:51Z"
		}
	]`,
	"create_tag": `{}`,
	"delete_tag": `{}`,
	"list_tags": `[
		{
			"name": "tageroo",
			"group_ids": [
				0
			],
			"created_at": "2017-10-12T13:56:16Z",
			"author_id": "smith@example.com"
		},
		{
			"name": "tagonanza",
			"group_ids": [
				1
			],
			"created_at": "2019-08-16T16:55:51Z",
			"author_id": "jones@example.com"
		}
	]`,
	"update_tag": `{}`,
	"list_groups_properties": `[
		{
			"id": 0,
			"properties": {
				"abc": {
					"a_property": "a"
				}
			}
		}
	]`,
	"reactivate_email":       `{}`,
	"update_company_details": `{}`,
	"create_canned_response": `{"id": 2137}`,
	"delete_canned_response": `{}`,
	"list_canned_responses": `
		{"canned_responses": [
			{"id": 1, "text": "text", "tags": ["tag"], "group_id": 0, "is_private": false, "author_id": "author_id"},
			{"id": 3, "text": "text2", "tags": ["tag2","tag3"], "group_id": 1, "is_private": true, "author_id": "author_id2"}
		],
		"found_canned_responses": 10,
		"next_page_id": "next_page=="
	}`,
	"update_translations": `{}`,
	"list_translations":   `{"phrases":{"hello":"cześć","bye":"do widzenia"}}`,
}

func getMockResponseOK(method string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(mockedResponses[method])),
		Header:     make(http.Header),
	}
}

func getMockResponseNOK(message string) *http.Response {
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

func newServerMock(t *testing.T, method string) *serverMock {
	return &serverMock{
		Method: method}
}

func validateRequestBody(t *testing.T, want string, got io.ReadCloser) {
	t.Helper()

	body, err := ioutil.ReadAll(got)
	if err != nil {
		t.Errorf("Error reading request body: %s", err)
	}

	if string(body) != want {
		t.Errorf("Request body mismatch\nwant: %s\ngot: %s", want, string(body))
	}
}

func TestRejectAPICreationWithoutTokenGetter(t *testing.T) {
	_, err := configuration.NewAPI(nil, nil, "client_id")
	if err == nil {
		t.Error("API should not be created without token getter")
	}
}

func TestRegisterWebhookShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(newServerMock(t, "register_webhook"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	cpf := configuration.NewChatPresenceFilter().WithMyBots().WithUserIDs([]string{"agent@smith.com"}, true)
	webhookID, rErr := api.RegisterWebhook(&configuration.Webhook{
		Filters: &configuration.WebhookFilters{
			ChatPresence: cpf,
		},
	}, nil)
	if rErr != nil {
		t.Errorf("RegisterWebhook failed: %v", rErr)
	}

	if webhookID != "pqi8oasdjahuakndw9nsad9na" {
		t.Errorf("Invalid webhookID: %v", webhookID)
	}
}

func TestListWebhooksShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(newServerMock(t, "list_webhooks"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	resp, rErr := api.ListWebhooks(nil)
	if rErr != nil {
		t.Errorf("ListWebhooks failed: %v", rErr)
	}

	if len(resp) != 1 || resp[0].ID != "pqi8oasdjahuakndw9nsad9na" {
		t.Errorf("Invalid webhooks config: %v", resp)
	}
}

func TestUnregisterWebhookShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(newServerMock(t, "unregister_webhook"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.UnregisterWebhook("pqi8oasdjahuakndw9nsad9na", nil)
	if rErr != nil {
		t.Errorf("UnregisterWebhook failed: %v", rErr)
	}
}

func TestCreateBotOK(t *testing.T) {
	serverMock := newServerMock(t, "create_bot")
	client := NewTestClient(serverMock)

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	checkAPIrespondedOK := func(t *testing.T, res *configuration.CreateBotResponse, rErr error) {
		t.Helper()
		if rErr != nil {
			t.Errorf("CreateBot failed: %v", rErr)
		}
		if res.BotID != "5c9871d5372c824cbf22d860a707a578" {
			t.Errorf("Invalid botID: %v", res.BotID)
		}
		if res.Secret != "laudla991lamda0pnoaa0" {
			t.Errorf("Invalid secret: %v", res.Secret)
		}
	}

	t.Run("Only required fields", func(t *testing.T) {
		res, rErr := api.CreateBot("John Doe", nil)
		wantReq := `{"name":"John Doe"}`

		checkAPIrespondedOK(t, res, rErr)
		validateRequestBody(t, wantReq, serverMock.LastRequest.Body)
	})

	t.Run("All optional fields", func(t *testing.T) {
		res, rErr := api.CreateBot("John Doe", &configuration.CreateBotRequestOptions{
			Avatar: "https://example.com/avatar.png",
			Groups: []configuration.GroupConfig{
				{ID: 6, Priority: "first"},
			},
			OwnerClientID: "dummy_client_id",
			WorkScheduler: &configuration.WorkScheduler{
				Timezone: "dummy/timezone",
				Schedule: []configuration.Schedule{
					{
						Enabled: true,
						Day:     "monday",
						Start:   "09:00",
						End:     "17:00",
					},
				},
			},
		})
		wantReq := `{"name":"John Doe","avatar":"https://example.com/avatar.png","groups":[{"id":6,"priority":"first"}],"owner_client_id":"dummy_client_id","work_scheduler":{"timezone":"dummy/timezone","schedule":[{"enabled":true,"day":"monday","start":"09:00","end":"17:00"}]}}`

		checkAPIrespondedOK(t, res, rErr)
		validateRequestBody(t, wantReq, serverMock.LastRequest.Body)
	})

	t.Run("No work scheduler provided", func(t *testing.T) {
		res, rErr := api.CreateBot("John Doe", &configuration.CreateBotRequestOptions{
			Avatar: "https://example.com/avatar.png",
		})
		wantReq := `{"name":"John Doe","avatar":"https://example.com/avatar.png"}`

		checkAPIrespondedOK(t, res, rErr)
		validateRequestBody(t, wantReq, serverMock.LastRequest.Body)
	})

}

func TestCreateBotShouldReturnErrorForInvalidInput(t *testing.T) {
	client := NewTestClient(newServerMock(t, "create_bot"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	_, rErr := api.CreateBot("John Doe", &configuration.CreateBotRequestOptions{Avatar: "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg", Groups: []configuration.GroupConfig{{6, "supervisor"}}, OwnerClientID: "dummy_client_id", WorkScheduler: &configuration.WorkScheduler{Timezone: "dummy/timezone"}})
	if rErr.Error() != "DoNotAssign priority is allowed only as default group priority" {
		t.Errorf("CreateBot failed: %v", rErr)
	}
}

func TestUpdateBotOK(t *testing.T) {
	serverMock := newServerMock(t, "update_bot")
	client := NewTestClient(serverMock)

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}
	t.Run("Only required fields", func(t *testing.T) {
		if rErr := api.UpdateBot("pqi8oasdjahuakndw9nsad9na", nil); rErr != nil {
			t.Errorf("UpdateBot failed: %v", rErr)
		}
		validateRequestBody(t, `{"id":"pqi8oasdjahuakndw9nsad9na"}`, serverMock.LastRequest.Body)
	})
	t.Run("No work scheduler provided", func(t *testing.T) {
		if rErr := api.UpdateBot("pqi8oasdjahuakndw9nsad9na", &configuration.UpdateBotRequestOptions{
			Avatar: "https://example.com/avatar.png",
		}); rErr != nil {
			t.Errorf("UpdateBot failed: %v", rErr)
		}
		validateRequestBody(t, `{"id":"pqi8oasdjahuakndw9nsad9na","avatar":"https://example.com/avatar.png"}`, serverMock.LastRequest.Body)
	})
	t.Run("All optional fields", func(t *testing.T) {
		if rErr := api.UpdateBot("pqi8oasdjahuakndw9nsad9na", &configuration.UpdateBotRequestOptions{
			Avatar: "https://example.com/avatar.png",
			Groups: []configuration.GroupConfig{
				{ID: 6, Priority: "first"},
			},
			OwnerClientID: "dummy_client_id",
			WorkScheduler: &configuration.WorkScheduler{
				Timezone: "dummy/timezone",
				Schedule: []configuration.Schedule{
					{
						Enabled: true,
						Day:     "monday",
						Start:   "09:00",
						End:     "17:00",
					},
				},
			},
		}); rErr != nil {
			t.Errorf("UpdateBot failed: %v", rErr)
		}
		validateRequestBody(t, `{"id":"pqi8oasdjahuakndw9nsad9na","avatar":"https://example.com/avatar.png","groups":[{"id":6,"priority":"first"}],"owner_client_id":"dummy_client_id","work_scheduler":{"timezone":"dummy/timezone","schedule":[{"enabled":true,"day":"monday","start":"09:00","end":"17:00"}]}}`, serverMock.LastRequest.Body)
	})

}

func TestUpdateBotShouldReturnErrorForInvalidInput(t *testing.T) {
	client := NewTestClient(newServerMock(t, "update_bot"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	groups := []configuration.GroupConfig{{Priority: "supervisor"}}
	rErr := api.UpdateBot("pqi8oasdjahuakndw9nsad9na", &configuration.UpdateBotRequestOptions{Groups: groups})
	if rErr.Error() != "DoNotAssign priority is allowed only as default group priority" {
		t.Errorf("UpdateBot failed: %v", rErr)
	}
}

func TestDeleteBotShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(newServerMock(t, "delete_bot"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.DeleteBot("pqi8oasdjahuakndw9nsad9na")
	if rErr != nil {
		t.Errorf("DeleteBot failed: %v", rErr)
	}
}

func TestListBotsShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(newServerMock(t, "list_bots"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	resp, rErr := api.ListBots(true, []string{})
	if rErr != nil {
		t.Errorf("ListBots failed: %v", rErr)
	}

	if len(resp) != 2 {
		t.Errorf("Invalid number of bots: %v", len(resp))
	}

	if resp[0].ID != "5c9871d5372c824cbf22d860a707a578" {
		t.Errorf("Invalid bot ID: %v", resp[0].ID)
	}

	if resp[1].ID != "8g1231ss112c013cbf11d530b595h987" {
		t.Errorf("Invalid bot ID: %v", resp[1].ID)
	}
}

func TestGetBotShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(newServerMock(t, "get_bot"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	resp, rErr := api.GetBot("5c9871d5372c824cbf22d860a707a578", []string{})
	if rErr != nil {
		t.Errorf("GetBot failed: %v", rErr)
	}

	if resp.ID != "5c9871d5372c824cbf22d860a707a578" {
		t.Errorf("Invalid bot: %v", resp)
	}
}

func TestResetBotSecretOK(t *testing.T) {
	serverMock := newServerMock(t, "reset_bot_secret")
	client := NewTestClient(serverMock)

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}
	res, rErr := api.ResetBotSecret("5c9871d5372c824cbf22d860a707a578", nil)
	if rErr != nil {
		t.Errorf("ResetBotSecret failed: %v", rErr)
	}
	if res != "laudla991lamda0pnoaa0" {
		t.Error("Invalid secret returned")
	}

	validateRequestBody(t, `{"id":"5c9871d5372c824cbf22d860a707a578"}`, serverMock.LastRequest.Body)
}

func TestIssueBotTokenOK(t *testing.T) {
	serverMock := newServerMock(t, "issue_bot_token")
	client := NewTestClient(serverMock)

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}
	res, rErr := api.IssueBotToken(configuration.IssueBotTokenRequest{
		BotID:          "5c9871d5372c824cbf22d860a707a578",
		BotSecret:      "laudla991lamda0pnoaa0",
		OrganizationID: "aaabbbcccdddeeefff",
		ClientID:       "dummy_client_id",
	})
	if rErr != nil {
		t.Errorf("IssueBotToken failed: %v", rErr)
	}
	if res != "bbd8924fcbcdbddbeaf60c19b238b0b0" {
		t.Error("Invalid token returned")
	}

	validateRequestBody(t, `{"bot_id":"5c9871d5372c824cbf22d860a707a578","bot_secret":"laudla991lamda0pnoaa0","organization_id":"aaabbbcccdddeeefff","client_id":"dummy_client_id"}`, serverMock.LastRequest.Body)
}

func TestCreateBotTemplateOK(t *testing.T) {
	serverMock := newServerMock(t, "create_bot_template")
	client := NewTestClient(serverMock)

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}
	t.Run("Only required fields", func(t *testing.T) {
		res, rErr := api.CreateBotTemplate("John Doe", nil)
		wantReq := `{"name":"John Doe"}`

		if rErr != nil {
			t.Errorf("CreateBotTemplate failed: %v", rErr)
		}
		if res.BotTemplateID != "5c9871d5372c824cbf22d860a707a578" {
			t.Errorf("Invalid botTemplateID: %v", res.BotTemplateID)
		}
		if res.Secret != "laudla991lamda0pnoaa0" {
			t.Errorf("Invalid secret: %v", res.Secret)
		}
		validateRequestBody(t, wantReq, serverMock.LastRequest.Body)
	})

	t.Run("All optional fields", func(t *testing.T) {
		maxChatsCount := uint(420)
		res, rErr := api.CreateBotTemplate("John Doe", &configuration.CreateBotTemplateRequestOptions{
			Avatar:                      "https://example.com/avatar.png",
			MaxChatsCount:               &maxChatsCount,
			DefaultGroupPriority:        "first",
			JobTitle:                    "b0t",
			OwnerClientID:               "dummy_client_id",
			AffectExistingInstallations: true})
		wantReq := `{"name":"John Doe","avatar":"https://example.com/avatar.png","max_chats_count":420,"default_group_priority":"first","job_title":"b0t","owner_client_id":"dummy_client_id","affect_existing_installations":true}`

		if rErr != nil {
			t.Errorf("CreateBotTemplate failed: %v", rErr)
		}
		if res.BotTemplateID != "5c9871d5372c824cbf22d860a707a578" {
			t.Errorf("Invalid botTemplateID: %v", res.BotTemplateID)
		}
		if res.Secret != "laudla991lamda0pnoaa0" {
			t.Errorf("Invalid secret: %v", res.Secret)
		}
		validateRequestBody(t, wantReq, serverMock.LastRequest.Body)
	})

}

func TestUpdateBotTemplateOK(t *testing.T) {
	serverMock := newServerMock(t, "update_bot_template")
	client := NewTestClient(serverMock)

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}
	t.Run("Only required fields", func(t *testing.T) {
		if rErr := api.UpdateBotTemplate("5c9871d5372c824cbf22d860a707a578", nil); rErr != nil {
			t.Errorf("UpdateBotTemplate failed: %v", rErr)
		}
		validateRequestBody(t, `{"id":"5c9871d5372c824cbf22d860a707a578"}`, serverMock.LastRequest.Body)
	})

	t.Run("All optional fields", func(t *testing.T) {
		maxChatsCount := uint(420)
		if rErr := api.UpdateBotTemplate("5c9871d5372c824cbf22d860a707a578", &configuration.UpdateBotTemplateRequestOptions{
			Avatar:                      "https://example.com/avatar.png",
			MaxChatsCount:               &maxChatsCount,
			DefaultGroupPriority:        "first",
			JobTitle:                    "b0t",
			OwnerClientID:               "dummy_client_id",
			AffectExistingInstallations: true}); rErr != nil {
			t.Errorf("UpdateBotTemplate failed: %v", rErr)
		}
		validateRequestBody(t, `{"id":"5c9871d5372c824cbf22d860a707a578","avatar":"https://example.com/avatar.png","max_chats_count":420,"default_group_priority":"first","job_title":"b0t","owner_client_id":"dummy_client_id","affect_existing_installations":true}`, serverMock.LastRequest.Body)
	})
}

func TestDeleteBotTemplateOK(t *testing.T) {
	serverMock := newServerMock(t, "delete_bot_template")
	client := NewTestClient(serverMock)

	t.Run("ID only", func(t *testing.T) {
		api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
		if err != nil {
			t.Error("API creation failed")
		}
		if rErr := api.DeleteBotTemplate("5c9871d5372c824cbf22d860a707a578", nil); rErr != nil {
			t.Errorf("DeleteBotTemplate failed: %v", rErr)
		}
		validateRequestBody(t, `{"id":"5c9871d5372c824cbf22d860a707a578","affect_existing_installations":false}`, serverMock.LastRequest.Body)
	})

	t.Run("reqeust with ID, ownerClientID and affectExistingInstallations", func(t *testing.T) {
		api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
		if err != nil {
			t.Error("API creation failed")
		}
		if rErr := api.DeleteBotTemplate("5c9871d5372c824cbf22d860a707a578", &configuration.DeleteBotTemplateRequestOptions{OwnerClientID: "dummy_client_id", AffectExistingInstallations: true}); rErr != nil {
			t.Errorf("DeleteBotTemplate failed: %v", rErr)
		}
		validateRequestBody(t, `{"id":"5c9871d5372c824cbf22d860a707a578","owner_client_id":"dummy_client_id","affect_existing_installations":true}`, serverMock.LastRequest.Body)
	})
}

func TestListBotTemplatesOK(t *testing.T) {
	serverMock := newServerMock(t, "list_bot_templates")
	client := NewTestClient(serverMock)

	t.Run("No optional fields", func(t *testing.T) {
		api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
		if err != nil {
			t.Error("API creation failed")
		}
		res, rErr := api.ListBotTemplates(nil)
		if rErr != nil {
			t.Errorf("ListBotTemplates failed: %v", rErr)
		}
		if len(res) != 2 {
			t.Errorf("Invalid number of bots: %v", len(res))
		}
		if res[0].ID != "5c9871d5372c824cbf22d860a707a578" {
			t.Errorf("Invalid botTemplateID: %v", res[0].ID)
		}
		if res[0].JobTitle != "b0t" {
			t.Errorf("Invalid job title: %v", res[0].JobTitle)
		}
		validateRequestBody(t, "{}", serverMock.LastRequest.Body)
	})
	t.Run("With optional fields", func(t *testing.T) {
		api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
		if err != nil {
			t.Error("API creation failed")
		}
		res, rErr := api.ListBotTemplates(&configuration.ListBotTemplatesRequestOptions{OwnerClientID: "dummy_client_id"})
		if rErr != nil {
			t.Errorf("ListBotTemplates failed: %v", rErr)
		}
		if len(res) != 2 {
			t.Errorf("Invalid number of bots: %v", len(res))
		}
		validateRequestBody(t, `{"owner_client_id":"dummy_client_id"}`, serverMock.LastRequest.Body)
	})
}

func TestRegisterPropertyShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(newServerMock(t, "register_property"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.RegisterProperty(&configuration.PropertyConfig{
		Name:          "dummy_property",
		OwnerClientID: "dummy_client_id",
		Type:          "int",
		Domain:        []interface{}{2, 1, 3, 7},
		Description:   "This is a dummy property",
		Access: map[string]*configuration.PropertyAccess{
			"chat": {
				Agent:    []string{"read", "write"},
				Customer: []string{"read"},
			},
		},
		DefaultValue: 7,
	})
	if rErr != nil {
		t.Errorf("RegisterProperty failed: %v", rErr)
	}
}

func TestUnregisterPropertyShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(newServerMock(t, "unregister_property"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.UnregisterProperty("dummy_property", "dummy_client_id")
	if rErr != nil {
		t.Errorf("UnregisterProperty failed: %v", rErr)
	}
}

func TestPublishPropertyShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(newServerMock(t, "publish_property"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.PublishProperty("dummy_property", "dummy_client_id", true, false)
	if rErr != nil {
		t.Errorf("PublishProperty failed: %v", rErr)
	}
}

func TestListPropertiesShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(newServerMock(t, "list_properties"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	resp, rErr := api.ListProperties("dummy_client_id")
	if rErr != nil {
		t.Errorf("ListProperties failed: %v", rErr)
	}

	_, exists := resp["dummy_property"]
	if !exists || len(resp) != 1 {
		t.Errorf("Invalid property configs: %v", resp)
	}
}

func TestListLicensePropertiesOK(t *testing.T) {
	serverMock := newServerMock(t, "list_license_properties")
	client := NewTestClient(serverMock)

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	checkAPIrespondedOK := func(t *testing.T, resp configuration.Properties, rErr error) {
		t.Helper()
		if rErr != nil {
			t.Errorf("ListLicenseProperties failed: %v", rErr)
		}
		if len(resp) != 1 {
			t.Errorf("Invalid license properties: %v", resp)
		}

		if resp[ExpectedPropertiesNamespace]["string_property"] != "string value" {
			t.Errorf("Invalid license property %s.string_property: %v", ExpectedPropertiesNamespace, resp[ExpectedPropertiesNamespace]["string_property"])
		}
	}

	t.Run("No optional fields", func(t *testing.T) {
		resp, rErr := api.ListLicenseProperties(nil)
		checkAPIrespondedOK(t, resp, rErr)
		validateRequestBody(t, "{}", serverMock.LastRequest.Body)
	})
	t.Run("With optional fields", func(t *testing.T) {
		resp, rErr := api.ListLicenseProperties(&configuration.ListLicensePropertiesRequestOptions{
			Namespace:  "namespace",
			NamePrefix: "prefix",
		})
		checkAPIrespondedOK(t, resp, rErr)
		validateRequestBody(t, `{"namespace":"namespace","name_prefix":"prefix"}`, serverMock.LastRequest.Body)
	})

}

func TestCreateAgentShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(newServerMock(t, "create_agent"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	agentID, rErr := api.CreateAgent("smith@example.com", &configuration.AgentFields{Name: "Agent Smith"})
	if rErr != nil {
		t.Errorf("CreateAgent failed: %v", rErr)
	}

	if agentID != "smith@example.com" {
		t.Errorf("Invalid agent ID: %v", agentID)
	}
}

func TestGetAgentShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(newServerMock(t, "get_agent"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	agent, rErr := api.GetAgent("smith@example.com", []string{})
	if rErr != nil {
		t.Errorf("CreateAgent failed: %v", rErr)
	}

	if agent.ID != "smith@example.com" {
		t.Errorf("Invalid agent ID: %v", agent.ID)
	}

	if agent.AccountID != "d24fa41e-bc16-41b8-a15b-9ca45ff7e0cf" {
		t.Errorf("Invalid agent account ID: %v", agent.AccountID)
	}

	if agent.Name != "Agent Smith" {
		t.Errorf("Invalid agent name: %v", agent.Name)
	}
}

func TestListAgentsShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(newServerMock(t, "list_agents"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	agents, rErr := api.ListAgents([]int32{0, 1}, []string{"last_logout"})
	if rErr != nil {
		t.Errorf("CreateAgent failed: %v", rErr)
	}

	if len(agents) != 2 {
		t.Errorf("Invalid number of agents: %v", len(agents))
	}

	if agents[0].ID != "smith@example.com" {
		t.Errorf("Invalid agent ID: %v", agents[0].ID)
	}

	if agents[0].AccountID != "d24fa41e-bc16-41b8-a15b-9ca45ff7e0cf" {
		t.Errorf("Invalid agent account ID: %v", agents[0].AccountID)
	}

	if agents[0].LastLogout != "2022-08-23T14:31:21.000000Z" {
		t.Errorf("Invalid agent last_logout: %v", agents[0].LastLogout)
	}

	if agents[1].ID != "adam@example.com" {
		t.Errorf("Invalid agent name: %v", agents[1].ID)
	}
}

func TestUpdateAgentShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(newServerMock(t, "update_agent"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.UpdateAgent("smith@example.com", &configuration.AgentFields{JobTitle: "Virus"})
	if rErr != nil {
		t.Errorf("UpdateAgent failed: %v", rErr)
	}
}

func TestDeleteAgentShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(newServerMock(t, "delete_agent"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.DeleteAgent("smith@example.com")
	if rErr != nil {
		t.Errorf("DeleteAgent failed: %v", rErr)
	}
}

func TestSuspendAgentShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(newServerMock(t, "suspend_agent"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.SuspendAgent("smith@example.com")
	if rErr != nil {
		t.Errorf("SuspendAgent failed: %v", rErr)
	}
}

func TestUnsuspendAgentShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(newServerMock(t, "unsuspend_agent"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.UnsuspendAgent("smith@example.com")
	if rErr != nil {
		t.Errorf("UnsuspendAgent failed: %v", rErr)
	}
}

func TestRequestAgentUnsuspensionShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(newServerMock(t, "request_agent_unsuspension"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.RequestAgentUnsuspension()
	if rErr != nil {
		t.Errorf("RequestAgentUnsuspension failed: %v", rErr)
	}
}

func TestApproveAgentShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(newServerMock(t, "approve_agent"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.ApproveAgent("smith@example.com")
	if rErr != nil {
		t.Errorf("ApproveAgent failed: %v", rErr)
	}
}

func TestCreateGroupOK(t *testing.T) {
	serverMock := newServerMock(t, "create_group")
	client := NewTestClient(serverMock)

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	checkAPIrespondedOK := func(t *testing.T, groupID int32, rErr error) {
		t.Helper()
		if rErr != nil {
			t.Errorf("CreateGroup failed: %v", rErr)
		}
		if groupID != ExpectedNewGroupID {
			t.Errorf("Invalid group id: %v", groupID)
		}
	}

	t.Run("Only required fields", func(t *testing.T) {
		groupID, rErr := api.CreateGroup("name", map[string]configuration.GroupPriority{}, nil)
		checkAPIrespondedOK(t, groupID, rErr)
		validateRequestBody(t, `{"name":"name","agent_priorities":{}}`, serverMock.LastRequest.Body)
	})

	t.Run("Required and optional fields", func(t *testing.T) {
		groupID, rErr := api.CreateGroup("name", map[string]configuration.GroupPriority{}, &configuration.CreateGroupRequestOptions{
			LanguageCode: "en",
		})
		checkAPIrespondedOK(t, groupID, rErr)
		validateRequestBody(t, `{"name":"name","agent_priorities":{},"language_code":"en"}`, serverMock.LastRequest.Body)
	})
}

func TestUpdateGroupOK(t *testing.T) {
	serverMock := newServerMock(t, "update_group")
	client := NewTestClient(serverMock)

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	t.Run("Only required fields", func(t *testing.T) {
		rErr := api.UpdateGroup(420, nil)
		if rErr != nil {
			t.Errorf("UpdateGroup failed: %v", rErr)
		}
		validateRequestBody(t, `{"id":420}`, serverMock.LastRequest.Body)
	})
	t.Run("Only name", func(t *testing.T) {
		rErr := api.UpdateGroup(420, &configuration.UpdateGroupRequestOptions{
			Name: "Foo",
		})
		if rErr != nil {
			t.Errorf("UpdateGroup failed: %v", rErr)
		}
		validateRequestBody(t, `{"id":420,"name":"Foo"}`, serverMock.LastRequest.Body)
	})
	t.Run("Required and optional fields", func(t *testing.T) {
		rErr := api.UpdateGroup(420, &configuration.UpdateGroupRequestOptions{
			Name:         "Foo",
			LanguageCode: "en",
			AgentPriorities: map[string]configuration.GroupPriority{
				"foo": configuration.First,
			},
		})
		if rErr != nil {
			t.Errorf("UpdateGroup failed: %v", rErr)
		}
		validateRequestBody(t, `{"id":420,"name":"Foo","language_code":"en","agent_priorities":{"foo":"first"}}`, serverMock.LastRequest.Body)
	})
}

func TestDeleteGroupShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(newServerMock(t, "delete_group"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.DeleteGroup(11)
	if rErr != nil {
		t.Errorf("DeleteGroup failed: %v", rErr)
	}
}

func TestListGroupsShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(newServerMock(t, "list_groups"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	groups, rErr := api.ListGroups([]string{})
	if rErr != nil {
		t.Errorf("DeleteGroup failed: %v", rErr)
	}

	if len(groups) != 2 {
		t.Errorf("Invalid groups length: %v", len(groups))
	}

	if groups[0].ID != 0 {
		t.Errorf("Invalid group ID: %v", groups[0].ID)
	}

	if groups[1].ID != ExpectedNewGroupID {
		t.Errorf("Invalid group ID: %v", groups[1].ID)
	}
}

func TestGetGroupShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(newServerMock(t, "get_group"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	resp, rErr := api.GetGroup(1)
	if rErr != nil {
		t.Errorf("GetGroup failed: %v", rErr)
	}

	if resp.ID != 1 {
		t.Errorf("Invalid group id: %v", resp.ID)
	}

	if resp.LanguageCode != "en" {
		t.Errorf("Invalid group language: %v", resp.LanguageCode)
	}
}

func TestListWebhookNamesShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(newServerMock(t, "list_webhook_names"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	resp, rErr := api.ListWebhookNames("3.7")
	if rErr != nil {
		t.Errorf("ListWebhookNames failed: %v", rErr)
	}

	if len(resp) != 2 {
		t.Errorf("Invalid response length: %v", len(resp))
	}

	if resp[0].Action != "chat_access_granted" {
		t.Errorf("Invalid action in first element: %v", resp[0].Action)
	}
}

func TestEnableLicenseWebhooksShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(newServerMock(t, "enable_license_webhooks"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	err = api.EnableLicenseWebhooks(nil)
	if err != nil {
		t.Errorf("EnableWebhooks failed: %v", err)
	}
}

func TestDisableLicenseWebhooksShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(newServerMock(t, "disable_license_webhooks"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	err = api.DisableLicenseWebhooks(nil)
	if err != nil {
		t.Errorf("DisableWebhooks failed: %v", err)
	}
}

func TestGetLicenseWebhooksStateShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(newServerMock(t, "get_license_webhooks_state"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	state, err := api.GetLicenseWebhooksState(nil)
	if err != nil {
		t.Errorf("GetWebhooksState failed: %v", err)
	}
	if !state.Enabled {
		t.Error("webhooks' state should be enabled'")
	}
}

func TestDeleteLicensePropertiesShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(newServerMock(t, "delete_license_properties"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	err = api.DeleteLicenseProperties(nil)
	if err != nil {
		t.Errorf("DeleteLicenseProperties failed: %v", err)
	}
}

func TestDeleteGroupPropertiesShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(newServerMock(t, "delete_group_properties"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	err = api.DeleteGroupProperties(0, nil)
	if err != nil {
		t.Errorf("DeleteGroupProperties failed: %v", err)
	}
}

func TestUpdateLicensePropertiesShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(newServerMock(t, "update_license_properties"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	err = api.UpdateLicenseProperties(nil)
	if err != nil {
		t.Errorf("DeleteLicensePropertie s failed: %v", err)
	}
}

func TestUpdateGroupPropertiesShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(newServerMock(t, "update_group_properties"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	err = api.UpdateGroupProperties(0, nil)
	if err != nil {
		t.Errorf("DeleteGroupProperties failed: %v", err)
	}
}

func TestAddAutoAccessOK(t *testing.T) {
	serverMock := newServerMock(t, "add_auto_access")
	client := NewTestClient(serverMock)

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}
	checkAPIrespondedOK := func(t *testing.T, resp string, rErr error) {
		t.Helper()
		if rErr != nil {
			t.Errorf("AddAutoAccess failed: %v", rErr)
		}
		if resp != ExpectedNewAutoAccessID {
			t.Errorf("Invalid new auto access ID obtained: %v", resp)
		}
	}
	t.Run("only required fields", func(t *testing.T) {
		resp, err := api.AddAutoAccess(configuration.Access{[]int{}}, configuration.AutoAccessConditions{}, nil)
		checkAPIrespondedOK(t, resp, err)
		validateRequestBody(t, `{"access":{"groups":[]},"conditions":{}}`, serverMock.LastRequest.Body)
	})
	t.Run("required and optional fields", func(t *testing.T) {
		resp, err := api.AddAutoAccess(configuration.Access{[]int{}}, configuration.AutoAccessConditions{}, &configuration.AddAutoAccessRequestOptions{Description: "foo"})
		checkAPIrespondedOK(t, resp, err)
		validateRequestBody(t, `{"access":{"groups":[]},"conditions":{},"description":"foo"}`, serverMock.LastRequest.Body)
	})

}

func addressOf(in string) *string {
	return &in
}

func TestUpdateAutoAccessOK(t *testing.T) {
	serverMock := newServerMock(t, "update_auto_access")
	client := NewTestClient(serverMock)

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	t.Run("only required fields", func(t *testing.T) {
		err := api.UpdateAutoAccess("foo", nil)
		if err != nil {
			t.Errorf("UpdateAutoAccess failed: %v", err)
		}
		validateRequestBody(t, `{"id":"foo"}`, serverMock.LastRequest.Body)
	})
	t.Run("only description", func(t *testing.T) {
		err := api.UpdateAutoAccess("foo", &configuration.UpdateAutoAccessRequestOptions{Description: addressOf("bar")})
		if err != nil {
			t.Errorf("UpdateAutoAccess failed: %v", err)
		}
		validateRequestBody(t, `{"id":"foo","description":"bar"}`, serverMock.LastRequest.Body)
	})
	t.Run("all optional fields", func(t *testing.T) {
		err := api.UpdateAutoAccess("foo", &configuration.UpdateAutoAccessRequestOptions{Description: addressOf("bar"), Access: &configuration.Access{Groups: []int{420}}, Conditions: &configuration.AutoAccessConditions{}, NextID: addressOf("baz")})
		if err != nil {
			t.Errorf("UpdateAutoAccess failed: %v", err)
		}
		validateRequestBody(t, `{"id":"foo","access":{"groups":[420]},"conditions":{},"description":"bar","next_id":"baz"}`, serverMock.LastRequest.Body)
	})
	t.Run("clearing next_id", func(t *testing.T) {
		err := api.UpdateAutoAccess("foo", &configuration.UpdateAutoAccessRequestOptions{NextID: addressOf("")})
		if err != nil {
			t.Errorf("UpdateAutoAccess failed: %v", err)
		}
		validateRequestBody(t, `{"id":"foo","next_id":""}`, serverMock.LastRequest.Body)
	})

	updateAutoAccessOKOpts := []*configuration.UpdateAutoAccessRequestOptions{
		{Description: addressOf("baz"), Access: &configuration.Access{}, Conditions: &configuration.AutoAccessConditions{}, NextID: addressOf("bar")},
		{Description: addressOf("baz")},
		nil,
	}

	for _, tt := range updateAutoAccessOKOpts {

		err = api.UpdateAutoAccess("foo", tt)
		if err != nil {
			t.Errorf("UpdateAutoAccess failed: %v", err)
		}
	}
}

func TestDeleteAutoAccessShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(newServerMock(t, "delete_auto_access"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	err = api.DeleteAutoAccess("foo")
	if err != nil {
		t.Errorf("DeleteAutoAccess failed: %v", err)
	}
}

func TestListAutoAccessesShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(newServerMock(t, "list_auto_accesses"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	resp, err := api.ListAutoAccesses()
	if err != nil {
		t.Errorf("ListAutoAccesses failed: %v", err)
	}

	if len(resp) != 1 {
		t.Errorf("Invalid response length: %v", len(resp))
	}

	if resp[0].ID != "1faad6f5f1d6e8fdf27e8af9839783b7" {
		t.Errorf("Invalid response id: %v", resp[0].ID)
	}

	if resp[0].NextID != "pqi8oasdjahuakndw9nsad9na" {
		t.Errorf("Invalid response next id: %v", resp[0].NextID)
	}

	if resp[0].Description != "Chats on livechat.com from United States" {
		t.Errorf("Invalid response description: %v", resp[0].ID)
	}

	if len(resp[0].Access.Groups) != 1 || resp[0].Access.Groups[0] != 0 {
		t.Errorf("Invalid response access groups: %v", resp[0].Access.Groups)
	}

	if len(resp[0].Conditions.Domain.Values) != 1 || !resp[0].Conditions.Domain.Values[0].ExactMatch || resp[0].Conditions.Domain.Values[0].Value != "livechat.com" {
		t.Errorf("Invalid response domain values: %v", resp[0].Conditions.Domain.Values)
	}

	if len(resp[0].Conditions.Geolocation.Values) != 1 || resp[0].Conditions.Geolocation.Values[0].Country != "United States" || resp[0].Conditions.Geolocation.Values[0].CountryCode != "US" {
		t.Errorf("Invalid response geolocation values: %v", resp[0].Conditions.Geolocation.Values)
	}
}

func TestCheckProductLimitsForPlanShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(newServerMock(t, "check_product_limits_for_plan"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	planLimits, rErr := api.CheckProductLimitsForPlan("starter")
	if rErr != nil {
		t.Errorf("CheckProductLimitsForPlan failed: %v", rErr)
	}

	if len(planLimits) != 3 {
		t.Errorf("Wrong number of limits: %v", len(planLimits))
	}

	if planLimits[0].Resource != "groups" {
		t.Errorf("Invalid limit resource: %v", planLimits[0].Resource)
	}

	if planLimits[1].LimitBalance != 2 {
		t.Errorf("Invalid limit balance: %v", planLimits[1].LimitBalance)
	}

	if planLimits[2].Id != "0" {
		t.Errorf("Invalid limit id: %v", planLimits[2].Id)
	}
}

func TestListChannelsShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(newServerMock(t, "list_channels"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	channels, rErr := api.ListChannels()
	if rErr != nil {
		t.Errorf("ListChannels failed: %v", rErr)
	}

	if len(channels) != 3 {
		t.Errorf("Wrong number of channels: %v", len(channels))
	}

	if channels[0].ChannelType != "code" {
		t.Errorf("Invalid channel type: %v", channels[0].ChannelType)
	}

	if channels[1].FirstActivityTimestamp != "2017-10-12T15:20:00Z" {
		t.Errorf("Invalid channel first activity timestamp: %v", channels[1].FirstActivityTimestamp)
	}

	if channels[2].ChannelSubtype != "c6e4f62e2a2dab12531235b12c5a2a6b" {
		t.Errorf("Invalid channel subtype: %v", channels[3].ChannelSubtype)
	}
}

func TestCreateTagShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(newServerMock(t, "create_tag"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	if rErr := api.CreateTag("tageroo", []int{0}); rErr != nil {
		t.Errorf("CreateTag failed: %v", rErr)
	}
}

func TestDeleteTagShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(newServerMock(t, "delete_tag"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	if rErr := api.DeleteTag("tageroo"); rErr != nil {
		t.Errorf("DeleteTag failed: %v", rErr)
	}
}

func TestListTagsShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(newServerMock(t, "list_tags"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	resp, err := api.ListTags([]int{})
	if err != nil {
		t.Errorf("ListTags failed: %v", err)
	}

	if len(resp) != 2 {
		t.Errorf("Invalid response length: %v", len(resp))
	}

	if resp[0].Name != "tageroo" {
		t.Errorf("Invalid response name: %v", resp[0].Name)
	}

	if len(resp[0].GroupIDs) != 1 || resp[0].GroupIDs[0] != 0 {
		t.Errorf("Invalid response group_ids: %v", resp[0].GroupIDs)
	}

	if resp[0].CreatedAt != "2017-10-12T13:56:16Z" {
		t.Errorf("Invalid response created_at: %v", resp[0].CreatedAt)
	}

	if resp[0].AuthorID != "smith@example.com" {
		t.Errorf("Invalid response author_id: %v", resp[0].AuthorID)
	}
}

func TestUpdateTagShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(newServerMock(t, "update_tag"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	if rErr := api.UpdateTag("tageroo", []int{0}); rErr != nil {
		t.Errorf("UpdateTag failed: %v", rErr)
	}
}

func TestListGroupsProperties(t *testing.T) {
	serverMock := newServerMock(t, "list_groups_properties")
	client := NewTestClient(serverMock)

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	checkAPIrespondedOK := func(t *testing.T, resp []configuration.GroupProperties, rErr error) {
		t.Helper()
		if rErr != nil {
			t.Errorf("ListGroupsProperties failed: %v", rErr)
		}
		if len(resp) != 1 {
			t.Errorf("Invalid response length: %v", len(resp))
		}
		if resp[0].ID != 0 {
			t.Errorf("Invalid response id: %v", resp[0].ID)
		}
		if resp[0].Properties["abc"]["a_property"] != "a" {
			t.Errorf("Invalid response: %v", resp[0].ID)
		}
	}

	t.Run("no options", func(t *testing.T) {
		resp, rErr := api.ListGroupsProperties([]int{0, 1}, nil)
		checkAPIrespondedOK(t, resp, rErr)
		validateRequestBody(t, `{"group_ids":[0,1]}`, serverMock.LastRequest.Body)
	})
	t.Run("with optional parameters", func(t *testing.T) {
		resp, rErr := api.ListGroupsProperties([]int{0, 1}, &configuration.ListGroupsPropertiesRequestOptions{
			Namespace:  "foo",
			NamePrefix: "bar",
		})
		checkAPIrespondedOK(t, resp, rErr)
		validateRequestBody(t, `{"group_ids":[0,1],"namespace":"foo","name_prefix":"bar"}`, serverMock.LastRequest.Body)
	})

}

func TestReactivateEmail(t *testing.T) {
	client := NewTestClient(newServerMock(t, "reactivate_email"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	if rErr := api.ReactivateEmail("agent_id"); rErr != nil {
		t.Errorf("ReactivateEmail failed: %v", rErr)
	}
}

func TestUpdateCompanyDetails(t *testing.T) {
	serverMock := newServerMock(t, "update_company_details")
	client := NewTestClient(serverMock)

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	company := "Text"
	url := ""
	if rErr := api.UpdateCompanyDetails(
		configuration.CompanyDetails{
			Company: &company,
			URL:     &url,
		},
		true,
	); rErr != nil {
		t.Errorf("UpdateCompanyDetails failed: %v", rErr)
	}

	validateRequestBody(t, `{"company":"Text","url":"","enrich":true}`, serverMock.LastRequest.Body)
}

func TestCreateCannedResponseShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	serverMock := newServerMock(t, "create_canned_response")
	client := NewTestClient(serverMock)

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	id, rErr := api.CreateCannedResponse("text", []string{"tag"}, 0, false)
	if rErr != nil {
		t.Errorf("CreateCannedResponse failed: %v", rErr)
	}
	if id != 2137 {
		t.Errorf("Invalid response id: %v", id)
	}

	validateRequestBody(t, `{"text":"text","tags":["tag"],"group_id":0,"is_private":false}`, serverMock.LastRequest.Body)
}

func TestDeleteCannedResponseShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	serverMock := newServerMock(t, "delete_canned_response")
	client := NewTestClient(serverMock)

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	if rErr := api.DeleteCannedResponse(2137); rErr != nil {
		t.Errorf("DeleteCannedResponse failed: %v", rErr)
	}

	validateRequestBody(t, `{"id":2137}`, serverMock.LastRequest.Body)
}

func TestListCannedResponsesShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	serverMock := newServerMock(t, "list_canned_responses")
	client := NewTestClient(serverMock)

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	limit := 2
	pageID := "some_page=="
	resp, rErr := api.ListCannedResponses(&configuration.ListCannedResponsesRequestOptions{
		GroupIDs:       []int32{0, 1},
		IncludePrivate: true,
		Limit:          &limit,
		PageID:         &pageID,
	})
	if rErr != nil {
		t.Errorf("ListCannedResponses failed: %v", rErr)
	}

	if len(resp.CannedResponses) != 2 {
		t.Errorf("Invalid response length: %v", len(resp.CannedResponses))
	}

	if resp.CannedResponses[0].ID != 1 {
		t.Errorf("Invalid response id: %v", resp.CannedResponses[0].ID)
	}

	if resp.CannedResponses[1].ID != 3 {
		t.Errorf("Invalid response id: %v", resp.CannedResponses[1].ID)
	}

	if resp.CannedResponses[0].Text != "text" {
		t.Errorf("Invalid response text: %v", resp.CannedResponses[0].Text)
	}

	validateRequestBody(t, `{"group_ids":[0,1],"include_private":true,"limit":2,"page_id":"some_page=="}`, serverMock.LastRequest.Body)
}

func TestAddAllowedDomainOK(t *testing.T) {
	serverMock := newServerMock(t, "add_allowed_domain")
	client := NewTestClient(serverMock)

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	err = api.AddAllowedDomain("example.com")
	if err != nil {
		t.Errorf("AddAllowedDomain failed: %v", err)
	}

	validateRequestBody(t, `{"domain":"example.com"}`, serverMock.LastRequest.Body)
}

func TestDeleteAllowedDomainOK(t *testing.T) {
	serverMock := newServerMock(t, "delete_allowed_domain")
	client := NewTestClient(serverMock)

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	err = api.DeleteAllowedDomain("example.com")
	if err != nil {
		t.Errorf("DeleteAllowedDomain failed: %v", err)
	}

	validateRequestBody(t, `{"domain":"example.com"}`, serverMock.LastRequest.Body)
}

func TestListAllowedDomainsShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(newServerMock(t, "list_allowed_domains"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	resp, err := api.ListAllowedDomains()
	if err != nil {
		t.Errorf("ListAllowedDomains failed: %v", err)
	}

	if len(resp) != 2 {
		t.Errorf("Invalid response length: %v", len(resp))
	}

	if resp[0] != "example.com" {
		t.Errorf("Invalid first domain: %v", resp[0])
	}

	if resp[1] != "test.com" {
		t.Errorf("Invalid second domain: %v", resp[1])
	}
}

func TestUpdateTranslationsOK(t *testing.T) {
	serverMock := newServerMock(t, "update_translations")
	client := NewTestClient(serverMock)

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	err = api.UpdateTranslations(1, "pl", map[string]string{
		"hello": "cześć",
	})
	if err != nil {
		t.Errorf("UpdateTranslations failed: %v", err)
	}

	validateRequestBody(t, `{"group_id":1,"language":"pl","phrases":{"hello":"cześć"}}`, serverMock.LastRequest.Body)
}

func TestListTranslationsShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	serverMock := newServerMock(t, "list_translations")
	client := NewTestClient(serverMock)

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	resp, err := api.ListTranslations(1, "pl")
	if err != nil {
		t.Errorf("ListTranslations failed: %v", err)
	}

	if len(resp) != 2 {
		t.Errorf("Invalid response length: %v", len(resp))
	}

	if resp["hello"] != "cześć" {
		t.Errorf("Invalid `hello` translation: %v", resp["hello"])
	}
}
