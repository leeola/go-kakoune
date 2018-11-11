// Tabnine plugin
package tnp

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/leeola/gokakoune/api"
	"github.com/leeola/gokakoune/tabnine"
)

const (
	lookBeforeMax      = 3000
	lookAfterMax       = 3000
	pluginConfigSubdir = "tabnine"
)

func Plugin(k *api.KakInit) error {
	if err := k.DeclareOption([]string{"-hidden"}, "completions", "tabnine_completions"); err != nil {
		return fmt.Errorf("declareoption: %v", err)
	}

	if err := k.DeclareOption([]string{"-hidden"}, "str", "tabnine_before"); err != nil {
		return fmt.Errorf("declareoption tabnine_before: %v", err)
	}

	if err := k.DeclareOption([]string{"-hidden"}, "str", "tabnine_after"); err != nil {
		return fmt.Errorf("declareoption tabnine_after: %v", err)
	}

	err := k.Hook(nil, "global", "WinCreate", ".*", k.Expansion(func(k api.Kak) error {
		// experimenting with full completion, disabling other completers
		// k.Command("set-option", "window", "completers", "option=tabnine_completions", "%opt{completers}")
		k.Command("set-option", "window", "completers", "option=tabnine_completions")

		k.Command("tabnine-start")

		k.Command("tabnine-prefetch")

		err := k.Hook(nil, "window", "InsertIdle", ".*", k.Expansion(func(k api.Kak) error {
			k.Debug("insert idle, tabnine")
			k.Command("tabnine")
			return nil
		}))
		if err != nil {
			return err
		}

		err = k.Hook(nil, "window", "BufWritePost", ".*", k.Expansion(func(k api.Kak) error {
			k.Command("tabnine-prefetch")
			return nil
		}))
		if err != nil {
			return err
		}

		return nil
	}))
	if err != nil {
		return fmt.Errorf("window hook: %v", err)
	}

	tabnineCallback := k.Callback(
		[]string{
			"cursor_line", "cursor_column", "timestamp", "bufname", "buffile",
			"config",
			"opt_tabnine_before", "opt_tabnine_after"},
		func(k api.Kak) error {
			buffile, err := k.Var("buffile")
			if err != nil {
				return err
			}

			if buffile == "*debug*" || buffile == "*scratch*" {
				return nil
			}

			bufname, err := k.Var("bufname")
			if err != nil {
				return err
			}
			line, err := k.VarInt("cursor_line")
			if err != nil {
				return err
			}
			col, err := k.VarInt("cursor_column")
			if err != nil {
				return err
			}
			timestamp, err := k.VarInt("timestamp")
			if err != nil {
				return err
			}
			before, err := k.Option("tabnine_before")
			if err != nil {
				return err
			}
			after, err := k.Option("tabnine_after")
			if err != nil {
				return err
			}

			kakConfigDir, err := k.Var("config")
			if err != nil {
				return err
			}

			tn, err := tabnine.New(tabnine.Config{
				ConfigDir: filepath.Join(kakConfigDir, pluginConfigSubdir),
			})
			if err != nil {
				return fmt.Errorf("tabnine new: %v", err)
			}

			res, err := tn.Autocomplete(tabnine.AutocompleteRequest{
				Filename:                buffile,
				Before:                  before,
				RegionIncludesBeginning: len(before) < lookBeforeMax,
				After:             after,
				RegionIncludesEnd: len(after) < lookAfterMax,
				MaxNumResults:     5,
			})
			if err != nil {
				return fmt.Errorf("tabnine autocomplete: %v", err)
			}

			if len(res.Results) == 0 {
				k.Debug("no results")
			}

			subLen := len(res.SuffixToSubstitute)
			subCol := col - subLen

			header := fmt.Sprintf("%d.%d+%d@%d", line, subCol, subLen, timestamp)
			args := []interface{}{"buffer=" + bufname, "tabnine_completions", header}
			for _, result := range res.Results {
				// trim the suffix substitute.
				//
				// This ensures that to complete `foo` with `foobarbaz`, you don't
				// result in `foofoobarbaz`.
				// text := strings.TrimPrefix(result.Result, res.SuffixToSubstitute)
				text := result.Result
				// escape pipes
				escapedText := strings.Replace(text, "|", "\\|", -1)
				escapedMenuText := strings.Replace(result.Result, "|", "\\|", -1)
				compl := escapedText + "||" + escapedMenuText
				args = append(args, compl)
			}

			k.Command("set-option", args...)

			return nil
		})

	opts := api.DefineCommandOptions{}
	err = k.DefineCommand("tabnine", opts, k.Expansion(func(k api.Kak) error {
		err := k.EvaluateCommands([]string{"-draft"}, k.Expansion(func(k api.Kak) error {
			// select the hundred lines before, and after.
			// the callback will handle the actual parsing of the selection.
			//
			// Note that the before and after selections are very important to the return
			// output of tabnine.
			//
			// It appears that before must include up until the cursor location,
			// including the last character typed.
			//
			// If the before/after split is done incorrectly, the tabnine functionality
			// breaks entirely.
			k.Printf("exec \"<space>;h%dH\"\n", lookBeforeMax)
			k.Command("set-option", "current", "tabnine_before", "%val{selection}")
			k.Printf("exec \"<a-;>;l%dL\"\n", lookAfterMax)
			k.Command("set-option", "current", "tabnine_after", "%val{selection}")
			return nil
		}))
		if err != nil {
			return err
		}

		return k.EvaluateCommands(nil, tabnineCallback)
	}),
	)
	if err != nil {
		return fmt.Errorf("define cmd tabnine: %v", err)
	}

	opts = api.DefineCommandOptions{}
	err = k.DefineCommand("tabnine-prefetch", opts, k.Callback(
		[]string{"buffile", "config"},
		func(k api.Kak) error {
			buffile, err := k.Var("buffile")
			if err != nil {
				return err
			}
			kakConfigDir, err := k.Var("config")
			if err != nil {
				return err
			}

			if buffile == "*debug*" || buffile == "*scratch*" {
				return nil
			}

			tn, err := tabnine.New(tabnine.Config{
				ConfigDir: filepath.Join(kakConfigDir, pluginConfigSubdir),
			})
			if err != nil {
				return fmt.Errorf("tabnine new: %v", err)
			}

			err = tn.Prefetch(tabnine.PrefetchRequest{Filename: buffile})
			if err != nil {
				return fmt.Errorf("tabnine prefetch: %v", err)
			}

			return nil
		}))
	if err != nil {
		return fmt.Errorf("define cmd tabnine-prefetch: %v", err)
	}

	opts = api.DefineCommandOptions{}
	err = k.DefineCommand("tabnine-start", opts, k.Callback(nil,
		func(k api.Kak) error {
			if len(os.Args) < 1 {
				return fmt.Errorf("os args missing exec value, first arg")
			}
			selfBin := os.Args[0]

			cmd := exec.Command(selfBin, "http-serve-background")
			return cmd.Run()
		}))
	if err != nil {
		return fmt.Errorf("define cmd tabnine-prefetch: %v", err)
	}

	return nil
}
