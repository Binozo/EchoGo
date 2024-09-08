package echohost

import (
	"errors"
	"github.com/Binozo/EchoGo/internal/mtk"
	"github.com/Binozo/EchoGo/internal/wsbindings"
	"github.com/Binozo/EchoGo/pkg/constants"
	"github.com/electricbubble/gadb"
	"log"
	"os"
	"path"
	"time"
)

type Alexa struct {
	adbClient *gadb.Client
	alexaCon  *gadb.Device
}

const waitForBootTimeout = time.Minute
const hostFileName = "host"
const remotePath = "/data/local/tmp"

func NewAlexa() (*Alexa, error) {
	adbClient, err := gadb.NewClient()
	if err != nil {
		return nil, err
	}
	return &Alexa{
		adbClient: &adbClient,
	}, nil
}

func (a *Alexa) GetAlexaAdbConnection() (*gadb.Device, error) {
	if a.alexaCon != nil && a.IsConnectionHealthy() {
		return a.alexaCon, nil
	}
	list, err := a.adbClient.DeviceList()
	if err != nil {
		return nil, err
	}

	for _, device := range list {
		device.DeviceInfo()
		product := device.DeviceInfo()["product"]
		if product == "biscuit_puffin" {
			a.alexaCon = &device
			return &device, nil
		}
	}
	return nil, ErrAlexaNotConnected
}

func (a *Alexa) IsConnected() (bool, error) {
	_, err := a.GetAlexaAdbConnection()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *Alexa) IsConnectionHealthy() bool {
	state, err := a.alexaCon.State()
	if err != nil {
		return false
	}
	return state == gadb.StateOnline
}

func (a *Alexa) Boot() error {
	err := mtk.Boot()
	if err != nil {
		return err
	}
	// Waiting for alexa to come online
	now := time.Now()
	for {
		time.Sleep(time.Millisecond * 100)
		connected, err := a.IsConnected()
		if err != nil && !errors.Is(err, ErrAlexaNotConnected) {
			return err
		}
		if connected {
			break
		}
		if time.Since(now) > waitForBootTimeout {
			return ErrAlexaBootTimeout
		}
	}
	// set the device variable
	_, err = a.GetAlexaAdbConnection()
	return err
}

func (a *Alexa) DeployServer() error {
	hostFile, err := os.Open(hostFileName)
	if err != nil {
		return err
	}
	defer hostFile.Close()
	if err = a.alexaCon.PushFile(hostFile, path.Join(remotePath, hostFileName)); err != nil {
		return err
	}

	if _, err = a.alexaCon.RunShellCommand("chmod", "+x", path.Join(remotePath, hostFileName)); err != nil {
		return err
	}

	port := constants.Port
	if err = a.alexaCon.Forward(port, port, false); err != nil {
		return err
	}

	go func() {
		command, err := a.alexaCon.RunShellCommand(path.Join(remotePath, hostFileName))
		if err != nil {
			log.Println("Failed to start remote server:", err.Error())
		}
		log.Println("Remote server stopped unexpectedly:", command)
	}()
	log.Println("Waiting for remote server to come online...")
	for {
		time.Sleep(time.Millisecond * 100)
		isOnline := wsbindings.CheckHealth()
		if isOnline {
			break
		}
	}
	log.Println("Remote server is online 🥳")
	return nil
}

func (a *Alexa) GetLedControl() (*wsbindings.LedControl, error) {
	ctrl := wsbindings.GetLedControl()
	return ctrl, ctrl.Connect()
}

func (a *Alexa) GetButtonListener() (*wsbindings.ButtonControl, error) {
	ctrl := wsbindings.GetButtonControl()
	return ctrl, ctrl.Connect()
}

func (a *Alexa) GetMicListener() (*wsbindings.MicControl, error) {
	ctrl := wsbindings.GetMicControl()
	return ctrl, ctrl.Connect()
}

func (a *Alexa) GetSpeakerControl() (*wsbindings.SpeakerControl, error) {
	ctrl := wsbindings.GetSpeakerControl()
	return ctrl, ctrl.Connect()
}
