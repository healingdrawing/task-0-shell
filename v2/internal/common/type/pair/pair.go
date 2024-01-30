// Released under an MIT license. See LICENSE.

// Package pair provides penishell's cons cell type.
package pair

import (
	"fmt"

	"penishell/internal/common"
	"penishell/internal/common/interface/cell"
	"penishell/internal/common/interface/literal"
)

const name = "cons"

//nolint:gochecknoglobals
var (
	// Null is the empty list. It is also used to mark the end of a list.
	Null cell.I
)

// T (pair) is a cons cell.
type T struct {
	car cell.I
	cdr cell.I
}

type pair = T

// Equal returns true if c is a pair with elements that are equal to p's.
func (p *pair) Equal(c cell.I) bool {
	if p == Null && c == Null {
		return true
	}

	return p.car.Equal(Car(c)) && p.cdr.Equal(Cdr(c))
}

// Literal returns the literal representation of the pair p.
func (p *pair) Literal() string {
	return p.string(literal.String)
}

// Name returns the name for a pair type.
func (p *pair) Name() string {
	return name
}

// String returns the text representation of the pair p.
func (p *pair) String() string {
	return p.string(common.String)
}

func (p *pair) string(toString func(cell.I) string) string {
	s := ""

	improper := false

	tail := Cdr(p)
	if !Is(tail) {
		improper = true
		s += "(|" + name + " "
	}

	sublist := false

	head := Car(p)
	if Is(head) && Is(Cdr(head)) {
		sublist = true
		s += "("
	}

	if head == nil {
		s += "()"
	} else if head != Null {
		s += toString(head)
	}

	if sublist {
		s += ")"
	}

	if !improper && tail == Null {
		return s
	}

	s += " "
	if tail == nil {
		s += "()"
	} else {
		s += toString(tail)
	}

	if improper {
		s += "|)"
	}

	return s
}

// Functions specific to pair.

// Car returns the car/head/first member of the pair c.
// If c is not a pair, this function will panic.
func Car(c cell.I) cell.I {
	return To(c).car
}

// Cdr returns the cdr/tail/rest member of the pair c.
// If c is not a pair, this function will panic.
func Cdr(c cell.I) cell.I {
	return To(c).cdr
}

// Caar returns the car of the car of the pair c.
// A non-pair value where a pair is expected will cause a panic.
func Caar(c cell.I) cell.I {
	return To(To(c).car).car
}

// Cadr returns the car of the cdr of the pair c.
// A non-pair value where a pair is expected will cause a panic.
func Cadr(c cell.I) cell.I {
	return To(To(c).cdr).car
}

// Cdar returns the cdr of the car of the pair c.
// A non-pair value where a pair is expected will cause a panic.
func Cdar(c cell.I) cell.I {
	return To(To(c).car).cdr
}

// Cddr returns the cdr of the cdr of the pair c.
// A non-pair value where a pair is expected will cause a panic.
func Cddr(c cell.I) cell.I {
	return To(To(c).cdr).cdr
}

// Caddr returns the car of the cdr of the cdr of the pair c.
// A non-pair value where a pair is expected will cause a panic.
func Caddr(c cell.I) cell.I {
	return To(To(To(c).cdr).cdr).car
}

// Cons conses h and t together to form a new pair.
func Cons(h, t cell.I) cell.I {
	return &pair{car: h, cdr: t}
}

// SetCar sets the car/head/first of the pair c to value.
// If c is not a pair, this function will panic.
func SetCar(c, value cell.I) {
	To(c).car = value
}

// SetCdr sets the cdr/tail/rest of the pair c to value.
// If c is not a pair, this function will panic.
func SetCdr(c, value cell.I) {
	To(c).cdr = value
}

// A compiler-checked list of interfaces this type satisfies. Never called.
func implements() { //nolint:deadcode,unused
	var t pair

	// The pair type is a cell.
	_ = cell.I(&t)

	// The pair type has a literal representation.
	_ = literal.I(&t)

	// The pair type is a stringer.
	_ = fmt.Stringer(&t)
}

func init() { //nolint:gochecknoinits
	pair := &pair{}
	pair.car = pair
	pair.cdr = pair

	Null = cell.I(pair)
}
