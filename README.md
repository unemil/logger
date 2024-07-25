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
export LOG_FORMAT=FORMAT
```

## Levels

- TRACE
- DEBUG
- INFO (default)
- WARN
- ERROR
- FATAL
- PANIC

## Formats

- TEXT (default)
- JSON

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
	logger.Errorf(
		logger.Context(context.Background(), "username", "unemil"),
		"test",
		errors.New(http.StatusText(http.StatusUnauthorized)),
		logger.Field("status", http.StatusUnauthorized),
	)

	// time=2024-07-25T13:25:53+03:00 level=ERROR source=test/main.go:12 msg=test error=Unauthorized status=401 username=unemil
}
```

[doc-img]: https://pkg.go.dev/badge/github.com/unemil/logger
[doc]: https://pkg.go.dev/github.com/unemil/logger
[slog-doc]: https://pkg.go.dev/log/slog
