package webhooks_validators

import (
	"fmt"

	"github.com/livechat/lc-sdk-go/webhooks"
)

func ChatUserAdded(licenseID int, payload interface{}) error {
	if licenseID != 100012582 {
		return fmt.Errorf("Invalid licenseID: %v", licenseID)
	}
	wh, ok := payload.(*webhooks.ChatUserAdded)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", payload)
	}

	var errors string
	PropEq("ChatID", wh.ChatID, "PJ0MRSHTDG", &errors)
	PropEq("UserType", wh.UserType, "agent", &errors)

	user := wh.User
	PropEq("User.ID", user.ID, "l.wojciechowski@livechatinc.com", &errors)
	PropEq("User.Type", user.Type, "agent", &errors)
	PropEq("User.Name", user.Name, "Łukasz Wojciechowski", &errors)
	PropEq("User.Email", user.Email, "l.wojciechowski@livechatinc.com", &errors)
	PropEq("User.Avatar", user.Avatar, "livechat.s3.amazonaws.com/default/avatars/a14.png", &errors)
	PropEq("User.Present", user.Present, true, &errors)
	PropEq("User.LastSeen", user.LastSeen.String(), "1970-01-01 01:00:00 +0100 CET", &errors)
	//	PropEq("User.RoutingStatus", user.RoutingStatus, "accepting_chats", &errors) ERROR - Incorrect type User instead of Customer or Agent, some data might be missing

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}
