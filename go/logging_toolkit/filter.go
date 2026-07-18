package logging_toolkit

import (
	"regexp"
)

// Filter determines whether a log entry should be recorded.
type Filter struct {
	MinLevel Level
	Pattern  *regexp.Regexp
}

// ShouldLog returns true if the entry at the given level and message passes the
// filter.
func (f *Filter) ShouldLog(level Level, msg string) bool {
	if level < f.MinLevel {
		return false
	}
	if f.Pattern != nil && !f.Pattern.MatchString(msg) {
		return false
	}
	return true
}
