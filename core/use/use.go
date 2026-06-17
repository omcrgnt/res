// Package use registers org system defaults (logger, telemetry) in res.Default.
//
// Import for side effects at the app composition root:
//
//	import _ "github.com/omcrgnt/res/core/use"
package use

import (
	_ "github.com/omcrgnt/logger/use"
	_ "github.com/omcrgnt/telemetry/use"
)
