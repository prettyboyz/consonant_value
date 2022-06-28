package txn

import (
	mrand "math/rand"
	"time"
)

var chaosEnabled = false
var chaosSetting Chaos

// Chaos holds parameters for the failure injection mechanism.
type Chaos struct {
	// KillChance is the 0.0 to 1.0 chance that a given checkpoint
	// within the algorithm will raise an interruption that will
	// stop the procedure.
	KillChance float64

	// SlowdownChance is the 0.0 to 1.0 chance that a given checkpoint
	// within the algorithm will be delayed by Slowdown before
	// continuing.
	SlowdownChance float64
	Slowdown       time.Duration

	// If Breakpoint is set, the above settings will only affect the
	// named breakpoint.
	Breakpoint string
}

// SetChaos sets the failure injection parameters to c.