// Released under an MIT license. See LICENSE.

// Package create provides helper functions for creating penishell types.
package create

import (
	"penishell/internal/common/interface/cell"
	"penishell/internal/common/type/pair"
	"penishell/internal/common/type/sym"
)

// Bool returns the penishell value corresponding to the value of the boolean a.
func Bool(a bool) cell.I {
	if a {
		return sym.True
	}

	return pair.Null
}
