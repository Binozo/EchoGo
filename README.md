# EchoGo
A Go SDK for your Echo Dot **2. Gen**

- [Features](#-features)
- [Requirements](#-requirements)
- [Development](#-Development)
- - [Bare metal](#native--development)
- - [Remote host controlled](#host--development)
- [Deployment](#-deployment) (Full tutorial)
- [Go example](#go-quickstart)

For a very quick quickstart (🚀🚀🚀) take a look at [Deployment](#-deployment)

## ⚡ Features
- 🚥 Full control over LEDs
- 🎤 Access to microphone, speaker and buttons
- ⚙️ Preconfigured compiler

## 👷‍♂️ Requirements
- A fully rooted Echo Dot **2. Gen** ([Dragon863's Project](https://github.com/Dragon863/EchoCLI))
- You need your `preloader_no_hdr.bin` file for booting

## 🧑‍💻 Development
You can use this project in two ways. Building your Go application using

- [[NATIVE](#native)] the prebuilt bindings and compile it using the prebuilt compiler (Depends on cgo and heavy compiler)
- [[HOST](#host)] a client api which runs on the host machine (Can be run without cgo and compiles much faster)

Pros/Cons about the [NATIVE] way:

| Pros                          | Cons                              |
|-------------------------------|-----------------------------------|
| App runs directly on the echo | A host machine is required anyway |
| More control about the echo   | Complex building process          |

Pros/Cons about the [HOST] way:

| Pros                      | Cons                                               |
|---------------------------|----------------------------------------------------|
| Easy to develop and build | Is "remote controlling"                            |
| Much better flexibility   | Is limited by the api and can't be extended easily |

## [NATIVE] ⚙️ Development
Develop using this strategy if you want to run your app 'bare metal'.
This required using cgo and the specified compiler but grants you much more control.

In this case you may need a direct Wi-Fi connection between your echo and your Wi-Fi network.
Please make sure you blocked Amazon's OTA servers, otherwise your echo will almost brick itself after a few weeks through OTA updates (already experienced that myself).  [More info](https://github.com/Dragon863/EchoCLI#notice)

### The prebuilt compiler
Either use the prebuilt compiler or build it yourself.

#### Building it yourself
Please note: Building the compiler will take a some time.
```shell
$ git clone https://github.com/Binozo/EchoGo && cd EchoGo
$ docker build -f compiler/Dockerfile -t echogosdkcompiler:latest .
$ docker build -t ghcr.io/binozo/echogo:latest .
```

#### Pull the prebuilt image
Run the following commands:
```shell
$ docker pull ghcr.io/binozo/echogo:latest
```

### Compilation
Run the following command:
```shell
$ docker run -v "$(pwd)":/EchoGoSDK ghcr.io/binozo/echogo:latest
```

> [!NOTE]
> By default the compiler tries to compile `cmd/main.go`

If you want to compile another target in the `cmd` directory you can pass an environment variable like this:
```shell
$ docker run --env TARGET="my_go_file.go" -v "$(pwd)":/EchoGoSDK ghcr.io/binozo/echogo:latest
```

This will generate the `main` executable file in the `build` directory.

Now jump to the deployment step.


## [HOST] 🖥️️ Development
By using this strategy you build your application in a way that it runs on the host machine which then remotely controls the echo.
So there is no need for a complex compiling process nor for cgo.

In this case a pre-programmed server executable will be deployed on the echo and acts as a websocket server to control the echo.
Build the server:
```shell
$ docker run --env TARGET="server.go" -v "$(pwd)":/EchoGoSDK ghcr.io/binozo/echogo:latest
```

Now deploy it:
```shell
$ chmod +x && adb push build/main /data/local/tmp/main
$ adb shell "./data/local/tmp/main"
```
Now the server is up and running.
In order to proxy the connection to your computer execute the following command:
```shell
$ adb forward tcp:8092 tcp:8092
```
Now the websocket server will be available on [localhost:8092](localhost:8092).

Those routes are available:
- `/`: Displays "EchoGo"
- `/buttons`: Websocket subscription to get button press events
- `/microphone` Websocket subscription to read microphone data. (Will be in 9 channels, 16kHz, PCM S24_3LE wav format)
- `/speaker` Websocket connection to play audio through the speaker. (Audio format must be wav, 2 channel 48kHz PCM S16_LE format)
- `/led` Websocket connection to control the LEDs. Example payload:
```json
{
    "leds":[
        {
            "led": 0,
            "r": 255,
            "g": 0,
            "b": 0
        },
        ...
        {
            "led": 11,
            "r": 255,
            "g": 0,
            "b": 0
        }
    ]
}
```

## 🚀 Deployment
This section will help you to get an example project up and running on your echo 😎

Prerequisites:
- Your own `preloader_no_hdr.bin` file for booting
- A music file by your choice

> [!TIP]
> You need that music file in a specific format. Here is a ffmpeg snippet that converts it right:
> `ffmpeg -i music.mp3 -ar 48000 -ac 2 -volumedetect -af "volume=0.5" -f s16le music.wav` (2 Channel 48kHz S16_LE )

Pull this repository:
```bash
$ git clone https://github.com/Binozo/EchoGo && cd EchoGo 
$ git submodule init && git submodule update
```

Follow the official guide [here](https://github.com/bkerler/mtkclient/tree/f338168caba2b100eca85e5e8dfcea78cb27f1e2?tab=readme-ov-file#install) until 'Grab files' (Only install system dependencies)
Now build mtkclient for booting:
```bash
$ cd mtkclient
$ python3 -m venv .venv
$ source .venv/bin/activate
$ pip3 install -r requirements.txt
$ pip3 install .
```
Now as a last step you need to set [those](https://github.com/bkerler/mtkclient/tree/f338168caba2b100eca85e5e8dfcea78cb27f1e2?tab=readme-ov-file#install-rules) rules.

Get back to the project's root directory and run the following command to build the echo's server app:
```bash
$ docker run --env TARGET="server.go" -v "$(pwd)":/EchoGoSDK ghcr.io/binozo/echogo:latest
```
This may take a while. The host app will be placed in the `build` directory. Run `cp build/main host` to place it correctly.

Finally build the host app:
```bash
$ go build -o app cmd/host.go
```
Now also place your music file (named `music.wav`) in this root directory.

Run 🥳 🚀:
```bash
$ ./app
```
This example app will boot your echo and run some example code. Take a look here: `cmd/host.go:35` to get started 😎

## Go Quickstart
### Install the package
```shell
$ go get -u github.com/Binozo/EchoGo
```

### [NATIVE] Example code
```go
package examples

import (
  "bytes"
  "fmt"
  "github.com/Binozo/EchoGo/v2/pkg/bindings/buttons"
  "github.com/Binozo/EchoGo/v2/pkg/bindings/led"
  "github.com/Binozo/EchoGo/v2/pkg/bindings/mic"
  "time"
)

func main() {
  // Init Buttons
  if err := buttons.Init(); err != nil {
    panic(err)
  }

  dotBtn := buttons.GetDotButton()
  fmt.Println("Waiting for dot button press")
  dotBtn.WaitForClick()

  // Init LEDs
  if err := led.Init(); err != nil {
    panic(err)
  }

  // Clear all LED lights
  if err := led.Clear(); err != nil {
    panic(err)
  }

  // Prepare microphone
  if err := mic.Init(); err != nil {
    panic(err)
  }

  // Record microphone for 5 seconds
  micDevice := mic.GetDevice()
  audioStream := make(chan []byte)
  go func() {
    err := micDevice.GetAudioStream(micDevice.DeviceConfig, audioStream)
    if err != nil {
      panic(err)
    }
  }()

  recordingSeconds := 5
  fmt.Println("Recording for", recordingSeconds, "seconds...")
  start := time.Now()
  dataBuffer := new(bytes.Buffer)
  for {
    audioData := <-audioStream
    dataBuffer.Write(audioData)
    if time.Now().Sub(start).Seconds() >= float64(recordingSeconds) {
      fmt.Println("Stopping!")
      break
    }
  }
  close(audioStream)
  fmt.Println("Recorded", len(dataBuffer.Bytes()), "bytes")

  // Playing Audio through speaker
  // speaker.GetDevice().SendAudioStream(<YOUR WAV DATA>)

  // Run a fancy RGB light animation
  if err := led.Fun(); err != nil {
    panic(err)
  }
}

```

## Reverse Engineering
I made some scripts to built `strace` and `ltrace` for the Echo Dot.
Take a look into the `binaries` directory.

## Todo
- [ ] Make use of `tinymix 61 100` audio control (100 for max volume)