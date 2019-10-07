package customer

// GroupStatus represents status of groups.
type GroupStatus int

// Possible values of GroupStatus.
const (
	// GroupStatusUnknown should never be returned by API Client.
	GroupStatusUnknown GroupStatus = iota
	GroupStatusOnline
	GroupStatusOffline
	GroupStatusOnlineForQueue
)

// FormType represents type of form templates.
type FormType string

// Possible values of FormType.
const (
	FormTypePrechat  FormType = "prechat"
	FormTypePostchat FormType = "postchat"
	FormTypeTicket   FormType = "ticket"
	FormTypeEmail    FormType = "email"
)

// Recipients represents possible event recipients.
type Recipients string

// Possible values of Recipients.
const (
	All    Recipients = "all"
	Agents Recipients = "agents"
)

func toGroupStatus(s string) GroupStatus {
	switch s {
	case "online":
		return GroupStatusOnline
	case "offline":
		return GroupStatusOffline
	case "online_for_queue":
		return GroupStatusOnlineForQueue
	default:
		return GroupStatusUnknown
	}
}
