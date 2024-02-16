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

func main() {
	// LOG_LEVEL=DEBUG

	var (
		ctx = context.Background()
		err = errors.New("test error")
	)

	logger.Info(ctx, "test")

	// {"time":"2024-02-15T16:32:52+03:00","level":"INFO","source":"test/main.go:18","msg":"test"}

	ctx = logger.With(ctx, "username", "unemil")
	logger.Debugf(ctx, "test", logger.Fields{
		"username":     nil,
		"error_levels": []logger.Level{logger.LevelError, logger.LevelFatal, logger.LevelPanic},
	})

	// {"time":"2024-02-15T16:32:52+03:00","level":"DEBUG","source":"test/main.go:23","msg":"test","error_levels":["ERROR","FATAL","PANIC"],"username":null}

	logger.Error(ctx, "test", err)

	// {"time":"2024-02-15T16:32:52+03:00","level":"ERROR","source":"test/main.go:30","msg":"test","error":"test error","username":"unemil"}
}
```

[doc-img]: https://pkg.go.dev/badge/github.com/unemil/logger
[doc]: https://pkg.go.dev/github.com/unemil/logger
[slog-doc]: https://pkg.go.dev/log/slog
