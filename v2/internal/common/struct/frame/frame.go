// Released under an MIT license. See LICENSE.

// Package frame provides penishell's call stack frame type.
package frame

import (
	"penishell/internal/common/interface/reference"
	"penishell/internal/common/interface/scope"
	"penishell/internal/common/struct/loc"
)

// T (frame) is stack frame or activation record.
type T struct {
	previous *frame
	scope    scope.I
	source   loc.T
}

type frame = T

// Dup creates a duplicate of the frame f with a new scope s.
func Dup(s scope.I, f *frame) *frame {
	dup := *f
	dup.scope = s

	return &dup
}

// New creates a new frame with the scope s and previous frame p.
func New(s scope.I, p *frame) *frame {
	f := &frame{scope: s}

	if p != nil {
		f.previous = p
		f.source = p.source
	}

	return f
}

// Loc returns the current location.
func (f *frame) Loc() *loc.T {
	return &f.source
}

// Previous returns the previous frame.
func (f *frame) Previous() *frame {
	return f.previous
}

// Resolve looks for a lexical and then dynamic resolution of k.
// The scope where the reference r was found is also returned.
func (f *frame) Resolve(k string) (s scope.I, r reference.I) {
	s = f.scope

	r = s.Lookup(k)
	if r != nil {
		return
	}

	for f = f.previous; f != nil; f = f.previous {
		for s = f.scope; s != nil; s = s.Enclosing() {
			r = s.Public().Get(k)
			if r != nil {
				return
			}
		}
	}

	return
}

// Scope returns the current frame's scope.
func (f *frame) Scope() scope.I {
	return f.scope
}

// Update sets the current lexical location.
func (f *frame) Update(source *loc.T) {
	f.source = *source
}
