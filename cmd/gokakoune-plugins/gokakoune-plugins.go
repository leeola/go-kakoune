package main

import (
	"github.com/leeola/gokakoune/api"
	"github.com/leeola/gokakoune/plugins/compilecheck"
	"github.com/leeola/gokakoune/plugins/jumpdef"
	"github.com/leeola/gokakoune/plugins/showdoc"
)

func main() {
	kak := api.New()

	opts := api.DefineCommandOptions{}
	kak.DefineCommand("gokakoune-compile-check", opts, compilecheck.CompileCheckSubprocs...)

	opts = api.DefineCommandOptions{}
	kak.DefineCommand("gokakoune-jump-def", opts, jumpdef.JumpDefSubprocs...)

	opts = api.DefineCommandOptions{}
	kak.DefineCommand("gokakoune-show-doc", opts, showdoc.ShowDocSubprocs...)

	//opts = api.DefineCommandOptions{}
	//kak.DefineCommand("gokakoune-rename", opts, rename.RenameSubprocs...)
}
