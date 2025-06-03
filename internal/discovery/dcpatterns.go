package dcpatterns

import (
	"fmt"
	"regexp"
	"strings"
)

type DCPattern struct {
	Key   string
	Label string
}

func GetRecommendedDC(hostname string, patterns map[string]DCPattern) string {
	type match struct {
		pattern string
		dc      DCPattern
	}
	var matches []match

	for pattern, dc := range patterns {
		var regex string
		if strings.ContainsAny(pattern, ".*+?()|[]{}^$") {
			regex = "(?i)" + pattern
		} else {
			regex = "(?i)^" + regexp.QuoteMeta(pattern) + ".*"
		}
		if matched, _ := regexp.MatchString(regex, hostname); matched {
			matches = append(matches, match{pattern: pattern, dc: dc})
		}
	}

	if len(matches) == 0 {
		fmt.Printf("DEBUG: No pattern matched for hostname %q\n", hostname)
		return ""
	}

	// Select the longest matching pattern for specificity
	var bestMatch match
	for _, m := range matches {
		if len(m.pattern) > len(bestMatch.pattern) {
			bestMatch = m
		}
	}

	fmt.Printf("DEBUG: Selected pattern %q for hostname %q, DC: %s:%s\n", bestMatch.pattern, hostname, bestMatch.dc.Key, bestMatch.dc.Label)
	return bestMatch.dc.Key + ":" + bestMatch.dc.Label
}

func GetUniqueDCLocations(patterns map[string]DCPattern) map[string]struct{} {
	unique := make(map[string]struct{})
	for _, dc := range patterns {
		unique[dc.Key+":"+dc.Label] = struct{}{}
	}
	return unique
}
