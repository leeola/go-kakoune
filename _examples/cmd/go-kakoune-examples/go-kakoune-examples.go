package main

import (
	"github.com/leeola/go-kakoune/_examples"
	"github.com/leeola/go-kakoune/api"
)

func main() {
	kak := api.New()

	err := kak.DefineCommand("go-kakoune-hello",
		api.DefineCommandOptions{}, _examples.Hello)
	if err != nil {
		panic(err)
	}

	kak.DefineCommand("go-kakoune-multicommand",
		api.DefineCommandOptions{}, _examples.MultiCommand...)
}
