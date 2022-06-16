package clog

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type CLog interface {
	Error(err error, msg string)
	Debug(msg string)
	Fatal(msg string)
	Info(msg string)
	Warn(msg string)
	Panic(msg string)
	Custom(msg map[string]string)
}

func Error(err error, msg string) {
	log.Error().Err(err).Msg(msg)
}

func Debug(msg string) {
	log.Debug().Msg(msg)
}

func Fatal(msg string) {
	log.Fatal().Msg(msg)
}

func Info(msg string) {
	log.Info().Msg(msg)
}

func Warn(msg string) {
	log.Warn().Msg(msg)
}

func Panic(msg string) {
	log.Panic().Msg(msg)
}

func Custom(msg map[string]interface{}) {
	log.Info().Fields(msg).Msg("")
}

func UseConsoleOutput() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
