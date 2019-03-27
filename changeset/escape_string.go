package changeset

import (
	"html"
)

// EscapeString escapes special characters like "<" to become "&lt;". this is helper for html.EscapeString
func EscapeString(ch *Changeset, field string) {
	ApplyString(ch, field, html.EscapeString)
}

// EscapeString escapes special characters like "<" to become "&lt;" for all changes on changeset. this is helper for html.EscapeString
func EscapeStringAll(ch *Changeset) {
	for field, _ := range ch.changes {
		ApplyString(ch, field, html.EscapeString)
	}
}
