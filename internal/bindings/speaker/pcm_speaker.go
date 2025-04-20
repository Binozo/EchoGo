//go:build server

package speaker

import (
	"context"
	"github.com/Binozo/GoTinyAlsa/pkg/pcm"
	"github.com/Binozo/GoTinyAlsa/pkg/tinyalsa"
	"os/exec"
	"time"
)

const cardNr = 0
const deviceNr = 23

type PcmSpeaker struct {
	device       *tinyalsa.AlsaDevice
	audioSession *tinyalsa.AudioSession

	noiseSuppressorCancel context.CancelFunc
}

// NewPcmSpeaker returns the pre-configured speaker alsa device
func NewPcmSpeaker() (*PcmSpeaker, error) {
	device := tinyalsa.NewDevice(cardNr, deviceNr, pcm.Config{
		Channels:    2,
		SampleRate:  48000,
		PeriodSize:  1024,
		PeriodCount: 2,
		Format:      tinyalsa.PCM_FORMAT_S16_LE,
	})
	speaker := &PcmSpeaker{
		device: &device,
	}
	if err := speaker.Init(); err != nil {
		return nil, err
	}
	return speaker, nil
}

func (p *PcmSpeaker) Init() error {
	cmd := exec.Command("stop", "mixer")
	if err := cmd.Run(); err != nil {
		return err
	}

	audioSession, err := p.device.NewAudioSession()
	if err != nil {
		return err
	}
	p.audioSession = &audioSession

	context, cancel := context.WithCancel(context.Background())
	p.noiseSuppressorCancel = cancel

	go func() {
		// This is very special. When we kill the mixer process and do not play any sound the speaker will emit random electronic noise
		// If you want to know how that hears like then simply put a return before this for loop
		// Do bypass this problem we randomly emit no noise every second
		// This bypass will be automatically disabled when you start playing sound yourself to ensure smooth playback
		for {
			time.Sleep(1 * time.Second)

			if context.Err() != nil {
				p.noiseSuppressorCancel = nil
				return
			}

			antiNoiseData := make([]byte, 512)
			if err := p.audioSession.Pump(antiNoiseData); err != nil {
				return
			}

		}
	}()
	return nil
}

func (p *PcmSpeaker) Pump(data []byte) error {
	if p.noiseSuppressorCancel != nil {
		p.noiseSuppressorCancel()
	}
	return p.audioSession.Pump(data)
}

func (p *PcmSpeaker) Close() {
	p.audioSession.Close()
}
