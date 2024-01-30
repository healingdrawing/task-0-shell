package commands

import (
	"penishell/internal/common/interface/cell"
	"penishell/internal/common/interface/integer"
	"penishell/internal/common/type/list"
	"penishell/internal/common/type/num"
	"penishell/internal/common/type/pair"
	"penishell/internal/common/validate"
)

// ListMethods returns a mapping of names to list methods.
func ListMethods() map[string]func(cell.I, cell.I) cell.I {
	return map[string]func(cell.I, cell.I) cell.I{
		"append":   appendMethod,
		"extend":   extend,
		"get":      get,
		"head":     head,
		"length":   length,
		"reverse":  reverse,
		"set-head": setHead,
		"set-tail": setTail,
		"slice":    slice,
		"tail":     tail,
	}
}

func appendMethod(s, args cell.I) cell.I {
	v := validate.Fixed(args, 1, 1)

	self := pair.To(s)

	return list.Append(self, v...)
}

func extend(s, args cell.I) cell.I {
	v := validate.Fixed(args, 1, 1)

	self := pair.To(s)

	return list.Join(self, v[0])
}

func get(s, args cell.I) cell.I {
	v, args := validate.Variadic(args, 0, 1)

	self := pair.To(s)

	i := int64(0)
	if len(v) != 0 {
		i = integer.Value(v[0])
	}

	var dflt cell.I
	if args != pair.Null {
		dflt = args
	}

	return pair.Car(list.Tail(self, i, dflt))
}

func head(s, _ cell.I) cell.I {
	return pair.Car(pair.To(s))
}

func length(s, args cell.I) cell.I {
	validate.Fixed(args, 0, 0)

	return num.Int(int(list.Length(pair.To(s))))
}

func reverse(s, args cell.I) cell.I {
	validate.Fixed(args, 0, 0)

	return list.Reverse(pair.To(s))
}

func setHead(s, args cell.I) cell.I {
	v := pair.Car(args)
	pair.SetCar(s, v)

	return v
}

func setTail(s, args cell.I) cell.I {
	v := pair.Car(args)
	pair.SetCdr(s, v)

	return v
}

func slice(s, args cell.I) cell.I {
	v := validate.Fixed(args, 1, 2)

	start := integer.Value(v[0])
	end := int64(0)

	if len(v) == 2 { //nolint:gomnd
		end = integer.Value(v[1])
	}

	return list.Slice(s, start, end)
}

func tail(s, _ cell.I) cell.I {
	return pair.Cdr(pair.To(s))
}
