package logger

import "log/slog"

const (
	levelTrace slog.Level = 4 * (iota - 2)
	levelDebug
	levelInfo
	levelWarn
	levelError
	levelFatal
	levelPanic
)

const (
	levelTraceName = "TRACE"
	levelDebugName = "DEBUG"
	levelInfoName  = "INFO"
	levelWarnName  = "WARN"
	levelErrorName = "ERROR"
	levelFatalName = "FATAL"
	levelPanicName = "PANIC"
)

var levels = map[slog.Level]string{
	levelTrace: levelTraceName,
	levelDebug: levelDebugName,
	levelInfo:  levelInfoName,
	levelWarn:  levelWarnName,
	levelError: levelErrorName,
	levelFatal: levelFatalName,
	levelPanic: levelPanicName,
}
