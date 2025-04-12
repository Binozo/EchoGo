package buttons

type Controller interface {
	Init() error
	SubscribeToButton(button Button, callback ButtonClickCallback) (*EventSubscription, error)
	GetDotButton() Button
	GetVolumeButton() Button
}

func NewButtonController() (Controller, error) {
	controller := &EvDevController{}
	if err := controller.Init(); err != nil {
		return nil, err
	}
	return controller, nil
}
