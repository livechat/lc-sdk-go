package webhooks_validators

import (
	"fmt"

	"github.com/livechat/lc-sdk-go/webhooks"
)

func IncomingChatThread(licenseID int, payload interface{}) error {
	if licenseID != 100012582 {
		return fmt.Errorf("Invalid licenseID: %v", licenseID)
	}
	wh, ok := payload.(*webhooks.IncomingChatThread)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", payload)
	}


	chat := wh.Chat

	var errors string
	PropEq("Chat.ID", chat.ID, "PS0X0L086G", &errors)
	PropEq("Chat.Access.GroupIDs", len(chat.Access.GroupIDs), 1, &errors)
	PropEq("Chat.Access.GroupIDs[0]", chat.Access.GroupIDs[0], 0, &errors)
	PropEq("Chat.Users length", len(chat.Users()), 2, &errors)

	PropEq("Chat.Customers", len(chat.Customers), 1, &errors)
	cid := "345f8235-d60d-433e-63c5-7f813a6ffe25"
	customer := chat.Customers[cid]
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

	PropEq("Chat.Agents.length", len(chat.Agents), 1, &errors)
	aid := "l.wojciechowski@livechatinc.com"
	agent := chat.Agents[aid]
	PropEq("Agent.ID", agent.ID, "l.wojciechowski@livechatinc.com", &errors)
	PropEq("Agent.Type", agent.Type, "agent", &errors)
	PropEq("Agent.Name", agent.Name, "Łukasz Wojciechowski", &errors)
	PropEq("Agent.Email", agent.Email, "l.wojciechowski@livechatinc.com", &errors)
	PropEq("Agent.Avatar", agent.Avatar, "livechat.s3.amazonaws.com/default/avatars/a14.png", &errors)
	PropEq("Agent.Present", agent.Present, true, &errors)
	PropEq("Agent.LastSeen", agent.LastSeen.String(), "1970-01-01 01:00:00 +0100 CET", &errors)
	PropEq("Agent.RoutingStatus", agent.RoutingStatus, "accepting_chats", &errors)

	PropEq("Chat.Threads.length", len(chat.Threads), 1, &errors)
	thread := chat.Threads[0]
	PropEq("Thread.ID", thread.ID, "PZ070E0W1B", &errors)
	PropEq("Thread.Timestamp", thread.Timestamp.String(), "2019-10-11 11:40:59 +0200 CEST", &errors)
	PropEq("Thread.Active", thread.Active, true, &errors)
	PropEq("Thread.UserIDs[0]", thread.UserIDs[0], "345f8235-d60d-433e-63c5-7f813a6ffe25", &errors)
	PropEq("Thread.UserIDs[1]", thread.UserIDs[1], "l.wojciechowski@livechatinc.com", &errors)
	PropEq("Thread.RestrictedAccess", thread.RestrictedAccess, false, &errors)
	PropEq("Thread.Order", thread.Order, 18, &errors)
	PropEq("Thread.Properties.routing.continuous", thread.Properties["routing"]["continuous"], false, &errors)
	PropEq("Thread.Properties.routing.idle", thread.Properties["routing"]["idle"], false, &errors)
	PropEq("Thread.Properties.routing.referrer", thread.Properties["routing"]["referrer"], "", &errors)
	PropEq("Thread.Properties.routing.start_url", thread.Properties["routing"]["start_url"], "https://cdn.livechatinc.com/labs/?license=100007977/", &errors)
	PropEq("Thread.Properties.routing.unassigned", thread.Properties["routing"]["unassigned"], false, &errors)
	PropEq("Thread.Access.GroupIDs", thread.Access.GroupIDs[0], 0, &errors)
	PropEq("Thread.Events.length", len(thread.Events), 2, &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}