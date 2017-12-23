package complete

import (
	"github.com/leeola/gokakoune/api"
	"github.com/leeola/gokakoune/api/vars"
)

var CompleteExpressions = api.Expansions{
	api.Func{
		ExportVars: []string{
			vars.Text,
			vars.BufFile,
			vars.CursorByteOffset,
		},
		Func: func(kak *api.Kak) error {
			kak.Echo("wee")

			return nil
		},
	},
}
