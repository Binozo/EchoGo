package buttons

type Button struct {
	internalName string
	Type         string
}

type ClickType uint16

const (
	DotClick        ClickType = 138
	VolumeUpClick   ClickType = 115
	VolumeDownClick ClickType = 114
)

func GetDotButton() *Button {
	return &Button{
		dotButton,
		"Dot",
	}
}

func GetVolumeButton() *Button {
	return &Button{
		volumeButtons,
		"Volume",
	}
}
