package tabnine

import (
	"fmt"
	"io"
)

// IO abstracts the raw TabNine communication protocol away from the API
// implementation.
//
// This is useful because at the time of writing, the named pipe
// implementation in this package did not work. So for now, a simpler
// and likely less efficient HTTP implementation is being used. This
// interface ensures that in the future if named pipes can be made to
// work, the API functionality just swaps IO implementations.
type IO interface {
	SendRecv(req io.Reader) (res io.ReadCloser, err error)

	// A named pipe solution may necessitate a close method on the IO,
	// but for now it's not needed.
	// Close() error
}

type Config struct {
	// TabnineBin is the path of the tabnine binary.
	TabnineBin string

	ConfigDir string

	// IO if supplied will be used for all IO to Tabnine.
	IO IO
}

type Tabnine struct {
	tabnineBin string
	configDir  string
	io         IO
}

func New(c Config) (*Tabnine, error) {
	// if empty, use a default
	if c.IO == nil {
		io, err := NewHTTPClient(c.ConfigDir)
		if err != nil {
			return nil, fmt.Errorf("newhttpclient: %v", err)
		}
		c.IO = io
	}

	return &Tabnine{
		tabnineBin: c.TabnineBin,
		configDir:  c.ConfigDir,
		io:         c.IO,
	}, nil
}
