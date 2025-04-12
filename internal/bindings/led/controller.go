package led

type Controller interface {
	Init() error
	GetNumLEDs() (int, error)
	SetLEDs(led ...Led) error
}

func NewDefaultController() (Controller, error) {
	controller := &I2CController{}

	if err := controller.Init(); err != nil {
		return nil, err
	}

	return controller, nil
}
