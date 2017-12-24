package errsig

import (
	"github.com/leeola/gokakoune/api"
	"github.com/leeola/gokakoune/api/vars"
)

var ErrSigExpressions = api.Expansions{
	api.Func{
		ExportVars: []string{
			vars.Text,
			vars.BufFile,
			vars.CursorByteOffset,
		},
		Func: func(kak *api.Kak) error {
			kak.Echo("wee")

			dostuff()

			return nil
		},
	},
}
