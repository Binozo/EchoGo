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

func (c *ClickType) String() string {
	switch *c {
	case DotClick:
		return "dot"
	case VolumeUpClick:
		return "volume_up"
	case VolumeDownClick:
		return "volume_down"
	default:
		return "unknown"
	}
}

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
