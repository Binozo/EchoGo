package led

// Clear all LEDs and turn them off
func Clear() error {
	return SetColorAll(0, 0, 0)
}

// Fun with colors
func Fun() error {
	colors := []struct {
		red   int
		green int
		blue  int
	}{
		{red: 255, green: 0, blue: 0},
		{red: 255, green: 119, blue: 0},
		{red: 255, green: 255, blue: 0},
		{red: 119, green: 255, blue: 0},
		{red: 0, green: 255, blue: 0},
		{red: 0, green: 255, blue: 119},
		{red: 0, green: 255, blue: 255},
		{red: 0, green: 119, blue: 255},
		{red: 0, green: 0, blue: 255},
		{red: 119, green: 0, blue: 255},
		{red: 255, green: 0, blue: 255},
		{red: 255, green: 0, blue: 119},
	}

	for {
		for j := 0; j < LedCount; j++ {
			for i := 0; i < LedCount; i++ {
				color := colors[i]
				err := SetColor((i+j)%LedCount, color.red, color.green, color.blue)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
