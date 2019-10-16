package customer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/livechat/lc-sdk-go/errors"
	"github.com/livechat/lc-sdk-go/internal/events"
)

const apiVersion = "v3.1"

type API struct {
	httpClient  *http.Client
	ApiURL      string
	tokenGetter func() *Token
}
type Token struct {
	License     string
	ClientID    string
	AccessToken string
	Region      string
}

type TokenGetter func() *Token

func NewAPI(t TokenGetter, client *http.Client) *API {
	if client == nil {
		client = &http.Client{
			Timeout: 20 * time.Second,
		}
	}

	return &API{
		tokenGetter: t,
		ApiURL:      "https://api.livechatinc.com/",
		httpClient:  client,
	}
}

type continuousChat struct {
	*Chat
	Continuous bool `json:"continuous"`
}

func (a *API) StartChat(c *Chat, continuous bool) (chatID, threadID string, err error) {
	cc := continuousChat{c, continuous}

	var resp struct {
		ChatID   string `json:"chat_id"`
		ThreadID string `json:"thread_id"`
	}

	return resp.ChatID, resp.ThreadID, a.call("start_chat", cc, &resp)
}

func (a *API) SendMessage(chatID, text string, whisper bool) (eventID string, err error) {
	recipients := "all"
	if whisper {
		recipients = "agents"
	}

	e := events.Message{
		Event: events.Event{
			Type:       "message",
			Recipients: recipients,
		},
		Text: text,
	}

	return a.SendEvent(chatID, e)
}

func (a *API) SendSystemMessage(chatID, text, messageType string) (eventID string, err error) {
	e := events.SystemMessage{
		Event: events.Event{
			Type: "system_message",
		},
		Text: text,
		Type: messageType,
	}

	return a.SendEvent(chatID, e)
}

func (a *API) SendEvent(chatID string, e interface{}) (eventID string, err error) {
	switch v := e.(type) {
	case events.Event:
	case events.Message:
	case events.SystemMessage:
	default:
		return "", fmt.Errorf("event type %s not supported", v)
	}

	var resp struct {
		EventID string `json:"event_id"`
	}
	err = a.call("send_event", map[string]interface{}{
		"chat_id": chatID,
		"event":   e,
	}, &resp)

	return resp.EventID, err
}

func (a *API) ActivateChat(chatID string) (threadID string, err error) {
	payload := map[string]interface{}{
		"chat": map[string]string{
			"id": chatID,
		},
	}

	var resp struct {
		ThreadID string `json:"thread_id"`
	}

	return resp.ThreadID, a.call("activate_chat", payload, &resp)
}
func (a *API) call(action string, request interface{}, response interface{}) error {
	rawBody, err := json.Marshal(request)
	if err != nil {
		return err
	}
	token := a.tokenGetter()

	url := fmt.Sprintf("%s/%s/customer/action/%s?license_id=%v", a.ApiURL, apiVersion, action, token.License)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(rawBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Set("User-agent", fmt.Sprintf("GO SDK Application %s", token.ClientID))
	req.Header.Set("X-Region", token.Region)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		apiErr := &errors.ErrAPI{}
		if err := json.Unmarshal(bodyBytes, apiErr); err != nil {
			return fmt.Errorf("couldn't unmarshal error response: %s (code: %d, raw body: %s)", err.Error(), resp.StatusCode, string(bodyBytes))
		}
		if apiErr.Error() == "" {
			return fmt.Errorf("couldn't unmarshal error response (code: %d, raw body: %s)", resp.StatusCode, string(bodyBytes))
		}
		return apiErr
	}

	if err != nil {
		return err
	}

	return json.Unmarshal(bodyBytes, response)
}
