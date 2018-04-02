package rename

import (
	"errors"
	"fmt"
	"strings"

	"github.com/leeola/gokakoune/api"
	"github.com/leeola/gokakoune/api/vars"
	"github.com/leeola/gokakoune/util"
)

const (
	gorenameBin = "gorename"
)

var RenameExpressions = api.Expansions{
	api.Prompt{
		Text: "rename: ",
		Expansions: api.Expansions{
			api.Func{
				ExportVars: []string{
					vars.Text,
					vars.BufFile,
					vars.CursorByteOffset,
				},
				Func: func(kak *api.Kak) error {
					text, err := kak.Var(vars.Text)
					if err != nil {
						return err
					}

					buffile, err := kak.Var(vars.BufFile)
					if err != nil {
						return err
					}

					cursorByteOffset, err := kak.Var(vars.CursorByteOffset)
					if err != nil {
						return err
					}

					stdout, stderr, exit, err := util.Exec(gorenameBin,
						"-offset", fmt.Sprintf("%s:#%s", buffile, cursorByteOffset),
						"-to", text)
					if err != nil {
						return err
					}

					// TODO(leeola): hook into a compile checker, to try and report
					// bad syntax, if possible.
					if exit != 0 {
						kak.Debugf("error: %q", stderr)
						return errors.New("bad exit, not compile checking yet...")
					}

					// TODO(leeola): unset any error code?

					kak.Command("edit!")

					// gorename reports what it changed, like how many files and how many
					// renames it did. So pass that report back to the user.
					stdout = strings.TrimSuffix(stdout, "\n")
					kak.Echo(stdout)

					return nil
				},
			},
		},
	},
}
