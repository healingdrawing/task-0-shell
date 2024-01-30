// Released under an MIT license. See LICENSE.

// Package literal defines the interface for penishell types that can be expressed as literals.
package literal

import (
	"penishell/internal/common/interface/cell"
)

// I (literal) is any type that can be expressed as a literal.
type I interface {
	Literal() string
}

// String returns the literal string representaition for a cell, if possible.
func String(c cell.I) string {
	l, ok := c.(I)
	if !ok {
		// Not all cell types can be expressed as literals.
		panic(c.Name() + " does not have a literal representation")
	}

	return l.Literal()
}
