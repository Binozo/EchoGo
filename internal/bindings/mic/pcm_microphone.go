//go:build server

package mic

import (
	"context"
	"errors"
	"github.com/Binozo/EchoGo/v3/pkg/mic"
	"github.com/Binozo/GoTinyAlsa/pkg/pcm"
	"github.com/Binozo/GoTinyAlsa/pkg/tinyalsa"
	"os/exec"
)

const cardNr = 0
const deviceNr = 24

type PcmMicrophone struct {
	device *tinyalsa.AlsaDevice
}

// NewMicrophone returns the pre-configured microphone alsa device
func NewMicrophone() (*PcmMicrophone, error) {
	device := tinyalsa.NewDevice(cardNr, deviceNr, pcm.Config{
		Channels:    9,
		SampleRate:  16000,
		PeriodSize:  512,
		PeriodCount: 5,
		Format:      tinyalsa.PCM_FORMAT_S24_3LE,
	})
	mic := &PcmMicrophone{
		device: &device,
	}

	if err := mic.Init(); err != nil {
		return nil, err
	}
	return mic, nil
}

// Init the microphone i/o.
// Stop the mixer process which blocks our operations
func (p *PcmMicrophone) Init() error {
	cmd := exec.Command("stop", "mixer")
	return cmd.Run()
}

func (p *PcmMicrophone) Listen(callback mic.AudioCallback, context context.Context) error {
	if callback == nil {
		return errors.New("callback can't be nil")
	}
	stream := make(chan []byte)

	go func() {
		defer recover()
		defer close(stream)

		for {
			if context.Err() != nil {
				return
			}
			audio := <-stream
			callback(audio)
		}
	}()

	return p.device.GetAudioStream(p.device.DeviceConfig, stream)
}
