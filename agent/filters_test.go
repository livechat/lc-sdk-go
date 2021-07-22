package agent_test

import (
	"testing"

	"github.com/livechat/lc-sdk-go/v4/agent"
)

func TestPropertyFilterTypeExistenceOnly(t *testing.T) {
	pft := agent.NewPropertyFilterType(true, nil, false)
	if *pft.Exists != true {
		t.Errorf("PropertyFilterType.Exists invalid: %v", pft.Exists)
	}

	if pft.Values != nil {
		t.Errorf("PropertyFilterType.Values should not be set: %v", pft.Values)
	}

	if pft.ExcludeValues != nil {
		t.Errorf("PropertyFilterType.ExcludeValues should not be set: %v", pft.ExcludeValues)
	}

	if pft.RequireEveryValue != nil {
		t.Errorf("PropertyFilterType.RequireEveryValue should not be set: %v", pft.RequireEveryValue)
	}
}

func TestPropertyFilterTypeMatchValues(t *testing.T) {
	values := make([]interface{}, 3)
	for i := range values {
		values[i] = i
	}
	pft := agent.NewPropertyFilterType(true, values, false)

	for i := range pft.Values {
		if pft.Values[i] != i {
			t.Errorf("PropertyFilterType.Values invalid: %v", pft.Values)
		}
	}

	if pft.Exists != nil {
		t.Errorf("PropertyFilterType.Exists should not be set: %v", *pft.Exists)
	}

	if pft.ExcludeValues != nil {
		t.Errorf("PropertyFilterType.ExcludeValues should not be set: %v", pft.ExcludeValues)
	}

	if *pft.RequireEveryValue {
		t.Errorf("PropertyFilterType.RequireEveryValue invalid: %v", pft.RequireEveryValue)
	}
}

func TestPropertyFilterTypeExcludeValues(t *testing.T) {
	values := make([]interface{}, 3)
	for i := range values {
		values[i] = i
	}
	pft := agent.NewPropertyFilterType(false, values, true)

	for i := range pft.ExcludeValues {
		if pft.ExcludeValues[i] != i {
			t.Errorf("PropertyFilterType.ExcludeValues invalid: %v", pft.ExcludeValues)
		}
	}

	if pft.Exists != nil {
		t.Errorf("PropertyFilterType.Exists should not be set: %v", *pft.Exists)
	}

	if pft.Values != nil {
		t.Errorf("PropertyFilterType.Values should not be set: %v", pft.Values)
	}

	if !*pft.RequireEveryValue {
		t.Errorf("PropertyFilterType.RequireEveryValue invalid: %v", pft.RequireEveryValue)
	}
}

func TestArchivesFiltersSimpleTypeFields(t *testing.T) {
	af := agent.NewArchivesFilters()
	af.ByAgents(true, []interface{}{"a"}, false).
		ByGroups([]uint{1}).
		ByQuery("query").
		FromDate("11-09-2001").
		ToDate("02-04-2137").
		ByEventTypes(true, []string{"filled_form", "file"}, true)

	if af.Agents.Values[0] != "a" {
		t.Errorf("ArchivesFilters.Agents invalid: %v", af.Agents)
	}

	if af.Agents.RequireEveryValue == nil || *af.Agents.RequireEveryValue {
		t.Errorf("ArchivesFilters.Agents invalid: %v", af.Agents)
	}

	if af.GroupIDs[0] != 1 {
		t.Errorf("ArchivesFilters.GroupIDs invalid: %v", af.GroupIDs)
	}

	if af.Query != "query" {
		t.Errorf("ArchivesFilters.Query invalid: %v", af.Query)
	}

	if af.From != "11-09-2001" {
		t.Errorf("ArchivesFilters.From invalid: %v", af.From)
	}

	if af.To != "02-04-2137" {
		t.Errorf("ArchivesFilters.To invalid: %v", af.To)
	}

	if af.EventTypes.Values[0] != "filled_form" || af.EventTypes.Values[1] != "file" {
		t.Errorf("ArchivesFilters.EventTypes.Values invalid: %v", af.EventTypes.Values)
	}

	if af.EventTypes == nil || !*af.EventTypes.RequireEveryValue {
		t.Errorf("ArchivesFilters.EventTypes invalid: %v", af.EventTypes)
	}
}

func TestArchiveFiltersPropertyFilterTypeFields(t *testing.T) {
	values := make([]interface{}, 3)
	for i := range values {
		values[i] = i
	}
	af := agent.NewArchivesFilters()
	af.ByTags(true, values, true)
	af.BySales(true, values, true)
	af.ByGoals(true, values, true)

	for i := range af.Tags.Values {
		if af.Tags.Values[i] != i {
			t.Errorf("ArchivesFilters.Tags.Values invalid: %v", af.Tags.Values)
		}
	}

	if af.Tags.Exists != nil {
		t.Errorf("ArchivesFilters.Tags.Exists should not be set: %v", *af.Tags.Exists)
	}

	if af.Tags.ExcludeValues != nil {
		t.Errorf("ArchivesFilters.Tags.ExcludeValues should not be set: %v", af.Tags.ExcludeValues)
	}

	if !*af.Tags.RequireEveryValue {
		t.Errorf("ArchivesFilters.Tags.RequireEveryValue should not be set: %v", af.Tags.RequireEveryValue)
	}

	for i := range af.Sales.Values {
		if af.Sales.Values[i] != i {
			t.Errorf("ArchivesFilters.Sales.Values invalid: %v", af.Sales.Values)
		}
	}

	if af.Sales.Exists != nil {
		t.Errorf("ArchivesFilters.Sales.Exists should not be set: %v", *af.Sales.Exists)
	}

	if af.Sales.ExcludeValues != nil {
		t.Errorf("ArchivesFilters.Sales.ExcludeValues should not be set: %v", af.Sales.ExcludeValues)
	}

	if !*af.Sales.RequireEveryValue {
		t.Errorf("ArchivesFilters.Sales.RequireEveryValue should not be set: %v", af.Sales.RequireEveryValue)
	}

	for i := range af.Goals.Values {
		if af.Goals.Values[i] != i {
			t.Errorf("ArchivesFilters.Goals.Values invalid: %v", af.Goals.Values)
		}
	}

	if af.Goals.Exists != nil {
		t.Errorf("ArchivesFilters.Goals.Exists should not be set: %v", *af.Goals.Exists)
	}

	if af.Goals.ExcludeValues != nil {
		t.Errorf("ArchivesFilters.Goals.ExcludeValues should not be set: %v", af.Goals.ExcludeValues)
	}

	if !*af.Goals.RequireEveryValue {
		t.Errorf("ArchivesFilters.Goals.RequireEveryValue should not be set: %v", af.Goals.RequireEveryValue)
	}
}

func TestArchiveFiltersByThreadsClearsOtherFilters(t *testing.T) {
	af := agent.NewArchivesFilters()
	af.ByQuery("query")
	af.ByThreads([]string{"thread"})

	if af.ThreadIDs[0] != "thread" {
		t.Errorf("ArchivesFilters.ThreadIDs invalid: %v", af.ThreadIDs)
	}

	if af.Query != "" {
		t.Errorf("ArchivesFilters.Query should not be set: %v", af.Query)
	}
}

func TestStringFilterMatchValues(t *testing.T) {
	sf := agent.NewStringFilter([]string{"value"}, true)

	if sf.Values[0] != "value" {
		t.Errorf("StringFilter.Values invalid: %v", sf.Values)
	}

	if sf.ExcludeValues != nil {
		t.Errorf("StringFilter.ExcludeValues should not be set: %v", sf.ExcludeValues)
	}
}

func TestStringFilterExcludeValues(t *testing.T) {
	sf := agent.NewStringFilter([]string{"value"}, false)

	if sf.ExcludeValues[0] != "value" {
		t.Errorf("StringFilter.Values invalid: %v", sf.ExcludeValues)
	}

	if sf.Values != nil {
		t.Errorf("StringFilter.ExcludeValues should not be set: %v", sf.Values)
	}
}

func TestCustomersFiltersStringFilterFields(t *testing.T) {
	cf := agent.NewCustomersFilters()
	cf.ByCountry([]string{"Wakanda"}, true).
		ByName([]string{"Pink Panther"}, true).
		ByEmail([]string{"e@mail"}, false).
		ByID([]string{"id"}, false)

	if cf.Country.Values[0] != "Wakanda" {
		t.Errorf("CustomersFilters.Country.Values invalid: %v", cf.Country.Values)
	}

	if cf.Country.ExcludeValues != nil {
		t.Errorf("CustomersFilters.Country.ExcludeValues should not be set: %v", cf.Country.ExcludeValues)
	}

	if cf.Name.Values[0] != "Pink Panther" {
		t.Errorf("CustomersFilters.Name.Values invalid: %v", cf.Name.Values)
	}

	if cf.Name.ExcludeValues != nil {
		t.Errorf("CustomersFilters.Name.ExcludeValues should not be set: %v", cf.Name.ExcludeValues)
	}

	if cf.Email.ExcludeValues[0] != "e@mail" {
		t.Errorf("CustomersFilters.Email.ExcludeValues invalid: %v", cf.Email.ExcludeValues)
	}

	if cf.Email.Values != nil {
		t.Errorf("CustomersFilters.Email.Values should not be set: %v", cf.Email.Values)
	}

	if cf.CustomerID.ExcludeValues[0] != "id" {
		t.Errorf("CustomersFilters.CustomerID.ExcludeValues invalid: %v", cf.CustomerID.ExcludeValues)
	}

	if cf.CustomerID.Values != nil {
		t.Errorf("CustomersFilters.CustomerID.Values should not be set: %v", cf.CustomerID.Values)
	}
}

func TestCustomersFiltersRangeFilterFields(t *testing.T) {
	rf := &agent.RangeFilter{
		LT: 5,
		GT: 2,
	}
	cf := agent.NewCustomersFilters()
	cf.ByChatsCount(rf).ByThreadsCount(rf).ByVisitsCount(rf)

	if cf.ChatsCount.LT != 5 || cf.ChatsCount.GT != 2 {
		t.Errorf("CustomersFilters.ChatsCount invalid: %v", cf.ChatsCount)
	}

	if cf.ThreadsCount.LT != 5 || cf.ThreadsCount.GT != 2 {
		t.Errorf("CustomersFilters.ThreadsCount invalid: %v", cf.ThreadsCount)
	}

	if cf.VisitsCount.LT != 5 || cf.VisitsCount.GT != 2 {
		t.Errorf("CustomersFilters.VisitsCount invalid: %v", cf.VisitsCount)
	}
}

func TestCustomersFiltersDateRangeFilterFields(t *testing.T) {
	drf := &agent.DateRangeFilter{
		GT: "11-09-2001",
		LT: "02-04-2137",
	}
	cf := agent.NewCustomersFilters()
	cf.ByCreationTime(drf).ByAgentsLastActivity(drf).ByCustomersLastActivity(drf)

	if cf.CreatedAt.LT != "02-04-2137" || cf.CreatedAt.GT != "11-09-2001" {
		t.Errorf("CustomersFilters.CreatedAt invalid: %v", cf.CreatedAt)
	}

	if cf.AgentLastEventCreatedAt.LT != "02-04-2137" || cf.AgentLastEventCreatedAt.GT != "11-09-2001" {
		t.Errorf("CustomersFilters.AgentLastEventCreatedAt invalid: %v", cf.AgentLastEventCreatedAt)
	}

	if cf.CustomerLastEventCreatedAt.LT != "02-04-2137" || cf.CustomerLastEventCreatedAt.GT != "11-09-2001" {
		t.Errorf("CustomersFilters.CustomerLastEventCreatedAt invalid: %v", cf.CustomerLastEventCreatedAt)
	}
}

func TestChatsFilters(t *testing.T) {
	cf := agent.NewChatsFilters()
	if !cf.IncludeActive {
		t.Errorf("ChatsFilters.IncludeActive should be true by default")
	}

	cf.WithoutActiveChats().ByGroups([]uint{1})
	if cf.IncludeActive {
		t.Errorf("ChatsFilters.IncludeActive should be toggled to false")
	}
	if cf.GroupIDs[0] != 1 {
		t.Errorf("ChatsFilters.GroupIDs invalid: %v", cf.GroupIDs)
	}
}

func TestThreadsFilters(t *testing.T) {
	tf := agent.NewThreadsFilters()

	tf.FromDate("11-09-2001").ToDate("02-04-2137")
	if tf.From != "11-09-2001" {
		t.Errorf("ThreadsFilters.From invalid: %v", tf.From)
	}

	if tf.To != "02-04-2137" {
		t.Errorf("ThreadsFilters.To invalid: %v", tf.To)
	}
}

func TestIntegerFilterMatchValues(t *testing.T) {
	intF := agent.NewIntegerFilter([]int64{12345678901234567}, true)

	if intF.Values[0] != 12345678901234567 {
		t.Errorf("IntegerFilter.Values invalid: %v", intF.Values)
	}

	if intF.ExcludeValues != nil {
		t.Errorf("IntegerFilter.ExcludeValues should not be set: %v", intF.ExcludeValues)
	}
}

func TestIntegerFilterExcludeValues(t *testing.T) {
	intF := agent.NewIntegerFilter([]int64{12345678901234567}, false)

	if intF.ExcludeValues[0] != 12345678901234567 {
		t.Errorf("IntegerFilter.Values invalid: %v", intF.ExcludeValues)
	}

	if intF.Values != nil {
		t.Errorf("IntegerFilter.ExcludeValues should not be set: %v", intF.Values)
	}
}

func TestCustomersFiltersIntegerFilterField(t *testing.T) {
	cf := agent.NewCustomersFilters()
	cf.ByChatGroupIDs([]int64{12345678901234567}, true)

	if cf.ChatGroupIDs.Values[0] != 12345678901234567 {
		t.Errorf("CustomersFilters.ChatGroupIDs.Values invalid: %v", cf.ChatGroupIDs.Values)
	}
}
