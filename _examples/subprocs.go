package _examples

import (
	"github.com/leeola/go-kakoune/api"
)

var Subprocs = []api.Subproc{
	{
		Func: func(kak *api.Kak) error {
			// TODO(leeola): add a set option here,
			// to demonstrate writing to kakoune from multiple
			// process calls.
			return nil
		},
	},
	{
		Func: func(kak *api.Kak) error {
			kak.Echo("hello from 2nd func")
			return nil
		},
	},
}
