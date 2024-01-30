// Released under an MIT license. See LICENSE.

package commands

import (
	"fmt"
	"strings"

	"penishell/internal/common"
	"penishell/internal/common/interface/cell"
	"penishell/internal/common/interface/integer"
	"penishell/internal/common/type/create"
	"penishell/internal/common/type/num"
	"penishell/internal/common/type/pair"
	"penishell/internal/common/type/str"
	"penishell/internal/common/validate"
)

// StringFunctions returns a mapping of names to string methods.
func StringFunctions() map[string]func(cell.I) cell.I {
	return map[string]func(cell.I) cell.I{
		"format":      sprintf,
		"length":      slength,
		"lower":       lower,
		"replace":     sreplace,
		"slice":       sslice,
		"trim-prefix": trimPrefix,
		"trim-suffix": trimSuffix,
		"upper":       upper,
	}
}

func isString(args cell.I) cell.I {
	v := validate.Fixed(args, 1, 1)

	return create.Bool(str.Is(v[0]))
}

func lower(args cell.I) cell.I {
	v := validate.Fixed(args, 1, 1)

	return str.New(strings.ToLower(common.String(v[0])))
}

func makeString(args cell.I) cell.I {
	v := validate.Fixed(args, 1, 1)

	return str.New(common.String(v[0]))
}

func slength(args cell.I) cell.I {
	v := validate.Fixed(args, 1, 1)

	return num.Int(len(common.String(v[0])))
}

func sslice(args cell.I) cell.I {
	v := validate.Fixed(args, 2, 3)

	s := common.String(v[0])

	length := int64(len(s))

	start := integer.Value(v[1])
	if start < 0 {
		panic("slice starts before first element")
	} else if start > length {
		start = length
	}

	end := length
	if len(v) == 3 { //nolint:gomnd
		end = integer.Value(v[2])
		if end > length {
			end = length
		} else if end < 0 {
			end = length + end
		}
	}

	if end < start {
		panic("end of slice before start")
	}

	return str.New(s[start:end])
}

func sreplace(args cell.I) cell.I {
	v := validate.Fixed(args, 3, 4)

	s := common.String(v[0])
	old := common.String(v[1])
	replacement := common.String(v[2])

	n := -1
	// The 4th argument, if passed, limits the number of replacements.
	if len(v) == 4 { //nolint:gomnd
		n = int(integer.Value(v[3]))
	}

	return str.New(strings.Replace(s, old, replacement, n))
}

// TODO: Extend penishell types to play nicer with fmt and Sprintf.
func sprintf(args cell.I) cell.I {
	v, args := validate.Variadic(args, 1, 1)

	argv := []interface{}{}

	for args != pair.Null {
		argv = append(argv, pair.Car(args))
		args = pair.Cdr(args)
	}

	return str.New(fmt.Sprintf(common.String(v[0]), argv...))
}

func trimPrefix(args cell.I) cell.I {
	v := validate.Fixed(args, 2, 2)

	return str.New(strings.TrimPrefix(common.String(v[0]), common.String(v[1])))
}

func trimSuffix(args cell.I) cell.I {
	v := validate.Fixed(args, 2, 2)

	return str.New(strings.TrimSuffix(common.String(v[0]), common.String(v[1])))
}

func upper(args cell.I) cell.I {
	v := validate.Fixed(args, 1, 1)

	return str.New(strings.ToUpper(common.String(v[0])))
}
