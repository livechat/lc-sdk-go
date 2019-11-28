package agent

type postback struct {
	ID      string `json:"id"`
	Toggled bool   `json:"toggled"`
}

type ban struct {
	Days uint `json:"days"`
}

// MulticastScopes aggregates Agent and Customer Scopes that multicast should be sent to
type MulticastScopes struct {
	Agents    *MulticastScopesAgents    `json:"agents,omitempty"`
	Customers *MulticastScopesCustomers `json:"customers,omitempty"`
}

// MulticastScopesAgents represents scopes for multicast to agents
type MulticastScopesAgents struct {
	Groups []uint   `json:"groups,omitempty"`
	IDs    []string `json:"ids,omitempty"`
	All    *bool    `json:"all,omitempty"`
}

// MulticastScopesCustomers represents scopes for multicast to customers
type MulticastScopesCustomers struct {
	IDs []string `json:"ids,omitempty"`
}

type transferTarget struct {
	Type string        `json:"type"`
	IDs  []interface{} `json:"ids"`
}
