package buttons

import (
	"context"
	evdev "github.com/gvalkov/golang-evdev"
)

type ButtonClickCallback func(button Button, clickType ClickType, down bool)

type EventSubscription struct {
	btn    Button
	device *evdev.InputDevice
	cancel context.CancelFunc
}

func (e *EventSubscription) Cancel() {
	e.cancel()
}
