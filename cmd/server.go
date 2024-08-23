package main

import (
	"github.com/Binozo/EchoGoSDK/internal/server"
	"github.com/Binozo/EchoGoSDK/pkg/buttons"
	"github.com/Binozo/EchoGoSDK/pkg/led"
	"github.com/Binozo/EchoGoSDK/pkg/mic"
	"log"
)

func main() {
	log.Println("Initializing")
	err := mic.Init()
	if err != nil {
		panic(err)
	}
	err = led.Init()
	if err != nil {
		panic(err)
	}
	err = led.Clear()
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
