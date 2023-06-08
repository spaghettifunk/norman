package entities

import (
	"testing"
	"time"
)

func TestNewDateTimeFormatSpec(t *testing.T) {
	// Valid EPOCH format
	format := "EPOCH|MILLISECONDS"
	expectedSpec := &DateTimeFormatSpec{
		Size:        1,
		UnitSpec:    time.Millisecond,
		PatternSpec: DateTimeFormatPatternSpec{Pattern: "EPOCH"},
	}
	spec, err := NewDateTimeFormatSpec(format)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !compareDateTimeFormatSpec(spec, expectedSpec) {
		t.Errorf("Unexpected DateTimeFormatSpec: %+v, expected: %+v", spec, expectedSpec)
	}

	// Valid SIMPLE_DATE_FORMAT format
	format = "SIMPLE_DATE_FORMAT|yyyy-MM-dd|America/Los_Angeles"
	expectedSpec = &DateTimeFormatSpec{
		Size:        1,
		UnitSpec:    time.Millisecond,
		PatternSpec: DateTimeFormatPatternSpec{Pattern: "yyyy-MM-dd", TimeZone: "America/Los_Angeles"},
	}
	spec, err = NewDateTimeFormatSpec(format)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !compareDateTimeFormatSpec(spec, expectedSpec) {
		t.Errorf("Unexpected DateTimeFormatSpec: %+v, expected: %+v", spec, expectedSpec)
	}

	// Valid TIMESTAMP format
	format = "TIMESTAMP"
	expectedSpec = &DateTimeFormatSpec{
		Size:        1,
		UnitSpec:    time.Millisecond,
		PatternSpec: DateTimeFormatPatternSpec{Pattern: "TIMESTAMP"},
	}
	spec, err = NewDateTimeFormatSpec(format)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !compareDateTimeFormatSpec(spec, expectedSpec) {
		t.Errorf("Unexpected DateTimeFormatSpec: %+v, expected: %+v", spec, expectedSpec)
	}

	// Invalid format
	format = "INVALID_FORMAT"
	_, err = NewDateTimeFormatSpec(format)
	if err == nil {
		t.Error("Expected an error for invalid format, but got nil")
	}
	expectedError := "invalid format: INVALID_FORMAT, must be of format 'EPOCH|<timeUnit>(|<size>)' or 'SIMPLE_DATE_FORMAT|<pattern>(|<timeZone>)' or 'TIMESTAMP'"
	if err.Error() != expectedError {
		t.Errorf("Unexpected error message: %v, expected: %v", err.Error(), expectedError)
	}

	format = "1:MILLISECONDS:EPOCH"
	_, err = NewDateTimeFormatSpec(format)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

// Helper function to compare DateTimeFormatSpec objects
func compareDateTimeFormatSpec(spec1, spec2 *DateTimeFormatSpec) bool {
	if spec1 == nil && spec2 == nil {
		return true
	}
	if spec1 == nil || spec2 == nil {
		return false
	}
	return spec1.Size == spec2.Size && spec1.UnitSpec == spec2.UnitSpec &&
		compareDateTimeFormatPatternSpec(spec1.PatternSpec, spec2.PatternSpec)
}

// Helper function to compare DateTimeFormatPatternSpec objects
func compareDateTimeFormatPatternSpec(patternSpec1, patternSpec2 DateTimeFormatPatternSpec) bool {
	return patternSpec1.Pattern == patternSpec2.Pattern && patternSpec1.TimeZone == patternSpec2.TimeZone
}
