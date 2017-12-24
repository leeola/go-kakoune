package main

import "github.com/leeola/gokakoune/api"

type Foo struct{}

func Zazzle(kak *api.Kak) (*Foo, error) {
	kak.Echo("wee")

	return &Foo{}, nil
}
