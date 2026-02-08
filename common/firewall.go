package common

type FirewallZoneCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Target      string `json:"target,omitempty"`
}

type FirewallZoneResponse struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Target      string   `json:"target,omitempty"`
	Services    []string `json:"services,omitempty"`
	Ports       []string `json:"ports,omitempty"`
	RichRules   []string `json:"rich_rules,omitempty"`
}

type FirewallZoneListResponse struct {
	Zones []FirewallZoneResponse `json:"zones"`
}
