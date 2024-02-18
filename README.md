# [![doc-img]][doc]

Structured logging based on [slog][slog-doc].

## Install

```sh
go get github.com/unemil/logger
```

**Compatibility:** go >= 1.21

## Usage

Set `LOG_LEVEL` using one of the [LEVEL](#levels) values (default: **INFO**).

```sh
export LOG_LEVEL=LEVEL
```

## Levels

- TRACE
- DEBUG
- INFO
- WARN
- ERROR
- FATAL
- PANIC

## Example

```go
package main

import (
	"context"
	"errors"

	"github.com/unemil/logger"
)

var (
	ctx    = context.Background()
	err    = errors.New("test error")
	fields = logger.Fields{
		Field("levels", []logger.Level{
			logger.LevelTrace,
			logger.LevelDebug,
			logger.LevelInfo,
			logger.LevelWarn,
			logger.LevelError,
			logger.LevelFatal,
			logger.LevelPanic,
		}),
		Field("username", nil),
	}
)

func Field(key logger.FieldKey, value logger.FieldValue) logger.Field {
	return logger.Field{Key: key, Value: value}
}

func main() {
	// LOG_LEVEL=DEBUG

	logger.Info(ctx, "test")

	// {"time":"2024-02-18T16:42:55+03:00","level":"INFO","source":"test/main.go:34","msg":"test"}

	ctx = logger.With(ctx, "username", "unemil")
	logger.Debugf(ctx, "test", fields...)

	// {"time":"2024-02-18T16:42:55+03:00","level":"DEBUG","source":"test/main.go:39","msg":"test","username":null,"levels":["TRACE","DEBUG","INFO","WARN","ERROR","FATAL","PANIC"]}

	logger.Error(ctx, "test", err)

	// {"time":"2024-02-18T16:42:55+03:00","level":"ERROR","source":"test/main.go:43","msg":"test","username":"unemil","error":"test error"}
}
```

[doc-img]: https://pkg.go.dev/badge/github.com/unemil/logger
[doc]: https://pkg.go.dev/github.com/unemil/logger
[slog-doc]: https://pkg.go.dev/log/slog
