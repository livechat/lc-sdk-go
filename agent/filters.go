package agent

// PropertiesFilters represents set of filters for Chat properties
type PropertiesFilters map[string]map[string]*propertyFilterType

type propertyFilterType struct {
	Exists        *bool         `json:"exists,omitempty"`
	Values        []interface{} `json:"values,omitempty"`
	ExcludeValues []interface{} `json:"exclude_values,omitempty"`
}

// NewPropertyFilterType creates new filter object for Chat properties
// If only first parameter is passed, filter will check only existence of property
// Passing additional parameters will check if property values match/exclude given values
// based on the first parameter
func NewPropertyFilterType(includes bool, vals ...interface{}) *propertyFilterType {
	pft := &propertyFilterType{}
	switch {
	case vals == nil:
		pft.Exists = &includes
	case includes:
		pft.Values = vals
	case !includes:
		pft.ExcludeValues = vals
	}
	return pft
}

// Archives filters

type archivesFilters struct {
	Agents     *propertyFilterType `json:"agents,omitempty"`
	GroupIDs   []uint              `json:"group_ids,omitempty"`
	From       string              `json:"from,omitempty"`
	To         string              `json:"to,omitempty"`
	Properties PropertiesFilters   `json:"properties,omitempty"`
	Tags       *propertyFilterType `json:"tags,omitempty"`
	Sales      *propertyFilterType `json:"sales,omitempty"`
	Goals      *propertyFilterType `json:"goals,omitempty"`
	Surveys    []SurveyFilter      `json:"surveys,omitempty"`
	ThreadIDs  []string            `json:"thread_ids,omitempty"`
	Query      string              `json:"query,omitempty"`
	Events     *eventsFilter       `json:"events,omitempty"`
}

type eventsFilter struct {
	Types []string `json:"types,omitempty"`
}

// SurveyFilter represents structure to match surveys when getting Archives
type SurveyFilter struct {
	Type     string `json:"type"`
	AnswerID string `json:"answer_id"`
}

// NewArchivesFilters creates empty structure to aggregate filters for ListArchives method
func NewArchivesFilters() *archivesFilters {
	return &archivesFilters{}
}

// ByAgents extends archives filter with agents specific property filter
// See NewPropertyFilterType definition for details of filter creation
func (af *archivesFilters) ByAgents(includes bool, vals ...interface{}) *archivesFilters {
	af.Agents = NewPropertyFilterType(includes, vals...)
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
func (af *archivesFilters) ByTags(includes bool, vals ...interface{}) *archivesFilters {
	af.Tags = NewPropertyFilterType(includes, vals...)
	return af
}

// BySales extends archives filter with sales specific property filter
// See NewPropertyFilterType definition for details of filter creation
func (af *archivesFilters) BySales(includes bool, vals ...interface{}) *archivesFilters {
	af.Sales = NewPropertyFilterType(includes, vals...)
	return af
}

// ByGoals extends archives filter with goals specific property filter
// See NewPropertyFilterType definition for details of filter creation
func (af *archivesFilters) ByGoals(includes bool, vals ...interface{}) *archivesFilters {
	af.Goals = NewPropertyFilterType(includes, vals...)
	return af
}

// ByEventTypes extends archives filter with event.types to match
func (af *archivesFilters) ByEventTypes(types ...string) *archivesFilters {
	af.Events = &eventsFilter{Types: types}
	return af
}

// Customer filters

type customersFilters struct {
	Country                    *stringFilter    `json:"country,omitempty"`
	Email                      *stringFilter    `json:"email,omitempty"`
	Name                       *stringFilter    `json:"name,omitempty"`
	CustomerID                 *stringFilter    `json:"customer_id,omitempty"`
	ChatsCount                 *RangeFilter     `json:"chats_count,omitempty"`
	ThreadsCount               *RangeFilter     `json:"threads_count,omitempty"`
	VisitsCount                *RangeFilter     `json:"visits_count,omitempty"`
	CreatedAt                  *DateRangeFilter `json:"created_at,omitempty"`
	AgentLastEventCreatedAt    *DateRangeFilter `json:"agent_last_event_created_at,omitempty"`
	CustomerLastEventCreatedAt *DateRangeFilter `json:"customer_last_event_created_at,omitempty"`
}

type stringFilter struct {
	Values        []string `json:"values,omitempty"`
	ExcludeValues []string `json:"exclude_values,omitempty"`
}

// NewStringFilter creates new filter for string values
// `inclusive` parameter controls if the filtered values should match or exclude given values
func NewStringFilter(values []string, inclusive bool) *stringFilter {
	sf := &stringFilter{}
	switch {
	case inclusive:
		sf.Values = values
	default:
		sf.ExcludeValues = values
	}
	return sf
}

// RangeFilter represents structure to define a range in which filtered numbers should be matched
// LTE - less than or equal
// LT  - less than
// GTE - greater than or equal
// GT  - greater than
// EQ  - equal
type RangeFilter struct {
	LTE int64 `json:"lte,omitempty"`
	LT  int64 `json:"lt,omitempty"`
	GTE int64 `json:"gte,omitempty"`
	GT  int64 `json:"gt,omitempty"`
	EQ  int64 `json:"eq,omitempty"`
}

// DateRangeFilter represents structure to define a range in which filtered dates should be matched
// Dates are represented in ISO 8601 format with microseconds resolution, e.g. 2017-10-12T15:19:21.010200+01:00 in specific timezone or 2017-10-12T14:19:21.010200Z in UTC.
// LTE - less than or equal
// LT  - less than
// GTE - greater than or equal
// GT  - greater than
// EQ  - equal
type DateRangeFilter struct {
	LTE string `json:"lte,omitempty"`
	LT  string `json:"lt,omitempty"`
	GTE string `json:"gte,omitempty"`
	GT  string `json:"gt,omitempty"`
	EQ  string `json:"eq,omitempty"`
}

// NewCustomersFilters creates empty structure to aggregate filters for customers in ListCustomers method
func NewCustomersFilters() *customersFilters {
	return &customersFilters{}
}

// ByCountry extends customers filters with string filter for customer's country
// See NewStringFilter definition for details of filter creation
func (cf *customersFilters) ByCountry(values []string, inclusive bool) *customersFilters {
	cf.Country = NewStringFilter(values, inclusive)
	return cf
}

// ByEmail extends customers filters with string filter for customer's email
// See NewStringFilter definition for details of filter creation
func (cf *customersFilters) ByEmail(values []string, inclusive bool) *customersFilters {
	cf.Email = NewStringFilter(values, inclusive)
	return cf
}

// ByName extends customers filters with string filter for customer's name
// See NewStringFilter definition for details of filter creation
func (cf *customersFilters) ByName(values []string, inclusive bool) *customersFilters {
	cf.Name = NewStringFilter(values, inclusive)
	return cf
}

// ByID extends customers filters with string filter for customer's ID
// See NewStringFilter definition for details of filter creation
func (cf *customersFilters) ByID(values []string, inclusive bool) *customersFilters {
	cf.CustomerID = NewStringFilter(values, inclusive)
	return cf
}

// ByChatsCount extends customers filters with range filter for customer's chats count
// See RangeFilter definition for details of filter creation
func (cf *customersFilters) ByChatsCount(ranges *RangeFilter) *customersFilters {
	cf.ChatsCount = ranges
	return cf
}

// ByThreadsCount extends customers filters with range filter for customer's threads count
// See RangeFilter definition for details of filter creation
func (cf *customersFilters) ByThreadsCount(ranges *RangeFilter) *customersFilters {
	cf.ThreadsCount = ranges
	return cf
}

// ByVisitsCount extends customers filters with range filter for customer's visits count
// See RangeFilter definition for details of filter creation
func (cf *customersFilters) ByVisitsCount(ranges *RangeFilter) *customersFilters {
	cf.VisitsCount = ranges
	return cf
}

// ByCreationTime extends customers filters with range filter for customer's creation time
// See DateRangeFilter definition for details of filter creation
func (cf *customersFilters) ByCreationTime(timeRange *DateRangeFilter) *customersFilters {
	cf.CreatedAt = timeRange
	return cf
}

// ByAgentsLastActivity extends customers filters with range filter for last agent's activity with customer
// See DateRangeFilter definition for details of filter creation
func (cf *customersFilters) ByAgentsLastActivity(timeRange *DateRangeFilter) *customersFilters {
	cf.AgentLastEventCreatedAt = timeRange
	return cf
}

// ByCustomersLastActivity extends customers filters with range filter for customer's last activity
// See DateRangeFilter definition for details of filter creation
func (cf *customersFilters) ByCustomersLastActivity(timeRange *DateRangeFilter) *customersFilters {
	cf.CustomerLastEventCreatedAt = timeRange
	return cf
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
