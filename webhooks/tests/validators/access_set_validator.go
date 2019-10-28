package webhooks_validators

import (
	"fmt"

	"github.com/livechat/lc-sdk-go/webhooks"
)

func AccessSet(licenseID int, payload interface{}) error {
	if licenseID != 100012582 {
		return fmt.Errorf("Invalid licenseID: %v", licenseID)
	}
	wh, ok := payload.(*webhooks.AccessSet)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", payload)
	}

	var errors string
	PropEq("ID", wh.ID, "PJ0MRSHTDG", &errors)
	PropEq("Resource", wh.Resource, "chat", &errors)
	PropEq("Access.GroupIDs.length", len(wh.Access.GroupIDs), 1, &errors)
	PropEq("Access.GroupIDs[0]", wh.Access.GroupIDs[0], 1, &errors)
	
	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}
