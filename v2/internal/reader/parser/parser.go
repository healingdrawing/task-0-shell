// Released under an MIT license. See LICENSE.

// Package parser provides a recursive descent parser for the penishell language.
package parser

import (
	"fmt"
	"strconv"
	"strings"

	"penishell/internal/common"
	"penishell/internal/common/interface/cell"
	"penishell/internal/common/struct/token"
	"penishell/internal/common/type/list"
	"penishell/internal/common/type/num"
	"penishell/internal/common/type/pair"
	"penishell/internal/common/type/status"
	"penishell/internal/common/type/str"
	"penishell/internal/common/type/sym"

	"github.com/michaelmacinnis/adapted"
)

// T holds the state of the parser.
type T struct {
	ahead int             // Lookahead count.
	emit  func(cell.I)    // Function to call to emit a parsed command.
	item  func() *token.T // Function to call to get another token.
	token *token.T        // Token lookahead.

	// Completion state.
	current cell.I // The command being parsed, so far.
}

// New creates a new parser.
// It connects a producer of tokens with a consumer of cells.
func New(emit func(cell.I), item func() *token.T) *T {
	return &T{emit: emit, item: item, current: pair.Null}
}

// Copy copies the current parser but replaces its emit and item functions.
func (p *T) Copy(emit func(cell.I), item func() *token.T) *T {
	c := *p

	c.emit = emit
	c.item = item

	return &c
}

// Current returns the command currently being parsed.
func (p *T) Current() cell.I {
	return p.current
}

// Parse consumes tokens and emits cells until there are no more tokens.
func (p *T) Parse() (err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		switch r := r.(type) {
		case error:
			err = r
		case string:
			err = common.Error(r)
		default:
			err = fmt.Errorf("unexpected error: %v", r) //nolint:goerr113
		}
	}()

	for t := p.peek(); t != nil; t = p.peek() {
		if t.Is('\n') {
			p.consume()

			continue
		}

		p.emit(p.possibleBackground())
	}

	return nil
}

func (p *T) consume() *token.T {
	if p.ahead == 0 {
		panic("nothing to consume.")
	}

	t := p.token

	p.ahead = 0
	p.token = nil

	return t
}

func (p *T) check(c cell.I) cell.I {
	if c == nil {
		t := p.peek()

		loc := t.Source()
		l := loc.Name
		x := strconv.Itoa(loc.Char)
		y := strconv.Itoa(loc.Line)

		panic(l + ":" + y + ":" + x + ": unexpected '" + t.Source().Text + "'")
	}

	return c
}

func (p *T) expect(cs ...token.Class) {
	if p.peek().Is(cs...) {
		p.consume()

		return
	}

	// Make a nice error message.
	n := len(cs)
	e := make([]string, n)

	for i, c := range cs[:n-1] {
		e[i] = c.String()
	}

	l := cs[n-1].String()
	if n > 2 { //nolint:gomnd
		l = ", or " + l
	} else if n > 1 {
		l = " or " + l
	}

	l = strings.Join(e, ", ") + l
	s := p.peek().Value()

	panic("expected " + l + ` got "` + s + `"`)
}

func (p *T) peek() *token.T {
	if p.ahead > 0 {
		return p.token
	}

	t := p.item()

	p.token = t
	p.ahead = 1

	return t
}

// T state functions.

// <possibleBackground> ::= <command> '&'?
func (p *T) possibleBackground() cell.I {
	c := p.command()

	t := p.peek()
	if t.Is(token.Background) {
		p.consume()

		c = list.New(sym.Token(t), c)
	}

	return c
}

// <command> ::= <possibleAndf> (Orf <possibleAndf>)* .
func (p *T) command() cell.I {
	c := p.possibleAndf()

	t := p.peek()
	if t.Is(token.Orf) {
		c = list.New(sym.Token(t), c)

		for p.peek().Is(token.Orf) {
			p.consume()
			c = list.Append(c, p.possibleAndf())
		}
	}

	return c
}

// <possibleAndf> ::= <possiblePipeline> (Andf <possiblePipeline>)* .
func (p *T) possibleAndf() cell.I {
	c := p.possiblePipeline()

	t := p.peek()
	if t.Is(token.Andf) {
		c = list.New(sym.Token(t), c)

		for p.peek().Is(token.Andf) {
			p.consume()
			c = list.Append(c, p.possiblePipeline())
		}
	}

	return c
}

// <possiblePipeline> ::= <possibleSequence> (Pipe <possiblePipeline>)?
func (p *T) possiblePipeline() cell.I {
	c := p.possibleSequence()

	if p.peek().Is(token.Pipe) {
		s := sym.Token(p.consume())

		c = pair.Cons(p.possiblePipeline(), c)
		c = pair.Cons(s, c)
	}

	return c
}

// <possibleSequence> ::= <possibleRedirection> (';' <possibleRedirection>)* .
func (p *T) possibleSequence() cell.I {
	c := p.possibleRedirection()

	if p.peek().Is(';') {
		c = list.New(sym.New("block"), c)

		for p.peek().Is(';') {
			p.consume()

			c = list.Append(c, p.possibleRedirection())
		}
	}

	return c
}

// <possibleRedirection> ::= <possibleSustitution> (Redirect <expression>)* .
func (p *T) possibleRedirection() cell.I {
	c := p.possibleSubstitution()

	for p.peek().Is(token.Redirect) {
		s := sym.Token(p.consume())
		c = list.New(s, p.check(p.implicitJoin(p.element())), c)

		for p.peek().Is(token.Space) {
			p.consume()
		}
	}

	return c
}

// <possibleSubstitution> ::= <statement> (Substitute <command> ')' <statement>?)* .
func (p *T) possibleSubstitution() cell.I {
	c := p.statement()
	if c == nil {
		return c
	}

	if p.peek().Is(token.Substitute) {
		c = pair.Cons(sym.New("process-substitution"), c)

		for p.peek().Is(token.Substitute) {
			s := sym.Token(p.consume())
			l := pair.Cons(s, p.element())
			c = list.Append(c, l)

			if !p.peek().Is(token.Substitute) {
				s := p.statement()
				if s != nil {
					c = list.Join(c, s)
				}
			}
		}
	}

	return c
}

func (p *T) braces() (c cell.I) {
	if p.peek().Is('{') {
		p.consume()

		n := p.peek()

		switch {
		case n.Is('\n'):
			p.consume()

			c = p.subStatement()
		case n.Is('{'):
			c = p.braces()
			p.expect('}')
		default:
			c = p.implicitJoin(p.element())
			c = pair.Cons(c, pair.Null)

			p.expect('}')
		}
	}

	return
}

func (p *T) assignments() (c, l cell.I) {
	l = pair.Null

	for {
		for p.peek().Is(token.Space) {
			p.consume()
		}

		c = p.braces()
		if c != nil {
			break
		}

		c = p.element()
		if c == nil {
			break
		}

		e := p.peek()
		if sym.Is(c) && e.Is(token.Symbol) && e.Value() == "=" {
			p.consume()

			v := p.check(p.implicitJoin(p.element()))

			l = list.Append(l, list.New(sym.New("export"), c, v))

			continue
		} else {
			c = p.implicitJoin(c)

			break
		}
	}

	return c, l
}

func (p *T) statement() (c cell.I) {
	// Reset current command.
	p.current = pair.Null

	c, l := p.assignments()
	if l != pair.Null {
		defer func() {
			if c != nil {
				c = list.Join(l, pair.Cons(c, pair.Null))
			} else {
				c = l
			}

			c = pair.Cons(sym.New("block"), c)
		}()
	}

	if c == nil {
		return
	}

	c = pair.Cons(c, pair.Null)

	for {
		p.current = c

		if p.peek().Is(token.Space) {
			p.consume()

			continue
		}

		t := p.braces()
		if t == nil {
			t = p.implicitJoin(p.element())
			if t == nil {
				break
			}

			t = pair.Cons(t, pair.Null)
		}

		c = list.Join(c, t)
	}

	return c
}

func (p *T) subStatement() cell.I {
	c := p.block()

	p.expect('}')

	for p.peek().Is(token.Space) {
		p.consume()
	}

	s := p.statement()
	if s != nil {
		c = list.Join(c, s)
	}

	return c
}

func (p *T) block() cell.I {
	c := pair.Null

	for !p.peek().Is('}') {
		if p.peek().Is('\n') {
			p.consume()

			continue
		}

		c = list.Append(c, p.check(p.possibleBackground()))
	}

	return c
}

func (p *T) implicitJoin(c cell.I) cell.I {
	if c == nil {
		return nil
	}

	l := list.New(c)

	for t := p.element(); t != nil; t = p.element() {
		l = list.Append(l, t)
	}

	if list.Length(l) == 1 {
		return c
	}

	l = pair.Cons(sym.New(""), l)

	return pair.Cons(sym.New("mend"), l)
}

func (p *T) element() cell.I {
	if p.peek().Is('`') {
		p.consume()

		c := p.check(p.value())

		c = pair.Cons(sym.New("capture"), list.New(c))
		c = list.New(sym.New("splice"), c)

		return c
	}

	return p.expression()
}

func (p *T) expression() cell.I {
	if p.peek().Is('$') {
		p.consume()

		c := p.braces()
		if c == nil {
			c = p.check(p.expression())
		} else {
			c = pair.Car(c)
		}

		return list.New(sym.New("resolve"), p.check(c))
	}

	return p.value()
}

func (p *T) meta(c cell.I) cell.I {
	t := pair.Car(c)

	if !sym.Is(t) {
		panic("meta command must start with a symbol not " + t.Name())
	}

	var create func(string) cell.I = nil

	switch sym.To(t).String() {
	case "cons":
		return pair.Cons(pair.Cadr(c), pair.Caddr(c))

	case "number":
		create = num.New

	case "status":
		create = status.New

	case "symbol":
		create = sym.New
	}

	if create == nil {
		panic("invalid meta command")
	}

	t = pair.Cadr(c)

	arg, ok := t.(fmt.Stringer)
	if ok {
		return create(arg.String())
	}

	// TODO: What case are we handling here?
	return num.New(arg.String())
}

func (p *T) value() cell.I {
	t := p.peek()

	meta := false
	if t.Is(token.MetaOpen) {
		meta = true
	} else if !t.Is('(') {
		return p.word()
	}

	p.consume()

	c := p.command()
	if c == nil {
		t := p.peek()
		if t.Is(')') {
			p.consume()

			return pair.Null
		}

		panic("unexpected '" + t.Source().Text + "'")
	}

	if meta {
		p.expect(token.MetaClose)

		return p.meta(c)
	}

	p.expect(')')

	return c
}

func (p *T) symbol(t *token.T) cell.I {
	if t.Is(token.Symbol) {
		p.consume()

		return sym.Token(t)
	}

	return nil
}

func (p *T) word() cell.I {
	t := p.peek()
	if t.Is(token.DollarSingleQuoted) {
		p.consume()

		text := t.Value()

		s, err := adapted.ActualBytes(text[2 : len(text)-1])
		if err != nil {
			panic(err.Error())
		}

		return str.New(s)
	}

	if t.Is(token.DoubleQuoted) {
		p.consume()

		text := t.Value()

		s, err := adapted.ActualBytes(text[1 : len(text)-1])
		if err != nil {
			panic(err.Error())
		}

		return list.New(sym.New("interpolate"), str.New(s))
	}

	if t.Is(token.SingleQuoted) {
		p.consume()

		s := t.Value()

		return str.New(s[1 : len(s)-1])
	}

	return p.symbol(t)
}
