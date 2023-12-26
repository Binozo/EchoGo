package mic

import (
	"github.com/Binozo/GoTinyAlsa/pkg/pcm"
	"github.com/Binozo/GoTinyAlsa/pkg/tinyalsa"
	"os/exec"
)

const CardNr = 0
const DeviceNr = 24

// Init the microphone i/o.
// Stop the mixer process which blocks our operations
func Init() error {
	cmd := exec.Command("stop", "mixer")
	return cmd.Run()
}

// GetDevice returns the pre-configured microphone alsa device
func GetDevice() tinyalsa.AlsaDevice {
	return tinyalsa.NewDevice(CardNr, DeviceNr, pcm.Config{
		Channels:    9,
		SampleRate:  16000,
		PeriodSize:  512,
		PeriodCount: 5,
		Format:      tinyalsa.PCM_FORMAT_S24_3LE,
	})
}
