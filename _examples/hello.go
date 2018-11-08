package _examples

import (
	"github.com/leeola/gokakoune/api"
)

var Hello = api.Exp{
	Content: func(kak *api.Kak) error {
		kak.Echo("hello world")
		return nil
	},
}
