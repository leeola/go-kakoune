package _examples

import (
	"github.com/leeola/gokakoune/api"
	"github.com/leeola/gokakoune/api/vars"
)

var Prompt = api.Expansions{
	api.Prompt{
		Text: "foo",
		Expansions: api.Expansions{
			api.Func{
				ExportVars: []string{
					vars.Text,
				},
				Func: func(kak *api.Kak) error {
					promptText, err := kak.Var(vars.Text)
					if err != nil {
						return err
					}

					kak.Echof("hello %q, from prompt!", promptText)
					return nil
				},
			},
		},
	},
}
