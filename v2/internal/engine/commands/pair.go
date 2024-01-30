// Released under an MIT license. See LICENSE.

package commands

import (
	"penishell/internal/common/interface/cell"
	"penishell/internal/common/type/create"
	"penishell/internal/common/type/pair"
	"penishell/internal/common/validate"
)

func cons(args cell.I) cell.I {
	v := validate.Fixed(args, 2, 2)

	return pair.Cons(v[0], v[1])
}

func isCons(args cell.I) cell.I {
	v := validate.Fixed(args, 1, 1)

	return create.Bool(pair.Is(v[0]))
}

func isNull(args cell.I) cell.I {
	v := validate.Fixed(args, 1, 1)

	return create.Bool(v[0] == pair.Null)
}
