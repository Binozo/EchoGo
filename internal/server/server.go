package server

import (
	"fmt"
	"github.com/Binozo/EchoGo/v2/internal"
	internalLed "github.com/Binozo/EchoGo/v2/internal/bindings/led"
	"github.com/Binozo/EchoGo/v2/pkg/buttons"
	"github.com/Binozo/EchoGo/v2/pkg/led"
	"github.com/Binozo/EchoGo/v2/pkg/mic"
	"github.com/Binozo/EchoGo/v2/pkg/speaker"
	"github.com/gin-gonic/gin"
	"golang.org/x/sys/unix"
	"log"
	"net/http"
	"os"
	"time"
)

const Port = internal.Port

type Server struct {
	router           *gin.Engine
	ledController    led.Controller
	buttonController buttons.Controller
	mic              mic.Microphone
	speaker          speaker.Speaker
}

func NewServer(buttonController buttons.Controller, microphone mic.Microphone, speaker speaker.Speaker) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	server := &Server{
		buttonController: buttonController,
		mic:              microphone,
		speaker:          speaker,
	}

	router.GET("/", server.rootHandler)
	router.GET("/kill", server.killHandler)
	router.GET("/ping", server.pingHandler)
	router.POST("/leds/set", server.ledsHandler)
	router.GET("/buttons", server.buttonHandler)
	router.GET("/microphone", server.microphoneHandler)
	router.POST("/speaker", server.speakerHandler)

	server.router = router

	go func() {
		uptime, err := getUptime()
		minUptime := time.Second * 90

		if err != nil || uptime < minUptime {
			// If we start too soon the native bootup from the echo will break (LEDs will spin forever)
			stillWait := minUptime - uptime
			log.Printf("Uptime is currently at %0.2fs, waiting %0.2fs for LED setup\n", uptime.Seconds(), stillWait.Seconds())
			time.Sleep(stillWait)
		}

		ledController, err := internalLed.NewDefaultController()
		if err != nil {
			log.Fatalf("Failed to initialize LED controller: %v", err)
		}

		server.ledController = ledController
		clearLeds(ledController)
	}()

	return server
}

func (s *Server) Serve() error {
	return s.router.Run(fmt.Sprintf(":%d", Port))
}

func (s *Server) rootHandler(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte("Echo up and running"))
}

func (s *Server) killHandler(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte("Bye bye"))
	c.Writer.Flush()
	go func() {
		os.Exit(0)
	}()
}

func (s *Server) pingHandler(c *gin.Context) {
	c.Status(http.StatusOK)

	if s.ledController == nil {
		return
	}
	go func() {
		numLEDs, err := s.ledController.GetNumLEDs()
		if err != nil {
			log.Printf("Error getting number of LEDs: %v\n", err)
			return
		}

		for i := 0; i < numLEDs; i++ {
			leds := make([]led.Led, numLEDs)

			for j := 0; j < numLEDs; j++ {
				if i == j {
					leds[j] = led.Led{
						ID: j,
						R:  255,
						G:  255,
						B:  255,
					}
				} else {
					leds[j] = led.Led{
						ID: j,
						R:  0,
						G:  0,
						B:  0,
					}
				}
			}

			if err := s.ledController.SetLEDs(leds...); err != nil {
				log.Printf("Error setting LEDs: %v\n", err)
				return
			}
			time.Sleep(time.Millisecond * 25)
		}

		// Turnoff
		for i := 255; i >= 0; i -= 25 {
			brightness := uint8(i)
			if brightness < 6 {
				brightness = 0
			}
			if err := s.ledController.SetLEDs(led.Led{
				ID: numLEDs - 1,
				R:  brightness,
				G:  brightness,
				B:  brightness,
			}); err != nil {
				log.Printf("Error setting LEDs: %v\n", err)
			}
			time.Sleep(time.Millisecond * 13)
		}

	}()
}

func clearLeds(ledController led.Controller) {
	// Clear LEDs
	numLEDs, err := ledController.GetNumLEDs()
	if err != nil {
		log.Fatalf("Failed to get number of LEDs: %v", err)
		return
	}

	leds := make([]led.Led, numLEDs)
	for i := 0; i < numLEDs; i++ {
		leds[i] = led.Led{
			ID: i,
			R:  0,
			G:  0,
			B:  0,
		}
	}
	if err = ledController.SetLEDs(leds...); err != nil {
		log.Fatalf("Failed to set LEDs: %v", err)
		return
	}
}

func getUptime() (time.Duration, error) {
	var info unix.Sysinfo_t
	if err := unix.Sysinfo(&info); err != nil {
		return time.Duration(0), err
	}
	return time.Second * time.Duration(info.Uptime), nil
}
