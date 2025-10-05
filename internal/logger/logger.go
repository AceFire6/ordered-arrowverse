package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/AceFire6/ordered-arrowverse/internal/build"
)

type Settings struct {
	BuildInfo      *build.Info
	ServiceName    string
	Environment    string
	IsDevelopment  bool
	GlobalLogLevel zerolog.Level
}

func SetupLogger(settings Settings) *zerolog.Logger {
	log.Logger = zerolog.New(os.Stdout).With().
		Timestamp().
		Caller().
		Str("role", settings.ServiceName).
		Str("environment", settings.Environment).
		Logger()

	if settings.BuildInfo != nil {
		log.Logger = log.Logger.With().Interface("build_info", settings.BuildInfo).Logger()
	}

	// This lets zerolog extract the stack trace from errors
	// requires an error log and the .Stack() call for that log
	// eg. log.Err(err).Stack().Msg("on no!")
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack //nolint:reassign // this is recommended by zerolog here: https://github.com/rs/zerolog?tab=readme-ov-file#error-logging-with-stacktrace

	if settings.IsDevelopment {
		// In development mode we output colorized human-readable output instead of JSON
		log.Logger = log.Output(zerolog.NewConsoleWriter())
	}

	// Make sure all logs are using the same log level
	zerolog.SetGlobalLevel(settings.GlobalLogLevel)

	log.Info().Msg("Logger setup ✨")

	return &log.Logger
}
