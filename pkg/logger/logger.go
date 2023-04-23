package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	configuration "github.com/spaghettifunk/norman/internal/common"
)

var levels = map[string]int8{
	"panic": 5,
	"fatal": 4,
	"error": 3,
	"warn":  2,
	"info":  1,
	"debug": 0,
	"trace": -1,
}

func InitLogger(lc configuration.Configuration) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// in case the user put some random capitals letter
	ll := strings.ToLower(lc.Logger.Level)
	l := zerolog.Level(levels[ll])
	zerolog.SetGlobalLevel(l)

	if lc.Logger.Pretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}
