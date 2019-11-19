package agent

type PropertiesFilters map[string]map[string]*PropertyFilterType

type PropertyFilterType struct {
	Exists        bool          `json:"exists,omitempty"`
	Values        []interface{} `json:"values,omitempty"`
	ExcludeValues []interface{} `json:"exclude_values,omitempty"`
}

//Archives filters

type ArchivesFilters struct {
	AgentIDs   []string           `json:"agent_ids"`
	GroupIDs   []int32            `json:"group_ids"`
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

func (af *ArchivesFilters) ByGroups(groupIDs []int32) *ArchivesFilters {
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
	pft := PropertyFilterType{}
	switch {
	case vals == nil:
		pft.Exists = includes
	case includes:
		pft.Values = vals
	case !includes:
		pft.ExcludeValues = vals
	}
	af.Tags = pft
	return af
}

func (af *ArchivesFilters) BySales(includes bool, vals ...interface{}) *ArchivesFilters {
	pft := PropertyFilterType{}
	switch {
	case vals == nil:
		pft.Exists = includes
	case includes:
		pft.Values = vals
	case !includes:
		pft.ExcludeValues = vals
	}
	af.Sales = pft
	return af
}

func (af *ArchivesFilters) ByGoals(includes bool, vals ...interface{}) *ArchivesFilters {
	pft := PropertyFilterType{}
	switch {
	case vals == nil:
		pft.Exists = includes
	case includes:
		pft.Values = vals
	case !includes:
		pft.ExcludeValues = vals
	}
	af.Goals = pft
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
	Values        []string `json:"values"`
	ExcludeValues []string `json:"exclude_values"`
}

type RangeFilter struct {
	LTE *int64 `json:"lte"`
	LT  *int64 `json:"lt"`
	GTE *int64 `json:"gte"`
	GT  *int64 `json:"gt"`
	EQ  *int64 `json:"eq"`
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
	if shouldExclude {
		cf.Country = &StringFilter{
			ExcludeValues: values,
		}
	} else {
		cf.Country = &StringFilter{
			Values: values,
		}
	}
	return cf
}

func (cf *CustomersFilters) ByEmail(values []string, shouldExclude bool) *CustomersFilters {
	if shouldExclude {
		cf.Email = &StringFilter{
			ExcludeValues: values,
		}
	} else {
		cf.Email = &StringFilter{
			Values: values,
		}
	}
	return cf
}

func (cf *CustomersFilters) ByName(values []string, shouldExclude bool) *CustomersFilters {
	if shouldExclude {
		cf.Name = &StringFilter{
			ExcludeValues: values,
		}
	} else {
		cf.Name = &StringFilter{
			Values: values,
		}
	}
	return cf
}

func (cf *CustomersFilters) ByID(values []string, shouldExclude bool) *CustomersFilters {
	if shouldExclude {
		cf.CustomerID = &StringFilter{
			ExcludeValues: values,
		}
	} else {
		cf.CustomerID = &StringFilter{
			Values: values,
		}
	}
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
	GroupIDs      []int32           `json:"group_ids"`
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

func (cf *ChatsFilters) ByGroups(groupIDs []int32) *ChatsFilters {
	cf.GroupIDs = groupIDs
	return cf
}

func (cf *ChatsFilters) ByProperties(propsFilters PropertiesFilters) *ChatsFilters {
	cf.Properties = propsFilters
	return cf
}
