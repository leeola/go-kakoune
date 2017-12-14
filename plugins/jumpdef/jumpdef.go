package jumpdef

import (
	"errors"
	"fmt"
	"strings"

	"github.com/leeola/gokakoune/api"
	"github.com/leeola/gokakoune/api/vars"
	"github.com/leeola/gokakoune/util"
)

// NOTE(leeola): this has plenty of limitations in implementation,
// improvements coming in the near future. this is just a hack.
// Eg, proper error handling, buffer saving, etc.
var JumpDefSubprocs = []api.Subproc{
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
				"guru", "definition", bufname+":#"+cursorByteOffset)
			if err != nil {
				return err
			}

			if exit != 0 {
				return fmt.Errorf("unexpected guru exit code: %d", exit)
			}

			split := strings.SplitN(stdout, ":", 4)
			if len(split) < 4 {
				return errors.New("unexpected guru stdout")
			}

			file, line, col := split[0], split[1], split[2]
			// desc comes with a newline, trim it.
			desc := strings.TrimSpace(split[3])

			// TODO(leeola): make this a native Go command.
			kak.Command("edit", file, line, col)

			kak.Echo(desc)

			return nil
		},
	},
}
