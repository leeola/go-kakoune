package api

type Expansion interface {
	Init(*Kak) error
	Funcs(*Kak) ([]Func, error)
}

type DefineCommand struct {
	Name string

	Options DefineCommandOptions

	Expansions []Expansion
}

type Func struct {
	ID string

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
}

func (c *DefineCommand) Init(k *Kak) error {
	return nil
}

func (c *DefineCommand) Run(k *Kak) error {
	return nil
}
