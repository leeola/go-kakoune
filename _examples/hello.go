package _examples

import (
	"github.com/leeola/gokakoune/api"
)

var Hello = api.Subproc{
	Func: func(kak *api.Kak) error {
		kak.Echo("hello world")
		return nil
	},
}
