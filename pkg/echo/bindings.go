package echo

import (
	"github.com/Binozo/EchoGo/v3/pkg/buttons"
	"github.com/Binozo/EchoGo/v3/pkg/led"
	"github.com/Binozo/EchoGo/v3/pkg/mic"
	"github.com/Binozo/EchoGo/v3/pkg/speaker"
)

func (e *Echo) GetButtonController() buttons.Controller {
	return buttons.NewDefaultHttpController()
}

func (e *Echo) GetLedController() led.Controller {
	return led.NewDefaultHttpController()
}

func (e *Echo) GetMicrophone() mic.Microphone {
	return mic.NewDefaultHttpMicrophone()
}

func (e *Echo) GetSpeaker() speaker.Speaker {
	return speaker.NewDefaultHttpSpeaker()
}
