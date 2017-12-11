package main

import (
	"github.com/leeola/go-kakoune/api"
	"github.com/leeola/go-kakoune/plugins/compilecheck"
)

func main() {
	kak := api.New()

	opts := api.DefineCommandOptions{}
	kak.DefineCommand("gokakoune-compile-check", opts, compilecheck.CompileCheck...)
}
