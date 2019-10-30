package webhooks_validators

import (
	"fmt"

	"github.com/livechat/lc-sdk-go/webhooks"
)

func AgentStatusChanged(licenseID int, payload interface{}) error {
	if licenseID != 100012582 {
		return fmt.Errorf("Invalid licenseID: %v", licenseID)
	}
	wh, ok := payload.(*webhooks.AgentStatusChanged)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", payload)
	}

	var errors string
	PropEq("AgentID", wh.AgentID, "5c9871d5372c824cbf22d860a707a578", &errors)
	PropEq("Status", wh.Status, "accepting chats", &errors)
	
	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}
