// Released under an MIT license. See LICENSE.

// Package cell defines the interface for all penishell types.
package cell

// I (cell) is the basic unit of storage in penishell.
type I interface {
	Equal(c I) bool
	Name() string
}
