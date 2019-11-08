package action

type Webhook string

const (
	IncomingChatThread          Webhook = "incoming_chat_thread"
	IncomingEvent               Webhook = "incoming_event"
	IncomingRichMessagePostback Webhook = "incoming_rich_message_postback"
	LastSeenTimestampUpdated    Webhook = "last_seen_timestamp_updated"
	ThreadClosed                Webhook = "thread_closed"
	ChatPropertiesUpdated       Webhook = "chat_properties_updated"
	ChatPropertiesDeleted       Webhook = "chat_properties_deleted"
	ChatThreadPropertiesUpdated Webhook = "chat_thread_properties_updated"
	ChatThreadPropertiesDeleted Webhook = "chat_thread_properties_deleted"
	ChatUserAdded               Webhook = "chat_user_added"
	ChatUserRemoved             Webhook = "chat_user_removed"
	ChatThreadTagged            Webhook = "chat_thread_tagged"
	ChatThreadUntagged          Webhook = "chat_thread_untagged"
	AgentStatusChanged          Webhook = "agent_status_changed"
	AgentDeleted                Webhook = "agent_deleted"
)
