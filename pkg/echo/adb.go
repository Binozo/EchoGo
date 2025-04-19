package echo

import (
	"errors"
	"github.com/electricbubble/gadb"
	"net"
	"os/exec"
)

func startAdbClient() (*gadb.Client, error) {
	adbClient, err := gadb.NewClient()
	if err != nil {
		opError := &net.OpError{}
		if errors.As(err, &opError) {
			// adb server isn't running
			if startErr := startAdbServer(); startErr != nil {
				// adb binary not found
				return nil, err
			}
			// Try again
			adbClient, err = gadb.NewClient()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return &adbClient, nil
}

func startAdbServer() error {
	_, err := exec.Command("adb", "start-server").CombinedOutput()
	return err
}

func (e *Echo) getAlexaAdbConnection() (*gadb.Device, error) {
	list, err := e.adbClient.DeviceList()
	if err != nil {
		return nil, err
	}

	for _, device := range list {
		device.DeviceInfo()
		product := device.DeviceInfo()["product"]
		if product == "biscuit_puffin" {
			e.adbDevice = &device
			return e.adbDevice, nil
		}
	}
	return nil, ErrAlexaNotConnected
}
