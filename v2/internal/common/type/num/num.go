// Released under an MIT license. See LICENSE.

// Package num provides penishell's rational number type.
package num

import (
	"fmt"
	"math/big"

	"penishell/internal/common/interface/boolean"
	"penishell/internal/common/interface/cell"
	"penishell/internal/common/interface/literal"
	"penishell/internal/common/interface/rational"
)

const name = "number"

// T (num) wraps Go's big.Rat type.
type T big.Rat

type num = T

// Int creates a num from the integer i.
func Int(i int) cell.I {
	return Rat(big.NewRat(int64(i), 1))
}

// New creates a new num from a string.
func New(s string) cell.I {
	v := &big.Rat{}

	if _, ok := v.SetString(s); !ok {
		panic("'" + s + "' is not a valid number")
	}

	return Rat(v)
}

// Rat creates wraps the *big.Rat r as a num.
func Rat(r *big.Rat) cell.I {
	return (*num)(r)
}

// Bool returns the boolean value of the num n.
func (n *num) Bool() bool {
	return n.Rat().Cmp(&big.Rat{}) != 0
}

// Equal returns true if c is the same number as the num n.
func (n *num) Equal(c cell.I) bool {
	return Is(c) && n.Rat().Cmp(To(c).Rat()) == 0
}

// Literal returns the literal representation of the num n.
func (n *num) Literal() string {
	return "(|" + name + " " + n.String() + "|)"
}

// Name returns the type name for the num n.
func (n *num) Name() string {
	return name
}

// Rat returns the value of the num n as a *big.Rat.
func (n *num) Rat() *big.Rat {
	return (*big.Rat)(n)
}

// String returns the text of the num n.
func (n *num) String() string {
	return n.Rat().RatString()
}

// A compiler-checked list of interfaces this type satisfies. Never called.
func implements() { //nolint:deadcode,unused
	var t num

	// The num type has a boolean value.
	_ = boolean.I(&t)

	// The num type is a cell.
	_ = cell.I(&t)

	// The num type has a literal representation.
	_ = literal.I(&t)

	// The num type is a rational.
	_ = rational.I(&t)

	// The num type is a stringer.
	_ = fmt.Stringer(&t)
}
