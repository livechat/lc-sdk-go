package configuration

type BotStatus string

const (
	AcceptingChats    BotStatus = "accepting chats"
	NotAcceptingChats BotStatus = "not accepting chats"
	Offline           BotStatus = "offline"
)

type WebhookAction string

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
)

type GroupPriority string

const (
	First       GroupPriority = "first"
	Normal      GroupPriority = "normal"
	Last        GroupPriority = "last"
	DoNotAssign GroupPriority = "supervisor"
)
