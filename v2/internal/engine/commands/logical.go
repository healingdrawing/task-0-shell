// Released under an MIT license. See LICENSE.

package commands

import (
	"penishell/internal/common/interface/boolean"
	"penishell/internal/common/interface/cell"
	"penishell/internal/common/type/create"
	"penishell/internal/common/validate"
)

func not(args cell.I) cell.I {
	v := validate.Fixed(args, 1, 1)

	return create.Bool(!boolean.Value(v[0]))
}
