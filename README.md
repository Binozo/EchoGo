# EchoGoSDK
A Go SDK for your Echo Dot **2. Gen**

## ⚡ Features
- 🚥 Full control over LEDs
- 🎤 Access to microphone + speaker
- ⚙️ Preconfigured compiler

## Setup
You need a fully rooted Echo Dot **2. Gen**. To perform that root you need to follow the steps in [Dragon863's Project](https://github.com/Dragon863/EchoCLI).
Make sure you preserve your `preloader_no_hdr.bin` file. You need it for booting.

> [!TIP]
> Make sure you blocked Amazon's OTA servers. [More info](https://github.com/Dragon863/EchoCLI#notice)

Now you need a host system for your Echo because it can't live on its own.
In this example, we are using a Raspberry Pi Zero 2 W.

Compile your Go executable. Take a look at [Compilation](README.md#Compilation) for more info.
Make sure your executable is named `echogo`.

Now copy the following 3 files into the home directory of your Raspberry Pi (`/home/pi/`):
- The `preloader_no_hdr.bin` file for booting
- The `boot.sh` file from this repository for installation ([Link](https://github.com/Binozo/EchoGoSDK/raw/master/boot.sh))
- Your `echogo` executable

SSH into your Pi's home directory and execute:
```shell
pi@raspberrypi:~ $ chmod +x boot.sh
pi@raspberrypi:~ $ sudo ./boot.sh
```

This will setup everything for you. The `boot.sh` script moves itself and the 2 other files you copied into `/opt/echogo/` and creates a systemd service.
The systemd service runs automatically and boots your Echo Dot including the Go executable you created.

> [!TIP]
> If you want to update your Go executable, just replace the `/opt/echogo/echogo` binary on your Pi.
 
> [!TIP]
> You can always run `journalctl -u echogo --no-pager -f` to view logs.


## Compilation
If you are compiling your Go project without any C bindings you can just run:
```shell
$ env GOOS=linux GOARCH=arm go build [...]
```


If your Go project is using C bindings (the EchoGoSDK does) you need to compile your project using the Android NDK toolchain.

You have two options:
- Build the compiler Docker image yourself
- Pull the prebuilt Docker image

### Build the compiler Docker image yourself
Run the following command.

##### Please note: Building the compiler will take a _long_ time.
```shell
$ docker build -f compiler/Dockerfile -t echogosdkcompiler:latest .
$ docker build -t ghcr.io/binozo/echogo:latest .
```

### Pull the prebuilt image
Run the following commands:
```shell
$ docker pull ghcr.io/binozo/echogo:latest
```

### Finally compiling your project
Run the following commands:
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

## Quickstart
### Install the package
```shell
$ go get -u github.com/Binozo/EchoGoSDK
```

### Example code
```go
package examples

import (
	"github.com/Binozo/EchoGoSDK/pkg/led"
	"github.com/Binozo/EchoGoSDK/pkg/mic"
	"bytes"
	"fmt"
	"time"
)

func main() {
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