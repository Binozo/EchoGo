package bridge

import (
	"errors"
	"fmt"
	"github.com/Binozo/EchoGo/v2/internal"
	"github.com/electricbubble/gadb"
	"log"
	"os"
	"path"
	"time"
)

const hostFileName = "build/server"
const remotePath = "/data/local/tmp"

type AdbBridge struct {
	client *gadb.Client
	device *gadb.Device
}

func NewAdbBridge() (Bridge, error) {
	adbClient, err := gadb.NewClient()
	if err != nil {
		return nil, err
	}

	// In case we restarted the host application we need to kill the server on the echo
	adbClient.DisconnectAll()

	return &AdbBridge{
		client: &adbClient,
	}, nil
}

func (a *AdbBridge) WaitForGettingOnline() error {
	for {
		connected, err := a.IsConnected()
		if err != nil {
			return err
		}
		if connected {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (a *AdbBridge) IsConnected() (bool, error) {
	_, err := a.getAlexaAdbConnection()
	if err != nil {
		if errors.Is(err, ErrAlexaNotConnected) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (a *AdbBridge) getAlexaAdbConnection() (*gadb.Device, error) {
	list, err := a.client.DeviceList()
	if err != nil {
		return nil, err
	}

	for _, device := range list {
		device.DeviceInfo()
		product := device.DeviceInfo()["product"]
		if product == "biscuit_puffin" {
			a.device = &device
			return a.device, nil
		}
	}
	return nil, ErrAlexaNotConnected
}

func (a *AdbBridge) Deploy(file *os.File) error {
	if err := a.device.PushFile(file, path.Join(remotePath, hostFileName)); err != nil {
		return err
	}

	if _, err := a.device.RunShellCommand("chmod", "+x", path.Join(remotePath, hostFileName)); err != nil {
		return err
	}
	return a.device.Forward(internal.Port, internal.Port, false)
}

func (a *AdbBridge) Run() error {
	out, err := a.device.RunShellCommand("nohup", fmt.Sprintf("\".%s\"", path.Join(remotePath, hostFileName)), ">/dev/null 2>&1", "&")
	log.Println(out)
	//out, err := a.device.RunShellCommand(fmt.Sprintf("./%s", path.Join(remotePath, hostFileName)))
	//fmt.Println(out)
	return err
}
