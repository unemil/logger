package logger

import (
	"log/slog"
	"os"
)

// Level represents of the levels of logging
type Level string

const (
	logLevel = "LOG_LEVEL"

	// LevelTrace represents the level that is required for the most detailed logging.
	LevelTrace Level = "TRACE"
	// LevelDebug represents the level that is used to get detailed debugging information.
	LevelDebug Level = "DEBUG"
	// LevelInfo represents the level that provides general information messages (default).
	LevelInfo Level = "INFO"
	// LevelWarn represents the level that indicates potential problems that should be corrected.
	LevelWarn Level = "WARN"
	// LevelError represents the level that indicates unexpected but recoverable issues.
	LevelError Level = "ERROR"
	// LevelFatal represents the level that indicates critical errors that cause the program to terminate.
	LevelFatal Level = "FATAL"
	// LevelPanic represents the level that indicates panics that cause the program to immediate termination.
	LevelPanic Level = "PANIC"

	levelTrace = slog.Level(slog.LevelDebug - (4 << 0))
	levelDebug = slog.LevelDebug
	levelInfo  = slog.LevelInfo
	levelWarn  = slog.LevelWarn
	levelError = slog.LevelError
	levelFatal = slog.Level(slog.LevelError + (4 << 0))
	levelPanic = slog.Level(slog.LevelError + (4 << 1))
)

var levels = map[slog.Level]Level{
	levelTrace: LevelTrace,
	levelDebug: LevelDebug,
	levelInfo:  LevelInfo,
	levelWarn:  LevelWarn,
	levelError: LevelError,
	levelFatal: LevelFatal,
	levelPanic: LevelPanic,
}

func envLogLevel() slog.Level {
	switch env := Level(os.Getenv(logLevel)); env {
	case LevelTrace:
		return levelTrace
	case LevelDebug:
		return levelDebug
	}

	return levelInfo
}
