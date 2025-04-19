package echo

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func (e *Echo) Run(serverPath string) error {
	file, err := os.Open(serverPath)
	if err != nil {
		return err
	}
	defer file.Close()

	alreadyRunning := true
	if e.adbDevice == nil {
		alreadyRunning = false
		_, err = e.getAlexaAdbConnection()
		if err != nil {
			return err
		}
	}

	if err = e.unlockFilesystem(); err != nil {
		return errors.Join(errors.New("couldn't unlock filesystem"), err)
	}
	defer e.lockFilesystem()

	if alreadyRunning {
		// Stop the service
		e.adbDevice.RunShellCommand("stop", serviceName)
	}
	if err = e.adbDevice.Push(file, remoteServerPath, time.Now(), 777); err != nil {
		return err
	}

	cmd := exec.Command("adb", "shell", fmt.Sprintf(".%s", remoteServerPath))
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
