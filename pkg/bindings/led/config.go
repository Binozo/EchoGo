package led

const ledCount = 13

// LedCount Defines how many LEDs are available
const LedCount = ledCount - 1

var leds []Led

func init() {
	leds = make([]Led, ledCount)
	for i := 0; i < ledCount; i++ {
		leds[i] = Led{
			Nr: i,
			R:  0,
			G:  0,
			B:  0,
		}
	}
}
