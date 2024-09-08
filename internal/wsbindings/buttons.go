package wsbindings

import (
	"github.com/Binozo/EchoGo/internal/payloads"
	"time"
)

type ButtonControl struct {
	baseConnection
}

func GetButtonControl() *ButtonControl {
	return &ButtonControl{}
}

func (b *ButtonControl) Connect() error {
	conn, err := connect("/buttons")
	if err != nil {
		return err
	}
	b.con = conn
	conn.SetWriteDeadline(time.Time{})
	conn.SetReadDeadline(time.Time{})

	return nil
}

func (b *ButtonControl) WaitForClickEvent() (payloads.ClickEvent, error) {
	if b.con == nil {
		if err := b.Connect(); err != nil {
			return payloads.ClickEvent{}, err
		}
	}

	var clickEvent payloads.ClickEvent
	if err := b.con.ReadJSON(&clickEvent); err != nil {
		return payloads.ClickEvent{}, err
	}
	return clickEvent, nil
}
