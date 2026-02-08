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

type FirewallRuleRequest struct {
	Zone     string `json:"zone"`
	RuleType string `json:"rule_type"` // "rich", "port", or "service"
	Rule     string `json:"rule,omitempty"`
	Port     string `json:"port,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Service  string `json:"service,omitempty"`
}

type FirewallRuleResponse struct {
	Zone     string `json:"zone"`
	RuleType string `json:"rule_type"`
	Rule     string `json:"rule,omitempty"`
	Port     string `json:"port,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Service  string `json:"service,omitempty"`
}
