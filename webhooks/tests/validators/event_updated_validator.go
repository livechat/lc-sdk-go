package webhooks_validators

import (
	"fmt"

	"github.com/livechat/lc-sdk-go/webhooks"
)

func EventUpdated(licenseID int, payload interface{}) error {
	if licenseID != 100012582 {
		return fmt.Errorf("Invalid licenseID: %v", licenseID)
	}
	wh, ok := payload.(*webhooks.EventUpdated)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", payload)
	}

	var errors string
	PropEq("ChatID", wh.ChatID, "123-123-123-123", &errors)
	PropEq("ThreadID", wh.ThreadID, "E2WDHA8A", &errors)
	
	e := wh.Event.Message()
	PropEq("Event.ID", e.ID, "PZ070E0W1B_3", &errors)
	PropEq("Event.Type", e.Type, "message", &errors)
	PropEq("Event.Text", e.Text, "14", &errors)
	PropEq("Event.CustomID", e.CustomID, "1dnepb4z00t", &errors)
	PropEq("Event.Recipients", e.Recipients, "all", &errors)
	PropEq("Event.CreatedAt", e.CreatedAt.String(), "2019-10-11 09:41:00.877 +0000 UTC", &errors)
	PropEq("Event.AuthorID", e.AuthorID, "345f8235-d60d-433e-63c5-7f813a6ffe25", &errors)
	
	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}
