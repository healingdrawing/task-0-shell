// Released under an MIT license. See LICENSE.

package commands

import (
	"penishell/internal/common/interface/cell"
	"penishell/internal/common/interface/integer"
	"penishell/internal/common/type/chn"
	"penishell/internal/common/type/create"
	"penishell/internal/common/validate"
)

func isChan(args cell.I) cell.I {
	v := validate.Fixed(args, 1, 1)

	return create.Bool(chn.Is(v[0]))
}

func makeChan(args cell.I) cell.I {
	v := validate.Fixed(args, 0, 1)

	n := int64(0)
	if len(v) > 0 {
		n = integer.Value(v[0])
	}

	return chn.New(n)
}
