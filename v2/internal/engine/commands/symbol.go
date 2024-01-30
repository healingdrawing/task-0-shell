// Released under an MIT license. See LICENSE.

package commands

import (
	"penishell/internal/common"
	"penishell/internal/common/interface/cell"
	"penishell/internal/common/type/create"
	"penishell/internal/common/type/sym"
	"penishell/internal/common/validate"
)

func isSymbol(args cell.I) cell.I {
	v := validate.Fixed(args, 1, 1)

	return create.Bool(sym.Is(v[0]))
}

func makeSymbol(args cell.I) cell.I {
	v := validate.Fixed(args, 1, 1)

	return sym.New(common.String(v[0]))
}
