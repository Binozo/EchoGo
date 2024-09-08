package led

import "github.com/Binozo/EchoGo/v2/pkg/constants"

var leds []Led

func init() {
	leds = make([]Led, constants.LedCount+1)
	for i := 0; i <= constants.LedCount; i++ {
		leds[i] = Led{
			Nr: i,
			R:  0,
			G:  0,
			B:  0,
		}
	}
}
