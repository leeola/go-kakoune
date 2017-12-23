package main

import (
	"github.com/leeola/gokakoune/api"
	"github.com/leeola/gokakoune/plugins/compilecheck"
	"github.com/leeola/gokakoune/plugins/complete"
	"github.com/leeola/gokakoune/plugins/jumpdef"
	"github.com/leeola/gokakoune/plugins/rename"
	"github.com/leeola/gokakoune/plugins/showdoc"
)

func main() {
	kak := api.New()

	opts := api.DefineCommandOptions{}
	kak.DefineCommand("gokakoune-compile-check", opts, compilecheck.CompileCheckExpressions...)

	opts = api.DefineCommandOptions{}
	kak.DefineCommand("gokakoune-jump-def", opts, jumpdef.JumpDefExpressions...)

	opts = api.DefineCommandOptions{}
	kak.DefineCommand("gokakoune-show-doc", opts, showdoc.ShowDocExpressions...)

	opts = api.DefineCommandOptions{}
	kak.DefineCommand("gokakoune-rename", opts, rename.RenameExpressions...)

	opts = api.DefineCommandOptions{}
	kak.DefineCommand("gokakoune-complete", opts, complete.CompleteExpressions...)

	kak.Command(`
    hook global WinCreate .* %{
      set-option window completers "option=gokakoune_completions:%opt{completers}"
    }

    #set-option window completers "option=gokakoune_completions:%opt{completers}"
  `)
}
