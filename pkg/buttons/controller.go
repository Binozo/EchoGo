package buttons

type Controller interface {
	Init() error
	SubscribeToButton(callback ButtonClickCallback) (*EventSubscription, error)
	GetDotButton() Button
	GetVolumeButton() Button
}
