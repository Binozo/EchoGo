package main

import (
	"github.com/Binozo/EchoGo/internal/server"
	"github.com/Binozo/EchoGo/pkg/bindings/buttons"
	"github.com/Binozo/EchoGo/pkg/bindings/led"
	"github.com/Binozo/EchoGo/pkg/bindings/mic"
	"github.com/Binozo/EchoGo/pkg/constants"
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
