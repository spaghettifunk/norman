package entities

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	colonSeparator            string = ":"
	colonFormatMinTokens      int    = 3
	colonFormatMaxTokens      int    = 4
	pipeSeparator             string = "|"
	pipeFormatMinTokens       int    = 1
	pipeFormatMaxTokens       int    = 3
	pipeTimezonePosition      int    = 2
	pipePatternPosition       int    = 1
	pipeTimeFormatPosition    int    = 0
	pipeSizePosition          int    = 2
	sizePosition              int    = 0
	timeUnitPosition          int    = 1
	timeFormatPosition        int    = 2
	patternPosition           int    = 3
	granularityNumberOfTokens int    = 2
)

var durationMap = map[string]time.Duration{
	"NANOSECONDS":  time.Nanosecond,
	"MICROSECONDS": time.Microsecond,
	"MILLISECONDS": time.Millisecond,
	"SECOND":       time.Second,
	"MINUTE":       time.Minute,
	"HOUR":         time.Hour,
	"DAY":          time.Hour * 24,
}

type TimeFormat string

const (
	EPOCH              TimeFormat = "EPOCH"
	TIMESTAMP          TimeFormat = "TIMESTAMP"
	SIMPLE_DATE_FORMAT TimeFormat = "SIMPLE_DATE_FORMAT"
)

type DateTimeFormatSpec struct {
	Size        int
	UnitSpec    time.Duration
	PatternSpec DateTimeFormatPatternSpec
}

func NewDateTimeFormatSpec(format string) (*DateTimeFormatSpec, error) {
	if format == "" {
		return nil, errors.New("must provide format")
	}
	// check if the first character of the string is a digit
	if format[0] >= '0' && format[0] <= '9' {
		return colonFormat(format)
	} else {
		// Pipe format
		return pipeFormat(format)
	}
}

func colonFormat(format string) (*DateTimeFormatSpec, error) {
	// Colon format
	formatTokens := strings.Split(format, colonSeparator)
	if len(formatTokens) < colonFormatMinTokens || len(formatTokens) > colonFormatMaxTokens {
		return nil, errors.New("invalid format: " + format + ", must be of format 'size:timeUnit:timeFormat(:patternWithTz)'")
	}

	timeFormat := TimeFormat(formatTokens[timeFormatPosition])
	size, err := strconv.Atoi(formatTokens[sizePosition])
	if err != nil {
		return nil, errors.New("invalid size: " + formatTokens[sizePosition] + " in format: " + format)
	}
	if size <= 0 {
		return nil, errors.New("invalid size: " + strconv.Itoa(size) + " in format: " + format + ", must be positive")
	}

	unitSpec, ok := durationMap[formatTokens[timeUnitPosition]]
	if !ok {
		return nil, errors.New("invalid time unit: " + formatTokens[timeUnitPosition] + " in format: " + format)
	}

	pattern := ""
	if len(formatTokens) > 3 {
		pattern = formatTokens[3]
	}

	patternSpec, err := NewDateTimeFormatPatternSpec(SIMPLE_DATE_FORMAT, pattern, "")
	if err != nil {
		return nil, errors.New("invalid SIMPLE_DATE_FORMAT pattern: " + pattern + " in format: " + format)
	}
	switch timeFormat {
	case EPOCH:
		return &DateTimeFormatSpec{
			Size:        size,
			UnitSpec:    unitSpec,
			PatternSpec: patternSpec,
		}, nil

	case TIMESTAMP:
		return &DateTimeFormatSpec{
			Size:        size,
			UnitSpec:    unitSpec,
			PatternSpec: patternSpec,
		}, nil

	case SIMPLE_DATE_FORMAT:
		return &DateTimeFormatSpec{
			Size:        size,
			UnitSpec:    unitSpec,
			PatternSpec: patternSpec,
		}, nil

	default:
		return nil, fmt.Errorf("invalid format: %s, must be of format 'EPOCH|<timeUnit>(|<size>)' or 'SIMPLE_DATE_FORMAT|<pattern>(|<timeZone>)' or 'TIMESTAMP'", string(timeFormat))
	}
}

func pipeFormat(format string) (*DateTimeFormatSpec, error) {
	tokens := strings.SplitN(format, string(pipeSeparator), pipeFormatMaxTokens)
	if len(tokens) < pipeFormatMinTokens || len(tokens) > pipeFormatMaxTokens {
		return nil, fmt.Errorf("invalid format: %s, must be of format 'EPOCH|<timeUnit>(|<size>)' or 'SIMPLE_DATE_FORMAT|<pattern>(|<timeZone>)' or 'TIMESTAMP'", format)
	}

	timeFormat := TimeFormat(tokens[pipeTimeFormatPosition])
	switch timeFormat {
	case EPOCH:
		size, unitSpec, err := getSizeUnit(format, tokens)
		if err != nil {
			return nil, err
		}
		return &DateTimeFormatSpec{
			Size:        size,
			UnitSpec:    unitSpec,
			PatternSpec: DateTimeFormatPatternSpec{Pattern: "EPOCH"},
		}, nil
	case TIMESTAMP:
		size, unitSpec, err := getSizeUnit(format, tokens)
		if err != nil {
			return nil, err
		}
		return &DateTimeFormatSpec{
			Size:        size,
			UnitSpec:    unitSpec,
			PatternSpec: DateTimeFormatPatternSpec{Pattern: "TIMESTAMP"},
		}, nil

	case SIMPLE_DATE_FORMAT:
		pattern := ""
		if len(tokens) > 1 {
			pattern = tokens[pipePatternPosition]
		}

		timeZone := ""
		if len(tokens) > 2 {
			timeZone = tokens[pipeTimezonePosition]
		}

		patternSpec, err := NewDateTimeFormatPatternSpec(SIMPLE_DATE_FORMAT, pattern, timeZone)
		if err != nil {
			return nil, fmt.Errorf("invalid SIMPLE_DATE_FORMAT pattern: %s, time zone: %s in format: %s", pattern, timeZone, format)
		}

		return &DateTimeFormatSpec{
			Size:        -1, // ignored
			UnitSpec:    0,  // ignored
			PatternSpec: patternSpec,
		}, nil

	default:
		return nil, fmt.Errorf("invalid format: %s, must be of format 'EPOCH|<timeUnit>(|<size>)' or 'SIMPLE_DATE_FORMAT|<pattern>(|<timeZone>)' or 'TIMESTAMP'", timeFormat)
	}
}

func getSizeUnit(format string, tokens []string) (int, time.Duration, error) {
	var us time.Duration
	size := 1
	if len(tokens) > 2 {
		var err error
		size, err = strconv.Atoi(tokens[pipeSizePosition])
		if err != nil {
			return size, time.Millisecond, fmt.Errorf("invalid size: %s in format: %s", tokens[pipeSizePosition], format)
		}
		if size <= 0 {
			return size, time.Millisecond, fmt.Errorf("invalid size: %s in format: %s, must be positive", strconv.Itoa(size), format)
		}
	}

	if len(tokens) > 1 {
		unitSpec, ok := durationMap[tokens[timeUnitPosition]]
		if !ok {
			return size, time.Millisecond, fmt.Errorf("invalid time unit: %s in format: %s", tokens[timeUnitPosition], format)
		}
		us = unitSpec
	}

	return size, us, nil
}

type DateTimeFormatPatternSpec struct {
	Pattern  string
	TimeZone string
}

func NewDateTimeFormatPatternSpec(format TimeFormat, pattern, timeZone string) (DateTimeFormatPatternSpec, error) {
	return DateTimeFormatPatternSpec{Pattern: pattern, TimeZone: timeZone}, nil
}
