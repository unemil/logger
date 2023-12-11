# [![doc-img]][doc]

Structured logging implementation using a standard library.

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

See the [documentation][doc] for more details.

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

	ctx := logger.WithContext(context.Background(), "username", "unemil")
	logger.Debug(ctx, "test")

	// {"time":"2023-12-11T03:29:23+03:00","level":"DEBUG","source":"test/main.go:13","msg":"test","username":"unemil"}
}
```

[doc-img]: https://pkg.go.dev/badge/github.com/unemil/logger
[doc]: https://pkg.go.dev/github.com/unemil/logger
