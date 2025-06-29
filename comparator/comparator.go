package comparator

type Comparator struct {
	Match       []string `json:"match,omitempty"`
	Ignore      []string `json:"ignore,omitempty"`
	MatchRegex  bool     `json:"matchRegex,omitempty"`
	IgnoreRegex bool     `json:"ignoreRegex,omitempty"`
	IgnoreAll   bool     `json:"ignoreAll,omitempty"`
}
