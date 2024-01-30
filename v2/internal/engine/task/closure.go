// Released under an MIT license. See LICENSE.

package task

import (
	"penishell/internal/common/interface/cell"
	"penishell/internal/common/interface/scope"
)

// Closure underlies the method, and syntax types.
type Closure struct {
	Body cell.I // Body of the routine.
	Labels
	Op
	Scope scope.I
}

// Labels hold the labels for a user-defined routine.
type Labels struct {
	Env    cell.I // Calling env label.
	Params cell.I // Param labels.
	Self   cell.I // Label for the env where this routine was found.
}
