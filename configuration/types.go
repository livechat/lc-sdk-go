package configuration

// BotStatus represents bot availability status
type BotStatus string

const (
	// AcceptingChats - Bot will be accepting chats
	AcceptingChats BotStatus = "accepting chats"
	// NotAcceptingChats - Bot will not accept any chats, yet it will still occupy seat
	NotAcceptingChats BotStatus = "not accepting chats"
	// Offline - Bot will not accept chats and will not use up seats
	Offline BotStatus = "offline"
)

// WebhookAction represents allowed values for action name
type WebhookAction string

// Following Webhook actions are supported
const (
	IncomingChatThread          WebhookAction = "incoming_chat_thread"
	IncomingEvent               WebhookAction = "incoming_event"
	IncomingRichMessagePostback WebhookAction = "incoming_rich_message_postback"
	LastSeenTimestampUpdated    WebhookAction = "last_seen_timestamp_updated"
	ThreadClosed                WebhookAction = "thread_closed"
	ChatPropertiesUpdated       WebhookAction = "chat_properties_updated"
	ChatPropertiesDeleted       WebhookAction = "chat_properties_deleted"
	ChatThreadPropertiesUpdated WebhookAction = "chat_thread_properties_updated"
	ChatThreadPropertiesDeleted WebhookAction = "chat_thread_properties_deleted"
	ChatUserAdded               WebhookAction = "chat_user_added"
	ChatUserRemoved             WebhookAction = "chat_user_removed"
	ChatThreadTagged            WebhookAction = "chat_thread_tagged"
	ChatThreadUntagged          WebhookAction = "chat_thread_untagged"
	AgentStatusChanged          WebhookAction = "agent_status_changed"
	AgentDeleted                WebhookAction = "agent_deleted"
	EventsMarkedAsSeen          WebhookAction = "events_marked_as_seen"
	AccessGranted               WebhookAction = "access_granted"
	AccessRevoked               WebhookAction = "access_revoked"
	AccessSet                   WebhookAction = "access_set"
	CustomerCreated             WebhookAction = "customer_created"
)

// GroupPriority represents priority of assigning chats in group
type GroupPriority string

const (
	// First - The highest chat routing priority. Agents with the first priority get chats before others from the same group, e.g. Bots can get chats before regular Agents.
	First GroupPriority = "first"
	// Normal - The medium chat routing priority. Agents with the normal priority get chats before those with the last priority, when there are no Agents with the first priority available with free slots in the group.
	Normal GroupPriority = "normal"
	// Last - The lowest chat routing priority. Agents with the last priority get chats when there are no Agents with the first or normal priority available with free slots in the group.
	Last GroupPriority = "last"
	// DoNotAssign - Bot will not be assigned to any chats. This can be used only in `default_group_priority`.
	DoNotAssign GroupPriority = "supervisor"
)
