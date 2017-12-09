package main

import (
	"github.com/leeola/go-kakoune/api"
	"github.com/leeola/go-kakoune/api/vars"
	"github.com/leeola/go-kakoune/plugins/compilecheck"
)

const (
	code_err = "code_err"
)

func main() {
	kak := api.New()

	opts := api.DefineCommandOptions{
		Vars: []string{
			vars.BufFile,
		},
	}

	kak.DefineCommand("go-kakoune-compile-check", opts, compilecheck.CompileCheck)
}
