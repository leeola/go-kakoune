package tabnine

import (
	"fmt"

	"github.com/leeola/gokakoune/api"
)

func Plugin(k *api.KakInit) error {
	if err := k.DeclareOption([]string{"-hidden"}, "completions", "tabnine_completions"); err != nil {
		return fmt.Errorf("declareoption: %v", err)
	}

	err := k.Hook("global", "WinCreate", ".*", k.Expansion(func(k api.Kak) error {
		// experimenting with full completion, disabling other completers
		// k.Command("set-option", "window", "completers", "option=tabnine_completions", "%opt{completers}")
		k.Command("set-option", "window", "completers", "option=tabnine_completions")

		return nil
	}))
	if err != nil {
		return fmt.Errorf("window hook: %v", err)
	}

	opts := api.DefineCommandOptions{}
	err = k.DefineCommand("tabnine", opts, k.Callback(
		[]string{"cursor_line", "cursor_column", "timestamp", "bufname"},
		func(k api.Kak) error {
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

			header := fmt.Sprintf("%d.%d@%d", line, col, timestamp)
			compl := "foo|bar|baz"
			k.Command("set-option", "buffer="+bufname, "tabnine_completions", header, compl)

			k.Echo("comple..")
			return nil
		}))
	if err != nil {
		return fmt.Errorf("define cmd tabnine: %v", err)
	}

	return nil
}
