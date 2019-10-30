package webhooks_validators

import (
	"fmt"

	"github.com/livechat/lc-sdk-go/webhooks"
)

func ChatThreadUntagged(licenseID int, payload interface{}) error {
	if licenseID != 100012582 {
		return fmt.Errorf("Invalid licenseID: %v", licenseID)
	}
	wh, ok := payload.(*webhooks.ChatThreadUntagged)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", payload)
	}

	var errors string
	PropEq("ChatID", wh.ChatID, "PJ0MRSHTDG", &errors)
	PropEq("ThreadID", wh.ThreadID, "K600PKZON8", &errors)
	PropEq("Tag", wh.Tag, "bug_report", &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}
