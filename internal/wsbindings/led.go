package wsbindings

import (
	"github.com/Binozo/EchoGo/v2/internal/payloads"
	"github.com/Binozo/EchoGo/v2/pkg/constants"
	"time"
)

type LedControl struct {
	baseConnection
}

func GetLedControl() *LedControl {
	return &LedControl{}
}

func (l *LedControl) Connect() error {
	conn, err := connect("/led")
	if err != nil {
		return err
	}
	l.con = conn
	conn.SetWriteDeadline(time.Time{})
	conn.SetReadDeadline(time.Time{})

	return nil
}

func (l *LedControl) SetColor(r, g, b int) error {
	if l.con == nil {
		if err := l.Connect(); err != nil {
			return err
		}
	}

	payload := payloads.LedsPayload{Leds: []payloads.LedPayload{}}

	for i := 0; i <= constants.LedCount; i++ {
		payload.Leds = append(payload.Leds, payloads.LedPayload{
			Led: i,
			R:   r,
			G:   g,
			B:   b,
		})
	}
	l.con.SetWriteDeadline(time.Now().Add(time.Second))
	defer l.con.SetWriteDeadline(time.Time{})
	return l.con.WriteJSON(payload)
}
