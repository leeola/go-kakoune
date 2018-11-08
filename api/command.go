package api

import (
	"fmt"
	"strings"
)

// Subproc executes Go code in a subproc of Kakoune.
//
// Each Subproc is effectively the same as the %sh{ .. } block found within
// a define-command command. Example:
//
//    define-command cmdName %{
//      %sh{
//        # do stuff in shell scope.
//      }
//    }
//
// The Subproc.Func is called from the shell expansion in the above example.
type Subproc struct {
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

	// Func is called within each subprocess specified in Kak.DefineCommand.
	//
	// It's important to understand that the function execution defines the
	// lifetime of the Kakoune command. Memory cannot be shared between
	// Subproc executions.
	//
	// To share memory/state between Func calls, set options within Kakoune
	// and retrieve them on future subprocs.
	Func func(*Kak) error
}

type DefineCommandOptions struct {
	Params int
}

// func (k *Kak) initCommand(name string, opts DefineCommandOptions, cs []Subproc) error {
// 	var blockStrs []string
// 	for i, c := range cs {
// 		var argStr string
// 		for i := 0; i < opts.Params; i++ {
// 			// prefix each item with a space!
// 			argStr += fmt.Sprintf(` "${%d}"`, i+1)
// 		}
//
// 		vars := make([]string, len(c.ExportVars))
// 		for i, v := range c.ExportVars {
// 			vars[i] = "$kak_" + v
// 		}
//
// 		blockStrs = append(blockStrs, fmt.Sprintf(`
//   %%sh{
//     # the following variables are being written in the def source
//     # code to make Kakoune export them to this shell scope. By doing
//     # so, they become available to the Go source code.
//     #
//     # Note that it appears Kakoune just uses regex on the codeblock,
//     # so the fact that the variables are commented out does not matter.
//     # It loads any kak variables specified in the code.
//     #
//     # %s
//
//     %s %q %d%s
//   }`,
// 			vars,
// 			k.bin, name, i, argStr))
// 	}
//
// 	// space omitted between %q%s on purpose,
// 	// see above loop code format.
// 	k.Printf(`
// define-command -params %d %s %%{
//   %s
// }
// `, opts.Params, name, strings.Join(blockStrs, "\n"),
// 	)
//
// 	return nil
// }

// func (k *Kak) runCommand(name string, opts DefineCommandOptions, cs []Subproc) error {
// 	if k.cmdBlockIndex > len(cs) {
// 		return fmt.Errorf("%s block unavailable: %d", name, k.cmdBlockIndex)
// 	}
//
// 	c := cs[k.cmdBlockIndex]
//
// 	// TODO(leeola): set the active command(s) so that we know what Vars[] should
// 	// be available.
// 	// k.activeCommands = c
//
// 	// NOTE(leeola): passing shared mutable references of the
// 	// params and vars to the user should be acceptable here.
// 	//
// 	// This is because no two commands will ever be called from
// 	// Kakoune within the same process, so technically all of
// 	// the memory of a single process should be owned by a single
// 	// kak-command regardless.
// 	if err := c.Func(k); err != nil {
// 		k.Failf("gokakoune: %s: %s", name, err.Error())
// 	}
//
// 	return nil
//
// }

func (k Kak) DefineCommand(name string, opts DefineCommandOptions, exp Expander) error {
	if k.isNop {
		return nil
	}

	_, err := fmt.Fprintf(k.writer, "define-command -params %d %s ", opts.Params, name)
	if err != nil {
		return fmt.Errorf("fprintf cmd start: %v", err)
	}

	// temp using empty context
	if err := wrapSh(exp).Expand(k); err != nil {
		return fmt.Errorf("expansion expand: %v", err)
	}

	return nil
}

// Command calls a kakoune command directly, escaping arguments
// automatically.
func (k *Kak) Command(name string, args ...interface{}) {
	v := make([]interface{}, len(args)+1)
	v[0] = name
	for i, a := range args {
		// // EscapeRune ensures that the double quote is escaped, but nothing
		// // else.
		// //
		// // This is because kakoune seems to have non-intuitive behavior with
		// // escaping. If we use something like `Sprintf("%q", a)`, newlines
		// // will be escaped in kakoune as well. We have to not escape newlines,
		// // but do escape the surrounding quotes to ensure it is read as a
		// // single argument.
		// //
		// // This feels a bit hacky, but i've not found a better way yet.
		// v[i+1] = fmt.Sprintf("\"%s\"", util.EscapeRune(a, '"'))

		s, ok := a.(string)
		if ok {
			s = escapeString(s)
		} else {
			s = fmt.Sprint(a)
		}

		v[i+1] = s
	}

	k.Println(v...)
}

func (k Kak) Prompt(promptMsg string, exp Expander) error {
	_, err := fmt.Fprintf(k.writer, "prompt %q ", promptMsg)
	if err != nil {
		return fmt.Errorf("fprintf cmd start: %v", err)
	}

	return wrapSh(exp).Expand(k)
}

func (k Kak) EvaluateCommands(exp Expander) error {
	_, err := fmt.Fprint(k.writer, "evaluate-commands ")
	if err != nil {
		return fmt.Errorf("fprintf eval: %v", err)
	}

	// do not wrap Sh. It doesn't need it, and will endlessly loop if you do ;)
	return exp.Expand(k)
}

func (k Kak) Hook(scope, hookName, filteringRegex string, exp Expander) error {
	_, err := fmt.Fprintf(k.writer, "hook %s %s %s ", scope, hookName, filteringRegex)
	if err != nil {
		return fmt.Errorf("print cmd: %v", err)
	}

	return wrapSh(exp).Expand(k)
}

func (k Kak) DeclareOption(flags []string, kakType, optName string) error {
	flagStr := strings.Join(flags, " ")
	if len(flags) != 0 {
		flagStr += " "
	}

	_, err := fmt.Fprintf(k.writer, "declare-option %s%s %s\n", flagStr, kakType, optName)
	if err != nil {
		return fmt.Errorf("print cmd: %v", err)
	}

	return nil
}

// wrapSh ensures that a shell expansion will be wrapped
// by a plain %{} expansion.
//
// Eg, %sh{} would become %{ eval %sh{} }. The function will not needlessly
// wrap a plain Expansion. Ie %{} does not become %{ %{} }.
//
// This was done because it was seen that %sh{} blocks used without being
// inside of a %{} expansion would not get access to $kak_bufname style
// vars, which Gokakoune depends upon for callback usage.
func wrapSh(exp Expander) Expander {
	switch v := exp.(type) {
	case Callback, Sh:
		exp = Expansion{
			Body: func(k Kak) error {
				return k.EvaluateCommands(v)
			},
		}
	default:
		// do nothing if it's an Expansion, or other type.
		// Other is allowed, as we assume it doesn't need
		// special treatment like callback or sh.
	}

	return exp
}
