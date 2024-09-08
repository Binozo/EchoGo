package server

import (
	"github.com/Binozo/EchoGo/pkg/bindings/mic"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func micHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failure:", err)
		return
	}
	defer c.Close()

	micDevice := mic.GetDevice()
	audioChan := make(chan []byte)
	go func() {
		err = micDevice.GetAudioStream(micDevice.DeviceConfig, audioChan)
		if err != nil {
			log.Println("GetAudioStream failure:", err)
			c.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			c.Close()
			return
		}
	}()
	defer close(audioChan)
	for {
		err := c.WriteMessage(websocket.BinaryMessage, <-audioChan)
		if err != nil {
			return // Client disconnected
		}
	}
}
