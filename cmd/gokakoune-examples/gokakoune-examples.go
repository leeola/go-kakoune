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

	err = kak.DefineCommand("gokakoune-subprocs",
		api.DefineCommandOptions{}, _examples.Subprocs...)
	if err != nil {
		panic(err)
	}

	// err = kak.DefineCommand("gokakoune-subproc-nesting",
	// 	api.DefineCommandOptions{}, _examples.SubprocNesting...)
	// if err != nil {
	// 	panic(err)
	// }
}
