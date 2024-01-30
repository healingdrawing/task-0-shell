// Released under an MIT license. See LICENSE.

// Package str provides penishell's string type.
package str

import (
	"fmt"

	"penishell/internal/common/interface/cell"
	"penishell/internal/common/interface/literal"

	"github.com/michaelmacinnis/adapted"
)

const name = "string"

// T (str) wraps Go's string type.
type T string

type str = T

// New creates a new str cell.
func New(v string) cell.I {
	s := str(v)

	return &s
}

// Equal returns true if the cell c wraps the same string and false otherwise.
func (s *str) Equal(c cell.I) bool {
	return Is(c) && s.String() == To(c).String()
}

// Literal returns the literal representation of the str s.
func (s *str) Literal() string {
	return adapted.CanonicalString(string(*s))
}

// Name returns the name of the str type.
func (s *str) Name() string {
	return name
}

// String returns the text of the str s.
func (s *str) String() string {
	return string(*s)
}

// A compiler-checked list of interfaces this type satisfies. Never called.
func implements() { //nolint:deadcode,unused
	var t str

	// The str type is a cell.
	_ = cell.I(&t)

	// The str type has a literal representation.
	_ = literal.I(&t)

	// The str type is a stringer.
	_ = fmt.Stringer(&t)
}
