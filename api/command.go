package api

import (
	"fmt"
	"strings"
)

type Subproc struct {
	Vars []string
	Func func(*Kak) error
}

type DefineCommandOptions struct {
	Params int
}

func (k *Kak) initCommand(name string, opts DefineCommandOptions, cs []Subproc) error {
	var blockStrs []string
	for i, c := range cs {
		var argStr string
		for i := 0; i < opts.Params; i++ {
			// prefix each item with a space!
			argStr += fmt.Sprintf(` "${%d}"`, i+1)
		}

		blockStrs = append(blockStrs, fmt.Sprintf(`
  %%sh{
    # the following variables are being written in the def source
    # code to make Kakoune export them to this shell scope. By doing
    # so, they become available to the Go source code.
    #
    # %s

    %s %q %d%s
  }`,
			c.Vars,
			k.bin, name, i, argStr))
	}

	// space omitted between %q%s on purpose,
	// see above loop code format.
	k.Printf(`
define-command -params %d %s %%{
  %s
}
`, opts.Params, name, strings.Join(blockStrs, "\n"),
	)

	return nil
}

func (k *Kak) runCommand(name string, opts DefineCommandOptions, cs []Subproc) error {
	if k.cmdBlockIndex > len(cs) {
		return fmt.Errorf("%s block unavailable: %d", name, k.cmdBlockIndex)
	}

	c := cs[k.cmdBlockIndex]

	// TODO(leeola): set the active command(s) so that we know what Vars[] should
	// be available.
	// k.activeCommands = c

	// NOTE(leeola): passing shared mutable references of the
	// params and vars to the user should be acceptable here.
	//
	// This is because no two commands will ever be called from
	// Kakoune within the same process, so technically all of
	// the memory of a single process should be owned by a single
	// kak-command regardless.
	if err := c.Func(k); err != nil {
		k.Failf("gokakoune: %s: %s", name, err.Error())
	}

	return nil

}

func (k *Kak) DefineCommand(name string, opts DefineCommandOptions, cs ...Subproc) error {
	if k.cmd == "" {
		return k.initCommand(name, opts, cs)
	}

	if k.cmd != name {
		return nil
	}

	return k.runCommand(name, opts, cs)
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
