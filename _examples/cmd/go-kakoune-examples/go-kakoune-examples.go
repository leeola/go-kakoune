package main

import (
	"github.com/leeola/go-kakoune/_examples"
	"github.com/leeola/go-kakoune/api"
)

func main() {
	kak := api.New()

	kak.DefineCommand("go-kakoune-hello", api.DefineCommandOptions{}, _examples.Hello)
}
