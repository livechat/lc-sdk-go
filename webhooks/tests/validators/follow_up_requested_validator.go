package webhooks_validators

import (
	"fmt"

	"github.com/livechat/lc-sdk-go/webhooks"
)

func FollowUpRequested(licenseID int, payload interface{}) error {
	if licenseID != 100012582 {
		return fmt.Errorf("Invalid licenseID: %v", licenseID)
	}
	wh, ok := payload.(*webhooks.FollowUpRequested)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", payload)
	}

	var errors string
	PropEq("ChatID", wh.ChatID, "XXXX", &errors)
	PropEq("ThreadID", wh.ThreadID, "YYYY", &errors)
	PropEq("CustomerID", wh.CustomerID, "AAA-BBB-CCC", &errors)
	
	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}
