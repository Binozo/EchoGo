package payloads

type LedsPayload struct {
	Leds []LedPayload `json:"leds"`
}

type LedPayload struct {
	Led int `json:"led"`
	R   int `json:"r"`
	G   int `json:"g"`
	B   int `json:"b"`
}
