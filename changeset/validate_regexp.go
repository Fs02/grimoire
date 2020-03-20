package changeset

import (
	"regexp"
	"strings"
)

// ValidateRegexpErrorMessage is the default error message for ValidateRegexp.
var ValidateRegexpErrorMessage = "{field}'s format is invalid"

// ValidateRegexp validates the value of given field to match given regexp.
func ValidateRegexp(ch *Changeset, field string, exp *regexp.Regexp, opts ...Option) {
	options := Options{
		message: ValidateRegexpErrorMessage,
	}
	options.apply(opts)

	val, exist := ch.changes[field]
	if !exist || contains(options.emptyValues, val) {
		return
	}

	if str, ok := val.(string); ok {
		match := exp.MatchString(str)
		if !match {
			msg := strings.Replace(options.message, "{field}", field, 1)
			AddError(ch, field, msg)
		}
		return
	}
}
