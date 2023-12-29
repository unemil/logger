# [![doc-img]][doc]

Structured logging based on [slog][slog-doc].

## Install

```sh
go get github.com/unemil/logger
```

**Compatibility:** go >= 1.21

## Usage

Set `LOG_LEVEL` one of the [LEVEL](#levels) values (default: **INFO**).

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

	"github.com/unemil/logger"
)

func main() {
	// LOG_LEVEL=DEBUG

	ctx := logger.With(context.Background(), "username", "unemil")
	logger.Debugf(ctx, "test", logger.Fields{"key": "value"})

	// {"time":"2023-12-27T14:42:09+03:00","level":"DEBUG","source":"test/main.go:13","msg":"test","key":"value","username":"unemil"}
}
```

[doc-img]: https://pkg.go.dev/badge/github.com/unemil/logger
[doc]: https://pkg.go.dev/github.com/unemil/logger
[slog-doc]: https://pkg.go.dev/log/slog
