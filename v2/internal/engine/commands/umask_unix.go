// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package commands

import (
	"fmt"

	"penishell/internal/common/interface/cell"
	"penishell/internal/common/interface/integer"
	"penishell/internal/common/type/sym"
	"penishell/internal/common/validate"
	"penishell/internal/system/process"
)

func umask(args cell.I) cell.I {
	v := validate.Fixed(args, 0, 1)

	nmask := int64(0)
	if len(v) == 1 {
		nmask = integer.Value(v[0])
	}

	omask := process.Umask(int(nmask))

	if nmask == 0 {
		process.Umask(omask)
	}

	return sym.New(fmt.Sprintf("0o%o", omask))
}
