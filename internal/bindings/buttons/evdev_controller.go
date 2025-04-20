package buttons

import (
	"context"
	"errors"
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

func (e *EvDevController) SubscribeToButton(callback ButtonClickCallback) (*EventSubscription, error) {
	if callback == nil {
		return nil, errors.New("callback can't be nil")
	}

	dotBtn := e.GetDotButton()
	volBtn := e.GetVolumeButton()
	dotDevice, err := evdev.Open(dotBtn.internalName)
	if err != nil {
		return nil, err
	}
	volDevice, err := evdev.Open(volBtn.internalName)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	eventSub := NewEventSubscription(cancel)

	readBtn := func(btn Button, btnDevice *evdev.InputDevice) {
		defer btnDevice.Release()

		beforeClickType := ClickType(0)
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

			clickType := ClickType(inputEvent.Code)
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
			callback(ButtonClickEvent{
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

func (e *EvDevController) GetVolumeButton() Button {
	return Button{
		internalName: volumeButton,
		Type:         VolumeButton,
	}
}

func (e *EvDevController) GetDotButton() Button {
	return Button{
		internalName: dotButton,
		Type:         DotButton,
	}
}
