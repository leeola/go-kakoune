package showdoc

import (
	"fmt"

	"github.com/leeola/gokakoune/api"
	"github.com/leeola/gokakoune/api/vars"
	"github.com/leeola/gokakoune/util"
)

const (
	gogetdocBin = "gogetdoc"
)

// BUG(leeola): multi-line escaping is not working yet.
var ShowDocSubprocs = []api.Subproc{
	{
		ExportVars: []string{
			vars.BufName,
			vars.CursorByteOffset,
		},
		Func: func(kak *api.Kak) error {
			bufname, err := kak.Var(vars.BufName)
			if err != nil {
				return err
			}

			cursorByteOffset, err := kak.Var(vars.CursorByteOffset)
			if err != nil {
				return err
			}

			stdout, _, exit, err := util.Exec(
				gogetdocBin, "-pos", bufname+":#"+cursorByteOffset)
			if err != nil {
				return err
			}

			if exit != 0 {
				return fmt.Errorf("unexpected %s exit code: %d", gogetdocBin, exit)
			}

			kak.Command("info", stdout)

			return nil
		},
	},
}
