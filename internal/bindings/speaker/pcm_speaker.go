package speaker

import (
	"github.com/Binozo/GoTinyAlsa/pkg/pcm"
	"github.com/Binozo/GoTinyAlsa/pkg/tinyalsa"
	"os/exec"
)

const cardNr = 0
const deviceNr = 23

type PcmSpeaker struct {
	device       *tinyalsa.AlsaDevice
	audioSession *tinyalsa.AudioSession
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
	return nil
}

func (p *PcmSpeaker) Pump(data []byte) error {
	return p.audioSession.Pump(data)
}

func (p *PcmSpeaker) Close() {
	p.audioSession.Close()
}
