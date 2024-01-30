// Released under an MIT license. See LICENSE.

package sym

import (
	"penishell/internal/common/interface/cell"
	"penishell/internal/common/struct/loc"
	"penishell/internal/common/struct/token"
)

// Plus is a symbol plus its lexical location.
type Plus struct {
	*sym
	source *loc.T
}

// Token creates a Plus from a token.T.
func Token(t *token.T) cell.I {
	p := symnew(t.Value())

	return &Plus{p, t.Source()}
}

// Source returns the lexical location for a sym that has it.
func (p *Plus) Source() *loc.T {
	return p.source
}
