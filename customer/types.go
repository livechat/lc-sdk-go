package customer

type GroupStatus int

const (
	GroupStatusUnknown GroupStatus = iota
	GroupStatusOnline
	GroupStatusOffline
	GroupStatusOnlineForQueue
)

type FormType string

const (
	FormTypePrechat  = "prechat"
	FormTypePostchat = "postchat"
	FormTypeTicket   = "ticket"
	FormTypeEmail    = "email"
)

type Recipients string

const (
	All    = "all"
	Agents = "agents"
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
