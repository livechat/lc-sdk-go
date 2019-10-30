package webhooks_validators

import (
	"fmt"

	"github.com/livechat/lc-sdk-go/webhooks"
)

func EventsMarkedAsSeen(licenseID int, payload interface{}) error {
	if licenseID != 100012582 {
		return fmt.Errorf("Invalid licenseID: %v", licenseID)
	}
	wh, ok := payload.(*webhooks.EventsMarkedAsSeen)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", payload)
	}

	var errors string
	PropEq("ChatID", wh.ChatID, "PJ0MRSHTDG", &errors)
	PropEq("UserID", wh.UserID, "b7eff798-f8df-4364-8059-649c35c9ed0c", &errors)
	PropEq("SeenUpTo", wh.SeenUpTo, "2017-10-12T15:19:21.010200Z", &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}
