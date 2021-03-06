package gree

type Command struct {
	Power           Power           `json:"power"`
	Mode            Mode            `json:"mode"`
	TemperatureUnit TemperatureUnit `json:"temperatureUnit"`
	Temperature     Temperature     `json:"temperature"`
	FanSpeed        FanSpeed        `json:"fanSpeed"`
	Air             Air             `json:"air"`
	Blow            Blow            `json:"blow"`
	Health          Health          `json:"health"`
	Sleep           Sleep           `json:"sleep"`
	Lights          Lights          `json:"lights"`
	SwingHor        SwingHor        `json:"swingHor"`
	SwingVert       SwingVert       `json:"swingVert"`
	Quiet           Quiet           `json:"quiet"`
	Turbo           Turbo           `json:"turbo"`
	PowerSave       PowerSave       `json:"powerSave"`
}
type Value struct {
	Off int `json:"off"`
	On  int `json:"on"`
}
type Power struct {
	Code  string         `json:"code"`
	Value map[string]int `json:"value"`
}

type Mode struct {
	Code  string         `json:"code"`
	Value map[string]int `json:"value"`
}

type TemperatureUnit struct {
	Code  string         `json:"code"`
	Value map[string]int `json:"value"`
}
type Temperature struct {
	Code string `json:"code"`
}

type FanSpeed struct {
	Code  string         `json:"code"`
	Value map[string]int `json:"value"`
}

type Air struct {
	Code  string         `json:"code"`
	Value map[string]int `json:"value"`
}
type Blow struct {
	Code  string         `json:"code"`
	Value map[string]int `json:"value"`
}
type Health struct {
	Code  string         `json:"code"`
	Value map[string]int `json:"value"`
}
type Sleep struct {
	Code  string         `json:"code"`
	Value map[string]int `json:"value"`
}
type Lights struct {
	Code  string         `json:"code"`
	Value map[string]int `json:"value"`
}

type SwingHor struct {
	Code  string         `json:"code"`
	Value map[string]int `json:"value"`
}

type SwingVert struct {
	Code  string         `json:"code"`
	Value map[string]int `json:"value"`
}

type Quiet struct {
	Code  string         `json:"code"`
	Value map[string]int `json:"value"`
}

type Turbo struct {
	Code  string         `json:"code"`
	Value map[string]int `json:"value"`
}
type PowerSave struct {
	Code  string         `json:"code"`
	Value map[string]int `json:"value"`
}

func Commands() Command {
	value := map[string]int{
		"off": 0,
		"on":  1,
	}
	modeValue := map[string]int{
		"auto":     0,
		"cool":     1,
		"dry":      2,
		"fan_only": 3,
		"heat":     4,
	}
	temperatureUnitValue := map[string]int{
		"celsius":    0,
		"fahrenheit": 1,
	}
	fanSpeedValue := map[string]int{
		"auto":       0,
		"low":        1,
		"mediumLow":  2,
		"medium":     3,
		"mediumHigh": 4,
		"high":       5,
	}
	airValue := map[string]int{
		"off":     0,
		"inside":  1,
		"outside": 2,
		"mode3":   3,
	}
	swingHorValue := map[string]int{
		"default":       0,
		"full":          1,
		"fixedLeft":     2,
		"fixedMidLeft":  3,
		"fixedMid":      4,
		"fixedMidRight": 5,
		"fixedRight":    6,
		"fullAlt":       7,
	}
	swingVertValue := map[string]int{
		"default":        0,
		"full":           1,
		"fixedTop":       2,
		"fixedMidTop":    3,
		"fixedMidBottom": 5,
		"fixedBottom":    6,
		"swingBottom":    7,
		"swingMidBottom": 8,
		"swingMid":       9,
		"swingMidTop":    10,
		"swingTop":       11,
	}
	quietValue := map[string]int{
		"off":   0,
		"mode1": 1,
		"mode2": 2,
		"mode3": 3,
	}
	power := Power{
		Code:  "Pow",
		Value: value,
	}
	mode := Mode{
		Code:  "Mod",
		Value: modeValue,
	}
	temperatureUnit := TemperatureUnit{
		Code:  "TemUn",
		Value: temperatureUnitValue,
	}
	temperature := Temperature{
		Code: "SetTem",
	}
	fanSpeed := FanSpeed{
		Code:  "WdSpd",
		Value: fanSpeedValue,
	}
	air := Air{
		Code:  "Air",
		Value: airValue,
	}
	blow := Blow{
		Code:  "Blo",
		Value: value,
	}
	health := Health{
		Code:  "Health",
		Value: value,
	}
	sleep := Sleep{
		Code:  "SwhSlp",
		Value: value,
	}
	lights := Lights{
		Code:  "Lig",
		Value: value,
	}
	swingHor := SwingHor{
		Code:  "SwingLfRig",
		Value: swingHorValue,
	}
	swingVert := SwingVert{
		Code:  "SwUpDn",
		Value: swingVertValue,
	}
	quiet := Quiet{
		Code:  "Quiet",
		Value: quietValue,
	}
	turbo := Turbo{
		Code:  "Tur",
		Value: value,
	}
	powerSave := PowerSave{
		Code:  "SvSt",
		Value: value,
	}

	return Command{
		Power:           power,
		Mode:            mode,
		TemperatureUnit: temperatureUnit,
		Temperature:     temperature,
		FanSpeed:        fanSpeed,
		Air:             air,
		Blow:            blow,
		Health:          health,
		Sleep:           sleep,
		Lights:          lights,
		SwingHor:        swingHor,
		SwingVert:       swingVert,
		Quiet:           quiet,
		Turbo:           turbo,
		PowerSave:       powerSave,
	}
}
