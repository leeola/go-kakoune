package api

import (
	"fmt"

	"github.com/leeola/go-kakoune/api/vars"
)

type Vars map[string]string

type Params []string

// TODO(leeola): a concurrent Kak param might be needed to safely
// define logic blocks with concurrent code. Hard to say at this time
// though. Concurrency doesn't *seem* likely to be used much, since
// and Kakoune interaction will still be non-concurrent as it is entirely
// over stdout.
type Command func(Params, Vars) error

type DefineCommandOptions struct {
	Params int
	Vars   []string
}

func (k *Kak) DefineCommand(name string, opts DefineCommandOptions, f Command) error {
	if k.cmd == name {
		// NOTE(leeola): passing shared mutable references of the
		// params and vars to the user should be acceptable here.
		//
		// This is because no two commands will ever be called from
		// Kakoune within the same process, so technically all of
		// the memory of a single process should be owned by a single
		// kak-command regardless.
		if err := f(k.args, k.vars); err != nil {
			k.Failf("go-kakoune: %s:", name, err.Error())
		}

		return nil
	}

	var argStr string
	for i := 0; i < opts.Params; i++ {
		// prefix each item with a space!
		argStr += fmt.Sprintf(` "${%d}"`, i+1)
	}

	// space omitted between %q%s on purpose,
	// see above loop code format.
	k.Printf(`
define-command -params %d %s %%{
  %%sh{

    # included variables
    # %s

    kak-go-kakoune-plugins %q%s
  }
}
`, opts.Params, name,
		opts.Vars,
		name, argStr,
	)

	return nil
}

// Command calls a kakoune command directly.
func (k *Kak) Command(name string, args ...string) {
	v := make([]interface{}, len(args)+1)
	v[0] = name
	for i, a := range args {
		// using %q to escape and quote the variable, to ensure that each
		// argument given to Command is an argument in the chosen command.
		v[i+1] = fmt.Sprintf("%q", a)
	}
	k.Println(v...)
}

func (v Vars) Option(o string) string {
	return v[vars.Option(o)]
}