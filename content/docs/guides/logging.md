---
title: Logging
weight: 43
---

Sphere's `log` package is a backend-agnostic structured logger. Package-level `Debug` / `Info` / `Warn` / `Error` calls go through a global logger that starts as a stdio backend (logfmt on stdout, Error+ on stderr). Call `log.InitWithBackends` early in `main`, or let `boot.WithLoggerBackend` do it.

Official templates still store a `zapx.Config` in `config.json` and install a zap backend at boot.

## Configuration

Configure the zap backend through your `config.json` file:

```json
{
  "log": {
    "level": "info",
    "console": {
      "disable": false
    },
    "file": {
      "file_name": "app.log",
      "max_size": 100,
      "max_backups": 3,
      "max_age": 28
    }
  }
}
```

### Configuration Options

- `level`: Minimum log level for zap (`debug`, `info`, `warn`, `error`). `log.WithMinLevel` does not apply to zapx; set this field instead.
- `console.disable`: Disable console output (default: `false`)
- `file`: File logging configuration (optional)
  - `file_name`: Log file path. Empty disables the file sink.
  - `max_size`: Max file size in MB before rotation
  - `max_backups`: Number of backup files to keep
  - `max_age`: Days to retain old files

## Initialization

```go
import (
    "github.com/go-sphere/sphere/core/boot"
    "github.com/go-sphere/sphere/log"
    "github.com/go-sphere/sphere/log/zapx"
)

func main() {
    backend := zapx.NewBackend(conf.Log, log.WithAttrs(map[string]any{
        "service": "user-api",
        "version": version,
    }))
    err := boot.Run(conf, app, boot.WithLoggerBackend(backend))
}
```

`WithLoggerBackend` installs the backend before start and syncs it after stop. If the backend implements `SlogLogger`, it also becomes the default `log/slog` handler.

`boot.WithLoggerInit(version, conf.Log)` still exists as a deprecated wrapper around the same path. New code should construct the backend explicitly.

To install the logger without boot:

```go
log.InitWithBackends(zapx.NewBackend(conf.Log))
defer log.Sync()
```

`InitWithBackends` does not close the previous backend. An empty or all-nil list keeps the current logger and warns on stderr. Pass `log.NewNopBackend()` to discard logs on purpose.

## Basic Usage

```go
log.Debug("Debug message")
log.Info("User created", log.String("user", "john"))
log.Warn("Warning message")
log.Error("Error occurred", log.Err(err))
```

Printf-style helpers exist (`Infof`, `Errorf`, …) but structured fields are preferred.

```go
log.Info("User operation",
    log.String("user_id", "123"),
    log.String("action", "login"),
    log.Duration("duration", time.Since(start)))
```

### Available Field Types

```go
log.String("key", "value")
log.Int("count", 42)
log.Int64("timestamp", time.Now().Unix())
log.Uint64("size", n)
log.Float64("score", 98.5)
log.Bool("success", true)
log.Duration("elapsed", duration)
log.Time("created_at", time.Now())
log.Err(err)                     // key="error"
log.Any("data", complexObject)
log.Group("req", log.String("id", id))
```

Fields are `slog.Attr` values. There is also a `log.Field` alias for compatibility.

## Logger Instances

```go
logger := log.With(
    log.WithName("user-service"),
    log.WithAttrs(map[string]any{"component": "api"}),
    log.AddCaller(),
)

logger.Info("Processing request")
```

Context-aware methods (`InfoContext`, `ErrorContext`, …) pass the caller's `context.Context` to the backend. `WrapBackendWithContextMerge` can inject attributes from that context.

### Options

- `log.WithName(name)` — logger name
- `log.AddCaller()` / `log.DisableCaller()` — file:line
- `log.WithAttrs(map[string]any{...})` — attributes on every line
- `log.WithMinLevel(level)` — drop entries below this level on `StdioBackend`. Ignored by zapx; set `zapx.Config.Level` instead.
- `log.WithStackAt(level)` — attach a stack at that level and above. It is **not** a level filter.

`StdioBackend` used to treat `WithStackAt` as a minimum-level filter. Code that relied on that should switch to `WithMinLevel`.

## HTTP Handler Logging

```go
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    logger := log.With(log.WithAttrs(map[string]any{
        "service": "UserService",
        "method":  "Create",
    }))

    user, err := s.repo.Create(ctx, req)
    if err != nil {
        logger.Error("Database error", log.Err(err))
        return nil, err
    }

    logger.Info("User created", log.String("user_id", user.ID))
    return user, nil
}
```

Keep logging in the service layer. Generated HTTP handlers already go through `httpz`; they do not need per-request `*gin.Context` log calls.

## Best Practices

- Use structured fields instead of `fmt.Sprintf` in the message
- Do not log passwords, tokens, or raw request bodies
- Include operation name and identifiers on error lines
- Prefer `log.Err(err)` over embedding the error in the message

```go
// Good
log.Info("User login",
    log.String("user_id", userID),
    log.String("ip", clientIP))

// Avoid
log.Info(fmt.Sprintf("User %s logged in from %s", userID, clientIP))
```

## Log Collection

File rotation is configured on `zapx.Config.File`. For development, [Logdy](https://github.com/logdyhq/logdy-core) can tail the file. For production, ship the JSON file with Promtail / Loki or any other collector; Sphere does not own that pipeline.

## Related

- [HTTP Runtime](http-runtime) — how `httpz` logs panics without leaking them to clients
- [Upgrading to v0.0.4](upgrading) — logger init and `WithStackAt` behavior changes
