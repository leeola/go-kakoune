package main

import (
	"github.com/leeola/go-kakoune/api"
	"github.com/leeola/go-kakoune/api/vars"
	"github.com/leeola/go-kakoune/commands/checksource"
)

func main() {
	kak := api.New()

	opts := api.DefineCommandOptions{
		Variables: []string{
			vars.BufFile,
		},
	}
	kak.DefineCommand("go-kakoune-kak-check-source", opts, func(p api.Params, v api.Vars) error {
		errLines, err := checksource.CheckSource(v[vars.BufFile])
		if err != nil {
			return err
		}

		for _, line := range errLines {
			kak.Command("set-code-err-line", line)
		}

		return nil
	})
}
