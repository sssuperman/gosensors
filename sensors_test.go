package gosensors

import (
	"fmt"
	"testing"
)

func TestPrint(t *testing.T) {
	sensors, err := NewFromSystem()
	// sensors, err := gosensors.NewFromFile("/path/to/log.txt")

	if err != nil {
		panic(err)
	}

	// Sensors implements Stringer interface,
	// so code below will print out JSON
	fmt.Println(sensors)

	// Iterate over chips
	for chip := range sensors.Chips {
		// Iterate over entries
		for key, value := range sensors.Chips[chip] {
			// If CPU or GPU, print out
      fmt.Println(key, ":", value)
		}
	}
}
