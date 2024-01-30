// Released under an MIT license. See LICENSE.

package validate

import (
	"fmt"

	"penishell/internal/common/interface/cell"
	"penishell/internal/common/type/list"
	"penishell/internal/common/type/pair"
)

// Variadic checks that there are at least min to max arguments and returns
// these as an array. Any remaining arguments are returned as a list.
func Variadic(actual cell.I, min, max int) ([]cell.I, cell.I) {
	expected := make([]cell.I, 0, max)

	for i := 0; i < max; i++ {
		if actual == pair.Null {
			if i < min {
				s := Count(min, "argument", "s")
				panic(fmt.Sprintf("expected %s, passed %d", s, i))
			}

			break
		}

		expected = append(expected, pair.Car(actual))

		actual = pair.Cdr(actual)
	}

	return expected, actual
}

// Fixed returns min to max arguments as an array.
func Fixed(actual cell.I, min, max int) []cell.I {
	expected, rest := Variadic(actual, min, max)
	if rest != pair.Null {
		s := Count(max, "argument", "s")
		n := int(list.Length(actual))

		panic(fmt.Sprintf("expected %s, passed %d", s, n))
	}

	return expected
}

// Count returns a human-readable count.
func Count(n int, label, p string) string {
	if n == 1 {
		p = ""
	}

	return fmt.Sprintf("%d %s%s", n, label, p)
}
