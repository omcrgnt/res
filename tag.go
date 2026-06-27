package res

// Tag is metadata on a registered resource ([Entry]), set by [Registry.AddWithTags].
// The registry does not interpret tags; callers read them from [Entry].
type Tag struct {
	name string
}

// TagRegular marks a normal application resource (not a library default).
var TagRegular = Tag{name: "regular"}

// TagReplaceable marks a resource as a fallback candidate when a caller
// chooses one of several matching entries.
var TagReplaceable = Tag{name: "replaceable"}

// TagFixed marks a resource that must not be replaced or deduplicated away.
// If another entry matches the same dependency type, sdi.Resolve fails.
var TagFixed = Tag{name: "fixed"}

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
