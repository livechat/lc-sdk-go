package configuration_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/livechat/lc-sdk-go/authorization"
	"github.com/livechat/lc-sdk-go/configuration"
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
	"remove_bot": `{}`,
	"list_bots": `{
    "bot_agents": [{
        "id": "5c9871d5372c824cbf22d860a707a578",
        "name": "John Doe",
        "avatar": "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg",
        "status": "accepting chats"
    }]
	}`,
	"get_bot": `{
    "bot_agent": {
        "id": "5c9871d5372c824cbf22d860a707a578",
        "name": "John Doe",
        "avatar": "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg",
        "status": "accepting chats",
        "application": {
            "client_id": "asXdesldiAJSq9padj"
        },
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
                "name": "incoming_chat_thread",
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
    }
	}`,
	"create_properties": `{
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
	"get_property_configs": `{
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

		if req.URL.String() != "https://api.livechatinc.com/v3.2/configuration/action/"+method+"?license_id=12345" {
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

func TestUnregisterWebhookShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
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

func TestCreateBotAgentShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "create_bot"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	botID, rErr := api.CreateBot("John Doe", "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg", "accepting chats", 6, "first", []*configuration.BotGroupConfig{}, &configuration.BotWebhooks{})
	if rErr != nil {
		t.Errorf("CreateBot failed: %v", rErr)
	}

	if botID != "5c9871d5372c824cbf22d860a707a578" {
		t.Errorf("Invalid botID: %v", botID)
	}
}

func TestCreateBotAgentShouldReturnErrorForInvalidInput(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "create_bot"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	groups := []*configuration.BotGroupConfig{&configuration.BotGroupConfig{Priority: "supervisor"}}
	_, rErr := api.CreateBot("John Doe", "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg", "accepting chats", 6, "first", groups, &configuration.BotWebhooks{})
	if rErr.Error() != "DoNotAssign priority is allowed only as default group priority" {
		t.Errorf("CreateBot failed: %v", rErr)
	}
}

func TestUpdateBotAgentShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "update_bot"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.UpdateBot("pqi8oasdjahuakndw9nsad9na", "John Doe", "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg", "accepting chats", 6, "first", []*configuration.BotGroupConfig{}, &configuration.BotWebhooks{})
	if rErr != nil {
		t.Errorf("UpdateBot failed: %v", rErr)
	}
}

func TestUpdateBotAgentShouldReturnErrorForInvalidInput(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "update_bot"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	groups := []*configuration.BotGroupConfig{&configuration.BotGroupConfig{Priority: "supervisor"}}
	rErr := api.UpdateBot("pqi8oasdjahuakndw9nsad9na", "John Doe", "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg", "accepting chats", 6, "first", groups, &configuration.BotWebhooks{})
	if rErr.Error() != "DoNotAssign priority is allowed only as default group priority" {
		t.Errorf("CreateBot failed: %v", rErr)
	}
}

func TestRemoveBotAgentShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "remove_bot"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.RemoveBot("pqi8oasdjahuakndw9nsad9na")
	if rErr != nil {
		t.Errorf("RemoveBot failed: %v", rErr)
	}
}

func TestGetBotAgentsShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "list_bots"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	resp, rErr := api.ListBots(true)
	if rErr != nil {
		t.Errorf("ListBots failed: %v", rErr)
	}

	if len(resp) != 1 || resp[0].ID != "5c9871d5372c824cbf22d860a707a578" {
		t.Errorf("Invalid bot agents: %v", resp)
	}
}

func TestGetBotAgentDetailsShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "get_bot"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	resp, rErr := api.GetBot("5c9871d5372c824cbf22d860a707a578")
	if rErr != nil {
		t.Errorf("GetBot failed: %v", rErr)
	}

	if resp.ID != "5c9871d5372c824cbf22d860a707a578" {
		t.Errorf("Invalid bot agents: %v", resp)
	}
}

func TestCreatePropertiesShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "create_properties"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.CreateProperties(map[string]*configuration.PropertyConfig{"foo": &configuration.PropertyConfig{Type: "string"}})
	if rErr != nil {
		t.Errorf("CreateProperties failed: %v", rErr)
	}
}

func TestGetPropertyConfigsShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "get_property_configs"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	resp, rErr := api.GetPropertyConfigs(true)
	if rErr != nil {
		t.Errorf("GetPropertyConfigs failed: %v", rErr)
	}

	if _, exists := resp["58737b5829e65621a45d598aa6f2ed8e"]; !exists || len(resp) != 1 {
		t.Errorf("Invalid property configs: %v", resp)
	}
}

func TestGetGroupShouldReturnDataReceivedFromConfigurationAPI(t *testing.T) {
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

func TestListLicensePropertiesShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
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

func TestListGroupPropertiesShouldReturnDataReceivedFromCustomerAPI(t *testing.T) {
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
