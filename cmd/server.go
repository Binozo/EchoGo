package main

import (
	"github.com/Binozo/EchoGoSDK/internal/server"
	"github.com/Binozo/EchoGoSDK/pkg/bindings/buttons"
	"github.com/Binozo/EchoGoSDK/pkg/bindings/led"
	"github.com/Binozo/EchoGoSDK/pkg/bindings/mic"
	"github.com/Binozo/EchoGoSDK/pkg/constants"
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
	log.Println("Listening on", constants.Port)

	err = server.Serve()
	if err != nil {
		panic(err)
	}
}
