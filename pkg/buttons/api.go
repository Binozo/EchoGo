package buttons

import (
	evdev "github.com/gvalkov/golang-evdev"
	"os/exec"
)

const dotButton = "/dev/input/event1"
const volumeButtons = "/dev/input/event2"

// Init the button listeners
// Kills alexa's native button functions
func Init() error {
	cmd := exec.Command("stop", "acebutton")
	return cmd.Run()
}

func listenForButtonClick(button string) (*evdev.InputEvent, error) {
	device, err := evdev.Open(button)
	if err != nil {
		return nil, err
	}
	inputEvent, err := device.ReadOne()
	if err != nil {
		return nil, err
	}
	return inputEvent, nil
}

func (b *Button) WaitForClick() (ClickType, error) {
	for {
		inputEvent, err := listenForButtonClick(b.internalName)
		if err != nil {
			return 0, err
		}
		if inputEvent.Value != 1 { // Not button down event
			continue
		}
		clickType := ClickType(inputEvent.Code)
		return clickType, nil
	}
}

func (b *Button) WaitForClickRelease() (ClickType, error) {
	for {
		inputEvent, err := listenForButtonClick(b.internalName)
		if err != nil {
			return 0, err
		}
		if inputEvent.Value != 0 { // Not button up event
			continue
		}
		clickType := ClickType(inputEvent.Code)
		return clickType, nil
	}
}
