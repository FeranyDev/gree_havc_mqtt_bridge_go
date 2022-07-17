package bemfa

type Command struct {
	Power       []string       `json:"power"`
	Mode        map[string]int `json:"mode"`
	Temperature int            `json:"temperature"`
	FanSpeed    map[string]int `json:"fanSpeed"`  // nolint:tagliatelle // Fuck YOU
	SwingHor    map[string]int `json:"swingHor"`  // nolint:tagliatelle // Fuck YOU
	SwingVert   map[string]int `json:"swingVert"` // nolint:tagliatelle // Fuck YOU
}

func Commands() Command {
	return Command{
		Power: []string{
			"off",
			"on",
		},
		Mode: map[string]int{
			"auto":      1,
			"cool":      2,
			"heat":      3,
			"fan_only":  4,
			"dry":       5,
			"sleep":     6,
			"powerSave": 7,
		},
		Temperature: 0,
		FanSpeed: map[string]int{
			"auto":       0,
			"low":        1,
			"mediumLow":  2,
			"medium":     3,
			"mediumHigh": 4,
			"high":       5,
		},
		SwingHor: map[string]int{
			"default": 0,
			"full":    1,
		},
		SwingVert: map[string]int{
			"default": 0,
			"full":    1,
		},
	}
}
