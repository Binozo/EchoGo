package server

import (
	"github.com/Binozo/EchoGoSDK/pkg/bindings/speaker"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func speakerHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failure:", err)
		return
	}
	defer c.Close()

	speakerDevice := speaker.GetDevice()
	session, err := speakerDevice.NewAudioSession()
	if err != nil {
		c.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}
	defer session.Close()

	for {
		_, binaryMusicData, err := c.ReadMessage()
		if err != nil {
			break
		}
		if err = session.Pump(binaryMusicData); err != nil {
			c.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			return
		}
	}
}
