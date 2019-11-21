package agent

type PropertiesFilters map[string]map[string]*PropertyFilterType

type PropertyFilterType struct {
	Exists        *bool         `json:"exists,omitempty"`
	Values        []interface{} `json:"values,omitempty"`
	ExcludeValues []interface{} `json:"exclude_values,omitempty"`
}

func NewPropertyFilterType(includes bool, vals ...interface{}) PropertyFilterType {
	pft := PropertyFilterType{}
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

type ArchivesFilters struct {
	AgentIDs   []string           `json:"agent_ids,omitempty"`
	GroupIDs   []uint             `json:"group_ids,omitempty"`
	DateFrom   string             `json:"date_from,omitempty"`
	DateTo     string             `json:"date_to,omitempty"`
	Properties PropertiesFilters  `json:"properties,omitempty"`
	Tags       PropertyFilterType `json:"tags,omitempty"`
	Sales      PropertyFilterType `json:"sales,omitempty"`
	Goals      PropertyFilterType `json:"goals,omitempty"`
	Surveys    []SurveyFilter     `json:"surveys,omitempty"`
	ThreadIDs  []string           `json:"thread_ids,omitempty"`
	Query      string             `json:"query,omitempty"`
}

type SurveyFilter struct {
	Type     string `json:"type"`
	AnswerID string `json:"answer_id"`
}

func NewArchivesFilters() *ArchivesFilters {
	return &ArchivesFilters{}
}

func (af *ArchivesFilters) ByAgents(agentIDs []string) *ArchivesFilters {
	af.AgentIDs = agentIDs
	return af
}

func (af *ArchivesFilters) ByGroups(groupIDs []uint) *ArchivesFilters {
	af.GroupIDs = groupIDs
	return af
}

func (af *ArchivesFilters) ByThreads(threadIDs []string) *ArchivesFilters {
	*af = ArchivesFilters{
		ThreadIDs: threadIDs,
	}
	return af
}

func (af *ArchivesFilters) ByQuery(query string) *ArchivesFilters {
	af.Query = query
	return af
}

func (af *ArchivesFilters) FromDate(date string) *ArchivesFilters {
	af.DateFrom = date
	return af
}

func (af *ArchivesFilters) ToDate(date string) *ArchivesFilters {
	af.DateTo = date
	return af
}

func (af *ArchivesFilters) ByProperties(propsFilters PropertiesFilters) *ArchivesFilters {
	af.Properties = propsFilters
	return af
}

func (af *ArchivesFilters) BySurveys(surveyFilters []SurveyFilter) *ArchivesFilters {
	af.Surveys = surveyFilters
	return af
}

func (af *ArchivesFilters) ByTags(includes bool, vals ...interface{}) *ArchivesFilters {
	af.Tags = NewPropertyFilterType(includes, vals...)
	return af
}

func (af *ArchivesFilters) BySales(includes bool, vals ...interface{}) *ArchivesFilters {
	af.Sales = NewPropertyFilterType(includes, vals...)
	return af
}

func (af *ArchivesFilters) ByGoals(includes bool, vals ...interface{}) *ArchivesFilters {
	af.Goals = NewPropertyFilterType(includes, vals...)
	return af
}

// Customer filters

type CustomersFilters struct {
	Country                    *StringFilter    `json:"country,omitempty"`
	Email                      *StringFilter    `json:"email,omitempty"`
	Name                       *StringFilter    `json:"name,omitempty"`
	CustomerID                 *StringFilter    `json:"customer_id,omitempty"`
	ChatsCount                 *RangeFilter     `json:"chats_count,omitempty"`
	ThreadsCount               *RangeFilter     `json:"threads_count,omitempty"`
	VisitsCount                *RangeFilter     `json:"visits_count,omitempty"`
	CreatedAt                  *DateRangeFilter `json:"created_at,omitempty"`
	AgentLastEventCreatedAt    *DateRangeFilter `json:"agent_last_event_created_at,omitempty"`
	CustomerLastEventCreatedAt *DateRangeFilter `json:"customer_last_event_created_at,omitempty"`
}

type StringFilter struct {
	Values        []string `json:"values,omitempty"`
	ExcludeValues []string `json:"exclude_values,omitempty"`
}

func NewStringFilter(values []string, inclusive bool) *StringFilter {
	sf := &StringFilter{}
	switch {
	case inclusive:
		sf.Values = values
	default:
		sf.ExcludeValues = values
	}
	return sf
}

type RangeFilter struct {
	LTE int `json:"lte,omitempty"`
	LT  int `json:"lt,omitempty"`
	GTE int `json:"gte,omitempty"`
	GT  int `json:"gt,omitempty"`
	EQ  int `json:"eq,omitempty"`
}

type DateRangeFilter struct {
	LTE string `json:"lte,omitempty"`
	LT  string `json:"lt,omitempty"`
	GTE string `json:"gte,omitempty"`
	GT  string `json:"gt,omitempty"`
	EQ  string `json:"eq,omitempty"`
}

func NewCustomersFilters() *CustomersFilters {
	return &CustomersFilters{}
}

func (cf *CustomersFilters) ByCountry(values []string, inclusive bool) *CustomersFilters {
	cf.Country = NewStringFilter(values, inclusive)
	return cf
}

func (cf *CustomersFilters) ByEmail(values []string, inclusive bool) *CustomersFilters {
	cf.Email = NewStringFilter(values, inclusive)
	return cf
}

func (cf *CustomersFilters) ByName(values []string, inclusive bool) *CustomersFilters {
	cf.Name = NewStringFilter(values, inclusive)
	return cf
}

func (cf *CustomersFilters) ByID(values []string, inclusive bool) *CustomersFilters {
	cf.CustomerID = NewStringFilter(values, inclusive)
	return cf
}

func (cf *CustomersFilters) ByChatsCount(ranges *RangeFilter) *CustomersFilters {
	cf.ChatsCount = ranges
	return cf
}

func (cf *CustomersFilters) ByThreadsCount(ranges *RangeFilter) *CustomersFilters {
	cf.ThreadsCount = ranges
	return cf
}

func (cf *CustomersFilters) ByVisitsCount(ranges *RangeFilter) *CustomersFilters {
	cf.VisitsCount = ranges
	return cf
}

func (cf *CustomersFilters) ByCreationTime(timeRange *DateRangeFilter) *CustomersFilters {
	cf.CreatedAt = timeRange
	return cf
}

func (cf *CustomersFilters) ByAgentsLastActivity(timeRange *DateRangeFilter) *CustomersFilters {
	cf.AgentLastEventCreatedAt = timeRange
	return cf
}

func (cf *CustomersFilters) ByCustomersLastActivity(timeRange *DateRangeFilter) *CustomersFilters {
	cf.CustomerLastEventCreatedAt = timeRange
	return cf
}

// Chats Filters
type ChatsFilters struct {
	IncludeActive bool              `json:"include_active,omitempty"`
	GroupIDs      []uint            `json:"group_ids,omitempty"`
	Properties    PropertiesFilters `json:"properties,omitempty"`
}

func NewChatsFilters() *ChatsFilters {
	return &ChatsFilters{
		IncludeActive: true,
	}
}

func (cf *ChatsFilters) WithoutActiveChats() *ChatsFilters {
	cf.IncludeActive = false
	return cf
}

func (cf *ChatsFilters) ByGroups(groupIDs []uint) *ChatsFilters {
	cf.GroupIDs = groupIDs
	return cf
}

func (cf *ChatsFilters) ByProperties(propsFilters PropertiesFilters) *ChatsFilters {
	cf.Properties = propsFilters
	return cf
}
