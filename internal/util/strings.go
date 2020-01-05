package util

import "strings"

// SplitAndTrimSpace splits s by sep and trims space.
func SplitAndTrimSpace(s, sep string) []string {
	if len(strings.TrimSpace(s)) == 0 {
		return []string{}
	}

	ret := strings.Split(s, sep)
	for i, v := range ret {
		ret[i] = strings.TrimSpace(v)
	}
	return ret
}
