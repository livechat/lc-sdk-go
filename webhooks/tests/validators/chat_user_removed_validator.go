package webhooks_validators

import (
	"fmt"

	"github.com/livechat/lc-sdk-go/webhooks"
)

func ChatUserRemoved(licenseID int, payload interface{}) error {
	if licenseID != 100012582 {
		return fmt.Errorf("Invalid licenseID: %v", licenseID)
	}
	wh, ok := payload.(*webhooks.ChatUserRemoved)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", payload)
	}

	var errors string
	PropEq("ChatID", wh.ChatID, "PS0X0L086G", &errors)
	PropEq("UserType", wh.UserType, "agent", &errors)
	PropEq("UserID", wh.UserID, "agent@livechatinc.com", &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}