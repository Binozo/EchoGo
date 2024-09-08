package led

import "github.com/Binozo/EchoGo/pkg/constants"

// Clear all LEDs and turn them off
func Clear() error {
	return SetColorAll(0, 0, 0)
}

// Fun with colors
func Fun() error {
	colors := []struct {
		red   int
		green int
		blue  int
	}{
		{red: 255, green: 0, blue: 0},
		{red: 255, green: 119, blue: 0},
		{red: 255, green: 255, blue: 0},
		{red: 119, green: 255, blue: 0},
		{red: 0, green: 255, blue: 0},
		{red: 0, green: 255, blue: 119},
		{red: 0, green: 255, blue: 255},
		{red: 0, green: 119, blue: 255},
		{red: 0, green: 0, blue: 255},
		{red: 119, green: 0, blue: 255},
		{red: 255, green: 0, blue: 255},
		{red: 255, green: 0, blue: 119},
	}

	for {
		for j := 0; j < constants.LedCount; j++ {
			for i := 0; i < constants.LedCount; i++ {
				color := colors[i]
				err := SetColor((i+j)%constants.LedCount, color.red, color.green, color.blue)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// LedPercent lights the LEDs respecting the percentual value
func LedPercent(percent float64) error {
	totalPercentPerLed := 1.0 / constants.LedCount
	targetMaxLED := (percent / totalPercentPerLed) - 1
	for i := 0; i < int(targetMaxLED); i++ {
		err := SetColor(i, 255, 255, 255)
		if err != nil {
			return err
		}
	}
	targetLedOpacity := totalPercentPerLed / (percent - ((targetMaxLED - 1) * totalPercentPerLed))
	opacity := int(255 * targetLedOpacity)
	return SetColor(int(targetMaxLED), opacity, opacity, opacity)
}
