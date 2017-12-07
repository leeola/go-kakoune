package checksource

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/leeola/go-kakoune/util"
)

func CheckSource(path string) ([]string, error) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return nil, errors.New("GOPATH not set")
	}

	srcpath := filepath.Join(gopath, "src")

	packagePath, err := filepath.Rel(srcpath, filepath.Dir(path))
	if err != nil {
		return nil, err
	}

	_, stderr, exit, err := util.Exec("go", "build", packagePath)
	if err != nil {
		return nil, err
	}

	// everything built fine.
	if exit == 0 {
		return nil, nil
	}

	lines := strings.Split(stderr, "\n")

	lenLines := len(lines)

	// the first and last line are not error reporting lines,
	// so make sure we actually got some lines.
	if lenLines <= 2 {
		return nil, errors.New("unexpected go build response")
	}

	return lines[1 : lenLines-1], nil
}
