package echo

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"github.com/Binozo/EchoGo/v2/internal/server"
	"log"
	"net/http"
	"os"
	"time"
)

const DefaultServerPath = "build/server"

const remoteServerPath = "/data/local/bin/server"
const remoteInitScriptPath = "/system/etc/init/echogo.rc"
const remoteServerScriptPath = "/data/local/bin/start_server.sh"
const serviceName = "echogo"

//go:embed echogo.rc
var initScript []byte

//go:embed start_server.sh
var serverScript []byte

func (e *Echo) Deploy(serverPath string) error {
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

	// Push init script
	if err = e.adbDevice.Push(bytes.NewReader(initScript), remoteInitScriptPath, time.Now(), 644); err != nil {
		return err
	}

	// Push server script
	if err = e.adbDevice.Push(bytes.NewReader(serverScript), remoteServerScriptPath, time.Now(), 644); err != nil {
		return err
	}
	if _, err = e.adbDevice.RunShellCommand("chcon", "u:object_r:system_file:s0", remoteServerScriptPath); err != nil {
		log.Fatal(err)
	}

	if err = e.disableSELinux(); err != nil {
		return err
	}

	if alreadyRunning {
		_, err = e.adbDevice.RunShellCommand("start", serviceName)
		return err
	}

	return nil
}

func (e *Echo) unlockFilesystem() error {
	_, err := e.adbDevice.RunShellCommand("mount", "-o", "rw,remount", "rootfs", "/") // mount -o rw,remount rootfs /
	return err
}

func (e *Echo) lockFilesystem() error {
	_, err := e.adbDevice.RunShellCommand("mount", "-o", "ro,remount", "rootfs", "/") // mount -o ro,remount rootfs /
	return err
}

func (e *Echo) disableSELinux() error {
	_, err := e.adbDevice.RunShellCommand("setenforce", "0")
	return err
}

func (e *Echo) forwardPorts() error {
	/*
		Maybe create a PR for gadb to support this:
		Syntax: adb forward tcp:HOST_BIND_IP:HOST_PORT tcp:DEVICE_PORT
		adb forward tcp:localhost:6996 tcp:6996
		adb forward tcp:127.0.0.1:6996 tcp:6996
	*/
	return e.adbDevice.Forward(server.Port, server.Port)
}

func (e *Echo) Ping() error {
	res, err := http.Get(fmt.Sprintf("http://localhost:%d", server.Port))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}
	return nil
}
