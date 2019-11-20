package agent

type Postback struct {
	ID      string `json:"id"`
	Toggled bool   `json:"toggled"`
}

type Ban struct {
	Days uint64 `json:"days"`
}

type MulticastScopes struct {
	Agents    *MulticastScopesAgents    `json:"agents,omitempty"`
	Customers *MulticastScopesCustomers `json:"customers,omitempty"`
}

type MulticastScopesAgents struct {
	Groups *[]uint64 `json:"groups,omitempty"`
	IDs    *[]string `json:"ids,omitempty"`
	All    *bool     `json:"all,omitempty"`
}

type MulticastScopesCustomers struct {
	IDs *[]string `json:"ids,omitempty"`
}

type TransferTarget struct {
	Type string `json:"type"`
	IDs  []uint `json:"ids"`
}
