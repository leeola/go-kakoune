package main

import (
	"github.com/leeola/gokakoune/api"
	"github.com/leeola/gokakoune/plugins/compilecheck"
)

func main() {
	kak := api.New()

	opts := api.DefineCommandOptions{}
	kak.DefineCommand("gokakoune-compile-check", opts, compilecheck.CompileCheck...)
}
