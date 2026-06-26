// Package respolicy wraps [res.Registry] with composable add-time policies.
//
// [res.Registry] stores entries as-is; respolicy runs [AddPolicy] before each
// Add or AddWithTags so callers can normalize or reject values at registration
// time without embedding rules in core res.
//
// Future work (not in this package yet):
//   - wire policies (registrant → build spec)
//   - Replace for materialize without a second Add
//   - composition-root wiring in app instead of raw [res.Global]
package respolicy
