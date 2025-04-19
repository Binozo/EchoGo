package echo

import (
	"github.com/Binozo/EchoGo/v2/internal/bindings/buttons"
	"github.com/Binozo/EchoGo/v2/internal/bindings/led"
)

func (e *Echo) GetButtonController() buttons.Controller {
	return buttons.NewDefaultHttpController()
}

func (e *Echo) GetLedController() led.Controller {
	return led.NewDefaultHttpController()
}
