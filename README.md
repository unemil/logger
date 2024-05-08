# [![doc-img]][doc]

Structured logging with context support based on a [slog][slog-doc].

## Install

```sh
go get github.com/unemil/logger
```

**Compatibility:** go >= 1.21

## Usage

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
	"net/http"

	"github.com/unemil/logger"
)

func main() {
	var (
		ctx = logger.Context(context.Background(), "username", "unemil")
		err = errors.New(http.StatusText(http.StatusUnauthorized))
	)

	logger.Errorf(ctx, "test", err, logger.Field("status", http.StatusUnauthorized))

	// {"time":"2024-05-08T20:15:05+03:00","level":"ERROR","source":"test/main.go:17","msg":"test","error":"Unauthorized","status":401,"username":"unemil"}
}
```

[doc-img]: https://pkg.go.dev/badge/github.com/unemil/logger
[doc]: https://pkg.go.dev/github.com/unemil/logger
[slog-doc]: https://pkg.go.dev/log/slog
