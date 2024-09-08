package server

import (
	"github.com/Binozo/EchoGo/v2/internal/payloads"
	"github.com/Binozo/EchoGo/v2/pkg/bindings/led"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func ledHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failure:", err)
		return
	}
	defer c.Close()

	for {
		var ledPayload payloads.LedsPayload

		err := c.ReadJSON(&ledPayload)
		if err != nil {
			c.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			return
		}

		for _, curLed := range ledPayload.Leds {
			led.SetColor(curLed.Led, curLed.R, curLed.G, curLed.B)
		}
	}
}
