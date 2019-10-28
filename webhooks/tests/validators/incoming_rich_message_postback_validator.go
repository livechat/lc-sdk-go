package webhooks_validators

import (
	"fmt"

	"github.com/livechat/lc-sdk-go/webhooks"
)

func IncomingRichMessagePostback(licenseID int, payload interface{}) error {
	if licenseID != 100012582 {
		return fmt.Errorf("Invalid licenseID: %v", licenseID)
	}
	wh, ok := payload.(*webhooks.IncomingRichMessagePostback)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", payload)
	}

	var errors string
	PropEq("ChatID", wh.ChatID, "PJ0MRSHTDG", &errors)
	PropEq("ThreadID", wh.ThreadID, "K600PKZON8", &errors)
	PropEq("UserID", wh.UserID, "b7eff798-f8df-4364-8059-649c35c9ed0c", &errors)
	PropEq("EventID", wh.EventID, "a0c22fdd-fb71-40b5-bfc6-a8a0bc3117f7", &errors)
	PropEq("Postback.ID", wh.Postback.ID, "action_yes", &errors)
	PropEq("Postback.Toggled", wh.Postback.Toggled, true, &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}
