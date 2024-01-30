// Released under an MIT license. See LICENSE.

// Package reference defines the interface for penishell's variable type.
package reference

import (
	"penishell/internal/common/interface/cell"
)

// I (reference) is anything that can hold a value.
type I interface {
	Copy() I
	Get() cell.I
	Set(cell.I)
}
