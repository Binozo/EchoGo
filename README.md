<br />
<div align="center">
<h3 align="center">EchoGo</h3>

  <p align="center">
    A Go SDK for your Echo Dot **2. Gen**
    <br />
  </p>
</div>

<img src="assets/thumbnail.gif">

## ⚡ Features
- 🚥 Full control over LEDs
- 🎤 Access to microphone, speaker and buttons
- ⚙️ Preconfigured compiler

## 🗣️ Quick pre yapping 

As I am building my privacy focused smart home I wanted an assistant like Alexa or Google Home to talk to. I already owned an Alexa Echo Dot 2. Gen so I decided to take full control of it.
Thanks to [Dragon863's Project](https://github.com/Dragon863/EchoCLI) I was able to do this and create an SDK for it.
With this SDK you take full control over your Echo and you can abuse it as you wish. You get full root access.

<details>
  <summary>Some interesting facts</summary>
  <ol>
    <li>🤫 The mute button seems to be a hardware button. If you access the microphone while the microphone is muted you will get an empty byte stream.</li>
    <li>🎨 The LEDs fully support RGB, you are not limited to Blue and Red colors.</li>
    <li>🛜 You can find saved Wi-Fi connections at `/data/misc/wifi/wpa_supplicant.conf`</li>
    <li>🔊 Uses TinyAlsa under the hood; Microphone has 9 channels while the speaker has 2 channels</li>
    <li>🤖 The echo is actually an Android 7 device</li>
  </ol>
  <p></p>
</details>


## 📋 Requirements
- A fully rooted Echo Dot **2. Gen** ([Dragon863's Project](https://github.com/Dragon863/EchoCLI))
- You need your `preloader_no_hdr.bin` file for booting
- A linux machine (This project has not been tested on Windows and Mac)
- If you wish to compile the server yourself you need Docker installed

> [!CAUTION]
> Make sure your echo is not able to connect to the Internet. Use a temporary hotspot for setup if needed, make sure the Echo can't reach the Wi-Fi network after that otherwise your echo may get bricked after a few days/weeks due to OTA updates.

## 👷‍♂️ Setup
First make sure you rooted your Echo. Take a look at this guide [here](https://danieldb.uk/posts/alexa-2/).

This project consists of two parts: The host controller and the echo server.

The echo server runs on the echo and allows clients to control the echo through HTTP.
The host controller takes care of booting the echo and setting it up. Additionally, it can control it too.

Now let's get to the `SHIP IT 🚀`-part.

### Server

To control the echo remotely you need to run a server on it.
Either you can download it prebuilt from GitHub or you compile it yourself.

For both refer to the [CLI](#-cli) section.

### Host
Open your terminal in your project and first import this package:
```shell
$ go get -u https://github.com/Binozo/EchoGo/v3
```

Also make sure you got your `preloader_no_hdr.bin` file in your project root or working directory.

Then install mtkclient:
```shell
$ git clone --depth=1 --branch 2.0.1 https://github.com/bkerler/mtkclient.git
$ git submodule update --recursive
# Install python dependencies
$ cd mtkclient
$ python -m venv .venv
$ source .venv/bin/activate
$ pip install -r requirements.txt
$ pip install .
```

Now we actually ship 🚀

Here is a full Go example. This shows the use of every component of the echo including LEDs, buttons, microphones and the speaker.

Press the dot button once to start recording. Then press the dot button again to stop recording and the echo will play your recording through the speakers.
This code must be run on the host machine. I recommend you to use this code as a starting point in your project as it already handles the booting process. Modify it as you wish.

> [!NOTE]  
> This code may crash on the first time because the echo server needs to be deployed first. The server must be started by the android system itself which only happens on boot. Refer to the [CLI](#-cli) section to deploy the server first, reboot and then run this code.

```go
package main

import (
	"bytes"
	"context"
	"errors"
	"github.com/Binozo/EchoGo/v2/internal/bindings/buttons"
	"github.com/Binozo/EchoGo/v2/internal/bindings/led"
	"github.com/Binozo/EchoGo/v2/pkg/echo"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)
	log.Println("Starting up...")

	alexa, err := echo.New()
	if err != nil {
		log.Fatal(err)
	}

	isOnline, err := alexa.IsOnline()
	if err != nil {
		log.Fatal(err)
	}

	if !isOnline {
		log.Println("Alexa is not online. Booting...")
		if err = alexa.Boot(echo.DefaultPreloaderPath); err != nil {
			log.Fatal(err)
		}
		// !!! We need to sleep here for 5 seconds because executing `alexa.Deploy()` disables SELinux enforcing
		// Doing this too early in the boot process breaks the microphone access functionality
		time.Sleep(5 * time.Second)
		log.Println("Deploying server app to alexa...")
		if err = alexa.Deploy(echo.DefaultServerPath); err != nil {
			log.Fatal(err)
		}
		log.Println("Bootup completed")
		time.Sleep(time.Second)
	} else {
		log.Println("Alexa is already online")

		// Check if server is online
		if err = alexa.Ping(); err != nil {
			log.Println("Server is not reachable. Starting...")
			if err = alexa.Deploy(echo.DefaultServerPath); err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second)
		}
	}
	
	// Custom code after this

	btn := alexa.GetButtonController()
	ledController := alexa.GetLedController()
	mic := alexa.GetMicrophone()
	speaker := alexa.GetSpeaker()

	log.Println("Press the Dot button to start recording")
	clearLeds(ledController)
	ledController.SetLEDs(led.Led{
		ID: 7,
		R:  255,
		G:  255,
		B:  255,
	}, led.Led{
		ID: 8,
		R:  255,
		G:  255,
		B:  255,
	})

	recording := false
	animatingCtx, cancel := context.WithCancel(context.Background())
	var recordingBuffer bytes.Buffer
	btnSub, err := btn.SubscribeToButton(func(event buttons.ButtonClickEvent) {
		if event.ClickType == buttons.DotClick && !event.Down {
			if recording {
				cancel()
				clearLeds(ledController)
				convertedAudio, err := convertRecordedAudio(recordingBuffer.Bytes())
				if err != nil {
					log.Fatal(err)
				}
				animatingCtx, cancel = context.WithCancel(context.Background())

				ledController.SetLEDs(led.Led{
					ID: 7,
					G:  255,
				}, led.Led{
					ID: 8,
					G:  255,
				})

				log.Println("Playing recording")
				if err = speaker.Pump(convertedAudio); err != nil {
					log.Fatal(err)
				}
				recording = false

				ledController.SetLEDs(led.Led{
					ID: 7,
					R:  255,
					G:  255,
					B:  255,
				}, led.Led{
					ID: 8,
					R:  255,
					G:  255,
					B:  255,
				})
			} else {
				// Start recording
				log.Println("Now recording")
				recording = true
				recordingBuffer.Reset()
				go animateRecording(ledController, animatingCtx)
				go func() {
					err := mic.Listen(func(audioData []byte) {
						recordingBuffer.Write(audioData)
					}, animatingCtx)
					if err != nil {
						if !errors.Is(err, context.Canceled) {
							log.Fatal(err)
						}
					}
				}()
			}
		}
	})
	defer btnSub.Cancel()
	if err != nil {
		log.Fatal(err)
	}

	// Keep running forever
	select {}
}

func convertRecordedAudio(data []byte) ([]byte, error) {
	log.Println("Converting", len(data), "bytes")
	cmd := exec.Command("ffmpeg", "-f", "s24le", "-ar", "16000", "-ac", "9", "-i", "pipe:0", "-af", "pan=stereo|c0=c0|c1=c1", "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	cmd.Stdin = bytes.NewReader(data)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	log.Println("Converted to", len(out.Bytes()), "bytes")
	return out.Bytes(), nil
}

func animateRecording(controller led.Controller, ctx context.Context) {
	for {
		if ctx.Err() != nil {
			return
		}

		controller.SetLEDs(led.Led{
			ID: 7,
			R:  255,
			G:  255,
			B:  255,
		}, led.Led{
			ID: 8,
			R:  255,
			G:  255,
			B:  255,
		})

		time.Sleep(time.Millisecond * 500)

		if ctx.Err() != nil {
			return
		}
		controller.SetLEDs(led.Led{
			ID: 7,
		}, led.Led{
			ID: 8,
		})
		time.Sleep(time.Millisecond * 500)
	}
}

func clearLeds(ledController led.Controller) error {
	numLeds, err := ledController.GetNumLEDs()
	if err != nil {
		return err
	}

	ledData := make([]led.Led, numLeds)
	for i := 0; i < numLeds; i++ {
		ledData[i] = led.Led{
			ID: i,
			R:  0,
			G:  0,
			B:  0,
		}
	}
	return ledController.SetLEDs(ledData...)
}

```

## 🤹 CLI

This project provides a handy cli to perform some routine tasks which is especially useful for debugging.

To use it, make sure you cloned this repo, ideally in your project folder.
Then run:

```shell
$ go run cmd/cli.go
NAME:
   EchoGo - Control your echo easily

USAGE:
   EchoGo [global options] [command [command options]]

COMMANDS:
   boot      Boot your echo
   shutdown  Shutdown your echo
   restart   Restart your echo
   compile   Compile the server app for your echo
   run       Run the server app on your echo for debugging purposes
   deploy    Deploy the server app to your echo
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

Here are some useful tips:
- Run `go run cmd/cli.go boot` to boot your echo
- Run `go run cmd/cli.go deploy -c` to compile and deploy the server to the echo
- Run `go run cmd/cli.go run -c` to compile and run the server to the echo

If you want to deploy the prebuilt server from GitHub:
```shell
$ wget -O build/server # TODO
$ go run cmd/cli.go deploy
```

> [!NOTE]  
> If you deploy the server for the first time you need to reboot your echo after a successful deployment. After that no further reboots are required.

## Debugging

Here are some useful commands for debugging the server

- Mount the filesystem read-write: `mount -o rw,remount rootfs /`
- Pushing the android init script: `adb push echogo.rc /system/etc/init/echogo.rc`
- Setting right permissions for the init shell script: `adb shell chcon u:object_r:system_file:s0 /data/local/bin/start_server.sh`
- Disable SELinux: `adb shell setenforce 0`
- Check for server failures: `adb shell dmesg | grep echogo`
- Check if SELinux annoys again: `adb shell dmesg | grep "avc: denied"`
- Check server status: `adb shell getprop init.svc.echogo`

## Todo
- [ ] Make use of `tinymix 61 100` audio control (100 for max volume)