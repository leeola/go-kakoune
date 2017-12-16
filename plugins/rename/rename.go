package rename

import (
	"errors"
	"fmt"

	"github.com/leeola/gokakoune/api"
	"github.com/leeola/gokakoune/api/vars"
	"github.com/leeola/gokakoune/util"
)

const (
	gorenameBin = "gorename"
)

// Unfinished implementation. We need to first prompt the user
// for input, which can't be done with gokakoune currently.
var RenameSubprocs = []api.Subproc{{
	ExportVars: []string{
		vars.BufFile,
		vars.CursorByteOffset,
	},
	Func: func(kak *api.Kak) error {
		buffile, err := kak.Var(vars.BufFile)
		if err != nil {
			return err
		}

		cursorByteOffset, err := kak.Var(vars.CursorByteOffset)
		if err != nil {
			return err
		}

		stdout, _, exit, err := util.Exec(
			gorenameBin, "-pos", fmt.Sprintf("%s:#%s", buffile, cursorByteOffset))
		if err != nil {
			return err
		}

		if exit != 0 {
			return errors.New("bad exit, not compile checking yet...")
		}

		// TODO(leeola): unset any error code?

		kak.Echo(stdout)
		return nil
	},
}}
