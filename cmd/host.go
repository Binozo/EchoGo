package main

import (
	"github.com/Binozo/EchoGo/v2/internal/bindings/buttons"
	"github.com/Binozo/EchoGo/v2/internal/bindings/led"
	"github.com/Binozo/EchoGo/v2/pkg/echo"
	"log"
	"os"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)
	log.Println("Starting up...")

	alexa, err := echo.New()
	if err != nil {
		log.Fatal(err)
	}

	isOnline, err := alexa.IsOnline()
	if err != nil {
		log.Fatal(err)
	}

	if !isOnline {
		log.Println("Alexa is not online. Booting...")
		if err = alexa.Boot(echo.DefaultPreloaderPath); err != nil {
			log.Fatal(err)
		}
		log.Println("Deploying server app to alexa...")
		if err = alexa.Deploy(echo.DefaultServerPath); err != nil {
			log.Fatal(err)
		}
		log.Println("Bootup completed")
	} else {
		log.Println("Alexa is already online")

		// Check if server is online
		if err = alexa.Ping(); err != nil {
			log.Println("Server is not reachable. Starting...")
			if err = alexa.Deploy(echo.DefaultServerPath); err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second)
		}
	}

	btn := alexa.GetButtonController()
	ledController := alexa.GetLedController()
	btnSub, err := btn.SubscribeToButton(func(clickEvent buttons.ButtonClickEvent) {
		color := uint8(255)
		if !clickEvent.Down {
			color = uint8(0)
		}
		switch clickEvent.ClickType {
		case buttons.DotClick:
			err := ledController.SetLEDs(led.Led{
				ID: 7,
				R:  color,
				G:  color,
				B:  color,
			}, led.Led{
				ID: 8,
				R:  color,
				G:  color,
				B:  color,
			})
			if err != nil {
				log.Fatal(err)
			}
			break
		case buttons.VolumeUpClick:
			err := ledController.SetLEDs(led.Led{
				ID: 4,
				R:  color,
				G:  color,
				B:  color,
			}, led.Led{
				ID: 5,
				R:  color,
				G:  color,
				B:  color,
			})
			if err != nil {
				log.Fatal(err)
			}
			break
		case buttons.VolumeDownClick:
			err := ledController.SetLEDs(led.Led{
				ID: 11,
				R:  color,
				G:  color,
				B:  color,
			}, led.Led{
				ID: 10,
				R:  color,
				G:  color,
				B:  color,
			})
			if err != nil {
				log.Fatal(err)
			}
			break
		}

	})
	if err != nil {
		log.Fatal(err)
	}
	defer btnSub.Cancel()

	// Keep running forever
	select {}
}
