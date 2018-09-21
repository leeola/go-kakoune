package compilecheck

import (
	"fmt"
	"strings"
	"time"

	"github.com/leeola/gokakoune/api"
	"github.com/leeola/gokakoune/api/vars"
	"github.com/leeola/gokakoune/errorlines"
)

var CompileCheckExpressions = api.Expansions{
	api.Func{
		ExportVars: []string{
			vars.BufFile,
		},
		Func: func(kak *api.Kak) error {
			buffile, err := kak.Var("buffile")
			if err != nil {
				return err
			}

			errLines, err := errorlines.GoTest(buffile)
			if err != nil {
				return err
			}

			// if there are no lines, clear the err output
			if len(errLines) == 0 {
				kak.Command("set-option", "buffer", "code_err", "false")
				kak.Command("remove-highlighter", "window/flag-lines_default_code_errors")
				return nil
			}

			var (
				first_file       string
				first_line       string
				first_desc       string
				code_errors_line []interface{}
			)

			for i, line := range errLines {
				cols := strings.SplitN(line, ":", 4)
				if len(cols) < 4 {
					return fmt.Errorf("incorrectly formatted error line: %q", line)
				}

				// Block left for reference.
				// file := cols[0]
				lineNo := cols[1]
				// col := cols[2]
				// desc := cols[3]

				if i == 0 {
					first_file = cols[0]
					first_line = lineNo
					first_desc = cols[3]
				}

				code_errors_line = append(code_errors_line, fmt.Sprintf("%s|{red,default}x", lineNo))
			}

			// TODO(leeola): store the file too? Pretty sure the reference implementation
			// is just assuming the same file, but that's faulty.

			// TODO(leeola): make all these commands native Go commands.
			kak.Command("set-option", "buffer", "code_err_file", first_file)
			kak.Command("set-option", "buffer", "code_err_line", first_line)
			kak.Command("set-option", "buffer", "code_err_desc", first_desc)

			// Clear previously assigned hightlighter. Otherwise kak fails.
			kak.Command("remove-highlighter", "window/flag-lines_default_code_errors")

			kak.Command("add-highlighter", "window/", "flag-lines", "default", "code_errors")
			kak.Command("set-option", "buffer", "code_err", "true")
			kak.Command("set-option", append([]interface{}{"global", "code_errors", time.Now().Unix()}, code_errors_line...)...)

			return nil
		},
	},
}
