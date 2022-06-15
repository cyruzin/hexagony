package clog

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var Clog CLog

type CLog interface {
	Error(err error, msg string)
	Debug(msg string)
	Fatal(msg string)
	Info(msg string)
	Warn(msg string)
	Panic(msg string)
}

func Error(err error, msg string) {
	log.Error().Err(err).Stack().Msg(msg)
}

func Debug(msg string) {
	log.Debug().Msg(msg)
}

func Fatal(msg string) {
	log.Fatal().Stack().Msg(msg)
}

func Info(msg string) {
	log.Info().Msg(msg)
}

func Warn(msg string) {
	log.Warn().Msg(msg)
}

func Panic(msg string) {
	log.Panic().Stack().Msg(msg)
}

func UseStack() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

func UseConsoleOutput() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func UseTimeFormatUnix() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}
