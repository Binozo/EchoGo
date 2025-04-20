package mic

import (
	"context"
)

type AudioCallback func(audioData []byte)

type Microphone interface {
	Init() error
	Listen(callback AudioCallback, context context.Context) error
}
