package webhooks_validators

import (
	"fmt"

	"github.com/livechat/lc-sdk-go/webhooks"
)

func ChatPropertiesUpdated(licenseID int, payload interface{}) error {
	if licenseID != 100012582 {
		return fmt.Errorf("Invalid licenseID: %v", licenseID)
	}
	wh, ok := payload.(*webhooks.ChatPropertiesUpdated)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", payload)
	}

	var errors string
	PropEq("ChatID", wh.ChatID, "PJ0MRSHTDG", &errors)

	PropEq("Properties.Rating.Score.Value", wh.Properties["rating"]["score"], float64(1), &errors)
	PropEq("Properties.Rating.Comment.Value", wh.Properties["rating"]["comment"], "Very good, veeeery good", &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}
