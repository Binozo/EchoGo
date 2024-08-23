package buttons

type ClickEvent struct {
	Button    *Button
	Down      bool
	ClickType ClickType
}
