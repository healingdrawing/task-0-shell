// Released under an MIT license. See LICENSE.

// Package rational defines the interface for penishell's numeric types.
package rational

import (
	"math/big"

	"penishell/internal/common/interface/cell"
)

// I (rational) is anything that can be treated as a rational number in penishell.
type I interface {
	Rat() *big.Rat
}

type rational = I

// Number returns the *big.Rat value for a cell, if possible.
func Number(c cell.I) *big.Rat {
	r, ok := c.(rational)
	if !ok {
		// Not all cell types can be treated as numbers.
		panic(c.Name() + " cannot be use in a numeric context")
	}

	return r.Rat()
}
