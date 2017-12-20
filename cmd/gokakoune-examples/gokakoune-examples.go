package main

import (
	"github.com/leeola/gokakoune/_examples"
	"github.com/leeola/gokakoune/api"
)

func main() {
	kak := api.New()

	err := kak.DefineCommand("gokakoune-hello",
		api.DefineCommandOptions{}, _examples.Hello)
	if err != nil {
		panic(err)
	}

	err = kak.DefineCommand("gokakoune-expressions",
		api.DefineCommandOptions{}, _examples.Expressions...)
	if err != nil {
		panic(err)
	}

	err = kak.DefineCommand("gokakoune-prompt",
		api.DefineCommandOptions{}, _examples.Prompt...)
	if err != nil {
		panic(err)
	}
}
