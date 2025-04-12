package led

const ledCount = 12

var leds []Led

func init() {
	leds = make([]Led, ledCount)
	for i := 0; i < ledCount; i++ {
		leds[i] = Led{
			ID: i,
			R:  0,
			G:  0,
			B:  0,
		}
	}
}
