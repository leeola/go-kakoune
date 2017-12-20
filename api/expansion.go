package api

import (
	"fmt"
	"strings"
)

type Expansion interface {
	Init(Context) (string, error)
	Children() []Expansion
}

type Runnable interface {
	Run(*Kak) error
}

type Context struct {
	BinName  string
	ID       int
	Children []string
}

type Expansions []Expansion

type DefineCommand struct {
	Name string

	Options DefineCommandOptions

	Expansions []Expansion
}

type Func struct {
	// ExportVars specifies the variables that Kakoune should export to the Subproc.
	//
	// Eg, if `[]string{"buffile"}` is the value of ExportVars, then the
	// environment variable `kak_buffile` will be exported to your subproc.
	// Retrieval of this variable can be done with `kak.Var("buffile")`,
	// also without the kak_ prefix.  All gokakoune functions will properly
	// prefix kak_ and kak_opt_ as needed.
	//
	// NOTE: these are not prefixed with `kak_`. Eg, to export `bufname` to a
	// subproc just specify the following, *without* the `kak_` prefix:
	//
	//    ExportVars: []string{"bufname"}
	//
	// Constants in the api/vars package are also available.
	ExportVars []string

	Func func(*Kak) error
}

type Sh struct {
}

type Prompt struct {
	Text       string
	Expansions []Expansion
}

func (e DefineCommand) Init(ctx Context) (string, error) {
	return fmt.Sprintf(`
define-command -params %d %s %%{
  %s
}`,
		e.Options.Params, e.Name,
		strings.Join(ctx.Children, "\n")), nil
}

func (e DefineCommand) Children() []Expansion {
	return e.Expansions
}

func (e Func) Init(ctx Context) (string, error) {
	// var argStr string
	// for i := 0; i < e.Params; i++ {
	// 	// TODO(leeola): use the new Go 1.10 string builder, i've not switched
	// 	// yet.
	// 	//
	// 	// prefix each item with a space!
	// 	argStr += fmt.Sprintf(` "${%d}"`, i+1)
	// }

	vars := make([]string, len(e.ExportVars))
	for i, v := range e.ExportVars {
		vars[i] = "$kak_" + v
	}

	return fmt.Sprintf(`
  %%sh{
    # the following variables are being written in the def source
    # code to make Kakoune export them to this shell scope. By doing
    # so, they become available to the Go source code.
    #
    # Note that it appears Kakoune just uses regex on the codeblock,
    # so the fact that the variables are commented out does not matter.
    # It loads any kak variables specified in the code.
    #
    # %s

    %s %d "$@"
  }
`,
		vars,
		ctx.BinName, ctx.ID), nil
}

func (e Func) Children() []Expansion {
	return nil
}

func (e Func) Run(k *Kak) error {
	return e.Func(k)
}

func (e Prompt) Init(ctx Context) (string, error) {
	return fmt.Sprintf(`
  prompt %q %%{
    %s
  }
`,
		e.Text,
		strings.Join(ctx.Children, "\n")), nil
}

func (e Prompt) Children() []Expansion {
	return e.Expansions
}
