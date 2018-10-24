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

	// gokakouneInit
	gokakouneInit bool

	// gokakouneBin is the name of the binary using this API, being called by
	// kakoune itself.
	gokakouneBin string

	funcArgs []string
	funcVars map[string]string

	// expansionID is passed in after the gokakouneBin and indicates a
	// func block to execute.
	expansionID    int
	expansionCount int

	// funcCalled will be true if a func was already called for this process.
	// If true, every action on Kakoune becomes a noop. This is because if
	// the func was already called, there is no other action that this
	// instance of gokakoune needs to do.
	//
	// NOTE(leeola): there is no goroutine locking/protection on this var.
	// This may be added in the future, but currently it isn't likely to
	// matter.
	funcCalled bool
}

func New() *Kak {
	var (
		notGokakouneInit bool
		gokakouneBin     string
		funcID           int
		funcArgs         []string
		funcVars         = map[string]string{}
	)

	// TODO(leeola): move this entire block of logic to some type
	// of init func? Because currently there is no way to inform
	// the caller that an error occured.
	lenArgs := len(os.Args)
	if lenArgs >= 1 {
		gokakouneBin = os.Args[0]
	} else {
		panic("cannot get plugin executable")
	}

	if lenArgs >= 2 {
		id, err := strconv.Atoi(os.Args[1])
		if err != nil {
			panic("expansionID is not valid int")
		}
		funcID = id
		notGokakouneInit = true
	}

	if lenArgs >= 3 {
		funcArgs = make([]string, len(os.Args[2:]))
		copy(funcArgs, os.Args[2:])
	}

	for _, env := range os.Environ() {
		if strings.HasPrefix(env, var_prefix) {
			kwargs := strings.SplitN(env, "=", 2)

			// TODO(leeola): Possibly return an error here?
			// This weird key could be a sign of something broken.
			if len(kwargs) != 2 {
				continue
			}

			funcVars[kwargs[0]] = kwargs[1]
		}
	}

	return &Kak{
		writer:        os.Stdout,
		gokakouneBin:  gokakouneBin,
		gokakouneInit: !notGokakouneInit,
		expansionID:   funcID,
		funcArgs:      funcArgs,
		funcVars:      funcVars,
	}
}

func (k *Kak) Debug(v ...interface{}) {
	// TODO(leeola): figure out the fastest way to print the v...
	// as if Sprintln did it, but WITHOUT the newline at the end.
	//
	// Apparently fmt.Sprint() and fmt.Sprintln() have different
	// behavior, Sprint puts spaces between arguments only if they're
	// not strings.. so i can't use Sprint() ...sadface.
	s := Escape(v...)
	k.Println("echo", "-debug", s)
}

func (k *Kak) Debugf(f string, v ...interface{}) {
	s := Escape(fmt.Sprintf(f, v...))
	k.Println("echo", "-debug", s)
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
