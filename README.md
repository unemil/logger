# [![doc-img]][doc]

Structured logging with context support based on a [slog][slog-doc].

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
	"math"
	"net/http"

	"github.com/unemil/logger"
)

var ctx = context.Background()

func main() {
	// LOG_LEVEL=DEBUG

	logger.Info(ctx, "test")

	// {"time":"2024-02-29T03:13:30+03:00","level":"INFO","source":"test/main.go:17","msg":"test"}

	ctx = logger.Context(ctx, "username", "unemil")
	logger.Debugf(ctx, "test",
		logger.Field("primes", []int{2, 3, 5, 7, 11}),
		logger.Field("pi", math.Pi),
		logger.Field("username", nil),
	)

	// {"time":"2024-02-29T03:13:30+03:00","level":"DEBUG","source":"test/main.go:22","msg":"test","primes":[2,3,5,7,11],"pi":3.141592653589793,"username":null}

	logger.Errorf(ctx, "test", errors.New("test error"), logger.Field("status", http.StatusInternalServerError))

	// {"time":"2024-02-29T03:13:30+03:00","level":"ERROR","source":"test/main.go:30","msg":"test","username":"unemil","error":"test error","status":500}
}
```

[doc-img]: https://pkg.go.dev/badge/github.com/unemil/logger
[doc]: https://pkg.go.dev/github.com/unemil/logger
[slog-doc]: https://pkg.go.dev/log/slog
