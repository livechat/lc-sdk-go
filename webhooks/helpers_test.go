package webhooks_test

import (
	"fmt"
	"time"

	"github.com/livechat/lc-sdk-go/v2/webhooks"
)

func propEq(propertyName string, actual, expected interface{}, validationAccumulator *string) {
	if actual != expected {
		*validationAccumulator += fmt.Sprintf("%s mismatch, actual: %v, expected: %v\n", propertyName, actual, expected)
	}
}

func incomingChat(wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.IncomingChat)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	chat := payload.Chat

	var errors string
	propEq("Chat.ID", chat.ID, "PS0X0L086G", &errors)
	propEq("Chat.Access.GroupIDs", len(chat.Access.GroupIDs), 1, &errors)
	propEq("Chat.Access.GroupIDs[0]", chat.Access.GroupIDs[0], 0, &errors)
	propEq("Chat.Users length", len(chat.Users()), 2, &errors)

	propEq("Chat.Customers", len(chat.Customers), 1, &errors)
	cid := "345f8235-d60d-433e-63c5-7f813a6ffe25"
	customer := chat.Customers[cid]
	propEq("Customer.ID", customer.ID, "345f8235-d60d-433e-63c5-7f813a6ffe25", &errors)
	propEq("Customer.Type", customer.Type, "customer", &errors)
	propEq("Customer.Name", customer.Name, "test", &errors)
	propEq("Customer.Email", customer.Email, "test@test.pl", &errors)
	propEq("Customer.Avatar", customer.Avatar, "", &errors)
	propEq("Customer.Present", customer.Present, true, &errors)
	propEq("Customer.EventsSeenUpTo", customer.EventsSeenUpTo.String(), "2019-10-08 13:56:53 +0000 UTC", &errors)

	lastVisit := customer.LastVisit
	propEq("LastVisit.IP", lastVisit.IP, "37.248.156.62", &errors)
	propEq("LastVisit.UserAgent", lastVisit.UserAgent, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36", &errors)
	propEq("LastVisit.StartedAt", lastVisit.StartedAt.String(), "2019-10-11 09:40:56.071345 +0000 UTC", &errors)

	geolocation := lastVisit.Geolocation
	propEq("Geolocation.Country", geolocation.Country, "Poland", &errors)
	propEq("Geolocation.CountryCode", geolocation.CountryCode, "PL", &errors)
	propEq("Geolocation.Region", geolocation.Region, "test", &errors)
	propEq("Geolocation.City", geolocation.City, "Wroclaw", &errors)
	propEq("Geolocation.Timezone", geolocation.Timezone, "test_timezone", &errors)

	propEq("LastPages", len(lastVisit.LastPages), 1, &errors)

	propEq("LastPages.OpenedAt", lastVisit.LastPages[0].OpenedAt.String(), "2019-10-11 09:40:56.071345 +0000 UTC", &errors)
	propEq("LastPages.URL", lastVisit.LastPages[0].URL, "https://cdn.livechatinc.com/labs/?license=100007977/", &errors)
	propEq("LastPages.Title", lastVisit.LastPages[0].Title, "LiveChat", &errors)

	statistics := customer.Statistics
	propEq("Statistics.VisistsCount", statistics.VisitsCount, 29, &errors)
	propEq("Statistics.ThreadsCount", statistics.ThreadsCount, 18, &errors)
	propEq("Statistics.ChatsCount", statistics.ChatsCount, 1, &errors)
	propEq("Statistics.PageViewsCount", statistics.PageViewsCount, 5, &errors)
	propEq("Statistics.GreetingsShownCount", statistics.GreetingsShownCount, 6, &errors)
	propEq("Statistics.GreetingsAcceptedCount", statistics.GreetingsAcceptedCount, 8, &errors)

	propEq("Customer.AgentLastEventCreatedAt", customer.AgentLastEventCreatedAt.String(), "2019-10-11 09:40:59.249 +0000 UTC", &errors)
	propEq("Customer.CustomerLastEventCreatedAt", customer.CustomerLastEventCreatedAt.String(), "2019-10-11 09:40:59.219001 +0000 UTC", &errors)

	propEq("Chat.Agents.length", len(chat.Agents), 1, &errors)
	aid := "l.wojciechowski@livechatinc.com"
	agent := chat.Agents[aid]
	propEq("Agent.ID", agent.ID, "l.wojciechowski@livechatinc.com", &errors)
	propEq("Agent.Type", agent.Type, "agent", &errors)
	propEq("Agent.Name", agent.Name, "Łukasz Wojciechowski", &errors)
	propEq("Agent.Email", agent.Email, "l.wojciechowski@livechatinc.com", &errors)
	propEq("Agent.Avatar", agent.Avatar, "livechat.s3.amazonaws.com/default/avatars/a14.png", &errors)
	propEq("Agent.Present", agent.Present, true, &errors)
	propEq("Agent.EventsSeenUpTo", agent.EventsSeenUpTo.String(), "1970-01-01 01:00:00 +0000 UTC", &errors)
	propEq("Agent.RoutingStatus", agent.RoutingStatus, "accepting_chats", &errors)

	propEq("Chat.Threads.length", len(chat.Threads), 1, &errors)
	thread := chat.Threads[0]
	propEq("Thread.ID", thread.ID, "PZ070E0W1B", &errors)
	propEq("Thread.Active", thread.Active, true, &errors)
	propEq("Thread.UserIDs[0]", thread.UserIDs[0], "345f8235-d60d-433e-63c5-7f813a6ffe25", &errors)
	propEq("Thread.UserIDs[1]", thread.UserIDs[1], "l.wojciechowski@livechatinc.com", &errors)
	propEq("Thread.RestrictedAccess", thread.RestrictedAccess, false, &errors)
	propEq("Thread.Properties.routing.continuous", thread.Properties["routing"]["continuous"], false, &errors)
	propEq("Thread.Properties.routing.idle", thread.Properties["routing"]["idle"], false, &errors)
	propEq("Thread.Properties.routing.referrer", thread.Properties["routing"]["referrer"], "", &errors)
	propEq("Thread.Properties.routing.start_url", thread.Properties["routing"]["start_url"], "https://cdn.livechatinc.com/labs/?license=100007977/", &errors)
	propEq("Thread.Properties.routing.unassigned", thread.Properties["routing"]["unassigned"], false, &errors)
	propEq("Thread.Access.GroupIDs", thread.Access.GroupIDs[0], 0, &errors)
	propEq("Thread.Events.length", len(thread.Events), 2, &errors)
	propEq("Thread.PreviousThreadID", thread.PreviousThreadID, "K600PKZOM8", &errors)
	propEq("Thread.NextThreadID", thread.NextThreadID, "K600PKZOO8", &errors)
	propEq("Thread.CreatedAt", thread.CreatedAt.String(), "2020-05-07 07:11:28.28834 +0000 UTC", &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}

func incomingEvent(wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.IncomingEvent)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var errors string
	propEq("ChatID", payload.ChatID, "PS0X0L086G", &errors)
	propEq("ThreadID", payload.ThreadID, "PZ070E0W1B", &errors)

	e := payload.Event.Message()
	propEq("Event.ID", e.ID, "PZ070E0W1B_3", &errors)
	propEq("Event.Type", e.Type, "message", &errors)
	propEq("Event.Text", e.Text, "14", &errors)
	propEq("Event.CustomID", e.CustomID, "1dnepb4z00t", &errors)
	propEq("Event.Recipients", e.Recipients, "all", &errors)
	propEq("Event.CreatedAt", e.CreatedAt.String(), "2019-10-11 09:41:00.877 +0000 UTC", &errors)
	propEq("Event.AuthorID", e.AuthorID, "345f8235-d60d-433e-63c5-7f813a6ffe25", &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}

func eventUpdated(wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.EventUpdated)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var errors string
	propEq("ChatID", payload.ChatID, "123-123-123-123", &errors)
	propEq("ThreadID", payload.ThreadID, "E2WDHA8A", &errors)

	e := payload.Event.Message()
	propEq("Event.ID", e.ID, "PZ070E0W1B_3", &errors)
	propEq("Event.Type", e.Type, "message", &errors)
	propEq("Event.Text", e.Text, "14", &errors)
	propEq("Event.CustomID", e.CustomID, "1dnepb4z00t", &errors)
	propEq("Event.Recipients", e.Recipients, "all", &errors)
	propEq("Event.CreatedAt", e.CreatedAt.String(), "2019-10-11 09:41:00.877 +0000 UTC", &errors)
	propEq("Event.AuthorID", e.AuthorID, "345f8235-d60d-433e-63c5-7f813a6ffe25", &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}

func incomingRichMessagePostback(wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.IncomingRichMessagePostback)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var errors string
	propEq("ChatID", payload.ChatID, "PJ0MRSHTDG", &errors)
	propEq("ThreadID", payload.ThreadID, "K600PKZON8", &errors)
	propEq("UserID", payload.UserID, "b7eff798-f8df-4364-8059-649c35c9ed0c", &errors)
	propEq("EventID", payload.EventID, "a0c22fdd-fb71-40b5-bfc6-a8a0bc3117f7", &errors)
	propEq("Postback.ID", payload.Postback.ID, "action_yes", &errors)
	propEq("Postback.Toggled", payload.Postback.Toggled, true, &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}

func chatDeactivated(wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.ChatDeactivated)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var errors string
	propEq("ChatID", payload.ChatID, "PS0X0L086G", &errors)
	propEq("ThreadID", payload.ThreadID, "PZ070E0W1B", &errors)
	propEq("UserID", payload.UserID, "l.wojciechowski@livechatinc.com", &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}

func chatPropertiesUpdated(wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.ChatPropertiesUpdated)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var errors string
	propEq("ChatID", payload.ChatID, "PJ0MRSHTDG", &errors)

	propEq("Properties.Rating.Score.Value", payload.Properties["rating"]["score"], float64(1), &errors)
	propEq("Properties.Rating.Comment.Value", payload.Properties["rating"]["comment"], "Very good, veeeery good", &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}

func threadPropertiesUpdated(wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.ThreadPropertiesUpdated)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var errors string
	propEq("ChatID", payload.ChatID, "PJ0MRSHTDG", &errors)
	propEq("ThreadID", payload.ThreadID, "K600PKZON8", &errors)

	propEq("Properties.Rating.Score.Value", payload.Properties["rating"]["score"], float64(1), &errors)
	propEq("Properties.Rating.Comment.Value", payload.Properties["rating"]["comment"], "Very good, veeeery good", &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}

func chatPropertiesDeleted(wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.ChatPropertiesDeleted)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var errors string
	propEq("ChatID", payload.ChatID, "PJ0MRSHTDG", &errors)

	propEq("Properties.Rating[0]", payload.Properties["rating"][0], "score", &errors)
	propEq("Properties.Rating[1]", payload.Properties["rating"][1], "comment", &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}

func threadPropertiesDeleted(wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.ThreadPropertiesDeleted)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var errors string
	propEq("ChatID", payload.ChatID, "PJ0MRSHTDG", &errors)
	propEq("ThreadID", payload.ThreadID, "K600PKZON8", &errors)

	propEq("Properties.Rating[0]", payload.Properties["rating"][0], "score", &errors)
	propEq("Properties.Rating[1]", payload.Properties["rating"][1], "comment", &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}

func chatUserAdded(wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.ChatUserAdded)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var errors string
	propEq("ChatID", payload.ChatID, "PJ0MRSHTDG", &errors)
	propEq("UserType", payload.UserType, "customer", &errors)

	customer := payload.User.Customer()
	propEq("Customer.ID", customer.ID, "345f8235-d60d-433e-63c5-7f813a6ffe25", &errors)
	propEq("Customer.Type", customer.Type, "customer", &errors)
	propEq("Customer.Name", customer.Name, "test", &errors)
	propEq("Customer.Email", customer.Email, "test@test.pl", &errors)
	propEq("Customer.Avatar", customer.Avatar, "", &errors)
	propEq("Customer.Present", customer.Present, true, &errors)
	propEq("Customer.EventsSeenUpTo", customer.EventsSeenUpTo.String(), "2019-10-08 11:56:53 +0000 UTC", &errors)

	lastVisit := customer.LastVisit
	propEq("LastVisit.IP", lastVisit.IP, "37.248.156.62", &errors)
	propEq("LastVisit.UserAgent", lastVisit.UserAgent, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36", &errors)
	propEq("LastVisit.StartedAt", lastVisit.StartedAt.String(), "2019-10-11 09:40:56.071345 +0000 UTC", &errors)

	geolocation := lastVisit.Geolocation
	propEq("Geolocation.Country", geolocation.Country, "Poland", &errors)
	propEq("Geolocation.CountryCode", geolocation.CountryCode, "PL", &errors)
	propEq("Geolocation.Region", geolocation.Region, "test", &errors)
	propEq("Geolocation.City", geolocation.City, "Wroclaw", &errors)
	propEq("Geolocation.Timezone", geolocation.Timezone, "test_timezone", &errors)

	propEq("LastPages", len(lastVisit.LastPages), 1, &errors)

	propEq("LastPages.OpenedAt", lastVisit.LastPages[0].OpenedAt.String(), "2019-10-11 09:40:56.071345 +0000 UTC", &errors)
	propEq("LastPages.URL", lastVisit.LastPages[0].URL, "https://cdn.livechatinc.com/labs/?license=100007977/", &errors)
	propEq("LastPages.Title", lastVisit.LastPages[0].Title, "LiveChat", &errors)

	statistics := customer.Statistics
	propEq("Statistics.VisistsCount", statistics.VisitsCount, 29, &errors)
	propEq("Statistics.ThreadsCount", statistics.ThreadsCount, 18, &errors)
	propEq("Statistics.ChatsCount", statistics.ChatsCount, 1, &errors)
	propEq("Statistics.PageViewsCount", statistics.PageViewsCount, 5, &errors)
	propEq("Statistics.GreetingsShownCount", statistics.GreetingsShownCount, 6, &errors)
	propEq("Statistics.GreetingsAcceptedCount", statistics.GreetingsAcceptedCount, 8, &errors)

	propEq("Customer.AgentLastEventCreatedAt", customer.AgentLastEventCreatedAt.String(), "2019-10-11 09:40:59.249 +0000 UTC", &errors)
	propEq("Customer.CustomerLastEventCreatedAt", customer.CustomerLastEventCreatedAt.String(), "2019-10-11 09:40:59.219001 +0000 UTC", &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}

func chatUserRemoved(wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.ChatUserRemoved)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var errors string
	propEq("ChatID", payload.ChatID, "PS0X0L086G", &errors)
	propEq("UserType", payload.UserType, "agent", &errors)
	propEq("UserID", payload.UserID, "agent@livechatinc.com", &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}

func threadTagged(wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.ThreadTagged)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var errors string
	propEq("ChatID", payload.ChatID, "PJ0MRSHTDG", &errors)
	propEq("ThreadID", payload.ThreadID, "K600PKZON8", &errors)
	propEq("Tag", payload.Tag, "bug_report", &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}

func threadUntagged(wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.ThreadUntagged)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var errors string
	propEq("ChatID", payload.ChatID, "PJ0MRSHTDG", &errors)
	propEq("ThreadID", payload.ThreadID, "K600PKZON8", &errors)
	propEq("Tag", payload.Tag, "bug_report", &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}

func agentDeleted(wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.AgentDeleted)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var errors string
	propEq("AgentID", payload.AgentID, "5c9871d5372c824cbf22d860a707a578", &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}

func eventsMarkedAsSeen(wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.EventsMarkedAsSeen)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var errors string
	propEq("ChatID", payload.ChatID, "PJ0MRSHTDG", &errors)
	propEq("UserID", payload.UserID, "b7eff798-f8df-4364-8059-649c35c9ed0c", &errors)
	propEq("SeenUpTo", payload.SeenUpTo, "2017-10-12T15:19:21.010200Z", &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}

func accessGranted(wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.AccessGranted)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var errors string
	propEq("ID", payload.ID, "PJ0MRSHTDX", &errors)
	propEq("Resource", payload.Resource, "thread", &errors)
	propEq("Access.GroupIDs.length", len(payload.Access.GroupIDs), 1, &errors)
	propEq("Access.GroupIDs[0]", payload.Access.GroupIDs[0], 2, &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}

func accessRevoked(wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.AccessRevoked)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var errors string
	propEq("ID", payload.ID, "PJ0MRSHTDV", &errors)
	propEq("Resource", payload.Resource, "chat", &errors)
	propEq("Access.GroupIDs.length", len(payload.Access.GroupIDs), 2, &errors)
	propEq("Access.GroupIDs[0]", payload.Access.GroupIDs[0], 3, &errors)
	propEq("Access.GroupIDs[1]", payload.Access.GroupIDs[1], 4, &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}

func accessSet(wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.AccessSet)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var errors string
	propEq("ID", payload.ID, "PJ0MRSHTDG", &errors)
	propEq("Resource", payload.Resource, "chat", &errors)
	propEq("Access.GroupIDs.length", len(payload.Access.GroupIDs), 1, &errors)
	propEq("Access.GroupIDs[0]", payload.Access.GroupIDs[0], 1, &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}

func customerCreated(wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.CustomerCreated)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}
	var errors string
	if payload.User == nil {
		return fmt.Errorf("`Customer.User` is nil")
	}
	propEq("User.ID", payload.User.ID, "baf3cf72-4768-42e4-6140-26dd36c962cc", &errors)
	t, err := time.Parse(time.RFC3339Nano, "2019-11-14T14:27:24.410018Z")
	if err != nil {
		return fmt.Errorf("Couldn't parse time: %v", err)
	}
	propEq("CreatedAt", payload.CreatedAt, t, &errors)
	propEq("Email", payload.Email, "customer1@example.com", &errors)
	propEq("Avatar", payload.Avatar, "https://example.com/avatars/1.jpg", &errors)
	propEq("SessionFields", len(payload.SessionFields), 2, &errors)
	propEq("SessionFields[0][some_key]", payload.SessionFields[0]["some_key"], "some_value", &errors)
	propEq("SessionFields[1][some_other_key]", payload.SessionFields[1]["some_other_key"], "some_other_value", &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}

func eventPropertiesDeleted(wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.EventPropertiesDeleted)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var errors string
	propEq("ChatID", payload.ChatID, "PJ0MRSHTDG", &errors)
	propEq("ThreadID", payload.ThreadID, "K600PKZON8", &errors)
	propEq("EventID", payload.EventID, "2_E2WDHA8A", &errors)

	propEq("Properties.Rating[0]", payload.Properties["rating"][0], "score", &errors)
	propEq("Properties.Rating[1]", payload.Properties["rating"][1], "comment", &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}

func eventPropertiesUpdated(wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.EventPropertiesUpdated)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var errors string
	propEq("ChatID", payload.ChatID, "PJ0MRSHTDG", &errors)
	propEq("ThreadID", payload.ThreadID, "K600PKZON8", &errors)
	propEq("EventID", payload.EventID, "2_E2WDHA8A", &errors)

	propEq("Properties.Rating.Score.Value", payload.Properties["rating"]["score"], float64(1), &errors)
	propEq("Properties.Rating.Comment.Value", payload.Properties["rating"]["comment"], "Very good, veeeery good", &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}

func routingStatusSet(wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.RoutingStatusSet)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", payload)
	}

	var errors string
	propEq("AgentID", payload.AgentID, "5c9871d5372c824cbf22d860a707a578", &errors)
	propEq("Status", payload.Status, "accepting chats", &errors)

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}
