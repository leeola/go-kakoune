package gurucomplete

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/leeola/gokakoune/util"
)

const guruBin = "guru"

const (
	lineUnknown int = iota
	lineMethod
	lineField
)

type Completion struct {
	Completion string
	File       string
	LineNo     int
	Column     int
}

// GuruComplete returns a set of completions for the given file/pos.
//
// NOTE(leeola): This function is not well optimized in any sense. It is parsing
// a slower output format from guru, not quite intended for code completion.
// Despite this, the function exists because i wanted code completion from syntax,
// and gocode wasn't working well enough for me. I'll likely be writing my own
// code completion eventually, unless Guru adds a native version of it.
func GuruComplete(filepath string, byteOffset int) ([]Completion, error) {
	stdout, _, exit, err := util.Exec(
		"guru", "describe", fmt.Sprintf("%s:#%d", filepath, byteOffset))
	if err != nil {
		return nil, err
	}

	if exit != 0 {
		return nil, fmt.Errorf("non-zero exit: %d", exit)
	}

	var (
		completions []Completion
		lineState   int
	)
	for i, line := range strings.Split(stdout, "\n") {
		if line == "" {
			break
		}

		split := strings.SplitN(line, ":", 3)
		if len(split) < 3 {
			return nil, fmt.Errorf("unexpected colum format in line: %d", i)
		}
		file, pos, desc := split[0], split[1], split[2]

		switch desc {
		case " Methods:":
			lineState = lineMethod
			continue
		}

		switch lineState {
		case lineMethod:
			desc = trimMethodPrefix(desc)

			// i believe [0] access is safe, i don't think Split will ever
			// return less than 1.
			sPos := strings.SplitN(pos, "-", 2)[0]

			sPosSplit := strings.SplitN(sPos, ".", 2)
			if len(sPosSplit) < 2 {
				return nil, fmt.Errorf("unexpected line.col format: %q", sPos)
			}

			lineNo, err := strconv.Atoi(sPosSplit[0])
			if err != nil {
				return nil, fmt.Errorf("failed to lineNo to int: %q", sPosSplit[0])
			}

			col, err := strconv.Atoi(sPosSplit[1])
			if err != nil {
				return nil, fmt.Errorf("failed to col to int: %q", sPosSplit[1])
			}

			completions = append(completions, Completion{
				// TODO(leeola): remove formatting on the desc for
				// methods. Guru adds lots of indentation, a method
				// prefix, etc.
				Completion: desc,
				File:       file,
				LineNo:     lineNo,
				Column:     col,
			})
		}
	}

	return completions, nil
}

func trimMethodPrefix(desc string) string {
	var spaceCount int
	for i, ru := range desc {
		switch ru {
		case ' ':
			spaceCount++
		}

		if spaceCount == 3 {
			return desc[i+1:]
		}
	}

	return ""
}
