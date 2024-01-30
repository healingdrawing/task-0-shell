// Released under an MIT license. See LICENSE.

package commands

import (
	"math"
	"math/rand"
	"time"

	"penishell/internal/common"
	"penishell/internal/common/interface/cell"
	"penishell/internal/common/interface/integer"
	"penishell/internal/common/interface/rational"
	"penishell/internal/common/type/create"
	"penishell/internal/common/type/num"
	"penishell/internal/common/validate"
)

func init() { //nolint:gochecknoinits
	rand.Seed(time.Now().UnixNano())
}

func isNumber(args cell.I) cell.I {
	v := validate.Fixed(args, 1, 1)

	return create.Bool(num.Is(v[0]))
}

func number(args cell.I) cell.I {
	v := validate.Fixed(args, 1, 1)

	if r, ok := v[0].(rational.I); ok {
		return num.Rat(r.Rat())
	}

	return num.New(common.String(v[0]))
}

func random(args cell.I) cell.I {
	v := validate.Fixed(args, 0, 1)

	n := math.MaxInt32
	if len(v) == 1 {
		n = int(integer.Value(v[0]))
	}

	return num.Int(rand.Intn(n)) //nolint:gosec
}
