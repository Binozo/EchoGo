package echo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"
)

const mtkClientLocation = "mtkclient"
const DefaultPreloaderPath = "preloader_no_hdr.bin"

func (e *Echo) Boot(preloaderPath string) error {
	if _, err := os.Stat(mtkClientLocation); errors.Is(err, os.ErrNotExist) {
		return errors.Join(fmt.Errorf("mtkclient is not installed at %s. Please install mtkclient", mtkClientLocation), err)
	}

	pythonPath := path.Join(mtkClientLocation, ".venv", "bin", "python3")
	if _, err := os.Stat(pythonPath); errors.Is(err, os.ErrNotExist) {
		return errors.Join(fmt.Errorf("virtual python environment at %s is missing. Please install mtkclient in a venv", pythonPath), err)
	}

	if _, err := os.Stat(preloaderPath); errors.Is(err, os.ErrNotExist) {
		return errors.Join(fmt.Errorf("preloader at %s not found", preloaderPath), err)
	}

	cmd := exec.Command(pythonPath, path.Join(mtkClientLocation, "mtk.py"), "plstage", fmt.Sprintf("--preloader=%s", preloaderPath))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	if err := e.WaitToComeOnline(ctx); err != nil {
		return err
	}

	// Disable SELinux policy because we don't care in this environment
	if err := e.disableSELinux(); err != nil {
		return err
	}

	return e.forwardPorts()
}

func (e *Echo) WaitToComeOnline(context context.Context) error {
	for {
		if context.Err() != nil {
			return context.Err()
		}

		_, err := e.getAlexaAdbConnection()
		if err != nil {
			if !errors.Is(err, ErrAlexaNotConnected) {
				return err
			}
		} else {
			// Connection found!
			return nil
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func (e *Echo) IsOnline() (bool, error) {
	_, err := e.getAlexaAdbConnection()
	if err != nil {
		if errors.Is(err, ErrAlexaNotConnected) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (e *Echo) Shutdown() error {
	_, err := e.getAlexaAdbConnection()
	if err != nil {
		if !errors.Is(err, ErrAlexaNotConnected) {
			return err
		}
	}

	_, err = e.adbDevice.RunShellCommand("reboot", "-p")
	return err
}

func (e *Echo) Restart(preloaderPath string) error {
	if err := e.Shutdown(); err != nil {
		return err
	}

	return e.Boot(preloaderPath)
}
