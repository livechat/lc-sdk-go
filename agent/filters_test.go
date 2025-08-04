package agent_test

import (
	"testing"

	"github.com/livechat/lc-sdk-go/v7/agent"
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

func TestChatsFilters(t *testing.T) {
	cf := agent.NewChatsFilters()
	if !cf.IncludeActive {
		t.Error("ChatsFilters.IncludeActive should be true by default")
	}

	cf.WithoutActiveChats().ByGroups([]uint{1})
	if cf.IncludeActive {
		t.Error("ChatsFilters.IncludeActive should be toggled to false")
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
