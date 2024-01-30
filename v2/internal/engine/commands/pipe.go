// Released under an MIT license. See LICENSE.

package commands

import (
	"penishell/internal/common/interface/cell"
	"penishell/internal/common/type/create"
	"penishell/internal/common/type/pipe"
	"penishell/internal/common/validate"
)

func isPipe(args cell.I) cell.I {
	v := validate.Fixed(args, 1, 1)

	return create.Bool(pipe.Is(v[0]))
}

func makePipe(args cell.I) cell.I {
	validate.Fixed(args, 0, 0)

	return pipe.New(nil, nil)
}
