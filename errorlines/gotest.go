package errorlines

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/leeola/gokakoune/util"
)

// GoTest the given path and return any error lines
//
// Path will be resolved from the GOPATH, to use a package
// directory.
func GoTest(path string) ([]string, error) {
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
		"go", "test", "-c", "-o", filepath.Join(tmpDir, "bin"), packagePath)
	if err != nil {
		return nil, err
	}

	// everything built fine.
	if exit == 0 {
		return nil, nil
	}

	lines := strings.Split(stderr, "\n")
	return cleanErrLines(lines), nil
}

// cleanErrLines cleans an go test stderr response to output one string per
// error line.
//
// foo.go:75:4: cannot use x (type Foo) as Y in assignment:
//      Foo does not implement Y (missing Baz method)
//
// See func test for more expected input.
func cleanErrLines(src []string) []string {
	var (
		dst      []string
		strBuild string
	)

	for _, s := range src {
		if s == "" {
			continue
		}

		// ignore all lines that start with a #, they're just package grouping
		// descriptions.
		if strings.HasPrefix(s, "#") {
			continue
		}

		if strings.HasPrefix(s, "\t") {
			strBuild = fmt.Sprintf("%s %s", strBuild, strings.TrimPrefix(s, "\t"))
			continue
		}

		if strBuild != "" {
			dst = append(dst, strBuild)
		}

		strBuild = s
	}

	return append(dst, strBuild)
}
