package webhooks_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/livechat/lc-sdk-go/v5/configuration"
	"github.com/livechat/lc-sdk-go/v5/webhooks"
)

func propEq(propertyName string, actual, expected interface{}, validationAccumulator *string) {
	if actual != expected {
		*validationAccumulator += fmt.Sprintf("%s mismatch, actual: %v, expected: %v\n", propertyName, actual, expected)
	}
}

func incomingChat(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.IncomingChat)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	chat := payload.Chat

	var propEqErrors string
	propEq("Chat.ID", chat.ID, "PS0X0L086G", &propEqErrors)
	propEq("Chat.Access.GroupIDs", len(chat.Access.GroupIDs), 1, &propEqErrors)
	propEq("Chat.Access.GroupIDs[0]", chat.Access.GroupIDs[0], 0, &propEqErrors)
	propEq("Chat.Users length", len(chat.Users()), 2, &propEqErrors)

	propEq("Chat.Customers", len(chat.Customers), 1, &propEqErrors)
	cid := "345f8235-d60d-433e-63c5-7f813a6ffe25"
	customer := chat.Customers[cid]
	propEq("Customer.ID", customer.ID, "345f8235-d60d-433e-63c5-7f813a6ffe25", &propEqErrors)
	propEq("Customer.Type", customer.Type, "customer", &propEqErrors)
	propEq("Customer.Name", customer.Name, "test", &propEqErrors)
	propEq("Customer.Email", customer.Email, "test@test.pl", &propEqErrors)
	propEq("Customer.Avatar", customer.Avatar, "", &propEqErrors)
	propEq("Customer.Present", customer.Present, true, &propEqErrors)
	propEq("Customer.EventsSeenUpTo", customer.EventsSeenUpTo.String(), "2019-10-08 13:56:53 +0000 UTC", &propEqErrors)

	lastVisit := customer.LastVisit
	propEq("LastVisit.IP", lastVisit.IP, "37.248.156.62", &propEqErrors)
	propEq("LastVisit.UserAgent", lastVisit.UserAgent, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36", &propEqErrors)
	propEq("LastVisit.StartedAt", lastVisit.StartedAt.String(), "2019-10-11 09:40:56.071345 +0000 UTC", &propEqErrors)

	geolocation := lastVisit.Geolocation
	propEq("Geolocation.Country", geolocation.Country, "Poland", &propEqErrors)
	propEq("Geolocation.CountryCode", geolocation.CountryCode, "PL", &propEqErrors)
	propEq("Geolocation.Region", geolocation.Region, "test", &propEqErrors)
	propEq("Geolocation.City", geolocation.City, "Wroclaw", &propEqErrors)
	propEq("Geolocation.Timezone", geolocation.Timezone, "test_timezone", &propEqErrors)

	propEq("LastPages", len(lastVisit.LastPages), 1, &propEqErrors)

	propEq("LastPages.OpenedAt", lastVisit.LastPages[0].OpenedAt.String(), "2019-10-11 09:40:56.071345 +0000 UTC", &propEqErrors)
	propEq("LastPages.URL", lastVisit.LastPages[0].URL, "https://cdn.livechatinc.com/labs/?license=100007977/", &propEqErrors)
	propEq("LastPages.Title", lastVisit.LastPages[0].Title, "LiveChat", &propEqErrors)

	statistics := customer.Statistics
	propEq("Statistics.VisistsCount", statistics.VisitsCount, 29, &propEqErrors)
	propEq("Statistics.ThreadsCount", statistics.ThreadsCount, 18, &propEqErrors)
	propEq("Statistics.ChatsCount", statistics.ChatsCount, 1, &propEqErrors)
	propEq("Statistics.PageViewsCount", statistics.PageViewsCount, 5, &propEqErrors)
	propEq("Statistics.GreetingsShownCount", statistics.GreetingsShownCount, 6, &propEqErrors)
	propEq("Statistics.GreetingsAcceptedCount", statistics.GreetingsAcceptedCount, 8, &propEqErrors)

	propEq("Customer.AgentLastEventCreatedAt", customer.AgentLastEventCreatedAt.String(), "2019-10-11 09:40:59.249 +0000 UTC", &propEqErrors)
	propEq("Customer.CustomerLastEventCreatedAt", customer.CustomerLastEventCreatedAt.String(), "2019-10-11 09:40:59.219001 +0000 UTC", &propEqErrors)

	propEq("Chat.Agents.length", len(chat.Agents), 1, &propEqErrors)
	aid := "l.wojciechowski@livechatinc.com"
	agent := chat.Agents[aid]
	propEq("Agent.ID", agent.ID, "l.wojciechowski@livechatinc.com", &propEqErrors)
	propEq("Agent.Type", agent.Type, "agent", &propEqErrors)
	propEq("Agent.Name", agent.Name, "Łukasz Wojciechowski", &propEqErrors)
	propEq("Agent.Email", agent.Email, "l.wojciechowski@livechatinc.com", &propEqErrors)
	propEq("Agent.Avatar", agent.Avatar, "livechat.s3.amazonaws.com/default/avatars/a14.png", &propEqErrors)
	propEq("Agent.Present", agent.Present, true, &propEqErrors)
	propEq("Agent.EventsSeenUpTo", agent.EventsSeenUpTo.String(), "1970-01-01 01:00:00 +0000 UTC", &propEqErrors)
	propEq("Agent.RoutingStatus", agent.RoutingStatus, "accepting_chats", &propEqErrors)

	propEq("Chat.Threads.length", len(chat.Threads), 1, &propEqErrors)
	thread := chat.Threads[0]
	propEq("Thread.ID", thread.ID, "PZ070E0W1B", &propEqErrors)
	propEq("Thread.Active", thread.Active, true, &propEqErrors)
	propEq("Thread.UserIDs[0]", thread.UserIDs[0], "345f8235-d60d-433e-63c5-7f813a6ffe25", &propEqErrors)
	propEq("Thread.UserIDs[1]", thread.UserIDs[1], "l.wojciechowski@livechatinc.com", &propEqErrors)
	propEq("Thread.RestrictedAccess", thread.RestrictedAccess, false, &propEqErrors)
	propEq("Thread.Properties.routing.continuous", thread.Properties["routing"]["continuous"], false, &propEqErrors)
	propEq("Thread.Properties.routing.idle", thread.Properties["routing"]["idle"], false, &propEqErrors)
	propEq("Thread.Properties.routing.referrer", thread.Properties["routing"]["referrer"], "", &propEqErrors)
	propEq("Thread.Properties.routing.start_url", thread.Properties["routing"]["start_url"], "https://cdn.livechatinc.com/labs/?license=100007977/", &propEqErrors)
	propEq("Thread.Properties.routing.unassigned", thread.Properties["routing"]["unassigned"], false, &propEqErrors)
	propEq("Thread.Access.GroupIDs", thread.Access.GroupIDs[0], 0, &propEqErrors)
	propEq("Thread.Events.length", len(thread.Events), 2, &propEqErrors)
	propEq("Thread.PreviousThreadID", thread.PreviousThreadID, "K600PKZOM8", &propEqErrors)
	propEq("Thread.NextThreadID", thread.NextThreadID, "K600PKZOO8", &propEqErrors)
	propEq("Thread.CreatedAt", thread.CreatedAt.String(), "2020-05-07 07:11:28.28834 +0000 UTC", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func incomingEvent(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.IncomingEvent)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var propEqErrors string
	propEq("ChatID", payload.ChatID, "PS0X0L086G", &propEqErrors)
	propEq("ThreadID", payload.ThreadID, "PZ070E0W1B", &propEqErrors)

	e := payload.Event.Message()
	propEq("Event.ID", e.ID, "PZ070E0W1B_3", &propEqErrors)
	propEq("Event.Type", e.Type, "message", &propEqErrors)
	propEq("Event.Text", e.Text, "14", &propEqErrors)
	propEq("Event.CustomID", e.CustomID, "1dnepb4z00t", &propEqErrors)
	propEq("Event.Visibility", e.Visibility, "all", &propEqErrors)
	propEq("Event.CreatedAt", e.CreatedAt.String(), "2019-10-11 09:41:00.877 +0000 UTC", &propEqErrors)
	propEq("Event.AuthorID", e.AuthorID, "345f8235-d60d-433e-63c5-7f813a6ffe25", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func eventUpdated(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.EventUpdated)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var propEqErrors string
	propEq("ChatID", payload.ChatID, "123-123-123-123", &propEqErrors)
	propEq("ThreadID", payload.ThreadID, "E2WDHA8A", &propEqErrors)

	e := payload.Event.Message()
	propEq("Event.ID", e.ID, "PZ070E0W1B_3", &propEqErrors)
	propEq("Event.Type", e.Type, "message", &propEqErrors)
	propEq("Event.Text", e.Text, "14", &propEqErrors)
	propEq("Event.CustomID", e.CustomID, "1dnepb4z00t", &propEqErrors)
	propEq("Event.Visibility", e.Visibility, "all", &propEqErrors)
	propEq("Event.CreatedAt", e.CreatedAt.String(), "2019-10-11 09:41:00.877 +0000 UTC", &propEqErrors)
	propEq("Event.AuthorID", e.AuthorID, "345f8235-d60d-433e-63c5-7f813a6ffe25", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func incomingRichMessagePostback(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.IncomingRichMessagePostback)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var propEqErrors string
	propEq("ChatID", payload.ChatID, "PJ0MRSHTDG", &propEqErrors)
	propEq("ThreadID", payload.ThreadID, "K600PKZON8", &propEqErrors)
	propEq("UserID", payload.UserID, "b7eff798-f8df-4364-8059-649c35c9ed0c", &propEqErrors)
	propEq("EventID", payload.EventID, "a0c22fdd-fb71-40b5-bfc6-a8a0bc3117f7", &propEqErrors)
	propEq("Postback.ID", payload.Postback.ID, "action_yes", &propEqErrors)
	propEq("Postback.Toggled", payload.Postback.Toggled, true, &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func chatDeactivated(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.ChatDeactivated)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var propEqErrors string
	propEq("ChatID", payload.ChatID, "PS0X0L086G", &propEqErrors)
	propEq("ThreadID", payload.ThreadID, "PZ070E0W1B", &propEqErrors)
	propEq("UserID", payload.UserID, "l.wojciechowski@livechatinc.com", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func chatPropertiesUpdated(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.ChatPropertiesUpdated)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var propEqErrors string
	propEq("ChatID", payload.ChatID, "PJ0MRSHTDG", &propEqErrors)

	propEq("Properties.Rating.Score.Value", payload.Properties["rating"]["score"], float64(1), &propEqErrors)
	propEq("Properties.Rating.Comment.Value", payload.Properties["rating"]["comment"], "Very good, veeeery good", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func threadPropertiesUpdated(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.ThreadPropertiesUpdated)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var propEqErrors string
	propEq("ChatID", payload.ChatID, "PJ0MRSHTDG", &propEqErrors)
	propEq("ThreadID", payload.ThreadID, "K600PKZON8", &propEqErrors)

	propEq("Properties.Rating.Score.Value", payload.Properties["rating"]["score"], float64(1), &propEqErrors)
	propEq("Properties.Rating.Comment.Value", payload.Properties["rating"]["comment"], "Very good, veeeery good", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func chatPropertiesDeleted(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.ChatPropertiesDeleted)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var propEqErrors string
	propEq("ChatID", payload.ChatID, "PJ0MRSHTDG", &propEqErrors)

	propEq("Properties.Rating[0]", payload.Properties["rating"][0], "score", &propEqErrors)
	propEq("Properties.Rating[1]", payload.Properties["rating"][1], "comment", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func threadPropertiesDeleted(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.ThreadPropertiesDeleted)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var propEqErrors string
	propEq("ChatID", payload.ChatID, "PJ0MRSHTDG", &propEqErrors)
	propEq("ThreadID", payload.ThreadID, "K600PKZON8", &propEqErrors)

	propEq("Properties.Rating[0]", payload.Properties["rating"][0], "score", &propEqErrors)
	propEq("Properties.Rating[1]", payload.Properties["rating"][1], "comment", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func userAddedToChat(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.UserAddedToChat)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var propEqErrors string
	propEq("ChatID", payload.ChatID, "PJ0MRSHTDG", &propEqErrors)
	propEq("ThreadID", payload.ThreadID, "K600PKZON8", &propEqErrors)
	propEq("Reason", payload.Reason, "manual", &propEqErrors)
	propEq("RequesterID", payload.RequesterID, "smith@example.com", &propEqErrors)
	propEq("UserType", payload.UserType, "customer", &propEqErrors)

	customer := payload.User.Customer()
	if customer == nil {
		return errors.New("`User.Customer` is nil")
	}
	propEq("Customer.ID", customer.ID, "345f8235-d60d-433e-63c5-7f813a6ffe25", &propEqErrors)
	propEq("Customer.Type", customer.Type, "customer", &propEqErrors)
	propEq("Customer.Name", customer.Name, "test", &propEqErrors)
	propEq("Customer.Email", customer.Email, "test@test.pl", &propEqErrors)
	propEq("Customer.Avatar", customer.Avatar, "", &propEqErrors)
	propEq("Customer.Present", customer.Present, true, &propEqErrors)
	propEq("Customer.EmailVerified", customer.EmailVerified, true, &propEqErrors)
	propEq("Customer.EventsSeenUpTo", customer.EventsSeenUpTo.String(), "2019-10-08 11:56:53 +0000 UTC", &propEqErrors)

	lastVisit := customer.LastVisit
	propEq("LastVisit.IP", lastVisit.IP, "37.248.156.62", &propEqErrors)
	propEq("LastVisit.UserAgent", lastVisit.UserAgent, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36", &propEqErrors)
	propEq("LastVisit.StartedAt", lastVisit.StartedAt.String(), "2019-10-11 09:40:56.071345 +0000 UTC", &propEqErrors)

	geolocation := lastVisit.Geolocation
	propEq("Geolocation.Country", geolocation.Country, "Poland", &propEqErrors)
	propEq("Geolocation.CountryCode", geolocation.CountryCode, "PL", &propEqErrors)
	propEq("Geolocation.Region", geolocation.Region, "test", &propEqErrors)
	propEq("Geolocation.City", geolocation.City, "Wroclaw", &propEqErrors)
	propEq("Geolocation.Timezone", geolocation.Timezone, "test_timezone", &propEqErrors)

	propEq("LastPages", len(lastVisit.LastPages), 1, &propEqErrors)

	propEq("LastPages.OpenedAt", lastVisit.LastPages[0].OpenedAt.String(), "2019-10-11 09:40:56.071345 +0000 UTC", &propEqErrors)
	propEq("LastPages.URL", lastVisit.LastPages[0].URL, "https://cdn.livechatinc.com/labs/?license=100007977/", &propEqErrors)
	propEq("LastPages.Title", lastVisit.LastPages[0].Title, "LiveChat", &propEqErrors)

	statistics := customer.Statistics
	propEq("Statistics.VisistsCount", statistics.VisitsCount, 29, &propEqErrors)
	propEq("Statistics.ThreadsCount", statistics.ThreadsCount, 18, &propEqErrors)
	propEq("Statistics.ChatsCount", statistics.ChatsCount, 1, &propEqErrors)
	propEq("Statistics.PageViewsCount", statistics.PageViewsCount, 5, &propEqErrors)
	propEq("Statistics.GreetingsShownCount", statistics.GreetingsShownCount, 6, &propEqErrors)
	propEq("Statistics.GreetingsAcceptedCount", statistics.GreetingsAcceptedCount, 8, &propEqErrors)

	propEq("Customer.AgentLastEventCreatedAt", customer.AgentLastEventCreatedAt.String(), "2019-10-11 09:40:59.249 +0000 UTC", &propEqErrors)
	propEq("Customer.CustomerLastEventCreatedAt", customer.CustomerLastEventCreatedAt.String(), "2019-10-11 09:40:59.219001 +0000 UTC", &propEqErrors)
	propEq("Customer.SessionFields[0][\"some_key\"]", customer.SessionFields[0]["some_key"], "some_value", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func userRemovedFromChat(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.UserRemovedFromChat)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var propEqErrors string
	propEq("ChatID", payload.ChatID, "PJ0MRSHTDG", &propEqErrors)
	propEq("ThreadID", payload.ThreadID, "K600PKZON8", &propEqErrors)
	propEq("Reason", payload.Reason, "manual", &propEqErrors)
	propEq("RequesterID", payload.RequesterID, "smith@example.com", &propEqErrors)
	propEq("UserType", payload.UserType, "agent", &propEqErrors)
	propEq("UserID", payload.UserID, "agent@livechatinc.com", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func threadTagged(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.ThreadTagged)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var propEqErrors string
	propEq("ChatID", payload.ChatID, "PJ0MRSHTDG", &propEqErrors)
	propEq("ThreadID", payload.ThreadID, "K600PKZON8", &propEqErrors)
	propEq("Tag", payload.Tag, "bug_report", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func threadUntagged(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.ThreadUntagged)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var propEqErrors string
	propEq("ChatID", payload.ChatID, "PJ0MRSHTDG", &propEqErrors)
	propEq("ThreadID", payload.ThreadID, "K600PKZON8", &propEqErrors)
	propEq("Tag", payload.Tag, "bug_report", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func agentCreated(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.AgentCreated)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var propEqErrors string
	propEq("ID", payload.ID, "smith@example.com", &propEqErrors)
	propEq("Name", payload.Name, "Agent Smith", &propEqErrors)
	propEq("Role", payload.Role, "viceowner", &propEqErrors)
	propEq("AwaitingApproval", payload.AwaitingApproval, false, &propEqErrors)
	propEq("Groups[0].ID", payload.Groups[0].ID, uint(5), &propEqErrors)
	propEq("Groups[1].Priority", payload.Groups[1].Priority, configuration.Last, &propEqErrors)
	propEq("Notifications[0]", payload.Notifications[0], "new_visitor", &propEqErrors)
	propEq("Notifications[1]", payload.Notifications[1], "new_goal", &propEqErrors)
	propEq("EmailSubscriptions[0]", payload.EmailSubscriptions[0], "weekly_summary", &propEqErrors)
	propEq("WorkScheduler.Timezone", payload.WorkScheduler.Timezone, "Europe/Warsaw", &propEqErrors)
	propEq("WorkScheduler.Schedule[0].Start", payload.WorkScheduler.Schedule[0].Start, "08:30", &propEqErrors)
	propEq("WorkScheduler.Schedule[0].Day", payload.WorkScheduler.Schedule[0].Day, configuration.Monday, &propEqErrors)
	propEq("WorkScheduler.Schedule[0].Enabled", payload.WorkScheduler.Schedule[0].Enabled, true, &propEqErrors)
	propEq("WorkScheduler.Schedule[0].Start", payload.WorkScheduler.Schedule[0].Start, "08:30", &propEqErrors)
	propEq("WorkScheduler.Schedule[0].End", payload.WorkScheduler.Schedule[0].End, "12:30", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func agentUpdated(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.AgentUpdated)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var propEqErrors string
	propEq("ID", payload.ID, "smith@example.com", &propEqErrors)
	propEq("WorkScheduler.Timezone", payload.WorkScheduler.Timezone, "Europe/Warsaw", &propEqErrors)
	propEq("WorkScheduler.Schedule[0].Day", payload.WorkScheduler.Schedule[0].Day, configuration.Monday, &propEqErrors)
	propEq("WorkScheduler.Schedule[0].Enabled", payload.WorkScheduler.Schedule[0].Enabled, true, &propEqErrors)
	propEq("WorkScheduler.Schedule[0].Start", payload.WorkScheduler.Schedule[0].Start, "08:30", &propEqErrors)
	propEq("WorkScheduler.Schedule[0].End", payload.WorkScheduler.Schedule[0].End, "12:30", &propEqErrors)
	propEq("WorkScheduler.Schedule[1].Day", payload.WorkScheduler.Schedule[1].Day, configuration.Friday, &propEqErrors)
	propEq("WorkScheduler.Schedule[1].Enabled", payload.WorkScheduler.Schedule[1].Enabled, true, &propEqErrors)
	propEq("WorkScheduler.Schedule[1].Start", payload.WorkScheduler.Schedule[1].Start, "07:30", &propEqErrors)
	propEq("WorkScheduler.Schedule[1].End", payload.WorkScheduler.Schedule[1].End, "21:30", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func agentDeleted(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.AgentDeleted)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var propEqErrors string
	propEq("ID", payload.ID, "smith@example.com", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func agentSuspended(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.AgentSuspended)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var propEqErrors string
	propEq("ID", payload.ID, "smith@example.com", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func agentUnsuspended(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.AgentUnsuspended)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var propEqErrors string
	propEq("ID", payload.ID, "smith@example.com", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func agentApproved(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.AgentApproved)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var propEqErrors string
	propEq("ID", payload.ID, "smith@example.com", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func eventsMarkedAsSeen(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.EventsMarkedAsSeen)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var propEqErrors string
	propEq("ChatID", payload.ChatID, "PJ0MRSHTDG", &propEqErrors)
	propEq("UserID", payload.UserID, "b7eff798-f8df-4364-8059-649c35c9ed0c", &propEqErrors)
	propEq("SeenUpTo", payload.SeenUpTo, "2017-10-12T15:19:21.010200Z", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func chatAccessUpdated(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.ChatAccessUpdated)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var propEqErrors string
	propEq("ID", payload.ID, "PJ0MRSHTDX", &propEqErrors)
	propEq("Access.GroupIDs.length", len(payload.Access.GroupIDs), 1, &propEqErrors)
	propEq("Access.GroupIDs[0]", payload.Access.GroupIDs[0], 2, &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func incomingCustomer(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.IncomingCustomer)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}
	var propEqErrors string
	if payload.User == nil {
		return errors.New("`Customer.User` is nil")
	}
	propEq("User.ID", payload.User.ID, "baf3cf72-4768-42e4-6140-26dd36c962cc", &propEqErrors)
	t, err := time.Parse(time.RFC3339Nano, "2019-11-14T14:27:24.410018Z")
	if err != nil {
		return fmt.Errorf("Couldn't parse time: %v", err)
	}
	propEq("CreatedAt", payload.CreatedAt, t, &propEqErrors)
	propEq("Email", payload.Email, "customer1@example.com", &propEqErrors)
	propEq("Avatar", payload.Avatar, "https://example.com/avatars/1.jpg", &propEqErrors)
	propEq("SessionFields", len(payload.SessionFields), 2, &propEqErrors)
	propEq("SessionFields[0][some_key]", payload.SessionFields[0]["some_key"], "some_value", &propEqErrors)
	propEq("SessionFields[1][some_other_key]", payload.SessionFields[1]["some_other_key"], "some_other_value", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func eventPropertiesDeleted(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.EventPropertiesDeleted)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var propEqErrors string
	propEq("ChatID", payload.ChatID, "PJ0MRSHTDG", &propEqErrors)
	propEq("ThreadID", payload.ThreadID, "K600PKZON8", &propEqErrors)
	propEq("EventID", payload.EventID, "2_E2WDHA8A", &propEqErrors)

	propEq("Properties.Rating[0]", payload.Properties["rating"][0], "score", &propEqErrors)
	propEq("Properties.Rating[1]", payload.Properties["rating"][1], "comment", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func eventPropertiesUpdated(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.EventPropertiesUpdated)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}

	var propEqErrors string
	propEq("ChatID", payload.ChatID, "PJ0MRSHTDG", &propEqErrors)
	propEq("ThreadID", payload.ThreadID, "K600PKZON8", &propEqErrors)
	propEq("EventID", payload.EventID, "2_E2WDHA8A", &propEqErrors)

	propEq("Properties.Rating.Score.Value", payload.Properties["rating"]["score"], float64(1), &propEqErrors)
	propEq("Properties.Rating.Comment.Value", payload.Properties["rating"]["comment"], "Very good, veeeery good", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func routingStatusSet(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.RoutingStatusSet)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", payload)
	}

	var propEqErrors string
	propEq("AgentID", payload.AgentID, "5c9871d5372c824cbf22d860a707a578", &propEqErrors)
	propEq("Status", payload.Status, "accepting chats", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func chatTransferred(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.ChatTransferred)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", payload)
	}

	var propEqErrors string
	propEq("ChatID", payload.ChatID, "PJ0MRSHTDG", &propEqErrors)
	propEq("ThreadID", payload.ThreadID, "K600PKZON8", &propEqErrors)
	propEq("RequesterID", payload.RequesterID, "5c9871d5372c824cbf22d860a707a578", &propEqErrors)
	propEq("Reason", payload.Reason, "manual", &propEqErrors)
	propEq("TransferredTo.GroupIDs.length", len(payload.TransferredTo.AgentIDs), 1, &propEqErrors)
	propEq("TransferredTo.GroupIDs[0]", payload.TransferredTo.AgentIDs[0], "l.wojciechowski@livechatinc.com", &propEqErrors)
	propEq("TransferredTo.GroupIDs.length", len(payload.TransferredTo.GroupIDs), 1, &propEqErrors)
	propEq("TransferredTo.GroupIDs[0]", payload.TransferredTo.GroupIDs[0], 2, &propEqErrors)
	propEq("Queue.Position", payload.Queue.Position, 42, &propEqErrors)
	propEq("Queue.WaitTime", payload.Queue.WaitTime, 1337, &propEqErrors)
	propEq("Queue.QueuedAt", payload.Queue.QueuedAt, "2019-12-09T12:01:18.909000Z", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func customerSessionFieldsUpdated(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.CustomerSessionFieldsUpdated)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}
	var propEqErrors string
	propEq("ID", payload.ID, "5280e68c-9692-4212-4ba9-85f7d8af55cd", &propEqErrors)
	propEq("ActiveChat.ChatID", payload.ActiveChat.ChatID, "PJ0MRSHTDG", &propEqErrors)
	propEq("ActiveChat.ThreadID", payload.ActiveChat.ThreadID, "K600PKZON8", &propEqErrors)
	propEq("ActiveChat.SessionFields[0][\"key\"]", payload.SessionFields[0]["key"], "value", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func groupCreated(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.GroupCreated)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}
	var propEqErrors string
	propEq("ID", payload.ID, 42, &propEqErrors)
	propEq("Name", payload.Name, "sales", &propEqErrors)
	propEq("LanguageCode", payload.LanguageCode, "en", &propEqErrors)
	propEq("AgentPriorities[\"agent@example.com\"]", payload.AgentPriorities["agent@example.com"], "normal", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func groupUpdated(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.GroupUpdated)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}
	var propEqErrors string
	propEq("ID", payload.ID, 42, &propEqErrors)
	propEq("Name", payload.Name, "sales", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func groupDeleted(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.GroupDeleted)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}
	var propEqErrors string
	propEq("ID", payload.ID, 42, &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func autoAccessAdded(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.AutoAccessAdded)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}
	var propEqErrors string
	propEq("ID", payload.ID, "pqi8oasdjahuakndw9nsad9na", &propEqErrors)
	propEq("Description", payload.Description, "Chats on livechat.com from United States", &propEqErrors)
	propEq("Access.Groups[0]", payload.Access.Groups[0], 1, &propEqErrors)
	propEq("Conditions.Domain.Values[0].Value", payload.Conditions.Domain.Values[0].Value, "livechat.com", &propEqErrors)
	propEq("Conditions.Domain.Values[0].ExactMatch", payload.Conditions.Domain.Values[0].ExactMatch, true, &propEqErrors)
	propEq("Conditions.Geolocation.Values[0].Country", payload.Conditions.Geolocation.Values[0].Country, "United States", &propEqErrors)
	propEq("Conditions.Geolocation.Values[0].CountryCode", payload.Conditions.Geolocation.Values[0].CountryCode, "US", &propEqErrors)
	propEq("NextID", payload.NextID, "1faad6f5f1d6e8fdf27e8af9839783b7", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func autoAccessUpdated(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.AutoAccessUpdated)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}
	var propEqErrors string
	propEq("ID", payload.ID, "pqi8oasdjahuakndw9nsad9na", &propEqErrors)
	propEq("Access.Groups[0]", payload.Access.Groups[0], 0, &propEqErrors)
	propEq("Access.Groups[1]", payload.Access.Groups[1], 42, &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func botCreated(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.BotCreated)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}
	var propEqErrors string
	propEq("ID", payload.ID, "5c9871d5372c824cbf22d860a707a578", &propEqErrors)
	propEq("Name", payload.Name, "Bot Name", &propEqErrors)
	propEq("DefaultGroupPriority", payload.DefaultGroupPriority, configuration.First, &propEqErrors)
	propEq("Groups[0].ID", payload.Groups[0].ID, uint(0), &propEqErrors)
	propEq("Groups[0].Priority", payload.Groups[0].Priority, configuration.Normal, &propEqErrors)
	propEq("OwnerClientID", payload.OwnerClientID, "asXdesldiAJSq9padj", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func botUpdated(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.BotUpdated)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}
	var propEqErrors string
	propEq("ID", payload.ID, "5c9871d5372c824cbf22d860a707a578", &propEqErrors)
	propEq("Name", payload.Name, "New Bot Name", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func botDeleted(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.BotDeleted)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}
	var propEqErrors string
	propEq("ID", payload.ID, "5c9871d5372c824cbf22d860a707a578", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}

func autoAccessDeleted(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.AutoAccessDeleted)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", wh.Payload)
	}
	var propEqErrors string
	propEq("ID", payload.ID, "pqi8oasdjahuakndw9nsad9na", &propEqErrors)

	if propEqErrors != "" {
		return errors.New(propEqErrors)
	}
	return nil
}
