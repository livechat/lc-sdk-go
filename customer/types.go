package customer

type GroupStatus int

const (
	GroupStatusUnknown GroupStatus = iota
	GroupStatusOnline
	GroupStatusOffline
	GroupStatusOnlineForQueue
)

type FormType int

const (
	FormTypePrechat = iota
	FormTypePostchat
	FormTypeTicket
	FormTypeEmail
)
