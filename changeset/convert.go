package changeset

// Convert a struct as changeset, every field's value will be treated as changes. Returns a new changeset.
func Convert(entity interface{}) *Changeset {
	ch := &Changeset{}
	ch.entity = entity
	ch.values = make(map[string]interface{})
	ch.changes, ch.types = mapSchema(ch.entity)

	return ch
}
