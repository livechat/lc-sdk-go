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
	return &authorization.Token{
		LicenseID:   12345,
		AccessToken: "access_token",
		Region:      "region",
	}
}

var mockedResponses = map[string]string{
	"register_webhook": `{
		"webhook_id": "pqi8oasdjahuakndw9nsad9na"
	}`,
	"get_webhooks_config": `[
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
	"create_bot_agent": `{
    "bot_agent_id": "5c9871d5372c824cbf22d860a707a578"
	}`,
	"update_bot_agent": `{}`,
	"remove_bot_agent": `{}`,
	"get_bot_agents": `{
    "bot_agents": [{
        "id": "5c9871d5372c824cbf22d860a707a578",
        "name": "John Doe",
        "avatar": "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg",
        "status": "accepting chats"
    }]
	}`,
	"get_bot_agent_details": `{
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

func TestGetWebhooksConfigShouldReturnDataReceivedFromConfApi(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "get_webhooks_config"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	resp, rErr := api.GetWebhooksConfig()
	if rErr != nil {
		t.Errorf("GetWebhooksConfig failed: %v", rErr)
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
	client := NewTestClient(createMockedResponder(t, "create_bot_agent"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	botID, rErr := api.CreateBotAgent("John Doe", "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg", "accepting chats", 6, "first", []*configuration.BotGroupConfig{}, &configuration.BotWebhooks{})
	if rErr != nil {
		t.Errorf("CreateBotAgent failed: %v", rErr)
	}

	if botID != "5c9871d5372c824cbf22d860a707a578" {
		t.Errorf("Invalid botID: %v", botID)
	}
}

func TestCreateBotAgentShouldReturnErrorForInvalidInput(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "create_bot_agent"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	groups := []*configuration.BotGroupConfig{&configuration.BotGroupConfig{Priority: "supervisor"}}
	_, rErr := api.CreateBotAgent("John Doe", "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg", "accepting chats", 6, "first", groups, &configuration.BotWebhooks{})
	if rErr.Error() != "DoNotAssign priority is allowed only as default group priority" {
		t.Errorf("CreateBotAgent failed: %v", rErr)
	}
}

func TestUpdateBotAgentShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "update_bot_agent"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.UpdateBotAgent("pqi8oasdjahuakndw9nsad9na", "John Doe", "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg", "accepting chats", 6, "first", []*configuration.BotGroupConfig{}, &configuration.BotWebhooks{})
	if rErr != nil {
		t.Errorf("UpdateBotAgent failed: %v", rErr)
	}
}

func TestUpdateBotAgentShouldReturnErrorForInvalidInput(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "update_bot_agent"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	groups := []*configuration.BotGroupConfig{&configuration.BotGroupConfig{Priority: "supervisor"}}
	rErr := api.UpdateBotAgent("pqi8oasdjahuakndw9nsad9na", "John Doe", "livechat.s3.amazonaws.com/1011121/all/avatars/bdd8924fcbcdbddbeaf60c19b238b0b0.jpg", "accepting chats", 6, "first", groups, &configuration.BotWebhooks{})
	if rErr.Error() != "DoNotAssign priority is allowed only as default group priority" {
		t.Errorf("CreateBotAgent failed: %v", rErr)
	}
}

func TestRemoveBotAgentShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "remove_bot_agent"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	rErr := api.RemoveBotAgent("pqi8oasdjahuakndw9nsad9na")
	if rErr != nil {
		t.Errorf("RemoveBotAgent failed: %v", rErr)
	}
}

func TestGetBotAgentsShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "get_bot_agents"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	resp, rErr := api.GetBotAgents(true)
	if rErr != nil {
		t.Errorf("GetBotAgents failed: %v", rErr)
	}

	if len(resp) != 1 || resp[0].ID != "5c9871d5372c824cbf22d860a707a578" {
		t.Errorf("Invalid bot agents: %v", resp)
	}
}

func TestGetBotAgentDetailsShouldReturnDataReceivedFromAgentAPI(t *testing.T) {
	client := NewTestClient(createMockedResponder(t, "get_bot_agent_details"))

	api, err := configuration.NewAPI(stubTokenGetter, client, "client_id")
	if err != nil {
		t.Errorf("API creation failed")
	}

	resp, rErr := api.GetBotAgentDetails("5c9871d5372c824cbf22d860a707a578")
	if rErr != nil {
		t.Errorf("GetBotAgentDetails failed: %v", rErr)
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
