package gosensors

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
)

// Sensors struct represents lm-sensors output.
// Content field contains string output.
// Chips field contains map[string]Entries.
// Example (JSON style):
//
//	"coretemp-isa-0000": {
//		"CPU": "+60.0°C",
//		"GPU": "+48.0°C",
//	}
type Sensors struct {
	Content string             `json:"-"`
	Chips   map[string]Entries `json:"chips"`
}

// Entries representing key, value pairs for chips.
// Example (JSON style):
// "GPU": "+56.0°C"
// "CPU": "+68.0°C"
type Entries map[string]float64

func parseTemperatureValue(line string) (string, float64) {

	parts := strings.Split(line, ":")
	value := strings.TrimRight(strings.TrimLeft(parts[1], " "), " ")
	temp_value := strings.Split(value, "°C")
	float_temp, _ := strconv.ParseFloat(temp_value[0], 32)
	return parts[0], float_temp

}

func construct(content string) *Sensors {
	s := &Sensors{}
	s.Content = content
	s.Chips = map[string]Entries{}
	parse_adpater := false

	lines := strings.Split(s.Content, "\n")

	var chip string
	for _, line := range lines {
		if len(line) > 0 {
			if !(strings.Contains(line, ":") || strings.Contains(line, "crit") || parse_adpater) {
				// Parse Chip
				fmt.Print()
				chip = line
				s.Chips[chip] = Entries{}
				parse_adpater = true
			} else if parse_adpater {
				parse_adpater = false

			} else if len(chip) > 0 {
				if strings.Contains(line, ":") && strings.Contains(line, "(") {
					// Sensor with threshold
					parts := strings.Split(line, ":")
					entry := parts[0]
					value := strings.TrimRight(strings.TrimLeft(parts[1], " "), " ")
					temp_value := strings.Split(value, " C")
					float_temp, _ := strconv.ParseFloat(temp_value[0], 32)
					s.Chips[chip][entry] = float_temp
					// fmt.Println(line)
				} else if strings.Contains(line, ":") {
					// Sensor without threshold
					parts := strings.Split(line, ":")
					entry := parts[0]
					value := strings.TrimRight(strings.TrimLeft(parts[1], " "), " ")
					temp_value := strings.Split(value, "°C")
					float_temp, _ := strconv.ParseFloat(temp_value[0], 32)
					s.Chips[chip][entry] = float_temp
				} else {
				}
			}
		}
	}

	return s
}

// NewFromSystem executes "sensors" system command and returns constructed Sensors struct.
// A successful call returns err == nil.
func NewFromSystem() (*Sensors, error) {
	out, err := exec.Command("sensors").Output()
	if err != nil {
		return &Sensors{}, errors.New("lm-sensors missing")
	}

	s := construct(string(out))

	return s, nil
}

// NewFromFile reads content from log file and returns constructed Sensors struct.
// A successful call returns err == nil.
func NewFromFile(path string) (*Sensors, error) {
	out, err := ioutil.ReadFile(path)
	if err != nil {
		return &Sensors{}, err
	}

	s := construct(string(out))
	return s, nil
}

// JSON returns JSON of Sensors.
func (s *Sensors) JSON() string {
	out, _ := json.Marshal(s)

	return string(out)
}

func (s *Sensors) String() string {
	return s.JSON()
}
