package main

import (
	"fmt"
	"github.com/Binozo/EchoGo/v2/internal/bindings/buttons"
	"github.com/Binozo/EchoGo/v2/internal/bindings/mic"
	"github.com/Binozo/EchoGo/v2/internal/bindings/speaker"
	"github.com/Binozo/EchoGo/v2/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)
	log.Println("Initializing")

	buttonController, err := buttons.NewButtonController()
	if err != nil {
		log.Fatalf("Failed to initialize Button controller: %v", err)
	}

	microphone, err := mic.NewMicrophone()
	if err != nil {
		log.Fatalf("Failed to initialize Microphone: %v", err)
	}

	pcmSpeaker, err := speaker.NewPcmSpeaker()
	if err != nil {
		log.Fatalf("Failed to initialize PCM Speaker: %v", err)
	}

	s := server.NewServer(buttonController, microphone, pcmSpeaker)
	log.Println("Starting server")

	if err := s.Serve(); err != nil {
		if strings.Contains(err.Error(), "address already in use") {
			log.Println("Server is already running, killing")
			response, err := http.Get(fmt.Sprintf("http://localhost:%d/kill", server.Port))
			if err != nil {
				log.Fatalf("Failed to send request to server: %v", err)
			}
			body, _ := io.ReadAll(response.Body)
			log.Println("Kill response from server:", string(body))
			response.Body.Close()
			time.Sleep(time.Millisecond * 100)
			log.Println("Now starting this instance")
			if err = s.Serve(); err != nil {
				log.Fatalf("Failed to start server: %v", err)
			}
		}
		log.Fatal(err)
	}
}
