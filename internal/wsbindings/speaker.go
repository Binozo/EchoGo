package wsbindings

import (
	"github.com/gorilla/websocket"
	"time"
)

type SpeakerControl struct {
	baseConnection
}

func GetSpeakerControl() *SpeakerControl {
	return &SpeakerControl{}
}

func (b *SpeakerControl) Connect() error {
	conn, err := connect("/speaker")
	if err != nil {
		return err
	}
	b.con = conn
	conn.SetWriteDeadline(time.Time{})
	conn.SetReadDeadline(time.Time{})

	return nil
}

func (b *SpeakerControl) Write(audiodata []byte) error {
	if b.con == nil {
		if err := b.Connect(); err != nil {
			return nil
		}
	}

	return b.con.WriteMessage(websocket.BinaryMessage, audiodata)
}
