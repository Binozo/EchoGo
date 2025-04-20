package buttons

import (
	"context"
	"errors"
	"github.com/Binozo/EchoGo/v3/pkg/buttons"
	evdev "github.com/gvalkov/golang-evdev"
	"os/exec"
)

const dotButton = "/dev/input/event1"
const volumeButton = "/dev/input/event2"

type EvDevController struct {
}

// Init the button listeners
// Kills alexa's native button functions
func (e *EvDevController) Init() error {
	cmd := exec.Command("stop", "acebutton")
	return cmd.Run()
}

func (e *EvDevController) SubscribeToButton(callback buttons.ButtonClickCallback) (*buttons.EventSubscription, error) {
	if callback == nil {
		return nil, errors.New("callback can't be nil")
	}

	dotBtn := e.GetDotButton()
	volBtn := e.GetVolumeButton()
	dotDevice, err := evdev.Open(dotButton)
	if err != nil {
		return nil, err
	}
	volDevice, err := evdev.Open(volumeButton)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	eventSub := buttons.NewEventSubscription(cancel)

	readBtn := func(btn buttons.Button, btnDevice *evdev.InputDevice) {
		defer btnDevice.Release()

		beforeClickType := buttons.ClickType(0)
		beforeDown := false

		for {
			if ctx.Err() != nil {
				return
			}

			inputEvent, err := btnDevice.ReadOne()
			if err != nil {
				// TODO: What to do now?
				return
			}

			clickType := buttons.ClickType(inputEvent.Code)
			if inputEvent.Code != 0 {
				beforeClickType = clickType
			} else {
				clickType = beforeClickType
			}

			down := inputEvent.Value == 1
			if beforeDown == down {
				continue
			}
			beforeDown = down
			callback(buttons.ButtonClickEvent{
				Button:    btn,
				ClickType: clickType,
				Down:      down,
			})
		}
	}

	go readBtn(dotBtn, dotDevice)
	go readBtn(volBtn, volDevice)

	return eventSub, nil
}

func (e *EvDevController) GetVolumeButton() buttons.Button {
	return buttons.Button{
		Type: buttons.VolumeButton,
	}
}

func (e *EvDevController) GetDotButton() buttons.Button {
	return buttons.Button{
		Type: buttons.DotButton,
	}
}

func NewButtonController() (buttons.Controller, error) {
	controller := &EvDevController{}
	if err := controller.Init(); err != nil {
		return nil, err
	}
	return controller, nil
}
