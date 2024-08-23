package led

import (
	"bytes"
	"fmt"
)

type Led struct {
	// Specifies the target LED on the Echo Dot. Ranges between 0 and 11 (12 LEDs).
	Nr int
	// R color value of the LED
	R int
	// G color value of the LED
	G int
	// B color value of the LED
	B int
}

// BuildArgument converts the defined R, G and B values into a hex format in order to send it to the i2C device
// Example with color blue: rgb(0, 0, 255) => 3030 3030 4646 (without empty spaces)
func (l *Led) BuildArgument() []byte {
	hexR := fmt.Sprintf("%02X", l.R)
	hexG := fmt.Sprintf("%02X", l.G)
	hexB := fmt.Sprintf("%02X", l.B)
	var buffer bytes.Buffer
	buffer.WriteString(hexR)
	buffer.WriteString(hexG)
	buffer.WriteString(hexB)
	return buffer.Bytes()
}

// SetColor sets the color of the LED struct
// Needs to be sent before calling BuildArgument
func (l *Led) SetColor(r int, g int, b int) {
	l.R = r
	l.G = g
	l.B = b
}
