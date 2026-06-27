// Package unique is a type-unique registry: at most one [res.Entry] per concrete type.
//
// System tag policy ([res.TagRegular], [res.TagReplaceable], [res.TagFixed]) is enforced on Add.
// Library use init registers defaults via [MustAddReplaceable] or [MustAddFixed] on [Global].
//
// App bootstrap API:
//
//	unique.Add(v)                              — materialized resource
//	unique.AddWithCustomTag(v, key, val)       — spec + custom tag (e.g. ecfg segment)
//	entry.ChangeValue(built)                   — materialize spec in-place (via WalkEntries)
//	unique.Merge(main, side)                   — merge app catalog into main
//
// [MustAddReplaceable] and [MustAddFixed] are the only public init entry points for Replaceable/Fixed.
package unique
