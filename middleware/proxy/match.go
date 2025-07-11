package proxy

import "strings"

type Match struct {
}

func NewMatch() *Match {
	return &Match{}
}

// Match
// return string(匹配前缀), bool（是否匹配）
func (m *Match) Match(path string, prefix []string) (string, bool) {
	for _, p := range prefix {
		if strings.HasPrefix(path, p) {
			return p, true
		}
	}
	return "", false
}
