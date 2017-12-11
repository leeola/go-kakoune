package main

import (
	"github.com/leeola/go-kakoune/_examples"
	"github.com/leeola/go-kakoune/api"
)

func main() {
	kak := api.New()

	err := kak.DefineCommand("gokakoune-hello",
		api.DefineCommandOptions{}, _examples.Hello)
	if err != nil {
		panic(err)
	}

	err = kak.DefineCommand("gokakoune-multicommand",
		api.DefineCommandOptions{}, _examples.MultiCommand...)
	if err != nil {
		panic(err)
	}
}
