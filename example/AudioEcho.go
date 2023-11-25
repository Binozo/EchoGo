package main

import (
	"bytes"
	"fmt"
	"github.com/Binozo/GoTinyAlsa/pkg/pcm"
	"github.com/Binozo/GoTinyAlsa/pkg/tinyalsa"
	"log"
	"os/exec"
	"time"
)

func main() {
	// Echo Dot's default output card is 0 and device is 3
	deviceNr := 24
	fmt.Println("Starting on", deviceNr)
	inputDevice := tinyalsa.NewDevice(0, deviceNr, pcm.Config{
		Channels:    9,
		SampleRate:  16000,
		PeriodSize:  512,
		PeriodCount: 5,
		Format:      tinyalsa.PCM_FORMAT_S24_3LE,
	})
	audioStream := make(chan []byte)
	go func() {
		err := inputDevice.GetAudioStream(inputDevice.DeviceConfig, audioStream)
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

	fmt.Println("Sending:")
	cmd := exec.Command("killall", "mixer")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	cmd = exec.Command("killall", "BTSinkPlayer")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	cmd = exec.Command("ledctrl", "-c")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	cmd = exec.Command("ledctrl", "-s", "zzz_rainbow")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Millisecond * 250)

	sendDevice := tinyalsa.NewDevice(0, 25, pcm.Config{
		Channels:         2,
		SampleRate:       8000,
		PeriodSize:       512,
		PeriodCount:      4,
		Format:           tinyalsa.PCM_FORMAT_S16_LE,
		SilenceThreshold: 512 * 4,
		StopThreshold:    512 * 4,
		StartThreshold:   512,
	})
	err := sendDevice.SendAudioStream(dataBuffer.Bytes())
	if err != nil {
		fmt.Println(err)
	}
}
