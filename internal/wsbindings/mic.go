package wsbindings

import "time"

type MicControl struct {
	baseConnection
}

func GetMicControl() *MicControl {
	return &MicControl{}
}

func (b *MicControl) Connect() error {
	conn, err := connect("/microphone")
	if err != nil {
		return err
	}
	b.con = conn
	conn.SetWriteDeadline(time.Time{})
	conn.SetReadDeadline(time.Time{})

	return nil
}

func (b *MicControl) Read() ([]byte, error) {
	if b.con == nil {
		if err := b.Connect(); err != nil {
			return nil, nil
		}
	}

	_, data, err := b.con.ReadMessage()
	return data, err
}
