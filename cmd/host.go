package main

import (
	"bytes"
	"context"
	"errors"
	"github.com/Binozo/EchoGo/v2/internal/bindings/led"
	"github.com/Binozo/EchoGo/v2/pkg/buttons"
	"github.com/Binozo/EchoGo/v2/pkg/echo"
	"log"
	"os"
	"os/exec"
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
		// !!! We need to sleep here for 5 seconds because executing `alexa.Deploy()` disables SELinux enforcing
		// Doing this too early in the boot process breaks the microphone access functionality
		time.Sleep(5 * time.Second)
		log.Println("Deploying server app to alexa...")
		if err = alexa.Deploy(echo.DefaultServerPath); err != nil {
			log.Fatal(err)
		}
		log.Println("Bootup completed")
		time.Sleep(time.Second)
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
	mic := alexa.GetMicrophone()
	speaker := alexa.GetSpeaker()

	log.Println("Press the Dot button to start recording")
	clearLeds(ledController)
	ledController.SetLEDs(led.Led{
		ID: 7,
		R:  255,
		G:  255,
		B:  255,
	}, led.Led{
		ID: 8,
		R:  255,
		G:  255,
		B:  255,
	})

	recording := false
	animatingCtx, cancel := context.WithCancel(context.Background())
	var recordingBuffer bytes.Buffer
	btnSub, err := btn.SubscribeToButton(func(event buttons.ButtonClickEvent) {
		if event.ClickType == buttons.DotClick && !event.Down {
			if recording {
				cancel()
				clearLeds(ledController)
				convertedAudio, err := convertRecordedAudio(recordingBuffer.Bytes())
				if err != nil {
					log.Fatal(err)
				}
				animatingCtx, cancel = context.WithCancel(context.Background())

				ledController.SetLEDs(led.Led{
					ID: 7,
					G:  255,
				}, led.Led{
					ID: 8,
					G:  255,
				})

				log.Println("Playing recording")
				if err = speaker.Pump(convertedAudio); err != nil {
					log.Fatal(err)
				}
				recording = false

				ledController.SetLEDs(led.Led{
					ID: 7,
					R:  255,
					G:  255,
					B:  255,
				}, led.Led{
					ID: 8,
					R:  255,
					G:  255,
					B:  255,
				})
			} else {
				// Start recording
				log.Println("Now recording")
				recording = true
				recordingBuffer.Reset()
				go animateRecording(ledController, animatingCtx)
				go func() {
					err := mic.Listen(func(audioData []byte) {
						recordingBuffer.Write(audioData)
					}, animatingCtx)
					if err != nil {
						if !errors.Is(err, context.Canceled) {
							log.Fatal(err)
						}
					}
				}()
			}
		}
	})
	defer btnSub.Cancel()
	if err != nil {
		log.Fatal(err)
	}

	// Keep running forever
	select {}
}

func convertRecordedAudio(data []byte) ([]byte, error) {
	log.Println("Converting", len(data), "bytes")
	cmd := exec.Command("ffmpeg", "-f", "s24le", "-ar", "16000", "-ac", "9", "-i", "pipe:0", "-af", "pan=stereo|c0=c0|c1=c1", "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	cmd.Stdin = bytes.NewReader(data)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	log.Println("Converted to", len(out.Bytes()), "bytes")
	return out.Bytes(), nil
}

func animateRecording(controller led.Controller, ctx context.Context) {
	for {
		if ctx.Err() != nil {
			return
		}

		controller.SetLEDs(led.Led{
			ID: 7,
			R:  255,
			G:  255,
			B:  255,
		}, led.Led{
			ID: 8,
			R:  255,
			G:  255,
			B:  255,
		})

		time.Sleep(time.Millisecond * 500)

		if ctx.Err() != nil {
			return
		}
		controller.SetLEDs(led.Led{
			ID: 7,
		}, led.Led{
			ID: 8,
		})
		time.Sleep(time.Millisecond * 500)
	}
}

func clearLeds(ledController led.Controller) error {
	numLeds, err := ledController.GetNumLEDs()
	if err != nil {
		return err
	}

	ledData := make([]led.Led, numLeds)
	for i := 0; i < numLeds; i++ {
		ledData[i] = led.Led{
			ID: i,
			R:  0,
			G:  0,
			B:  0,
		}
	}
	return ledController.SetLEDs(ledData...)
}
