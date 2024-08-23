package server

import (
	"github.com/Binozo/EchoGoSDK/pkg/bindings/led"
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

		var jsonPayload struct {
			Leds []struct {
				Led int `json:"led"`
				R   int `json:"r"`
				G   int `json:"g"`
				B   int `json:"b"`
			} `json:"leds"`
		}

		err := c.ReadJSON(&jsonPayload)
		if err != nil {
			c.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			return
		}

		for _, curLed := range jsonPayload.Leds {
			led.SetColor(curLed.Led, curLed.R, curLed.G, curLed.B)
		}
	}
}
