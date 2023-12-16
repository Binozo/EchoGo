package speaker

import (
	"github.com/Binozo/GoTinyAlsa/pkg/pcm"
	"github.com/Binozo/GoTinyAlsa/pkg/tinyalsa"
)

const CardNr = 0
const DeviceNr = 25

// GetDevice returns the pre-configured speaker alsa device
func GetDevice() tinyalsa.AlsaDevice {
	return tinyalsa.NewDevice(CardNr, DeviceNr, pcm.Config{
		Channels:    2,
		SampleRate:  16000,
		PeriodSize:  512,
		PeriodCount: 4,
		Format:      tinyalsa.PCM_FORMAT_S16_LE,
	})
}
