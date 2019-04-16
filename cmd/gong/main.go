package main

import (
	"github.com/ibudiallo/gong"
	"github.com/ibudiallo/gong/server"
)

func main() {

	env := &gong.Env{
		Port: ":8080",
		Root: "html",
	}

	server.Init(env)

}
