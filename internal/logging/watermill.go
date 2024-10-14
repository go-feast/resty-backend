package logging

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type WatermillLoggerAdapter struct {
	Log zerolog.Logger
}

func NewWatermillLogger() *WatermillLoggerAdapter {
	return &WatermillLoggerAdapter{
		Log: log.Logger,
	}
}

func (w *WatermillLoggerAdapter) Error(msg string, err error, fields watermill.LogFields) {
	w.Log.Err(err).Any("fields", fields).Msg(msg)
}

func (w *WatermillLoggerAdapter) Info(msg string, fields watermill.LogFields) {
	w.Log.Info().Any("fields", fields).Msg(msg)
}

func (w *WatermillLoggerAdapter) Debug(msg string, fields watermill.LogFields) {
	w.Log.Debug().Any("fields", fields).Msg(msg)
}

func (w *WatermillLoggerAdapter) Trace(msg string, fields watermill.LogFields) {
	w.Log.Trace().Any("fields", fields).Msg(msg)
}

func (w *WatermillLoggerAdapter) With(fields watermill.LogFields) watermill.LoggerAdapter {
	return &WatermillLoggerAdapter{w.Log.With().Any("fields", fields).Logger()}
}
