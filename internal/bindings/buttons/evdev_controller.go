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

func (e *EvDevController) SubscribeToButton(button Button, callback ButtonClickCallback) (*EventSubscription, error) {
	if callback == nil {
		return nil, errors.New("callback can't be nil")
	}

	device, err := evdev.Open(button.internalName)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	eventSub := &EventSubscription{
		btn:    button,
		device: device,
		cancel: cancel,
	}

	go func() {
		defer func() {
			device.Release()
		}()

		beforeClickType := ClickType(0)
		beforeDown := false

		for {
			if ctx.Err() != nil {
				return
			}
			inputEvent, err := device.ReadOne()
			if err != nil {
				cancel()
				return // TODO: Return error?
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
			callback(button, clickType, down)
		}
	}()
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
