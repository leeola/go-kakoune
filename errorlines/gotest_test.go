package errorlines

import (
	"strings"
	"testing"
)

func TestCleanErrLines(t *testing.T) {
	src := strings.Split(`# github.com/foo/bar/baz
foo.go:75:4: cannot use x (type Foo) as Y in assignment:
	Foo does not implement Y (missing Baz method)
bar.go:75:4: cannot use x (type Foo) as Y in assignment:
	Foo does not implement Y (missing Baz method)
	Foo does not implement Z (missing Baz method)
# github.com/foo/bar/baz
bar.go:75:4: cannot use x (type Foo) as Y in assignment
baz.go:75:4: unknown field 'Wee' in struct literal of type RollerCoaster`, "\n")

	dst := cleanErrLines(src)

	if got, want := len(dst), 4; got != want {
		t.Fatalf("unexpected result length. got:%d, want:%d", got, want)
	}

	gotLine := dst[0]
	wantLine := "foo.go:75:4: cannot use x (type Foo) as Y in assignment: Foo does not implement Y (missing Baz method)"
	if gotLine != wantLine {
		t.Errorf("unexpected line 0.\n  got:%q\n want:%q", gotLine, wantLine)
	}

	gotLine = dst[1]
	wantLine = "bar.go:75:4: cannot use x (type Foo) as Y in assignment: Foo does not implement Y (missing Baz method) Foo does not implement Z (missing Baz method)"
	if gotLine != wantLine {
		t.Errorf("unexpected line 1.\n  got:%q\n want:%q", gotLine, wantLine)
	}

	gotLine = dst[2]
	wantLine = "bar.go:75:4: cannot use x (type Foo) as Y in assignment"
	if gotLine != wantLine {
		t.Errorf("unexpected line 2.\n  got:%q\n want:%q", gotLine, wantLine)
	}

	gotLine = dst[3]
	wantLine = "baz.go:75:4: unknown field 'Wee' in struct literal of type RollerCoaster"
	if gotLine != wantLine {
		t.Errorf("unexpected line 3.\n  got:%q\n want:%q", gotLine, wantLine)
	}
}
