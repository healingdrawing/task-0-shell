// Released under an MIT license. See LICENSE.

package commands

import (
	"penishell/internal/common/interface/cell"
	"penishell/internal/common/interface/rational"
	"penishell/internal/common/type/pair"
	"penishell/internal/common/type/sym"
	"penishell/internal/common/validate"
)

func equal(args cell.I) cell.I {
	v, rest := validate.Variadic(args, 2, 2)

	for {
		if !v[0].Equal(v[1]) {
			return pair.Null
		}

		if rest == pair.Null {
			return sym.True
		}

		v[0] = v[1]
		v[1] = pair.Car(rest)

		rest = pair.Cdr(rest)
	}
}

func ge(args cell.I) cell.I {
	v, rest := validate.Variadic(args, 2, 2)

	prev := rational.Number(v[0])
	curr := rational.Number(v[1])

	for {
		if prev.Cmp(curr) < 0 {
			return pair.Null
		}

		if rest == pair.Null {
			return sym.True
		}

		prev = curr
		curr = rational.Number(pair.Car(rest))

		rest = pair.Cdr(rest)
	}
}

func gt(args cell.I) cell.I {
	v, rest := validate.Variadic(args, 2, 2)

	prev := rational.Number(v[0])
	curr := rational.Number(v[1])

	for {
		if prev.Cmp(curr) <= 0 {
			return pair.Null
		}

		if rest == pair.Null {
			return sym.True
		}

		prev = curr
		curr = rational.Number(pair.Car(rest))

		rest = pair.Cdr(rest)
	}
}

func le(args cell.I) cell.I {
	v, rest := validate.Variadic(args, 2, 2)

	prev := rational.Number(v[0])
	curr := rational.Number(v[1])

	for {
		if prev.Cmp(curr) > 0 {
			return pair.Null
		}

		if rest == pair.Null {
			return sym.True
		}

		prev = curr
		curr = rational.Number(pair.Car(rest))

		rest = pair.Cdr(rest)
	}
}

func lt(args cell.I) cell.I {
	v, rest := validate.Variadic(args, 2, 2)

	prev := rational.Number(v[0])
	curr := rational.Number(v[1])

	for {
		if prev.Cmp(curr) >= 0 {
			return pair.Null
		}

		if rest == pair.Null {
			return sym.True
		}

		prev = curr
		curr = rational.Number(pair.Car(rest))

		rest = pair.Cdr(rest)
	}
}
