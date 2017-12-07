package api

import (
	"fmt"
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
	Params    int
	Variables []string
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
		opts.Variables,
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

func (k *Kak) Echo(v ...interface{}) {
	// TODO(leeola): figure out the fastest way to print the v...
	// as if Sprintln did it, but WITHOUT the newline at the end.
	//
	// Apparently fmt.Sprint() and fmt.Sprintln() have different
	// behavior, Sprint puts spaces between arguments only if they're
	// not strings.. so i can't use Sprint() ...sadface.
	s := fmt.Sprintln(v...)
	l := len(s)
	s = s[:l-1]
	k.Printf("echo %q\n", s)
}

func (k *Kak) Echof(f string, v ...interface{}) {
	k.Printf("echo %q\n", fmt.Sprintf(f, v...))
}

func (k *Kak) Fail(v ...interface{}) {
	// TODO(leeola): figure out the fastest way to print the v...
	// as if Sprintln did it, but WITHOUT the newline at the end.
	//
	// Apparently fmt.Sprint() and fmt.Sprintln() have different
	// behavior, Sprint puts spaces between arguments only if they're
	// not strings.. so i can't use Sprint() ...sadface.
	s := fmt.Sprintln(v...)
	l := len(s)
	s = s[:l-1]
	k.Printf("fail %q\n", s)
}

func (k *Kak) Failf(f string, v ...interface{}) {
	k.Println("fail", fmt.Sprintf(f, v...))
}

// Print to the internal writer.
//
// This is a lower level interface, allowing you to send arbitrary
// commands to Kakoune. Use with caution.
func (k *Kak) Print(v ...interface{}) {
	fmt.Fprint(k.writer, v...)
}

// Println to the internal writer.
//
// This is a lower level interface, allowing you to send arbitrary
// commands to Kakoune. Use with caution.
func (k *Kak) Println(v ...interface{}) {
	fmt.Fprintln(k.writer, v...)
}

// Printf to the internal writer.
//
// This is a lower level interface, allowing you to send arbitrary
// commands to Kakoune. Use with caution.
func (k *Kak) Printf(f string, v ...interface{}) {
	fmt.Fprintf(k.writer, f, v...)
}
