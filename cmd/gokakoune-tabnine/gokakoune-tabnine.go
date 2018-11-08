package main

import (
	"fmt"

	"github.com/leeola/gokakoune/api"
	"github.com/leeola/gokakoune/plugins/tabnine"
)

func main() {
	k := api.New()
	if err := tabnine.Plugin(k); err != nil {
		panic(fmt.Sprintf("tabnine: %v", err))
	}
}
