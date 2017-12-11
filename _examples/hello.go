package _examples

import (
	"github.com/leeola/go-kakoune/api"
)

var Hello = api.Command{
	Func: func(kak *api.Kak) error {
		kak.Echo("hello")
		return nil
	},
}
