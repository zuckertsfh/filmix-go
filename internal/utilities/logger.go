package utilities

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

func InitLogger(env string) {
	if env == "prod" {
		// JSON structured logs
		Logger = log.Output(os.Stdout).With().Timestamp().Logger()
		return
	}

	// Winston-style for dev
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02 Jan 2006 15:04:05",
		FormatLevel: func(i any) string {
			if lvl, ok := i.(string); ok {
				return "[" + lvl + "]:"
			}
			return "[info]:"
		},
		FormatMessage: func(i any) string {
			if msg, ok := i.(string); ok {
				return msg
			}
			return ""
		},
		FormatFieldName:  func(i any) string { return "" },
		FormatFieldValue: func(i any) string { return "" },
	}

	Logger = log.Output(output).With().Timestamp().Logger()
}
