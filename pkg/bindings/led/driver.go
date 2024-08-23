package led

import (
	"bytes"
	"os"
	"os/exec"
)

// i2C device that sets the current led
const ledCurrentPath = "/sys/devices/soc/11007000.i2c/i2c-0/0-003f/led_current"

// i2C device that seems to control brightness
const privacyBrightnessPath = "/sys/devices/soc/10010000.keypad/amz_privacy/privacy_brightness"

// i2C device that controls the actual LEDs
const ledFrame = "/sys/devices/soc/11007000.i2c/i2c-0/0-003f/frame"

// file permission we need to access the i2C device
const perm = os.FileMode(0644)

// Init the basic setup for setting LED colors
// Also stops the ledcontroller service
func Init() error {
	err := basicSetup()
	if err != nil {
		return err
	}

	// ledcontroller may overwrite our led config
	// solution: let android kill it
	cmd := exec.Command("stop", "ledcontroller")
	return cmd.Run()
}

// SetColor Sets the specified color on the specified LED nr
// You must provide a valid ledNr, otherwise nothing will happen
// ledNr consists of range 0 - LedCount
func SetColor(ledNr int, r int, g int, b int) error {
	var targetColor bytes.Buffer
	for index, led := range leds {
		if led.Nr == ledNr {
			led.SetColor(r, g, b)
			leds[index] = led
		}
		targetColor.Write(led.BuildArgument())
	}
	return os.WriteFile(ledFrame, targetColor.Bytes(), perm)
}

// SetColorAll Sets the specified color on all available LEDs
func SetColorAll(r int, g int, b int) error {
	var targetColor bytes.Buffer
	for index, led := range leds {
		led.SetColor(r, g, b)
		leds[index] = led
		targetColor.Write(led.BuildArgument())
	}
	return os.WriteFile(ledFrame, targetColor.Bytes(), perm)
}

// Initialises the LED i2C device in order for us to control it
func basicSetup() error {
	ledCurrentPacket := []byte{48}
	privacyBrightnessPacket := []byte{48, 0}

	err := os.WriteFile(ledCurrentPath, ledCurrentPacket, perm)
	if err != nil {
		return err
	}

	return os.WriteFile(privacyBrightnessPath, privacyBrightnessPacket, perm)
}
