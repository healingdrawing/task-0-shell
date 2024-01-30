// Released under an MIT license. See LICENSE.

// Package boot provides what is necessary for bootstrapping penishell.
package boot

import _ "embed" // Blank import required by embed.

//go:embed boot.penishell
var script string //nolint:gochecknoglobals

// Script returns the boot script for penishell.
func Script() string { //nolint:funlen
	return script
}
