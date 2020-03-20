package changeset

import (
	"regexp"
	"strings"
)

// ValidatePatternErrorMessage is the default error message for ValidatePattern.
var ValidatePatternErrorMessage = "{field}'s format is invalid"

// ValidatePattern validates the value of given field to match given pattern.
func ValidatePattern(ch *Changeset, field string, pattern string, opts ...Option) {
	options := Options{
		message: ValidatePatternErrorMessage,
	}
	options.apply(opts)

	val, exist := ch.changes[field]
	if !exist || contains(options.emptyValues, val) {
		return
	}

	if str, ok := val.(string); ok {
		match, _ := regexp.MatchString(pattern, str)
		if !match {
			msg := strings.Replace(options.message, "{field}", field, 1)
			AddError(ch, field, msg)
		}
		return
	}
}
