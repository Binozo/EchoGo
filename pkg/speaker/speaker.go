package speaker

import (
	"github.com/Binozo/GoTinyAlsa/pkg/pcm"
	"github.com/Binozo/GoTinyAlsa/pkg/tinyalsa"
)

const CardNr = 0
const DeviceNr = 23

// GetDevice returns the pre-configured speaker alsa device
func GetDevice() tinyalsa.AlsaDevice {
	return tinyalsa.NewDevice(CardNr, DeviceNr, pcm.Config{
		Channels:    1,
		SampleRate:  48000,
		PeriodSize:  1024,
		PeriodCount: 2,
		Format:      tinyalsa.PCM_FORMAT_S16_LE,
	})
}
