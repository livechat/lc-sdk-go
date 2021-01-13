package configuration_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/livechat/lc-sdk-go/v2/authorization"
	"github.com/livechat/lc-sdk-go/v2/configuration"
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
		"webhook_id": "pqi8oasdjahuakndw9nsad9na"
	}`,
	"list_registered_webhooks": `[
    {
      "webhook_id": "pqi8oasdjahuakndw9nsad9na",
      "url": "http://myservice.com/webhooks",
      "description": "Test webhook",
      "action": "thread_closed",
      "secret_key": "laudla991lamda0pnoaa0",
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
    "bot_agent_id": "5c9871d5372c824cbf22d860a707a578"
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
	"register_properties": `{
		"58737b5829e65621a45d598aa6f2ed8e": {
			"greeting": {
				"type": "string",
				"locations": {
					"chat": {
						"access": {
							"agent": {
								"read": true,
								"write": false
							},
							"customer": {
								"read": true,
								"write": true
							}
						}
					}
				},
				"domain": [
					"hello",
					"hi"
				]
			},
			"scoring": {
				"type": "int",
				"locations": {
					"event": {
						"access": {
							"agent": {
								"read": true,
								"write": true
							}
						}
					}
				},
				"range": {
					"from": 0,
					"to": 10
				}
			}
		}
	}`,
	"list_registered_properties": `{
		"58737b5829e65621a45d598aa6f2ed8e": {
			"greeting": {
				"type": "string",
				"locations": {
					"chat": {
						"access": {
							"agent": {
								"read": true,
								"write": false
							},
							"customer": {
								"read": true,
								"write": true
							}
						}
					}
				},
				"domain": [
					"hello",
					"hi"
				]
			},
			"scoring": {
				"type": "int",
				"locations": {
					"event": {
						"access": {
							"agent": {
								"read": true,
								"write": true
							}
						}
					}
				},
				"range": {
					"from": 0,
					"to": 10
				}
			}
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

		if req.URL.String() != "https://api.livechatinc.com/v3.2/configuration/action/"+method {
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

	webhookID, rErr := api.RegisterWebhook(&configuration.Webhook{})
	if rErr != nil {
		t.Errorf("RegisterWebhook failed: %v", rErr)
	}

	if webhookID != "pqi8oasdjahuakndw9nsad9na" {
		t.Errorf("Invalid webhookID: %v", webhookID)
	}
}

func TestListRegisteredWebhooksShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "list_registered_webhooks"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	resp, rErr := api.ListRegisteredWebhooks()
	if rErr != nil {
		t.Errorf("ListRegisteredWebhooks failed: %v", rErr)
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

	rErr := api.UnregisterWebhook("pqi8oasdjahuakndw9nsad9na")
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

	botID, rErr := api.CreateBot("John Doe", "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg", "accepting chats", 6, "first", []*configuration.GroupConfig{}, &configuration.BotWebhooks{})
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
	_, rErr := api.CreateBot("John Doe", "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg", "accepting chats", 6, "first", groups, &configuration.BotWebhooks{})
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

	rErr := api.UpdateBot("pqi8oasdjahuakndw9nsad9na", "John Doe", "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg", "accepting chats", 6, "first", []*configuration.GroupConfig{}, &configuration.BotWebhooks{})
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
	rErr := api.UpdateBot("pqi8oasdjahuakndw9nsad9na", "John Doe", "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg", "accepting chats", 6, "first", groups, &configuration.BotWebhooks{})
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

func TestRegisterPropertiesShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "register_properties"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.RegisterProperties(map[string]*configuration.PropertyConfig{"foo": &configuration.PropertyConfig{Type: "string"}})
	if rErr != nil {
		t.Errorf("RegisterProperties failed: %v", rErr)
	}
}

func TestListRegisteredPropertiesShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "list_registered_properties"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	resp, rErr := api.ListRegisteredProperties(true)
	if rErr != nil {
		t.Errorf("ListRegisteredProperties failed: %v", rErr)
	}

	if _, exists := resp["58737b5829e65621a45d598aa6f2ed8e"]; !exists || len(resp) != 1 {
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
