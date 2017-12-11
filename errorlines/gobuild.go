package errorlines

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/leeola/gokakoune/util"
)

// GoBuild the given path and return any error lines
//
// Path will be resolved from the GOPATH, to use a package
// directory.
func GoBuild(path string) ([]string, error) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return nil, errors.New("GOPATH not set")
	}

	srcpath := filepath.Join(gopath, "src")

	packagePath, err := filepath.Rel(srcpath, filepath.Dir(path))
	if err != nil {
		return nil, err
	}

	tmpDir, err := ioutil.TempDir("", "gokakoune")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	_, stderr, exit, err := util.Exec(
		"go", "build", "-o", filepath.Join(tmpDir, "bin"), packagePath)
	if err != nil {
		return nil, err
	}

	// everything built fine.
	if exit == 0 {
		return nil, nil
	}

	lines := strings.Split(stderr, "\n")

	lenLines := len(lines)
	// collapse any lines that are just details to the previous line.
	// Eg, lines like:
	//
	// foo.go:75:4: cannot use x (type Foo) as Y in assignment:
	//      Foo does not implement Y (missing Baz method)
	//
	// NOTE(leeola): index starts as 1, because we can't collapse
	// the 0th line.
	for i := 1; i < lenLines; i++ {
		line := lines[i]
		if !strings.HasPrefix(line, "\t") {
			continue
		}

		prevI := i - 1
		lines[prevI] = lines[prevI] + " " + strings.TrimPrefix(line, "\t")

		// remove the collapsed element from the slice.
		lines = append(lines[:i], lines[i+1:]...)
		lenLines--

		// make sure to bump the index number down by 1, since we modified
		// the slice bounds.
		i--
	}

	// the first and last line are not error reporting lines,
	// so make sure we actually got some lines.
	if lenLines <= 2 {
		return nil, errors.New("unexpected go build response")
	}

	return lines[1 : lenLines-1], nil
}
