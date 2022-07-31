package clog

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Error is used for error level logs.
func Error(err error, msg string) {
	log.Error().Err(err).Msg(msg)
}

// Debug is used for debug level logs.
func Debug(msg string) {
	log.Debug().Msg(msg)
}

// Fatal is used for fatal level logs.
// The os.Exit(1) function is called after the log.
func Fatal(msg string) {
	log.Fatal().Msg(msg)
}

// Info is used for info level logs.
func Info(msg string) {
	log.Info().Msg(msg)
}

// Warn is used for warn level logs.
func Warn(msg string) {
	log.Warn().Msg(msg)
}

// Panic is used for panic level logs.
// The message is also sent to the panic function.
func Panic(msg string) {
	log.Panic().Msg(msg)
}

// Custom accepts a map[string] interface as a parameter.
// Useful for more detailed info level logging.
func Custom(msg map[string]interface{}) {
	log.Info().Fields(msg).Msg("")
}

// UseConsoleOutput writes the output to the console.
func UseConsoleOutput() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
