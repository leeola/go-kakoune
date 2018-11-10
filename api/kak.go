package api

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Kak struct {
	writer io.Writer

	// gokakouneBin is the name of the binary using this API, being called by
	// kakoune itself.
	gokakouneBin string

	isCallback   bool
	callbackName string
	callbackArgs []string
	callbackVars map[string]string

	isNop bool
}

type KakInit struct {
	Kak

	// callbackCount is used to create a name for callbacks without a name.
	callbackCount int
}

// New returns a new instance of Kakoune, using os.Args and os.Environ,
// and panics of any error is encountered.
//
// For a safe alternative, use NewSafe()
func New() *KakInit {
	k, err := NewSafe(os.Stdout, os.Args, os.Environ())
	if err != nil {
		panic(fmt.Sprintf("gokakoune new: %v", err))
	}
	return k
}

func NewSafe(w io.Writer, args, envs []string) (*KakInit, error) {
	var (
		isCallback   bool
		gokakouneBin string
		callbackName string
		callbackArgs []string
		callbackVars = map[string]string{}
	)

	// TODO(leeola): move this entire block of logic to some type
	// of init func? Because currently there is no way to inform
	// the caller that an error occured.
	lenArgs := len(os.Args)
	if lenArgs >= 1 {
		gokakouneBin = os.Args[0]
	} else {
		return nil, fmt.Errorf("cannot get gokakoune executable")
	}

	if lenArgs >= 2 {
		callbackName = os.Args[1]
		isCallback = true
	}

	if lenArgs >= 3 {
		callbackArgs = make([]string, len(os.Args[2:]))
		copy(callbackArgs, os.Args[2:])
	}

	for _, env := range envs {
		if strings.HasPrefix(env, var_prefix) {
			kwargs := strings.SplitN(env, "=", 2)

			// TODO(leeola): Possibly return an error here?
			// This weird key could be a sign of something broken.
			if len(kwargs) != 2 {
				continue
			}

			callbackVars[kwargs[0]] = kwargs[1]
		}
	}

	return &KakInit{
		Kak: Kak{
			writer:       w,
			gokakouneBin: gokakouneBin,
			isCallback:   isCallback,
			isNop:        isCallback,
			callbackName: callbackName,
			callbackArgs: callbackArgs,
			callbackVars: callbackVars,
		},
	}, nil
}

func (k *KakInit) Callback(exportVars []string, f func(Kak) error) Callback {
	name := strconv.Itoa(k.callbackCount)
	return k.CallbackWithName(name, exportVars, f)
}

func (k *KakInit) CallbackWithName(name string, exportVars []string, f func(Kak) error) Callback {
	k.callbackCount++

	// if the callback matches, call it.
	if k.isCallback && name == k.callbackName {
		cbKak := k.Kak
		// allow the cb to print commands and etc
		cbKak.isNop = false
		if err := f(cbKak); err != nil {
			cbKak.Failf("gokakoune error: callback %s: %v", name, err)
		}

		// TODO(leeola): we could potentially reduce work by exiting the process
		// immediately after calling the callback.
		//
		// Ie, when a callback is being requested, there is no other work that
		// can be done with gokakoune. Yet, a possibly large amount of code
		// may still be run. We may want to configurably os.Exit(0) after
		// the callback is called (and maybe even os.Exit(1) if error).
		//
		// This is a bit scary to run from inside some random lib, so i may just
		// put some type of hook instead, and let the API user define
		// kak.OnDone(func()) themselves, which seems like a better solution.
	}

	return Callback{
		Name:       name,
		ExportVars: exportVars,
	}
}

func (k Kak) Debug(v ...interface{}) {
	if k.isNop {
		return
	}

	// TODO(leeola): figure out the fastest way to print the v...
	// as if Sprintln did it, but WITHOUT the newline at the end.
	//
	// Apparently fmt.Sprint() and fmt.Sprintln() have different
	// behavior, Sprint puts spaces between arguments only if they're
	// not strings.. so i can't use Sprint() ...sadface.
	s := Escape(v...)
	k.Println("echo", "-debug", s)
}

func (k Kak) Debugf(f string, v ...interface{}) {
	if k.isNop {
		return
	}

	s := Escape(fmt.Sprintf(f, v...))
	k.Println("echo", "-debug", s)
}

func (k Kak) Echo(v ...interface{}) {
	if k.isNop {
		return
	}

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

func (k Kak) Echof(f string, v ...interface{}) {
	if k.isNop {
		return
	}

	k.Printf("echo %q\n", fmt.Sprintf(f, v...))
}

func (k Kak) Fail(v ...interface{}) {
	if k.isNop {
		return
	}

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

func (k Kak) Failf(f string, v ...interface{}) {
	if k.isNop {
		return
	}

	k.Println("fail", fmt.Sprintf(f, v...))
}

// Print to the internal writer.
//
// This is a lower level interface, allowing you to send arbitrary
// commands to Kakoune. Use with caution.
func (k Kak) Print(v ...interface{}) error {
	if k.isNop {
		return nil
	}

	_, err := fmt.Fprint(k.writer, v...)
	return err
}

// Println to the internal writer.
//
// This is a lower level interface, allowing you to send arbitrary
// commands to Kakoune. Use with caution.
func (k Kak) Println(v ...interface{}) error {
	if k.isNop {
		return nil
	}

	_, err := fmt.Fprintln(k.writer, v...)
	return err
}

// Printf to the internal writer.
//
// This is a lower level interface, allowing you to send arbitrary
// commands to Kakoune. Use with caution.
func (k Kak) Printf(f string, v ...interface{}) error {
	if k.isNop {
		return nil
	}

	_, err := fmt.Fprintf(k.writer, f, v...)
	return err
}
