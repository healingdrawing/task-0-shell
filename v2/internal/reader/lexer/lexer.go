// Released under an MIT license. See LICENSE.

// Package lexer provides a lexical scanner for the penishell language.
//
// The penishell lexer adapts the state function approach used by Go's text/template
// lexer and described in detail in Rob Pike's talk "Lexical Scanning in Go".
// See https://talks.golang.org/2011/lex.slide for more information.
package lexer

import (
	"strings"
	"unicode/utf8"

	"penishell/internal/common/struct/loc"
	"penishell/internal/common/struct/token"
)

// T holds the state of the scanner.
type T struct {
	expected []string // Completion candidates.

	bytes string   // Buffer being scanned.
	first int      // Index of the current token's first byte.
	index int      // Index of the current byte.
	queue []string // Buffers waiting to be scanned.
	runes int      // Runes scanned on the current line.
	saved action   // Escaped action.
	state action   // Current action.

	source loc.T

	tokens chan *token.T
}

const qsize = 2

// New creates a new lexer/scanner. Label can be a file name or other identifier.
func New(label string) *T {
	l := &T{
		source: loc.T{
			Char: 1,
			Line: 1,
			Name: label,
		},
	}

	l.state = skipWhitespace

	return l
}

// Copy makes a copy of the lexer with its own token channel.
// A copy is useful for doing partial parses for command completion.
func (l *T) Copy() *T {
	c := *l

	copy(c.queue, l.queue)

	c.tokens = make(chan *token.T, qsize)

	return &c
}

// Expected returns the list of expected strings. (Command completion).
func (l *T) Expected() []string {
	return l.expected
}

// Scan passes a text buffer to the lexer for scanning.
// If a buffer is currently being scanned, the new buffer will
// be appended to the list of buffers waiting to be scanned.
func (l *T) Scan(text string) {
	l.queue = append(l.queue, text)
}

// Text is used to return the text corresponding to the current token.
func (l *T) Text() string {
	return l.bytes[l.first:l.index]
}

// Token returns the next scanned token, or nil if no token is available.
func (l *T) Token() *token.T {
	for {
		l.gather()

		if len(l.bytes) == 0 {
			return nil
		}

		select {
		case t := <-l.tokens:
			return t
		default:
			state := l.state(l)
			if state != nil {
				l.state = state
			} else {
				close(l.tokens)
			}
		}
	}
}

type action func(*T) action

const eof = -1

func (l *T) accept(r token.Class, w int) {
	if r == '\n' {
		// Because we update lines here, if we emit a newline
		// it will be reported as being part of the next line.
		// We fix this when emitting the newline.
		l.source.Line++
		l.runes = 1
	} else {
		l.runes++
	}

	l.index += w
}

func (l *T) emit(c token.Class, v string) {
	source := l.source
	if c == '\n' {
		// Report newline as part of previous line.
		source.Line--
	}

	source.Text = strings.TrimRight(l.bytes[l.first:], "\n")

	t := token.New(c, v, &source)

	l.tokens <- t
	l.skip()
}

func (l *T) escape(escaped, a action) action {
	l.saved = escaped

	return a
}

func (l *T) gather() {
	if len(l.queue) == 0 {
		return
	}

	length := len(l.bytes)
	bytes := strings.Join(l.queue, "")

	if length > 0 && l.first < length {
		// Prepend leftover to new bytes.
		bytes = l.bytes[l.first:] + bytes
	} else {
		l.source.Char = 1
		l.runes = 1
	}

	l.queue = nil
	l.bytes = bytes
	l.index -= l.first
	l.first = 0
	l.tokens = make(chan *token.T, qsize)
}

func (l *T) next() token.Class {
	r, w := l.peek()
	l.accept(r, w)

	return r
}

func (l *T) peek() (token.Class, int) {
	r, w := rune(eof), 0
	if l.index < len(l.bytes) {
		r, w = utf8.DecodeRuneInString(l.bytes[l.index:])
	}

	return token.Class(r), w
}

func (l *T) resume() action {
	resumed := l.saved
	l.saved = nil

	return resumed
}

func (l *T) skip() {
	l.source.Char = l.runes
	l.first = l.index
}

// T states.

func afterAmpersand(l *T) action {
	r, w := l.peek()

	l.expected = []string{" ", "& "}

	switch r {
	case eof:
		return nil

	case '&':
		l.accept(r, w)
		l.emit(token.Andf, operator(l.Text()))

		return skipWhitespace
	default:
		l.emit(token.Background, operator(l.Text()))

		return collectHorizontalSpace
	}
}

func afterDollar(l *T) action {
	r, w := l.peek()

	// TODO: Indicate other completions are possible.
	l.expected = []string{"'"}

	switch r {
	case eof:
		return nil
	case '$': // Special-case to recognize $$.
		l.emit(r, l.Text())
		l.accept(r, w)
		l.emit(token.Symbol, l.Text())
	case '\'':
		l.accept(r, w)

		return scanDollarSingleQuoted
	case '\t', '\n', ' ', '"', '#', '&', '(',
		')', ';', '<', '>', '`', '|', '}':
		l.emit(token.Symbol, l.Text())
	default:
		l.emit('$', l.Text())
	}

	return collectHorizontalSpace
}

func afterDoubleGreaterThan(l *T) action {
	r, w := l.peek()

	l.expected = []string{" ", "& "}

	switch r {
	case eof:
		return nil
	case '&':
		l.accept(r, w)
		l.emit(token.Redirect, operator(l.Text()))
	default:
		l.emit(token.Redirect, operator(l.Text()))
	}

	return skipHorizontalSpace
}

func afterGreaterThan(l *T) action {
	r, w := l.peek()

	l.expected = []string{" ", "& ", "> ", ">& ", ">&| ", "| "}

	switch r {
	case eof:
		return nil
	case '&':
		l.accept(r, w)

		return afterGreaterThanAmpersand
	case '>':
		l.accept(r, w)

		return afterDoubleGreaterThan
	case '|':
		l.accept(r, w)
		l.emit(token.Redirect, operator(l.Text()))
	default:
		l.emit(token.Redirect, operator(l.Text()))
	}

	return skipHorizontalSpace
}

func afterGreaterThanAmpersand(l *T) action {
	r, w := l.peek()

	l.expected = []string{" ", "| "}

	switch r {
	case eof:
		return nil
	case '|':
		l.accept(r, w)
		l.emit(token.Redirect, operator(l.Text()))
	default:
		l.emit(token.Redirect, operator(l.Text()))
	}

	return skipHorizontalSpace
}

func afterOpenParen(l *T) action {
	r, w := l.peek()

	switch r {
	case eof:
		return nil
	case '|':
		l.accept(r, w)
		l.emit(token.MetaOpen, l.Text())
	default:
		l.emit('(', l.Text())
	}

	return skipWhitespace
}

func afterPipe(l *T) action {
	r, w := l.peek()

	l.expected = []string{" ", "& ", "<(", ">(", "| "}

	switch r {
	case eof:
		return nil
	case '&':
		l.accept(r, w)
		l.emit(token.Pipe, operator(l.Text()))

		return skipWhitespace
	case ')':
		l.accept(r, w)
		l.emit(token.MetaClose, l.Text())

		return collectHorizontalSpace
	case '<':
		l.accept(r, w)
		l.emit(token.Substitute, operator(l.Text()))

		return skipHorizontalSpace
	case '>':
		l.accept(r, w)
		l.emit(token.Substitute, operator(l.Text()))

		return skipHorizontalSpace
	case '|':
		l.accept(r, w)
		l.emit(token.Orf, operator(l.Text()))

		return skipWhitespace
	default:
		l.emit(token.Pipe, operator(l.Text()))

		return skipWhitespace
	}
}

func collectHorizontalSpace(l *T) action {
	for {
		r, w := l.peek()

		switch r {
		case eof:
			return nil
		case '\n':
			l.accept(r, w)
			l.emit(r, l.Text())

			return skipWhitespace
		case '#':
			l.accept(r, w)

			return skipComment
		case '\t', ' ':
			l.accept(r, w)

			continue
		default:
			s := l.Text()
			if len(s) > 0 {
				l.emit(token.Space, s)
			}

			return skipHorizontalSpace
		}
	}
}

func escapeNewline(l *T) action {
	r, w := l.peek()

	switch r {
	case eof:
		return nil
	case '\n':
		l.accept(r, w)
		l.skip()
	default:
		l.accept(r, w)
		l.saved = nil

		return scanSymbol
	}

	return l.resume()
}

func escapeNextCharacter(l *T) action {
	r := l.next()

	if r == eof {
		return nil
	}

	return l.resume()
}

func scanDollarSingleQuoted(l *T) action {
	for {
		c := l.next()

		switch c {
		case eof:
			return nil

		case '\'':
			l.emit(token.DollarSingleQuoted, l.Text())

			return collectHorizontalSpace

		case '\\':
			return l.escape(scanDollarSingleQuoted, escapeNextCharacter)

		default: // Continue and get next character.
		}
	}
}

func scanDoubleQuoted(l *T) action {
	for {
		c := l.next()

		switch c {
		case eof:
			return nil

		case '"':
			l.emit(token.DoubleQuoted, l.Text())

			return collectHorizontalSpace

		case '\\':
			return l.escape(scanDoubleQuoted, escapeNextCharacter)

		default: // Continue and get next character.
		}
	}
}

func scanSingleQuoted(l *T) action {
	for {
		r := l.next()

		switch r {
		case eof:
			return nil

		case '\'':
			l.emit(token.SingleQuoted, l.Text())

			return collectHorizontalSpace

		default: // Continue and get next character.
		}
	}
}

func scanSymbol(l *T) action {
	// Characters that can be in a symbol:
	// '!', '%', '*', '+', '-',
	// '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	// '?',
	// 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
	// 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	// '[', ']', '^', '_',
	// 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	// 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	//
	// And a trailing '$'.
	// (A trailing '$' should not be used in a variable name).
	//
	// The characters ',', '.', '/', ':', '=', '@', and '~' result in a
	// symbol of one character. (These characters should not be used in
	// variable names).
	for {
		r, w := l.peek()

		switch r {
		case eof:
			return nil

		case '\t', '\n', ' ', '"', '#', '&', '\'', '(',
			')', ';', '<', '>', '`', '{', '|', '}':
			l.emit(token.Symbol, l.Text())

			return collectHorizontalSpace

		case ',', '.', '/', ':', '=', '@', '~':
			s := l.Text()
			if len(s) > 0 {
				l.emit(token.Symbol, s)
			}

			l.accept(r, w)
			l.emit(token.Symbol, l.Text())

			return collectHorizontalSpace

		case '$':
			s := l.Text()
			if len(s) > 0 {
				l.emit(token.Symbol, s)
			}

			l.accept(r, w)

			return afterDollar

		case '\\':
			l.accept(r, w)

			return l.escape(scanSymbol, escapeNextCharacter)

		default:
			l.accept(r, w)
		}
	}
}

func skipComment(l *T) action {
	for {
		r := l.next()

		switch r {
		case eof:
			return nil

		case '\n':
			l.emit('\n', l.Text())

			return skipWhitespace

		default: // Continue and get next character.
		}
	}
}

func skipHorizontalSpace(l *T) action {
	return startState(l, skipHorizontalSpace, "\t ")
}

func skipWhitespace(l *T) action {
	return startState(l, skipWhitespace, "\n\t ")
}

func startState(l *T, state action, ignore string) action {
	l.expected = []string{}

	for {
		r := l.next()

		if strings.ContainsRune(ignore, rune(r)) {
			l.skip()

			continue
		}

		switch r {
		case eof:
			return nil

		case '\n', ')', ';', '`', '{', '}':
			l.emit(r, l.Text())

			return collectHorizontalSpace

		case '<':
			l.emit(token.Redirect, operator(l.Text()))

			return skipHorizontalSpace

		case '\\':
			return l.escape(state, escapeNewline)

		case ',', '.', '/', ':', '=', '@', '~':
			l.emit(token.Symbol, l.Text())

			return collectHorizontalSpace

		default:
			return initial(r)
		}
	}
}

// Helper functions.

func initial(r token.Class) action {
	a, ok := map[token.Class]action{
		'"':  scanDoubleQuoted,
		'#':  skipComment,
		'$':  afterDollar,
		'&':  afterAmpersand,
		'\'': scanSingleQuoted,
		'(':  afterOpenParen,
		'>':  afterGreaterThan,
		'|':  afterPipe,
	}[r]

	if ok {
		return a
	}

	return scanSymbol
}

func operator(s string) string {
	return map[string]string{
		"&":   "spawn",
		"&&":  "and",
		"<":   "input-from",
		">":   "output-to",
		">&":  "output-errors-to",
		">&|": "output-errors-clobbers",
		">>":  "append-output-to",
		">>&": "append-output-errors-to",
		">|":  "output-clobbers",
		"|":   "pipe-output-to",
		"|&":  "pipe-output-errors-to",
		"|<":  "-named-pipe-input-from",
		"|>":  "-named-pipe-output-to",
		"||":  "or",
	}[s]
}
