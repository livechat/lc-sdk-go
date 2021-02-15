package configuration_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/livechat/lc-sdk-go/v3/authorization"
	"github.com/livechat/lc-sdk-go/v3/configuration"
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
	licenseID := 12345
	return &authorization.Token{
		LicenseID:   &licenseID,
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
				"chat_member_ids": {
					"agents_any": ["johndoe@mail.com"]
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
		}],
		"webhooks": {
			"url": "http://myservice.com/webhooks",
			"secret_key": "JSauw0Aks8l-asAa",
			"actions": [{
				"name": "incoming_chats",
				"filters": {
					"chat_properties": {
						"source": {
							"type": {
								"values": ["facebook", "twitter"]
							}
						}
					}
				}
			},{
				"name": "incoming_event",
				"additional_data": ["chat_properties"]
			}]
		}
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
			]
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
		"name": "Agent Smith",
		"avatar_path": "https://domain.com/avatar.image.jpg",
		"role": "administrator",
		"login_status": "accepting chats"
	}`,
	"list_agents": `[
		{
			"id": "smith@example.com",
			"job_title": "Support Hero",
			"max_chats_count": 5,
			"summaries": [
				"daily_summary",
				"weekly_summary"
			]
		},
		{
			"id": "adam@example.com",
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
				"chat_properties"
			]
		},
		{
			"action": "event_properties_deleted",
			"filters": [
				"chat_member_ids",
				"only_my_chats"
			],
			"additional_data": [
				"chat_properties"
			]
		}
	]`,
	"enable_license_webhooks":  `{}`,
	"disable_license_webhooks": `{}`,
	"get_license_webhooks_state": `{
		"license_webhooks_enabled": true
	}`,
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

		if req.URL.String() != "https://api.livechatinc.com/v3.3/configuration/action/"+method {
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
			Body:       ioutil.NopCloser(bytes.NewBufferString(mockedResponses[method])),
			Header:     make(http.Header),
		}
	}
}

func TestRejectAPICreationWithoutTokenGetter(t *testing.T) {
	_, err := configuration.NewAPI(nil, nil, "client_id")
	if err == nil {
		t.Errorf("API should not be created without token getter")
	}
}

func TestRegisterWebhookShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "register_webhook"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	webhookID, rErr := api.RegisterWebhook(&configuration.Webhook{}, nil)
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
		t.Errorf("API creation failed")
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
		t.Errorf("API creation failed")
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
		t.Errorf("API creation failed")
	}

	botID, rErr := api.CreateBot("John Doe", "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg", 6, "first", []*configuration.GroupConfig{}, &configuration.BotWebhooks{})
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
		t.Errorf("API creation failed")
	}

	groups := []*configuration.GroupConfig{&configuration.GroupConfig{Priority: "supervisor"}}
	_, rErr := api.CreateBot("John Doe", "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg", 6, "first", groups, &configuration.BotWebhooks{})
	if rErr.Error() != "DoNotAssign priority is allowed only as default group priority" {
		t.Errorf("CreateBot failed: %v", rErr)
	}
}

func TestUpdateBotShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "update_bot"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.UpdateBot("pqi8oasdjahuakndw9nsad9na", "John Doe", "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg", 6, "first", []*configuration.GroupConfig{}, &configuration.BotWebhooks{})
	if rErr != nil {
		t.Errorf("UpdateBot failed: %v", rErr)
	}
}

func TestUpdateBotShouldReturnErrorForInvalidInput(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "update_bot"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	groups := []*configuration.GroupConfig{&configuration.GroupConfig{Priority: "supervisor"}}
	rErr := api.UpdateBot("pqi8oasdjahuakndw9nsad9na", "John Doe", "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg", 6, "first", groups, &configuration.BotWebhooks{})
	if rErr.Error() != "DoNotAssign priority is allowed only as default group priority" {
		t.Errorf("CreateBot failed: %v", rErr)
	}
}

func TestDeleteBotShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "delete_bot"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
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
		t.Errorf("API creation failed")
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
		t.Errorf("API creation failed")
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
		t.Errorf("API creation failed")
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
	})
	if rErr != nil {
		t.Errorf("RegisterProperty failed: %v", rErr)
	}
}

func TestUnregisterPropertyShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "unregister_property"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
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
		t.Errorf("API creation failed")
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
		t.Errorf("API creation failed")
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

func TestListGroupPropertiesShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "list_group_properties"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
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

func TestCreateAgentShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "create_agent"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
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
		t.Errorf("API creation failed")
	}

	agent, rErr := api.GetAgent("smith@example.com", []string{})
	if rErr != nil {
		t.Errorf("CreateAgent failed: %v", rErr)
	}

	if agent.ID != "smith@example.com" {
		t.Errorf("Invalid agent ID: %v", agent.ID)
	}

	if agent.Name != "Agent Smith" {
		t.Errorf("Invalid agent name: %v", agent.Name)
	}
}

func TestListAgentsShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "list_agents"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	agents, rErr := api.ListAgents([]int32{0, 1}, []string{})
	if rErr != nil {
		t.Errorf("CreateAgent failed: %v", rErr)
	}

	if len(agents) != 2 {
		t.Errorf("Invalid number of agents: %v", len(agents))
	}

	if agents[0].ID != "smith@example.com" {
		t.Errorf("Invalid agent ID: %v", agents[0].ID)
	}

	if agents[1].ID != "adam@example.com" {
		t.Errorf("Invalid agent name: %v", agents[1].ID)
	}
}

func TestUpdateAgentShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "update_agent"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
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
		t.Errorf("API creation failed")
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
		t.Errorf("API creation failed")
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
		t.Errorf("API creation failed")
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
		t.Errorf("API creation failed")
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
		t.Errorf("API creation failed")
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
		t.Errorf("API creation failed")
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
		t.Errorf("API creation failed")
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
		t.Errorf("API creation failed")
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
		t.Errorf("API creation failed")
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
		t.Errorf("API creation failed")
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
		t.Errorf("API creation failed")
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
		t.Errorf("API creation failed")
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
		t.Errorf("API creation failed")
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
		t.Errorf("API creation failed")
	}

	state, err := api.GetLicenseWebhooksState(nil)
	if err != nil {
		t.Errorf("GetWebhooksState failed: %v", err)
	}
	if !state.Enabled {
		t.Errorf("webhooks' state should be enabled'")
	}
}
