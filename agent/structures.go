package agent

import "github.com/livechat/lc-sdk-go/v2/objects"

type postback struct {
	ID      string `json:"id"`
	Toggled bool   `json:"toggled"`
}

type ban struct {
	Days uint `json:"days"`
}

type InitialChat struct {
	objects.InitialChat
	Users []*objects.User `json:"users,omitempty"`
}

// MulticastRecipients aggregates Agent and Customer recipients that multicast should be sent to
type MulticastRecipients struct {
	Agents    *MulticastRecipientsAgents    `json:"agents,omitempty"`
	Customers *MulticastRecipientsCustomers `json:"customers,omitempty"`
}

// MulticastRecipientsAgents represents recipients for multicast to agents
type MulticastRecipientsAgents struct {
	Groups []uint   `json:"groups,omitempty"`
	IDs    []string `json:"ids,omitempty"`
	All    *bool    `json:"all,omitempty"`
}

// MulticastRecipientsCustomers represents recipients for multicast to customers
type MulticastRecipientsCustomers struct {
	IDs []string `json:"ids,omitempty"`
}

type transferTarget struct {
	Type string        `json:"type"`
	IDs  []interface{} `json:"ids"`
}

type AgentsForTransfer []struct {
	AgentID          string `json:"agent_id"`
	TotalActiveChats uint   `json:"total_active_chats"`
}
