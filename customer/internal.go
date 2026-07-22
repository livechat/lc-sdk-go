package customer

type startChatRequest struct {
	Chat       *InitialChat `json:"chat,omitempty"`
	Continuous bool         `json:"continuous,omitempty"`
	Active     bool         `json:"active"`
}

type startChatResponse struct {
	ChatID   string   `json:"chat_id"`
	ThreadID string   `json:"thread_id"`
	EventIDs []string `json:"event_ids"`
}

type sendEventRequest struct {
	ChatID             string      `json:"chat_id"`
	Event              interface{} `json:"event"`
	AttachToLastThread *bool       `json:"attach_to_last_thread,omitempty"`
}

type sendEventResponse struct {
	EventID string `json:"event_id"`
}

type resumeChatRequest struct {
	Chat       *InitialChat `json:"chat"`
	Continuous bool         `json:"continuous,omitempty"`
	Active     bool         `json:"active"`
}

type resumeChatResponse struct {
	ThreadID string   `json:"thread_id"`
	EventIDs []string `json:"event_ids"`
}

type listChatsRequest struct {
	*hashedPaginationRequest
}

type listChatsResponse struct {
	hashedPaginationResponse
	ChatsInfo  []ChatInfo `json:"chats_info"`
	TotalChats uint       `json:"total_chats"`
}

type getChatRequest struct {
	ChatID   string `json:"chat_id"`
	ThreadID string `json:"thread_id,omitempty"`
}

type listThreadsRequest struct {
	*hashedPaginationRequest
	ChatID         string `json:"chat_id"`
	MinEventsCount uint   `json:"min_events_count,omitempty"`
}

type listThreadsResponse struct {
	hashedPaginationResponse
	Threads      []Thread `json:"threads"`
	FoundThreads uint     `json:"found_threads"`
}

type deactivateChatRequest struct {
	ID string `json:"id"`
}

type sendRichMessagePostbackRequest struct {
	ChatID   string   `json:"chat_id"`
	ThreadID string   `json:"thread_id"`
	EventID  string   `json:"event_id"`
	Postback postback `json:"postback"`
}

type postback struct {
	ID          string             `json:"id"`
	Toggled     bool               `json:"toggled"`
	ButtonType  string             `json:"button_type"`
	ButtonValue string             `json:"button_value"`
	Ecommerce   *PostbackEcommerce `json:"ecommerce,omitempty"`
}

type sendSneakPeekRequest struct {
	ChatID        string `json:"chat_id"`
	SneakPeekText string `json:"sneak_peek_text"`
}

type updateChatPropertiesRequest struct {
	ID         string     `json:"id"`
	Properties Properties `json:"properties"`
}

type deleteChatPropertiesRequest struct {
	ID         string              `json:"id"`
	Properties map[string][]string `json:"properties"`
}

type updateThreadPropertiesRequest struct {
	ChatID     string     `json:"chat_id"`
	ThreadID   string     `json:"thread_id"`
	Properties Properties `json:"properties"`
}

type deleteThreadPropertiesRequest struct {
	ChatID     string              `json:"chat_id"`
	ThreadID   string              `json:"thread_id"`
	Properties map[string][]string `json:"properties"`
}

type updateEventPropertiesRequest struct {
	ChatID     string     `json:"chat_id"`
	ThreadID   string     `json:"thread_id"`
	EventID    string     `json:"event_id"`
	Properties Properties `json:"properties"`
}

type deleteEventPropertiesRequest struct {
	ChatID     string              `json:"chat_id"`
	ThreadID   string              `json:"thread_id"`
	EventID    string              `json:"event_id"`
	Properties map[string][]string `json:"properties"`
}

type deleteEventRequest struct {
	ChatID   string `json:"chat_id"`
	ThreadID string `json:"thread_id"`
	EventID  string `json:"event_id"`
}

type updateCustomerRequest struct {
	Name          *string             `json:"name,omitempty"`
	NameIsDefault *bool               `json:"name_is_default,omitempty"`
	Email         *string             `json:"email,omitempty"`
	Avatar        *string             `json:"avatar,omitempty"`
	PhoneNumber   *string             `json:"phone_number,omitempty"`
	SessionFields []map[string]string `json:"session_fields,omitempty"`
	Address       *AddressUpdate      `json:"address,omitempty"`
}

type setCustomerSessionFieldsRequest struct {
	SessionFields []map[string]string `json:"session_fields"`
}

type listGroupStatusesRequest struct {
	All      bool  `json:"all,omitempty"`
	GroupIDs []int `json:"group_ids,omitempty"`
}

type listGroupStatusesResponse struct {
	Status map[int]string `json:"groups_status"`
}

type checkGoalsRequest struct {
	PageURL        string            `json:"page_url"`
	GroupID        int               `json:"group_id"`
	CustomerFields map[string]string `json:"customer_fields"`
}

type getFormRequest struct {
	GroupID int    `json:"group_id"`
	Type    string `json:"type"`
}

type getFormResponse struct {
	Form    *Form `json:"form"`
	Enabled bool  `json:"enabled"`
}

type getURLInfoRequest struct {
	URL string `json:"url"`
}

type markEventsAsSeenRequest struct {
	ChatID   string `json:"chat_id"`
	SeenUpTo string `json:"seen_up_to"`
}

type emptyResponse struct{}

type acceptGreetingRequest struct {
	GreetingID int    `json:"greeting_id"`
	UniqueID   string `json:"unique_id"`
}

type cancelGreetingRequest struct {
	UniqueID string `json:"unique_id"`
}

type hashedPaginationRequest struct {
	PageID    string `json:"page_id,omitempty"`
	Limit     uint   `json:"limit,omitempty"`
	SortOrder string `json:"sort_order,omitempty"`
}

type hashedPaginationResponse struct {
	PreviousPageID string `json:"previous_page_id,omitempty"`
	NextPageID     string `json:"next_page_id,omitempty"`
}

type requestEmailVerificationRequest struct {
	CallbackURI string `json:"callback_uri"`
}

type requestWelcomeMessageRequest struct {
	ID      string `json:"id,omitempty"`
	GroupID *int   `json:"group_id,omitempty"`
}

type requestWelcomeMessageResponse struct {
	ID             string          `json:"id"`
	PredictedAgent *PredictedAgent `json:"predicted_agent"`
	Queue          bool            `json:"queue"`
}
