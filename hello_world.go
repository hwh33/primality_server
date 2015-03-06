package main

import (
	"github.com/go-martini/martini"
)

func main() {
	m := martini.Classic()
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Get("/secretmessage", func() string {
		return "Hey bitch face"
	})
	m.Run()
}
