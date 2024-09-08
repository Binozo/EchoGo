package payloads

import "fmt"

type ClickEvent struct {
	Button string `json:"button"`
	Down   bool   `json:"down"`
	Type   string `json:"type"`
}

func (c *ClickEvent) String() string {
	return fmt.Sprintf("ClickEvent [%s] -> %s (Down: %v)", c.Button, c.Type, c.Down)
}
