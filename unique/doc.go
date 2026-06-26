// Package unique is a type-unique registry: at most one [res.Entry] per concrete type.
// Tag behavior ([res.TagRegular], [res.TagReplaceable], [res.TagFixed]) is enforced on Add.
// Library use init may register defaults via [MustAddReplaceable] or [MustAddFixed] on [Global].
package unique
