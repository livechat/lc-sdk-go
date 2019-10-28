package webhooks_validators

import (
	"fmt"

	"github.com/livechat/lc-sdk-go/webhooks"
)

func ChatThreadPropertiesDeleted(licenseID int, payload interface{}) error {
	if licenseID != 100012582 {
		return fmt.Errorf("Invalid licenseID: %v", licenseID)
	}
	wh, ok := payload.(*webhooks.ChatThreadPropertiesDeleted)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", payload)
	}

	var errors string
	PropEq("ChatID", wh.ChatID, "PJ0MRSHTDG", &errors)
	PropEq("ThreadID", wh.ThreadID, "K600PKZON8", &errors)

	PropEq("Properties.Rating[0]", wh.Properties["rating"][0], "score", &errors)
	PropEq("Properties.Rating[1]", wh.Properties["rating"][1], "comment", &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}
