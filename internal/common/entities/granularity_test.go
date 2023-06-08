package entities

import (
	"fmt"
	"testing"
)

func TestNewGranularitySpec(t *testing.T) {
	var formats = []string{
		"5:SECOND",
		"350:MILLISECONDS",
		"10:MINUTE",
		"1:HOUR",
		"12:SECOND",
		"4:MICROSECONDS",
		"3:NANOSECONDS",
	}
	for _, format := range formats {
		testname := fmt.Sprintf("format: %s", format)
		t.Run(testname, func(t *testing.T) {
			_, err := NewGranularitySpec(format)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
