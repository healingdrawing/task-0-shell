// Released under an MIT license. See LICENSE.

//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris
// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package history

import (
	"os"
	"path"
)

func file(op func(string) (*os.File, error)) (*os.File, error) {
	s, ok := os.LookupEnv("PENISHELL_HISTORY")
	if !ok {
		s = path.Join(os.Getenv("HOME"), ".penishell-history")
	}

	return op(s)
}
