package app

type AutoGenerated struct {
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
type ModeValue struct {
	Auto    int `json:"auto"`
	Cool    int `json:"cool"`
	Dry     int `json:"dry"`
	FanOnly int `json:"fan_only"`
	Heat    int `json:"heat"`
}
type Mode struct {
	Code  string         `json:"code"`
	Value map[string]int `json:"value"`
}
type TemperatureUnitValue struct {
	Celsius    int `json:"celsius"`
	Fahrenheit int `json:"fahrenheit"`
}
type TemperatureUnit struct {
	Code  string         `json:"code"`
	Value map[string]int `json:"value"`
}
type Temperature struct {
	Code string `json:"code"`
}
type FanSpeedValue struct {
	Auto       int `json:"auto"`
	Low        int `json:"low"`
	MediumLow  int `json:"mediumLow"`
	Medium     int `json:"medium"`
	MediumHigh int `json:"mediumHigh"`
	High       int `json:"high"`
}
type FanSpeed struct {
	Code  string         `json:"code"`
	Value map[string]int `json:"value"`
}
type AirValue struct {
	Off     int `json:"off"`
	Inside  int `json:"inside"`
	Outside int `json:"outside"`
	Mode3   int `json:"mode3"`
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
type SwingHorValue struct {
	Default       int `json:"default"`
	Full          int `json:"full"`
	FixedLeft     int `json:"fixedLeft"`
	FixedMidLeft  int `json:"fixedMidLeft"`
	FixedMid      int `json:"fixedMid"`
	FixedMidRight int `json:"fixedMidRight"`
	FixedRight    int `json:"fixedRight"`
	FullAlt       int `json:"fullAlt"`
}
type SwingHor struct {
	Code  string         `json:"code"`
	Value map[string]int `json:"value"`
}
type SwingVertValue struct {
	Default        int `json:"default"`
	Full           int `json:"full"`
	FixedTop       int `json:"fixedTop"`
	FixedMidTop    int `json:"fixedMidTop"`
	FixedMid       int `json:"fixedMid"`
	FixedMidBottom int `json:"fixedMidBottom"`
	FixedBottom    int `json:"fixedBottom"`
	SwingBottom    int `json:"swingBottom"`
	SwingMidBottom int `json:"swingMidBottom"`
	SwingMid       int `json:"swingMid"`
	SwingMidTop    int `json:"swingMidTop"`
	SwingTop       int `json:"swingTop"`
}
type SwingVert struct {
	Code  string         `json:"code"`
	Value map[string]int `json:"value"`
}
type QuietValue struct {
	Off   int `json:"off"`
	Mode1 int `json:"mode1"`
	Mode2 int `json:"mode2"`
	Mode3 int `json:"mode3"`
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

func Commands() AutoGenerated {
	value := map[string]int{
		"Off":            0,
		"On":             1,
		"Auto":           0,
		"Cool":           1,
		"Dry":            2,
		"FanOnly":        3,
		"Heat":           4,
		"Celsius":        1,
		"Fahrenheit":     2,
		"Low":            1,
		"MediumLow":      2,
		"Medium":         3,
		"MediumHigh":     4,
		"High":           5,
		"Inside":         1,
		"Outside":        2,
		"Default":        0,
		"Full":           1,
		"FixedLeft":      2,
		"FixedMidLeft":   3,
		"FixedMid":       4,
		"FixedMidRight":  5,
		"FixedRight":     6,
		"FullAlt":        7,
		"FixedTop":       2,
		"FixedMidTop":    3,
		"FixedMidBottom": 5,
		"FixedBottom":    6,
		"SwingBottom":    7,
		"SwingMidBottom": 8,
		"SwingMid":       9,
		"SwingMidTop":    10,
		"SwingTop":       11,
		"Mode1":          1,
		"Mode2":          2,
		"Mode3":          3,
	}
	power := Power{
		Code:  "Pow",
		Value: value,
	}
	mode := Mode{
		Code:  "Mod",
		Value: value,
	}
	temperatureUnit := TemperatureUnit{
		Code:  "TemUn",
		Value: value,
	}
	temperature := Temperature{
		Code: "SetTem",
	}
	fanSpeed := FanSpeed{
		Code:  "WdSpd",
		Value: value,
	}
	air := Air{
		Code:  "Air",
		Value: value,
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
		Value: value,
	}
	swingVert := SwingVert{
		Code:  "SwUpDn",
		Value: value,
	}
	quiet := Quiet{
		Code:  "Quiet",
		Value: value,
	}
	turbo := Turbo{
		Code:  "Tur",
		Value: value,
	}
	powerSave := PowerSave{
		Code:  "SvSt",
		Value: value,
	}
	return AutoGenerated{
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
