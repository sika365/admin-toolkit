package utils

import "regexp"

func FindStringSubmatch(reg *regexp.Regexp, s string) map[string]string {
	matchGroup := reg.FindStringSubmatch(s)
	return MatchGroupMap(reg, matchGroup)
}

// MatchGroupMap converts the match string array into a map of group keys to group values
func MatchGroupMap(reg *regexp.Regexp, match []string) map[string]string {
	ret := make(map[string]string)
	for i, name := range reg.SubexpNames() {
		if i != 0 && name != "" {
			ret[name] = match[i]
		}
	}
	return ret
}
