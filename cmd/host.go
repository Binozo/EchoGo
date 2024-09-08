package main

import (
	"errors"
	"github.com/Binozo/EchoGoSDK/pkg/client/echohost"
	"log"
	"os"
)

func main() {
	log.Println("Starting up...")

	alexa, err := echohost.NewAlexa()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Checking connection")
	connected, err := alexa.IsConnected()
	if err != nil && !errors.Is(err, echohost.ErrAlexaNotConnected) {
		log.Fatalln(err)
	}
	if !connected {
		log.Println("Not connected, booting...")
		if err = alexa.Boot(); err != nil {
			log.Fatalln(err)
		}
	}
	log.Println("Built connection with alexa")
	log.Println("Deploying server")
	if err = alexa.DeployServer(); err != nil {
		log.Fatalln(err)
	}
	log.Println("Remote server up and running")

	///////// Your code here /////////
	// Debug Light
	ledCtrl, err := alexa.GetLedControl()
	if err != nil {
		log.Println("Couldn't access led control")
		log.Fatalln(err)
	}
	if err = ledCtrl.SetColor(0, 255, 0); err != nil {
		log.Println("Couldn't set led")
		log.Fatalln(err)
	}
	ledCtrl.Close()

	// Debug Button listener
	btnCtrl, err := alexa.GetButtonListener()
	if err != nil {
		log.Println("Couldn't access buttons")
		log.Fatalln(err)
	}
	event, err := btnCtrl.WaitForClickEvent()
	if err != nil {
		log.Println("Couldn't wait for click")
		log.Fatalln(err)
	}
	log.Println("Received click event:", event.String())
	btnCtrl.Close()

	speakerControl, err := alexa.GetSpeakerControl()
	if err != nil {
		log.Println("Couldn't access speaker control")
		log.Fatalln(err)
	}

	// This specific operation is non-blocking
	// Your wav file MUST be in 48000kHz 2 channel S16_LE format otherwise your ears will suffer
	myWavFile, _ := os.ReadFile("music.wav")
	if err = speakerControl.Write(myWavFile); err != nil {
		log.Println("Couldn't write speaker control")
		log.Fatalln(err)
	}

	// Listen mic
	// Keep in mind that you get raw pcm data with 9 channels, 16kHz and S24_3LE format
	micCtrl, err := alexa.GetMicListener()
	if err != nil {
		log.Println("Couldn't access mic")
		log.Fatalln(err)
	}
	for {
		read, err := micCtrl.Read()
		if err != nil {
			panic(err)
		}
		log.Println("Read", len(read), "bytes from alexa microphone")
	}
}
