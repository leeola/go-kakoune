package api

import (
	"io"
	"os"
	"strings"
)

type Kak struct {
	writer io.Writer
	cmd    string
	args   []string
	vars   map[string]string
}

func New() *Kak {
	var (
		cmd  string
		args []string
	)

	switch l := len(os.Args); {
	case l == 2:
		cmd = os.Args[1]
	case l >= 3:
		cmd = os.Args[1]
		args = make([]string, len(os.Args[2:]))
		copy(args, os.Args[2:])
	}

	vars := map[string]string{}
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "kak_") {
			kwargs := strings.SplitN(env, "=", 2)

			// TODO(leeola): Possibly return an error here?
			// This weird key could be a sign of something broken.
			if len(kwargs) != 2 {
				continue
			}

			vars[kwargs[0]] = kwargs[1]
			// vars[kwargs[0]] = kwargs[1]
		}
	}

	return &Kak{
		writer: os.Stdout,
		cmd:    cmd,
		args:   args,
		vars:   vars,
	}
}
