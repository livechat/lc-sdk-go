package configuration

// WebhookAction represents allowed values for action name
type WebhookAction string

// Following Webhook actions are supported
const (
	IncomingChat                 WebhookAction = "incoming_chat"
	IncomingEvent                WebhookAction = "incoming_event"
	EventUpdated                 WebhookAction = "event_updated"
	IncomingRichMessagePostback  WebhookAction = "incoming_rich_message_postback"
	ChatDeactivated              WebhookAction = "chat_deactivated"
	ChatPropertiesUpdated        WebhookAction = "chat_properties_updated"
	ThreadPropertiesUpdated      WebhookAction = "thread_properties_updated"
	ChatPropertiesDeleted        WebhookAction = "chat_properties_deleted"
	ThreadPropertiesDeleted      WebhookAction = "thread_properties_deleted"
	UserAddedToChat              WebhookAction = "user_added_to_chat"
	UserRemovedFromChat          WebhookAction = "user_removed_from_chat"
	ThreadTagged                 WebhookAction = "thread_tagged"
	ThreadUntagged               WebhookAction = "thread_untagged"
	AgentCreated                 WebhookAction = "agent_created"
	AgentUpdated                 WebhookAction = "agent_updated"
	AgentDeleted                 WebhookAction = "agent_deleted"
	AgentSuspended               WebhookAction = "agent_suspended"
	AgentUnsuspended             WebhookAction = "agent_unsuspended"
	AgentApproved                WebhookAction = "agent_approved"
	EventsMarkedAsSeen           WebhookAction = "events_marked_as_seen"
	ChatAccessUpdated            WebhookAction = "chat_access_updated"
	EventPropertiesUpdated       WebhookAction = "event_properties_updated"
	EventPropertiesDeleted       WebhookAction = "event_properties_deleted"
	RoutingStatusSet             WebhookAction = "routing_status_set"
	ChatTransferred              WebhookAction = "chat_transferred"
	IncomingCustomer             WebhookAction = "incoming_customer"
	CustomerSessionFieldsUpdated WebhookAction = "customer_session_fields_updated"
	GroupCreated                 WebhookAction = "group_created"
	GroupUpdated                 WebhookAction = "group_updated"
	GroupDeleted                 WebhookAction = "group_deleted"
	AutoAccessAdded              WebhookAction = "auto_access_added"
	AutoAccessUpdated            WebhookAction = "auto_access_updated"
	AutoAccessDeleted            WebhookAction = "auto_access_deleted"
	BotCreated                   WebhookAction = "bot_created"
	BotUpdated                   WebhookAction = "bot_updated"
	BotDeleted                   WebhookAction = "bot_deleted"
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
