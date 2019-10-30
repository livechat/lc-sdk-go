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
	PropEq("UserType", wh.UserType, "customer", &errors)

	customer := wh.User.Customer()
	PropEq("Customer.ID", customer.ID, "345f8235-d60d-433e-63c5-7f813a6ffe25", &errors)
	PropEq("Customer.Type", customer.Type, "customer", &errors)
	PropEq("Customer.Name", customer.Name, "test", &errors)
	PropEq("Customer.Email", customer.Email, "test@test.pl", &errors)
	PropEq("Customer.Avatar", customer.Avatar, "", &errors)
	PropEq("Customer.Present", customer.Present, true, &errors)
	PropEq("Customer.LastSeen", customer.LastSeen.String(), "2019-10-08 13:56:53 +0200 CEST", &errors)

	lastVisit := customer.LastVisit
	PropEq("LastVisit.IP", lastVisit.IP, "37.248.156.62", &errors)
	PropEq("LastVisit.UserAgent", lastVisit.UserAgent, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36", &errors)
	PropEq("LastVisit.StartedAt", lastVisit.StartedAt.String(), "2019-10-11 09:40:56.071345 +0000 UTC", &errors)

	geolocation := lastVisit.Geolocation
	PropEq("Geolocation.Country", geolocation.Country, "Poland", &errors)
	PropEq("Geolocation.CountryCode", geolocation.CountryCode, "PL", &errors)
	PropEq("Geolocation.Region", geolocation.Region, "test", &errors)
	PropEq("Geolocation.City", geolocation.City, "Wroclaw", &errors)
	PropEq("Geolocation.Timezone", geolocation.Timezone, "test_timezone", &errors)

	PropEq("LastPages", len(lastVisit.LastPages), 1, &errors)

	PropEq("LastPages.OpenedAt", lastVisit.LastPages[0].OpenedAt.String(), "2019-10-11 09:40:56.071345 +0000 UTC", &errors)
	PropEq("LastPages.URL", lastVisit.LastPages[0].URL, "https://cdn.livechatinc.com/labs/?license=100007977/", &errors)
	PropEq("LastPages.Title", lastVisit.LastPages[0].Title, "LiveChat", &errors)

	statistics := customer.Statistics
	PropEq("Statistics.VisistsCount", statistics.VisitsCount, 29, &errors)
	PropEq("Statistics.ThreadsCount", statistics.ThreadsCount, 18, &errors)
	PropEq("Statistics.ChatsCount", statistics.ChatsCount, 1, &errors)
	PropEq("Statistics.PageViewsCount", statistics.PageViewsCount, 5, &errors)
	PropEq("Statistics.GreetingsShownCount", statistics.GreetingsShownCount, 6, &errors)
	PropEq("Statistics.GreetingsAcceptedCount", statistics.GreetingsAcceptedCount, 8, &errors)

	PropEq("Customer.AgentLastEventCreatedAt", customer.AgentLastEventCreatedAt.String(), "2019-10-11 09:40:59.249 +0000 UTC", &errors)
	PropEq("Customer.CustomerLastEventCreatedAt", customer.CustomerLastEventCreatedAt.String(), "2019-10-11 09:40:59.219001 +0000 UTC", &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}
