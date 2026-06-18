package res

// Tag is metadata on a registered resource ([Entry]), set by [Registry.AddWithTags].
// The registry does not interpret tags; callers read them from [Entry].
type Tag struct {
	name string
}

// TagReplaceable marks a resource as a fallback candidate when a caller
// chooses one of several matching entries.
var TagReplaceable = Tag{name: "replaceable"}

// TagFixed marks a resource that must not be replaced or deduplicated away.
// If another entry matches the same dependency type, sdi.Resolve fails.
var TagFixed = Tag{name: "fixed"}

// Has reports whether e has tag.
func (e Entry) Has(tag Tag) bool {
	_, ok := e.tags[tag]
	return ok
}

// Replaceable reports whether e has [TagReplaceable].
func (e Entry) Replaceable() bool {
	return e.Has(TagReplaceable)
}

// Fixed reports whether e has [TagFixed].
func (e Entry) Fixed() bool {
	return e.Has(TagFixed)
}

// Tags returns a copy of tags attached to e.
func (e Entry) Tags() []Tag {
	if len(e.tags) == 0 {
		return nil
	}
	out := make([]Tag, 0, len(e.tags))
	for tag := range e.tags {
		out = append(out, tag)
	}
	return out
}

type tagSet map[Tag]struct{}

func newTagSet(tags ...Tag) tagSet {
	if len(tags) == 0 {
		return nil
	}
	s := make(tagSet, len(tags))
	for _, tag := range tags {
		s[tag] = struct{}{}
	}
	return s
}
