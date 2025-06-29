package rules

type Rule struct {
	MatchRule          MatchRule   `json:"match"`
	CompareRule        CompareRule `json:"compare,omitempty"`
	RequestOperations  []Operation `json:"req,omitempty"`
	ResponseOperations []Operation `json:"res,omitempty"`
}

func NewRule(match MatchRule, compare CompareRule, req_operations []Operation, res_operations []Operation) *Rule {
	return &Rule{
		MatchRule:          match,
		CompareRule:        compare,
		RequestOperations:  req_operations,
		ResponseOperations: res_operations,
	}
}
