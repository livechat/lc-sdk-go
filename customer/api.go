package customer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/pkg/errors"
)

type API struct {
	clientID   string
	license    int32
	HTTPClient *http.Client
	ApiURL     string
	token      string
	tokenLock  sync.Mutex
}

func NewAPI(clientID string, license int32) *API {
	return &API{
		clientID:   clientID,
		license:    license,
		ApiURL:     "https://api.livechatinc.com/",
		HTTPClient: http.DefaultClient,
	}
}

func (a *API) SetToken(t string) {
	a.tokenLock.Lock()
	defer a.tokenLock.Unlock()

	a.token = t
}

type continuousChat struct {
	*Chat
	Continuous bool `json:"continuous"`
}

func (a *API) StartChat(c *Chat, continuous bool) (chatID, threadID string, err error) {
	if c.ID != "" {
		return "", "", fmt.Errorf("chat %s already started", c.ID)
	}

	cc := continuousChat{c, continuous}
	body, err := a.call("start_chat", cc)

	if err != nil {
		return "", "", err
	}

	resp := struct {
		ChatID   string `json:"chat_id"`
		ThreadID string `json:"thread_id"`
	}{}
	err = json.Unmarshal(body, &resp)

	if err != nil {
		return "", "", err
	}

	return resp.ChatID, resp.ThreadID, nil
}

func (a *API) SendMessage(chatID, text string, whisper bool) (eventID string, err error) {
	recipients := "all"
	if whisper {
		recipients = "agents"
	}
	payload := map[string]interface{}{
		"chat_id": chatID,
		"message:": map[string]string{
			"type":       "message",
			"text":       text,
			"recepients": recipients,
		},
	}

	body, err := a.call("send_event", payload)

	if err != nil {
		return "", err
	}

	resp := struct {
		EventID string `json:"event_id"`
	}{}
	err = json.Unmarshal(body, &resp)

	if err != nil {
		return "", err
	}

	return resp.EventID, err
}

func (a *API) SendSystemMessage() {

}

func (a *API) call(action string, payload interface{}) (json.RawMessage, error) {
	a.tokenLock.Lock()
	defer a.tokenLock.Unlock()

	rawBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/v3.1/customer/actions/%s?license_id=%v", a.ApiURL, action, a.license)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(rawBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.token))
	req.Header.Set("User-agent", fmt.Sprintf("GO SDK Application %s", a.clientID))

	resp, err := a.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("authorization error for token '%v'", a.token)
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: " + resp.Status + ", " + string(bodyBytes))
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return bodyBytes, nil
}
