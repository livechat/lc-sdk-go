package webhooks_test

import (
	"testing"

	"github.com/livechat/lc-sdk-go/webhooks"
)

const incomingChatThread = `{"webhook_id":"adaa18f4abe65dec0f86a5047ed32d2c","secret_key":"1234567890","action":"incoming_chat_thread","payload":{"chat":{"id":"PS0X0L086G","users":[{"id":"345f8235-d60d-433e-63c5-7f813a6ffe25","name":"test","email":"test@test.pl","present":true,"last_seen_timestamp":1570535813,"type":"customer","created_at":"2019-06-11T11:00:10.329000Z","last_visit":{"ip":"37.248.156.62","user_agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36","geolocation":{"country":"Poland","country_code":"PL","region":"test","city":"Wroclaw","timezone":"test_timezone"},"started_at":"2019-10-11T09:40:56.071345Z","last_pages":[{"opened_at":"2019-10-11T09:40:56.071345Z","url":"https://cdn.livechatinc.com/labs/?license=100007977/","title":"LiveChat"}]},"statistics":{"visits_count":29,"threads_count":18,"chats_count":1,"page_views_count":5,"greetings_shown_count":6,"greetings_accepted_count":8},"agent_last_event_created_at":"2019-10-11T09:40:59.249000Z","customer_last_event_created_at":"2019-10-11T09:40:59.219001Z"},{"id":"l.wojciechowski@livechatinc.com","name":"\u0141ukasz Wojciechowski","email":"l.wojciechowski@livechatinc.com","present":true,"last_seen_timestamp":0,"type":"agent","avatar":"livechat.s3.amazonaws.com/default/avatars/a14.png","routing_status":"accepting_chats"}],"thread":{"id":"PZ070E0W1B","timestamp":1570786859,"active":true,"order":18,"properties":{"routing":{"continuous":false,"idle":false,"referrer":"","start_url":"https://cdn.livechatinc.com/labs/?license=100007977/","unassigned":false},"source":{"client_id":"c6e4f62e2a2dab12531235b12c5a2a6b"}},"user_ids":["345f8235-d60d-433e-63c5-7f813a6ffe25","l.wojciechowski@livechatinc.com"],"events":[{"type":"filled_form","fields":[{"label":"Name:","type":"name","value":"rewrew"},{"label":"E-mail:","type":"email","value":""}],"id":"PZ070E0W1B_1","custom_id":"3g4a8a2p6e2","recipients":"all","created_at":"2019-10-11T09:40:59.219001Z","author_id":"345f8235-d60d-433e-63c5-7f813a6ffe25","properties":{"lc2":{"form_type":"prechat"}}},{"type":"message","text":"Hello. How may I help you?","id":"PZ070E0W1B_2","custom_id":"","recipients":"all","created_at":"2019-10-11T09:40:59.249000Z","author_id":"l.wojciechowski@livechatinc.com","properties":{"lc2":{"welcome_message":true}}}],"access":{"group_ids":[0]}},"properties":{"routing":{"continuous":false,"pinned":false},"source":{"client_id":"c6e4f62e2a2dab12531235b12c5a2a6b"},"supervising":{"agent_ids":""}},"access":{"group_ids":[0]}}},"additional_data":{}}`

func eq(actual interface{}, expected interface{}, t *testing.T) {
	if expected != actual {
		t.Errorf("Values does not match. Expected %s got %s", expected, actual)
	}
}

func TestParseIncomingChatThreadPayload(t *testing.T) {

	p, err := webhooks.ParseIncomingChatThreadPayload([]byte(incomingChatThread))

	if err != nil {
		t.Error(err)
	}

	eq(p.WebhookID, "adaa18f4abe65dec0f86a5047ed32d2c", t)
	eq(p.SecretKey, "1234567890", t)
	eq(p.Action, "incoming_chat_thread", t)

	eq(p.Payload.Chat.Access.GroupIDs[0], 0, t)
	eq(len(p.Payload.Chat.Access.GroupIDs), 1, t)
	eq(len(p.Payload.Chat.Users()), 2, t)

	eq(len(p.Payload.Chat.Customers), 1, t)
	cid := "345f8235-d60d-433e-63c5-7f813a6ffe25"
	eq(p.Payload.Chat.ID, "PS0X0L086G", t)
	eq(p.Payload.Chat.Customers[cid].ID, "345f8235-d60d-433e-63c5-7f813a6ffe25", t)
	eq(p.Payload.Chat.Customers[cid].Type, "customer", t)
	eq(p.Payload.Chat.Customers[cid].Name, "test", t)
	eq(p.Payload.Chat.Customers[cid].Email, "test@test.pl", t)
	eq(p.Payload.Chat.Customers[cid].Avatar, "", t)
	eq(p.Payload.Chat.Customers[cid].Present, true, t)
	eq(p.Payload.Chat.Customers[cid].LastSeen.String(), "2019-10-08 13:56:53 +0200 CEST", t)

	eq(p.Payload.Chat.Customers[cid].LastVisit.IP, "37.248.156.62", t)
	eq(p.Payload.Chat.Customers[cid].LastVisit.UserAgent, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36", t)
	eq(p.Payload.Chat.Customers[cid].LastVisit.StartedAt.String(), "2019-10-11 09:40:56.071345 +0000 UTC", t)

	eq(p.Payload.Chat.Customers[cid].LastVisit.Geolocation.Country, "Poland", t)
	eq(p.Payload.Chat.Customers[cid].LastVisit.Geolocation.CountryCode, "PL", t)
	eq(p.Payload.Chat.Customers[cid].LastVisit.Geolocation.Region, "test", t)
	eq(p.Payload.Chat.Customers[cid].LastVisit.Geolocation.City, "Wroclaw", t)
	eq(p.Payload.Chat.Customers[cid].LastVisit.Geolocation.Timezone, "test_timezone", t)

	eq(len(p.Payload.Chat.Customers[cid].LastVisit.LastPages), 1, t)
	eq(p.Payload.Chat.Customers[cid].LastVisit.LastPages[0].OpenedAt.String(), "2019-10-11 09:40:56.071345 +0000 UTC", t)
	eq(p.Payload.Chat.Customers[cid].LastVisit.LastPages[0].URL, "https://cdn.livechatinc.com/labs/?license=100007977/", t)
	eq(p.Payload.Chat.Customers[cid].LastVisit.LastPages[0].Title, "LiveChat", t)

	eq(p.Payload.Chat.Customers[cid].Statistics.VisitsCount, 29, t)
	eq(p.Payload.Chat.Customers[cid].Statistics.ThreadsCount, 18, t)
	eq(p.Payload.Chat.Customers[cid].Statistics.ChatsCount, 1, t)
	eq(p.Payload.Chat.Customers[cid].Statistics.PageViewsCount, 5, t)
	eq(p.Payload.Chat.Customers[cid].Statistics.GreetingsShownCount, 6, t)
	eq(p.Payload.Chat.Customers[cid].Statistics.GreetingsAcceptedCount, 8, t)

	eq(p.Payload.Chat.Customers[cid].AgentLastEventCreatedAt.String(), "2019-10-11 09:40:59.249 +0000 UTC", t)
	eq(p.Payload.Chat.Customers[cid].CustomerLastEventCreatedAt.String(), "2019-10-11 09:40:59.219001 +0000 UTC", t)

	eq(len(p.Payload.Chat.Agents), 1, t)
	aid := "l.wojciechowski@livechatinc.com"
	eq(p.Payload.Chat.Agents[aid].ID, "l.wojciechowski@livechatinc.com", t)
	eq(p.Payload.Chat.Agents[aid].Type, "agent", t)
	eq(p.Payload.Chat.Agents[aid].Name, "Łukasz Wojciechowski", t)
	eq(p.Payload.Chat.Agents[aid].Email, "l.wojciechowski@livechatinc.com", t)
	eq(p.Payload.Chat.Agents[aid].Avatar, "livechat.s3.amazonaws.com/default/avatars/a14.png", t)
	eq(p.Payload.Chat.Agents[aid].Present, true, t)
	eq(p.Payload.Chat.Agents[aid].LastSeen.String(), "1970-01-01 01:00:00 +0100 CET", t)
	eq(p.Payload.Chat.Agents[aid].RoutingStatus, "accepting_chats", t)

	eq(len(p.Payload.Chat.Threads), 1, t)
	eq(p.Payload.Chat.Threads[0].ID, "PZ070E0W1B", t)
	eq(p.Payload.Chat.Threads[0].Timestamp.String(), "2019-10-11 11:40:59 +0200 CEST", t)
	eq(p.Payload.Chat.Threads[0].Active, true, t)
	eq(p.Payload.Chat.Threads[0].UserIDs[0], "345f8235-d60d-433e-63c5-7f813a6ffe25", t)
	eq(p.Payload.Chat.Threads[0].UserIDs[1], "l.wojciechowski@livechatinc.com", t)
	eq(p.Payload.Chat.Threads[0].RestrictedAccess, false, t)
	eq(p.Payload.Chat.Threads[0].Order, 18, t)
	eq(p.Payload.Chat.Threads[0].Properties["routing"]["continuous"], false, t)
	eq(p.Payload.Chat.Threads[0].Properties["routing"]["idle"], false, t)
	eq(p.Payload.Chat.Threads[0].Properties["routing"]["referrer"], "", t)
	eq(p.Payload.Chat.Threads[0].Properties["routing"]["start_url"], "https://cdn.livechatinc.com/labs/?license=100007977/", t)
	eq(p.Payload.Chat.Threads[0].Properties["routing"]["unassigned"], false, t)
	eq(p.Payload.Chat.Threads[0].Access.GroupIDs[0], 0, t)
	eq(len(p.Payload.Chat.Threads[0].Events), 2, t)

	e := p.Payload.Chat.Threads[0].Events[0].FilledForm()
	eq(e.ID, "PZ070E0W1B_1", t)
	eq(e.CreatedAt.String(), "2019-10-11 09:40:59.219001 +0000 UTC", t)
	eq(e.AuthorID, "345f8235-d60d-433e-63c5-7f813a6ffe25", t)
	eq(e.Recipients, "all", t)
	eq(e.Properties["lc2"]["form_type"], "prechat", t)
	eq(e.Type, "filled_form", t)
	eq(len(e.Fields), 2, t)
	eq(e.Fields[0].Label, "Name:", t)
	eq(e.Fields[0].Type, "name", t)
	eq(e.Fields[0].Value, "rewrew", t)
	eq(e.Fields[1].Label, "E-mail:", t)
	eq(e.Fields[1].Type, "email", t)
	eq(e.Fields[1].Value, "", t)

	e1 := p.Payload.Chat.Threads[0].Events[1].Message()
	eq(e1.ID, "PZ070E0W1B_2", t)
	eq(e1.Text, "Hello. How may I help you?", t)
}

func BenchmarkParseIncomingChatThreadPayload(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = webhooks.ParseIncomingChatThreadPayload([]byte(incomingChatThread))
	}
}
