package server

import (
	"github.com/Binozo/EchoGoSDK/internal/payloads"
	"github.com/Binozo/EchoGoSDK/pkg/bindings/buttons"
	"log"
	"net/http"
)

func buttonHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failure:", err)
		return
	}
	defer c.Close()

	dotButton := buttons.GetDotButton()
	dotBtnChan := dotButton.ListenForEvents()
	defer close(dotBtnChan)
	volumeButton := buttons.GetVolumeButton()
	volumeBtnChan := volumeButton.ListenForEvents()
	defer close(volumeBtnChan)

	notifyClickEvent := func(event buttons.ClickEvent) error {
		payload := payloads.ClickEvent{
			Button: event.Button.Type,
			Down:   event.Down,
			Type:   event.ClickType.String(),
		}
		return c.WriteJSON(payload)
	}

	for {
		select {
		case clickEvent := <-dotBtnChan:
			if err := notifyClickEvent(clickEvent); err != nil {
				return
			}
			break
		case clickEvent := <-volumeBtnChan:
			if err := notifyClickEvent(clickEvent); err != nil {
				return
			}
			break
		}
	}
}
