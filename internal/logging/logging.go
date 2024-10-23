package logging

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() { //nolint:gochecknoinits
	log.Logger = log.Level(zerolog.InfoLevel)
}
