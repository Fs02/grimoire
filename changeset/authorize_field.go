package changeset

// AuthorizeFieldErrorMessage is the default error message for AuthorizeField
var AuthorizeFieldErrorMessage = "Forbidden"

// AuthorizeField check the field based on permission
func AuthorizedField(ch *Changeset, fields []string, allowed bool, opts ...Option) {
	options := Options{
		message: AuthorizeFieldErrorMessage,
	}
	options.apply(opts)

	for _, element := range fields {
		_, exist := ch.changes[element]
		if !exist {
			return
		}

		if !allowed {
			AddError(ch, element, options.message)
		}
	}
}
