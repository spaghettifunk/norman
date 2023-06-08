package entities

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	separator                 string = ":"
	sizePosition              int    = 0
	timeUnitPosition          int    = 1
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

const (
	colonSeparator       = ':'
	colonFormatMinTokens = 3
	colonFormatMaxTokens = 4

	pipeSeparator       = '|'
	pipeFormatMinTokens = 1
	pipeFormatMaxTokens = 3
)

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

	if format[0] >= '0' && format[0] <= '9' {
		// Colon format
		formatTokens := strings.SplitN(format, string(colonSeparator), colonFormatMaxTokens)
		if len(formatTokens) < colonFormatMinTokens || len(formatTokens) > colonFormatMaxTokens {
			return nil, errors.New("invalid format: " + format + ", must be of format 'size:timeUnit:timeFormat(:patternWithTz)'")
		}

		timeFormat := TimeFormat(formatTokens[2])
		size, err := strconv.Atoi(formatTokens[0])
		if err != nil {
			return nil, errors.New("invalid size: " + formatTokens[0] + " in format: " + format)
		}
		if size <= 0 {
			return nil, errors.New("invalid size: " + strconv.Itoa(size) + " in format: " + format + ", must be positive")
		}

		unitSpec, ok := durationMap[formatTokens[1]]
		if !ok {
			return nil, errors.New("invalid time unit: " + formatTokens[1] + " in format: " + format)
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

	} else {
		// Pipe format
		tokens := strings.SplitN(format, string(pipeSeparator), pipeFormatMaxTokens)
		if len(tokens) < pipeFormatMinTokens || len(tokens) > pipeFormatMaxTokens {
			return nil, errors.New("invalid format: " + format + ", must be of format 'EPOCH|<timeUnit>(|<size>)' or 'SIMPLE_DATE_FORMAT|<pattern>(|<timeZone>)' or 'TIMESTAMP'")
		}

		timeFormat := TimeFormat(tokens[0])
		switch timeFormat {
		case EPOCH:
			size := 1
			if len(tokens) > 2 {
				var err error
				size, err = strconv.Atoi(tokens[2])
				if err != nil {
					return nil, errors.New("invalid size: " + tokens[2] + " in format: " + format)
				}
				if size <= 0 {
					return nil, errors.New("invalid size: " + strconv.Itoa(size) + " in format: " + format + ", must be positive")
				}
			}

			if len(tokens) > 1 {
				unitSpec, ok := durationMap[tokens[1]]
				if !ok {
					return nil, errors.New("invalid time unit: " + tokens[1] + " in format: " + format)
				}
			}

			return &DateTimeFormatSpec{
				Size:        size,
				UnitSpec:    time.Millisecond,
				PatternSpec: DateTimeFormatPatternSpec{Pattern: "EPOCH"},
			}, nil

		case TIMESTAMP:
			return &DateTimeFormatSpec{
				Size:        1,
				UnitSpec:    time.Millisecond,
				PatternSpec: DateTimeFormatPatternSpec{Pattern: "TIMESTAMP"},
			}, nil

		case SIMPLE_DATE_FORMAT:
			pattern := ""
			if len(tokens) > 1 {
				pattern = tokens[1]
			}

			timeZone := ""
			if len(tokens) > 2 {
				timeZone = tokens[2]
			}

			patternSpec, err := NewDateTimeFormatPatternSpec(SIMPLE_DATE_FORMAT, pattern, timeZone)
			if err != nil {
				return nil, errors.New("invalid SIMPLE_DATE_FORMAT pattern: " + pattern + ", time zone: " + timeZone + " in format: " + format)
			}

			return &DateTimeFormatSpec{
				Size:        1,
				UnitSpec:    DateTimeFormatUnitSpec{Unit: "MILLISECONDS"},
				PatternSpec: patternSpec,
			}, nil

		default:
			return nil, fmt.Errorf("invalid format: %s, must be of format 'EPOCH|<timeUnit>(|<size>)' or 'SIMPLE_DATE_FORMAT|<pattern>(|<timeZone>)' or 'TIMESTAMP'", string(timeFormat))
		}
	}
}

type DateTimeFormatPatternSpec struct {
	Pattern  string
	TimeZone string
}

func NewDateTimeFormatPatternSpec(format TimeFormat, pattern, timeZone string) (DateTimeFormatPatternSpec, error) {
	// Add your implementation to handle different patterns and time zones as required
	return DateTimeFormatPatternSpec{Pattern: pattern, TimeZone: timeZone}, nil
}
