package _examples

import (
	"github.com/leeola/gokakoune/api"
)

var Expressions = api.Expansions{
	api.Func{
		Func: func(kak *api.Kak) error {
			kak.Echo("hello from 1st func")

			// TODO(leeola): add a set option here,
			// to demonstrate writing to kakoune from multiple
			// process calls.
			return nil
		},
	},
	api.Func{
		Func: func(kak *api.Kak) error {
			kak.Echo("hello from 2nd func")
			return nil
		},
	},
}
