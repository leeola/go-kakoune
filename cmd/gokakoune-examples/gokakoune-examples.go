package main

import (
	"fmt"

	"github.com/leeola/gokakoune/api"
)

func main() {
	kak := api.New()

	err := kak.DefineCommand("gokakoune-example-hello", api.DefineCommandOptions{}, api.Expansion{
		Body: func(kak api.Kak) error {
			kak.Echo("hello world!")
			return nil
		},
	})
	if err != nil {
		panic(err)
	}

	err = kak.DefineCommand("gokakoune-example-callback", api.DefineCommandOptions{},
		kak.Callback(nil, func(k api.Kak) error {
			k.Echo("hello world, from callback")
			return nil
		}))
	if err != nil {
		panic(err)
	}

	callbackB := kak.Callback(nil, func(kak api.Kak) error {
		kak.Echo("hello world, from nested callback")
		return nil
	})
	callbackA := kak.Callback(nil, func(k api.Kak) error {
		k.Echo("adding new command, dynamically and declaring callback")
		return k.DefineCommand("gokakoune-example-nested-callback-b", api.DefineCommandOptions{}, callbackB)
	})
	err = kak.DefineCommand("gokakoune-example-nested-callback", api.DefineCommandOptions{}, callbackA)
	if err != nil {
		panic(err)
	}

	promptCallback := kak.Callback([]string{"text"}, func(k api.Kak) error {
		promptText, err := k.Var("text")
		if err != nil {
			return fmt.Errorf("var: %v", err)
		}
		kak.Echof("hello %s!", promptText)
		return nil
	})
	err = kak.DefineCommand("gokakoune-example-prompt", api.DefineCommandOptions{},
		kak.Expansion(func(k api.Kak) error {
			k.Prompt("What is your name?", promptCallback)
			return nil
		}),
	)
	if err != nil {
		panic(err)
	}
}
