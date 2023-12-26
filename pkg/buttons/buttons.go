package buttons

import (
	evdev "github.com/gvalkov/golang-evdev"
	"log"
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

func ListenForDotButton(callback func()) {
	go func() {
		for {
			_, err := listenForButtonClick(dotButton)
			if err != nil {
				log.Fatalf("Somehow couldn't fetch input events for %s: %s", dotButton, err)
			}
			callback()
		}
	}()
}

func ListenForVolumeButtons(volumeUpCallback func(), volumeDownCallback func()) {
	go func() {
		for {
			inputEvent, err := listenForButtonClick(volumeButtons)
			if err != nil {
				log.Fatalf("Somehow couldn't fetch input events for %s: %s", volumeButtons, err)
			}
			if inputEvent.Code == 114 {
				volumeDownCallback()
			} else if inputEvent.Code == 115 {
				volumeUpCallback()
			} else {
				log.Printf("Unknown input event code %d. Please create a pull request!", inputEvent.Code)
			}
		}
	}()
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
