package mic

import (
	"github.com/Binozo/GoTinyAlsa/pkg/pcm"
	"github.com/Binozo/GoTinyAlsa/pkg/tinyalsa"
	"os/exec"
	"time"
)

const CardNr = 0
const DeviceNr = 24

// Init the microphone i/o.
// We need to kill the mixer process first.
// We are literally fighting against the Android System.
func Init() error {
	cmd := exec.Command("killall", "mixer")
	err := cmd.Run()
	if err != nil {
		return err
	}
	// We need a cool down to release the alsa device
	time.Sleep(time.Millisecond * 250)
	return nil
}

// GetDevice returns the pre-configured microphone alsa device
func GetDevice() tinyalsa.AlsaDevice {
	return tinyalsa.NewDevice(CardNr, DeviceNr, pcm.Config{
		Channels:    9,
		SampleRate:  16000,
		PeriodSize:  512,
		PeriodCount: 5,
		Format:      tinyalsa.PCM_FORMAT_S24_LE,
	})
}
