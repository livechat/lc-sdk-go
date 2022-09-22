package configuration_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/livechat/lc-sdk-go/v5/authorization"
	"github.com/livechat/lc-sdk-go/v5/configuration"
)

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
		"id": "5c9871d5372c824cbf22d860a707a578"
	}`,
	"update_bot": `{}`,
	"delete_bot": `{}`,
	"list_bots": `[
		{
			"id": "5c9871d5372c824cbf22d860a707a578",
			"name": "John Doe",
			"avatar_path": "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg"
		},
		{
			"id": "8g1231ss112c013cbf11d530b595h987",
			"name": "Jason Brown",
			"avatar_path": "livechat.s3.amazonaws.com/1011121/all/avatars/wff9482gkdjanzjgdsf88a184jsskaz1.jpg"
		}
	]`,
	"get_bot": `{
		"id": "5c9871d5372c824cbf22d860a707a578",
		"name": "John Doe",
		"avatar_path": "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg",
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
	"create_agent": `{
		"id": "smith@example.com"
	}`,
	"get_agent": `{
		"id": "smith@example.com",
		"account_id": "d24fa41e-bc16-41b8-a15b-9ca45ff7e0cf",
		"name": "Agent Smith",
		"avatar_path": "https://domain.com/avatar.image.jpg",
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
	"create_group": `{
		"id": 19
	}`,
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
	"add_auto_access":           `{ "id": "pqi8oasdjahuakndw9nsad9na" }`,
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

		if req.URL.String() != "https://api.livechatinc.com/v3.5/configuration/action/"+method {
			t.Errorf("Invalid URL for Configuration API request: %s", req.URL.String())
			return createServerError("Invalid URL")
		}

		if req.Method != "POST" {
			t.Errorf("Invalid method: %s for Configuration API action: %s", req.Method, method)
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

func TestRejectAPICreationWithoutTokenGetter(t *testing.T) {
	_, err := configuration.NewAPI(nil, nil, "client_id")
	if err == nil {
		t.Error("API should not be created without token getter")
	}
}

func TestRegisterWebhookShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "register_webhook"))

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
	client := NewTestClient(createMockedResponder(t, "list_webhooks"))

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
	client := NewTestClient(createMockedResponder(t, "unregister_webhook"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.UnregisterWebhook("pqi8oasdjahuakndw9nsad9na", nil)
	if rErr != nil {
		t.Errorf("UnregisterWebhook failed: %v", rErr)
	}
}

func TestCreateBotShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "create_bot"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	botID, rErr := api.CreateBot("John Doe", "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg", 6, "first", []*configuration.GroupConfig{}, "dummy_client_id", "dummy/timezone", configuration.WorkScheduler{})
	if rErr != nil {
		t.Errorf("CreateBot failed: %v", rErr)
	}

	if botID != "5c9871d5372c824cbf22d860a707a578" {
		t.Errorf("Invalid botID: %v", botID)
	}
}

func TestCreateBotShouldReturnErrorForInvalidInput(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "create_bot"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	groups := []*configuration.GroupConfig{{Priority: "supervisor"}}
	_, rErr := api.CreateBot("John Doe", "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg", 6, "first", groups, "dummy_client_id", "dummy/timezone", configuration.WorkScheduler{})
	if rErr.Error() != "DoNotAssign priority is allowed only as default group priority" {
		t.Errorf("CreateBot failed: %v", rErr)
	}
}

func TestUpdateBotShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "update_bot"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.UpdateBot("pqi8oasdjahuakndw9nsad9na", "John Doe", "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg", 6, "first", []*configuration.GroupConfig{}, "dummy/timezone", configuration.WorkScheduler{})
	if rErr != nil {
		t.Errorf("UpdateBot failed: %v", rErr)
	}
}

func TestUpdateBotShouldReturnErrorForInvalidInput(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "update_bot"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	groups := []*configuration.GroupConfig{{Priority: "supervisor"}}
	rErr := api.UpdateBot("pqi8oasdjahuakndw9nsad9na", "John Doe", "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg", 6, "first", groups, "dummy/timezone", configuration.WorkScheduler{})
	if rErr.Error() != "DoNotAssign priority is allowed only as default group priority" {
		t.Errorf("CreateBot failed: %v", rErr)
	}
}

func TestDeleteBotShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "delete_bot"))

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
	client := NewTestClient(createMockedResponder(t, "list_bots"))

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
	client := NewTestClient(createMockedResponder(t, "get_bot"))

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

func TestRegisterPropertyShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "register_property"))

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
	client := NewTestClient(createMockedResponder(t, "unregister_property"))

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
	client := NewTestClient(createMockedResponder(t, "publish_property"))

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
	client := NewTestClient(createMockedResponder(t, "list_properties"))

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

func TestListLicensePropertiesShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "list_license_properties"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
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

func TestListGroupPropertiesShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "list_group_properties"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
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

func TestCreateAgentShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "create_agent"))

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
	client := NewTestClient(createMockedResponder(t, "get_agent"))

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
	client := NewTestClient(createMockedResponder(t, "list_agents"))

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
	client := NewTestClient(createMockedResponder(t, "update_agent"))

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
	client := NewTestClient(createMockedResponder(t, "delete_agent"))

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
	client := NewTestClient(createMockedResponder(t, "suspend_agent"))

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
	client := NewTestClient(createMockedResponder(t, "unsuspend_agent"))

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
	client := NewTestClient(createMockedResponder(t, "request_agent_unsuspension"))

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
	client := NewTestClient(createMockedResponder(t, "approve_agent"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.ApproveAgent("smith@example.com")
	if rErr != nil {
		t.Errorf("ApproveAgent failed: %v", rErr)
	}
}

func TestCreateGroupShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "create_group"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	groupID, rErr := api.CreateGroup("name", "en", map[string]configuration.GroupPriority{})
	if rErr != nil {
		t.Errorf("GetGroup failed: %v", rErr)
	}

	if groupID != 19 {
		t.Errorf("Invalid group id: %v", groupID)
	}
}

func TestUpdateGroupShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "update_group"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	rErr := api.UpdateGroup(11, "name", "en", map[string]configuration.GroupPriority{})
	if rErr != nil {
		t.Errorf("UpdateGroup failed: %v", rErr)
	}
}

func TestDeleteGroupShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "delete_group"))

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
	client := NewTestClient(createMockedResponder(t, "list_groups"))

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

	if groups[1].ID != 19 {
		t.Errorf("Invalid group ID: %v", groups[1].ID)
	}
}

func TestGetGroupShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "get_group"))

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
	client := NewTestClient(createMockedResponder(t, "list_webhook_names"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	resp, rErr := api.ListWebhookNames("3.2")
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
	client := NewTestClient(createMockedResponder(t, "enable_license_webhooks"))

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
	client := NewTestClient(createMockedResponder(t, "disable_license_webhooks"))

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
	client := NewTestClient(createMockedResponder(t, "get_license_webhooks_state"))

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
	client := NewTestClient(createMockedResponder(t, "delete_license_properties"))

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
	client := NewTestClient(createMockedResponder(t, "delete_group_properties"))

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
	client := NewTestClient(createMockedResponder(t, "update_license_properties"))

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
	client := NewTestClient(createMockedResponder(t, "update_group_properties"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	err = api.UpdateGroupProperties(0, nil)
	if err != nil {
		t.Errorf("DeleteGroupProperties failed: %v", err)
	}
}

func TestAddAutoAccessShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "add_auto_access"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	resp, err := api.AddAutoAccess([]int{}, nil, nil, nil, "", "")
	if err != nil {
		t.Errorf("AddAutoAccess failed: %v", err)
	}

	if resp != "pqi8oasdjahuakndw9nsad9na" {
		t.Errorf("Invalid response: %v", resp)
	}
}

func TestUpdateAutoAccessShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "update_auto_access"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	err = api.UpdateAutoAccess("foo", []int{}, nil, nil, nil, "", "")
	if err != nil {
		t.Errorf("UpdateAutoAccess failed: %v", err)
	}
}

func TestDeleteAutoAccessShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "delete_auto_access"))

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
	client := NewTestClient(createMockedResponder(t, "list_auto_accesses"))

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
	client := NewTestClient(createMockedResponder(t, "check_product_limits_for_plan"))

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
	client := NewTestClient(createMockedResponder(t, "list_channels"))

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
	client := NewTestClient(createMockedResponder(t, "create_tag"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	if rErr := api.CreateTag("tageroo", []int{0}); rErr != nil {
		t.Errorf("CreateTag failed: %v", rErr)
	}
}

func TestDeleteTagShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "delete_tag"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	if rErr := api.DeleteTag("tageroo"); rErr != nil {
		t.Errorf("DeleteTag failed: %v", rErr)
	}
}

func TestListTagsShouldReturnDataReceivedFromConfAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "list_tags"))

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
	client := NewTestClient(createMockedResponder(t, "update_tag"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	if rErr := api.UpdateTag("tageroo", []int{0}); rErr != nil {
		t.Errorf("UpdateTag failed: %v", rErr)
	}
}

func TestListGroupsProperties(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "list_groups_properties"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Error("API creation failed")
	}

	resp, rErr := api.ListGroupsProperties("namespace", "prefix,", []int{0, 1})

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
