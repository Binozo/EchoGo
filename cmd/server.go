package main

import (
	"github.com/Binozo/EchoGoSDK/internal/server"
	"github.com/Binozo/EchoGoSDK/pkg/bindings/buttons"
	led2 "github.com/Binozo/EchoGoSDK/pkg/bindings/led"
	"github.com/Binozo/EchoGoSDK/pkg/bindings/mic"
	"log"
)

func main() {
	log.Println("Initializing")
	err := mic.Init()
	if err != nil {
		panic(err)
	}
	err = led2.Init()
	if err != nil {
		panic(err)
	}
	err = led2.Clear()
	if err != nil {
		panic(err)
	}
	err = buttons.Init()
	if err != nil {
		panic(err)
	}
	log.Println("Listening on", server.Port)

	err = server.Serve()
	if err != nil {
		panic(err)
	}
}
