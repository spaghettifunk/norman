package entities

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type GranularitySpec struct {
	Size     int64
	UnitSpec time.Duration
}

func NewGranularitySpec(format string) (*GranularitySpec, error) {
	granularityTokens := strings.Split(format, colonSeparator)
	if len(granularityTokens) != granularityNumberOfTokens {
		return nil, fmt.Errorf("wrong amount of tokens in string")
	}

	size, err := strconv.ParseInt(granularityTokens[sizePosition], 10, 64)
	if err != nil {
		return nil, err
	}

	timeUnit, ok := durationMap[granularityTokens[timeUnitPosition]]
	if !ok {
		return nil, fmt.Errorf("wrong duration time unit. %s does not exist", granularityTokens[timeUnitPosition])
	}

	return &GranularitySpec{
		Size:     size,
		UnitSpec: timeUnit,
	}, nil
}
