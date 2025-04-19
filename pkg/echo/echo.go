package echo

import (
	"errors"
	"github.com/electricbubble/gadb"
)

type Echo struct {
	adbClient *gadb.Client
	adbDevice *gadb.Device
}

func New() (*Echo, error) {
	// Init adb client
	adbClient, err := startAdbClient()
	if err != nil {
		return nil, err
	}

	echo := &Echo{
		adbClient: adbClient,
	}

	// Check if alexa is already connected
	_, err = echo.getAlexaAdbConnection()
	if err != nil && !errors.Is(err, ErrAlexaNotConnected) {
		return nil, err
	}

	return echo, nil
}
