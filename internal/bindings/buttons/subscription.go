package buttons

import (
	"context"
)

type ButtonClickCallback func(event ButtonClickEvent)

type ButtonClickEvent struct {
	Button    Button    `json:"button"`
	ClickType ClickType `json:"clickType"`
	Down      bool      `json:"down"`
}

type EventSubscription struct {
	btn    Button
	cancel context.CancelFunc
}

func (e *EventSubscription) Cancel() {
	e.cancel()
}
