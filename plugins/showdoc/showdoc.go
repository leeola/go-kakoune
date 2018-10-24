package showdoc

import (
	"fmt"
	"strings"

	"github.com/leeola/gokakoune/api"
	"github.com/leeola/gokakoune/api/vars"
	"github.com/leeola/gokakoune/util"
)

const (
	gogetdocBin = "gogetdoc"
)

var ShowDocExpressions = api.Expansions{
	api.Func{
		ExportVars: []string{
			vars.BufName,
			vars.CursorByteOffset,
			vars.WindowHeight,
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

			windowHeight, err := kak.VarInt(vars.WindowHeight)
			if err != nil {
				return err
			}

			stdout, _, exit, err := util.Exec(
				gogetdocBin, "-pos", bufname+":#"+cursorByteOffset)
			if err != nil {
				return err
			}

			if exit != 0 {
				kak.Debugf(gogetdocBin, "output:", stdout)
				return fmt.Errorf("unexpected %s exit code: %d", gogetdocBin, exit)
			}

			// info seems to have a bug in which it causes a silent failure
			// if the input is too large compared to the window height. So,
			// we need to trim the stdout by lines.
			//
			// TODO(leeola): i'm sure this is a slow (perf) way to implement
			// this, but i'm busy atm. This should be benched.
			split := strings.SplitN(stdout, "\n", windowHeight)
			if len(split) == windowHeight {
				lineCap := int(float32(windowHeight) * 0.75)
				split = split[:lineCap]
			}
			stdout = strings.Join(split, "\n")

			kak.Command("info", stdout)

			return nil
		},
	},
}
