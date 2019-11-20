package agent

type PropertiesFilters map[string]map[string]*PropertyFilterType

type PropertyFilterType struct {
	Exists        bool          `json:"exists,omitempty"`
	Values        []interface{} `json:"values,omitempty"`
	ExcludeValues []interface{} `json:"exclude_values,omitempty"`
}

func NewPropertyFilterType(includes bool, vals ...interface{}) PropertyFilterType {
	pft := PropertyFilterType{}
	switch {
	case vals == nil:
		pft.Exists = includes
	case includes:
		pft.Values = vals
	case !includes:
		pft.ExcludeValues = vals
	}
	return pft
}

// Archives filters

type ArchivesFilters struct {
	AgentIDs   []string           `json:"agent_ids"`
	GroupIDs   []uint             `json:"group_ids"`
	DateFrom   string             `json:"date_from"`
	DateTo     string             `json:"date_to"`
	Properties PropertiesFilters  `json:"properties"`
	Tags       PropertyFilterType `json:"tags"`
	Sales      PropertyFilterType `json:"sales"`
	Goals      PropertyFilterType `json:"goals"`
	Surveys    []SurveyFilter     `json:"surveys"`
	ThreadIDs  []string           `json:"thread_ids"`
	Query      string             `json:"query"`
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
	af = &ArchivesFilters{
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
	af.Tags = NewPropertyFilterType(includes, vals)
	return af
}

func (af *ArchivesFilters) BySales(includes bool, vals ...interface{}) *ArchivesFilters {
	af.Sales = NewPropertyFilterType(includes, vals)
	return af
}

func (af *ArchivesFilters) ByGoals(includes bool, vals ...interface{}) *ArchivesFilters {
	af.Goals = NewPropertyFilterType(includes, vals)
	return af
}

// Customer filters

type CustomersFilters struct {
	Country                    *StringFilter    `json:"country"`
	Email                      *StringFilter    `json:"email"`
	Name                       *StringFilter    `json:"name"`
	CustomerID                 *StringFilter    `json:"customer_id"`
	ChatsCount                 *RangeFilter     `json:"chats_count"`
	ThreadsCount               *RangeFilter     `json:"threads_count"`
	VisitsCount                *RangeFilter     `json:"visits_count"`
	CreatedAt                  *DateRangeFilter `json:"created_at"`
	AgentLastEventCreatedAt    *DateRangeFilter `json:"agent_last_event_created_at"`
	CustomerLastEventCreatedAt *DateRangeFilter `json:"customer_last_event_created_at"`
}

type StringFilter struct {
	Values        []string `json:"values,omitempty"`
	ExcludeValues []string `json:"exclude_values,omitempty"`
}

func NewStringFilter(values []string, shouldExclude bool) *StringFilter {
	sf := &StringFilter{}
	switch {
	case shouldExclude:
		sf.ExcludeValues = values
	default:
		sf.Values = values
	}
	return sf
}

type RangeFilter struct {
	LTE int `json:"lte"`
	LT  int `json:"lt"`
	GTE int `json:"gte"`
	GT  int `json:"gt"`
	EQ  int `json:"eq"`
}

type DateRangeFilter struct {
	LTE string `json:"lte"`
	LT  string `json:"lt"`
	GTE string `json:"gte"`
	GT  string `json:"gt"`
	EQ  string `json:"eq"`
}

func NewCustomersFilters() *CustomersFilters {
	return &CustomersFilters{}
}

func (cf *CustomersFilters) ByCountry(values []string, shouldExclude bool) *CustomersFilters {
	cf.Country = NewStringFilter(values, shouldExclude)
	return cf
}

func (cf *CustomersFilters) ByEmail(values []string, shouldExclude bool) *CustomersFilters {
	cf.Email = NewStringFilter(values, shouldExclude)
	return cf
}

func (cf *CustomersFilters) ByName(values []string, shouldExclude bool) *CustomersFilters {
	cf.Name = NewStringFilter(values, shouldExclude)
	return cf
}

func (cf *CustomersFilters) ByID(values []string, shouldExclude bool) *CustomersFilters {
	cf.CustomerID = NewStringFilter(values, shouldExclude)
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
	IncludeActive bool              `json:"include_active"`
	GroupIDs      []uint            `json:"group_ids"`
	Properties    PropertiesFilters `json:"properties"`
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
