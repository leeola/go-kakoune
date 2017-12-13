package main

import (
	"github.com/leeola/gokakoune/api"
	"github.com/leeola/gokakoune/plugins/compilecheck"
	"github.com/leeola/gokakoune/plugins/jumpdef"
)

func main() {
	kak := api.New()

	opts := api.DefineCommandOptions{}
	kak.DefineCommand("gokakoune-compile-check", opts, compilecheck.CompileCheck...)

	opts = api.DefineCommandOptions{}
	kak.DefineCommand("gokakoune-jump-def", opts, jumpdef.JumpDef...)
}
