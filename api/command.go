package api

import (
	"fmt"
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

func (k *Kak) DefineCommand(name string, opts DefineCommandOptions, exps ...Expansion) error {
	cd := DefineCommand{
		Name:       name,
		Options:    opts,
		Expansions: exps,
	}

	return k.Expansion(cd)
}

// func (k *Kak) RegisterFunc(f *Func) error {
// 	// noop if func was already called
// 	if k.funcCalled {
// 		return nil
// 	}
//
// 	// automatically assign func id if one is not provided.
// 	if f.ID == "" {
// 		f.ID = strconv.Itoa(k.funcAutoID)
// 		k.funcAutoID++
// 	}
//
// 	// if gokakoune is initing, there's no need to do any registering.
// 	if k.gokakouneInit {
// 		return nil
// 	}
//
// 	// noop if the id doesn't match.
// 	if f.ID != k.funcID {
// 		return nil
// 	}
//
// 	return f.Func(k)
// }

// Expansion processes a Gokakoune expansion.
//
// Printing commands, registering functions, etc.
func (k *Kak) Expansion(exp Expansion) error {
	// noop if func was already called
	if k.funcCalled {
		return nil
	}

	if k.gokakouneInit {
		init, err := k.initExpansion(exp)
		if err != nil {
			return err
		}

		k.Println(init)

		// returning here ensures we don't run expansions when gokakoune
		// should be initializing.
		return nil
	}

	return k.runExpansion(exp)
}

func (k *Kak) runExpansion(exp Expansion) error {
	// noop if func was already called
	if k.funcCalled {
		return nil
	}

	expansionCount := k.expansionCount
	k.expansionCount++

	for _, cExp := range exp.Children() {
		if err := k.runExpansion(cExp); err != nil {
			return err
		}
	}

	if expansionCount != k.expansionID {
		return nil
	}

	k.funcCalled = true

	runnable, ok := exp.(Runnable)
	if !ok {
		return fmt.Errorf("not runnable expansion: %d", expansionCount)
	}

	// NOTE(leeola): this behavior of consuming the error and printing
	// it to the user is debatable. My thought process though is that
	// gokakoune isn't failing, so the main k.Expansion() or k.DefineCommand()
	// API shouldn't be failing. The user command is failing, but that's
	// on their side.
	//
	// This differs from the error above, where an expansion is not runnable,
	// that's clearly related to a gokakoune error.
	if err := runnable.Run(k); err != nil {
		// report the error to the user
		k.Fail(err)
	}

	return nil
}

func (k *Kak) initExpansion(exp Expansion) (string, error) {
	expansionCount := k.expansionCount
	k.expansionCount++

	var childInits []string
	for _, cExp := range exp.Children() {
		init, err := k.initExpansion(cExp)
		if err != nil {
			return "", err
		}

		childInits = append(childInits, init)
	}

	return exp.Init(Context{
		BinName:  k.gokakouneBin,
		ID:       expansionCount,
		Children: childInits,
	})
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
