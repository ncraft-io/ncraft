package logs

import (
	"github.com/mojo-lang/mojo/go/pkg/logs"
)

type Level = logs.Level

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = logs.DebugLevel
	// InfoLevel is the default logging priority.
	InfoLevel = logs.InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = logs.WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = logs.ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel = logs.DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel = logs.PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = logs.FatalLevel
)

func ParseLogLevel(val string) Level {
	return logs.ParseLogLevel(val)
}
