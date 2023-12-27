package logger

import "log/slog"

const (
	logLevel = "LOG_LEVEL"

	logLevelTrace = "TRACE"
	logLevelDebug = "DEBUG"
	logLevelInfo  = "INFO"
	logLevelWarn  = "WARN"
	logLevelError = "ERROR"
	logLevelFatal = "FATAL"
	logLevelPanic = "PANIC"

	levelTrace = slog.Level(slog.LevelDebug - (4 << 0))
	levelDebug = slog.LevelDebug
	levelInfo  = slog.LevelInfo
	levelWarn  = slog.LevelWarn
	levelError = slog.LevelError
	levelFatal = slog.Level(slog.LevelError + (4 << 0))
	levelPanic = slog.Level(slog.LevelError + (4 << 1))
)

var levels = map[slog.Leveler]string{
	levelTrace: logLevelTrace,
	levelDebug: logLevelDebug,
	levelInfo:  logLevelInfo,
	levelWarn:  logLevelWarn,
	levelError: logLevelError,
	levelFatal: logLevelFatal,
	levelPanic: logLevelPanic,
}
