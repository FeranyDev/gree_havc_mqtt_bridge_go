package bemfa

type Command struct {
	Power       []string       `json:"power"`
	Mode        map[string]int `json:"mode"`
	Temperature int            `json:"temperature"`
	FanSpeed    map[string]int `json:"fanSpeed"`
	SwingHor    []int          `json:"swingHor"`
	SwingVert   []int          `json:"swingVert"`
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
		SwingHor: []int{
			0,
			1,
		},
		SwingVert: []int{
			0,
			1,
		},
	}
}
