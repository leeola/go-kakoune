package api

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Kak struct {
	writer        io.Writer
	cmd           string
	cmdBlockIndex int
	args          []string
	vars          map[string]string
}

func New() *Kak {
	var (
		cmd           string
		cmdBlockIndex int
		args          []string
	)

	// TODO(leeola): move this entire block of logic to some type
	// of init func? Because currently there is no way to inform
	// the caller that an error occured.
	lenArgs := len(os.Args)
	if lenArgs >= 3 {
		cmd = os.Args[1]

		cbi, err := strconv.Atoi(os.Args[2])
		if err != nil {
			panic(err)
		}
		cmdBlockIndex = cbi
	}

	if lenArgs >= 4 {
		args = make([]string, len(os.Args[2:]))
		copy(args, os.Args[2:])
	}

	vars := map[string]string{}
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, var_prefix) {
			kwargs := strings.SplitN(env, "=", 2)

			// TODO(leeola): Possibly return an error here?
			// This weird key could be a sign of something broken.
			if len(kwargs) != 2 {
				continue
			}

			vars[kwargs[0]] = kwargs[1]
		}
	}

	return &Kak{
		writer:        os.Stdout,
		cmd:           cmd,
		cmdBlockIndex: cmdBlockIndex,
		args:          args,
		vars:          vars,
	}
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
