package agent

// PropertiesFilters represents set of filters for Chat properties
type PropertiesFilters map[string]map[string]*propertyFilterType

type propertyFilterType struct {
	Exists            *bool         `json:"exists,omitempty"`
	Values            []interface{} `json:"values,omitempty"`
	ExcludeValues     []interface{} `json:"exclude_values,omitempty"`
	RequireEveryValue *bool         `json:"require_every_value,omitempty"`
}

// NewPropertyFilterType creates new filter object for Chat properties
// If the first parameter is passed along with nil values then the last parameter will be ignored and the filter will check only existence of property
// Otherwise it will check if property values match/exclude given values based on the first parameter
// The last parameter modifies the filter behavior so that it matches only those Chats that have or don't have all values in the property the filter relates to
func NewPropertyFilterType(includes bool, vals []interface{}, requireEveryValue bool) *propertyFilterType {
	pft := &propertyFilterType{}
	switch {
	case vals == nil:
		pft.Exists = &includes
	case includes:
		pft.Values = vals
		pft.RequireEveryValue = &requireEveryValue
	case !includes:
		pft.ExcludeValues = vals
		pft.RequireEveryValue = &requireEveryValue
	}
	return pft
}

// Archives filters

type archivesFilters struct {
	Agents        *propertyFilterType `json:"agents,omitempty"`
	AgentTypes    *AgentTypesFilter   `json:"agent_types,omitempty"`
	GroupIDs      []uint              `json:"group_ids,omitempty"`
	From          string              `json:"from,omitempty"`
	To            string              `json:"to,omitempty"`
	Properties    PropertiesFilters   `json:"properties,omitempty"`
	Tags          *propertyFilterType `json:"tags,omitempty"`
	Sales         *propertyFilterType `json:"sales,omitempty"`
	Goals         *propertyFilterType `json:"goals,omitempty"`
	Surveys       []SurveyFilter      `json:"surveys,omitempty"`
	ThreadIDs     []string            `json:"thread_ids,omitempty"`
	Query         string              `json:"query,omitempty"`
	EventTypes    *eventTypesFilter   `json:"event_types,omitempty"`
	Greetings     *GreetingsFilter    `json:"greetings,omitempty"`
	CustomerID    string              `json:"customer_id,omitempty"`
	CustomerEmail string              `json:"customer_email,omitempty"`
}

type eventTypesFilter struct {
	Values            []string `json:"values,omitempty"`
	ExcludeValues     []string `json:"exclude_values,omitempty"`
	RequireEveryValue *bool    `json:"require_every_value,omitempty"`
}

// AgentTypesFilter represents structure to match agent types when getting Archives
type AgentTypesFilter struct {
	AllValues     []string `json:"all_values,omitempty"`
	AnyValues     []string `json:"any_values,omitempty"`
	ExcludeValues []string `json:"exclude_values,omitempty"`
}

// SurveyFilter represents structure to match surveys when getting Archives
type SurveyFilter struct {
	Type     string `json:"type"`
	AnswerID string `json:"answer_id"`
}

// GreetingsFilter represents structure to match greetings when getting Archives
type GreetingsFilter struct {
	From          string         `json:"from,omitempty"`
	To            string         `json:"to,omitempty"`
	Values        []int64        `json:"values,omitempty"`
	ExcludeValues []int64        `json:"exclude_values,omitempty"`
	Exists        *bool          `json:"exists,omitempty"`
	Groups        *integerFilter `json:"groups,omitempty"`
}

// NewArchivesFilters creates empty structure to aggregate filters for ListArchives method
func NewArchivesFilters() *archivesFilters {
	return &archivesFilters{}
}

// ByAgents extends archives filter with agents specific property filter
// See NewPropertyFilterType definition for details of filter creation
func (af *archivesFilters) ByAgents(includes bool, vals []interface{}, requireEveryValue bool) *archivesFilters {
	af.Agents = NewPropertyFilterType(includes, vals, requireEveryValue)
	return af
}

// ByGroups extends archives filter with list of group IDs to match
func (af *archivesFilters) ByGroups(groupIDs []uint) *archivesFilters {
	af.GroupIDs = groupIDs
	return af
}

// ByThreads extends archives filter with list of thread IDs to match
// This method clears previously set filters as this type of filter cannot be used in combination with others
func (af *archivesFilters) ByThreads(threadIDs []string) *archivesFilters {
	*af = archivesFilters{
		ThreadIDs: threadIDs,
	}
	return af
}

// ByQuery extends archives filter with query to match
func (af *archivesFilters) ByQuery(query string) *archivesFilters {
	af.Query = query
	return af
}

// FromDate extends archives filter to exclude entries before given date
func (af *archivesFilters) FromDate(date string) *archivesFilters {
	af.From = date
	return af
}

// FromDate extends archives filter to exclude entries after given date
func (af *archivesFilters) ToDate(date string) *archivesFilters {
	af.To = date
	return af
}

// ByProperties extends archives filter with Chat properties to match
func (af *archivesFilters) ByProperties(propsFilters PropertiesFilters) *archivesFilters {
	af.Properties = propsFilters
	return af
}

// BySurveys extends archives filter with surveys to match
func (af *archivesFilters) BySurveys(surveyFilters []SurveyFilter) *archivesFilters {
	af.Surveys = surveyFilters
	return af
}

// ByTags extends archives filter with tags specific property filter
// See NewPropertyFilterType definition for details of filter creation
func (af *archivesFilters) ByTags(includes bool, vals []interface{}, requireEveryValue bool) *archivesFilters {
	af.Tags = NewPropertyFilterType(includes, vals, requireEveryValue)
	return af
}

// BySales extends archives filter with sales specific property filter
// See NewPropertyFilterType definition for details of filter creation
func (af *archivesFilters) BySales(includes bool, vals []interface{}, requireEveryValue bool) *archivesFilters {
	af.Sales = NewPropertyFilterType(includes, vals, requireEveryValue)
	return af
}

// ByGoals extends archives filter with goals specific property filter
// See NewPropertyFilterType definition for details of filter creation
func (af *archivesFilters) ByGoals(includes bool, vals []interface{}, requireEveryValue bool) *archivesFilters {
	af.Goals = NewPropertyFilterType(includes, vals, requireEveryValue)
	return af
}

// ByEventTypes extends archives filter with event_types.values to match if first parameter true
// Otherwise it extends archives filter with event_types.exclude_values
func (af *archivesFilters) ByEventTypes(includes bool, vals []string, requireEveryValue bool) *archivesFilters {
	if includes {
		af.EventTypes = &eventTypesFilter{Values: vals}
	} else {
		af.EventTypes = &eventTypesFilter{ExcludeValues: vals}
	}

	af.EventTypes.RequireEveryValue = &requireEveryValue
	return af
}

// ByGreetings extends archives filter with greetings to match
func (af *archivesFilters) ByGreetings(filters *GreetingsFilter) *archivesFilters {
	af.Greetings = filters
	return af
}

// ByAgentTypes extends archives filter with agent types to match
func (af *archivesFilters) ByAgentTypes(filter *AgentTypesFilter) *archivesFilters {
	af.AgentTypes = filter
	return af
}

type integerFilter struct {
	Values        []int64 `json:"values,omitempty"`
	ExcludeValues []int64 `json:"exclude_values,omitempty"`
}

// NewIntegerFilter creates new filter for integer values
// `inclusive` parameter controls if the filtered values should match or exclude given values
func NewIntegerFilter(values []int64, inclusive bool) *integerFilter {
	intF := &integerFilter{}
	switch {
	case inclusive:
		intF.Values = values
	default:
		intF.ExcludeValues = values
	}
	return intF
}

// Chats Filters
type chatsFilters struct {
	IncludeActive              bool              `json:"include_active,omitempty"`
	IncludeChatsWithoutThreads bool              `json:"include_chats_without_threads,omitempty"`
	GroupIDs                   []uint            `json:"group_ids,omitempty"`
	Properties                 PropertiesFilters `json:"properties,omitempty"`
}

// NewChatsFilters creates empty structure to aggregate filters for Chats in ListChats method
// By default filters include also active chats
func NewChatsFilters() *chatsFilters {
	return &chatsFilters{
		IncludeActive: true,
	}
}

// WithoutActiveChats extends chat filters to not include active chats
func (cf *chatsFilters) WithoutActiveChats() *chatsFilters {
	cf.IncludeActive = false
	return cf
}

// WithChatsWithoutThreads extends chat filters to include chats without threads
func (cf *chatsFilters) WithChatsWithoutThreads() *chatsFilters {
	cf.IncludeChatsWithoutThreads = true
	return cf
}

// ByGroups extends chat filters with group IDs to match
func (cf *chatsFilters) ByGroups(groupIDs []uint) *chatsFilters {
	cf.GroupIDs = groupIDs
	return cf
}

// ByProperties extends chat filters with Chat properties to match
func (cf *chatsFilters) ByProperties(propsFilters PropertiesFilters) *chatsFilters {
	cf.Properties = propsFilters
	return cf
}

// Threads Filters
type threadsFilters struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}

// NewThreadsFilters creates empty structure to aggregate filters for Threads in ListThreads method
func NewThreadsFilters() *threadsFilters {
	return &threadsFilters{}
}

// FromDate extends threads filter to exclude entries before given date
func (tf *threadsFilters) FromDate(date string) *threadsFilters {
	tf.From = date
	return tf
}

// FromDate extends threads filter to exclude entries after given date
func (tf *threadsFilters) ToDate(date string) *threadsFilters {
	tf.To = date
	return tf
}
