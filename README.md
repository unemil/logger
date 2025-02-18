# [![doc-img]][doc]

[doc-img]: https://pkg.go.dev/badge/github.com/unemil/logger
[doc]: https://pkg.go.dev/github.com/unemil/logger
[slog-doc]: https://pkg.go.dev/log/slog

Structured logging with context support based on [slog][slog-doc].

## Install

```sh
go get github.com/unemil/logger
```

**Compatibility:** go >= 1.21

## Usage

```sh
export LOG_LEVEL=LEVEL
export LOG_FORMAT=FORMAT
export LOG_FILE=FILE
```

## Levels

- TRACE
- DEBUG
- INFO
- WARN
- ERROR
- FATAL
- PANIC

## Formats

- TEXT
- JSON
