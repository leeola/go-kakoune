package _examples

import (
	"github.com/leeola/go-kakoune/api"
)

var Hello = api.Command{
	Func: func(kak *api.Kak, p api.Params, v api.Vars) error {
		kak.Echo("hello")
		return nil
	},
}
