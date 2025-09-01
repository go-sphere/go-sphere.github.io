---
title: Logging
weight: 43
---

Sphere provides a flexible logging system built on `go.uber.org/zap` for high-performance structured logging. It supports console and file outputs for both development and production environments.

## Configuration

Configure logging through your `config.json` file:

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

- `level`: Minimum log level (`debug`, `info`, `warn`, `error`)
- `console.disable`: Disable console output (default: `false`)
- `file`: File logging configuration (optional)
  - `file_name`: Log file path
  - `max_size`: Max file size in MB before rotation
  - `max_backups`: Number of backup files to keep
  - `max_age`: Days to retain old files

## Basic Usage

### Standard Logging

Use global logging functions with structured fields:

```go
package main

import "github.com/go-sphere/sphere/log"

func main() {
    log.Debug("Debug message")
    log.Info("User created", log.String("user", "john"))
    log.Warn("Warning message")
    log.Error("Error occurred", log.Err(err))
}
```

### Printf-style Logging

You can also use formatted logging:

```go
log.Infof("User %s created with ID %d", username, userID)
log.Errorf("Failed to connect to %s: %v", host, err)
```

### Structured Logging

Add context using structured fields:

```go
log.Info("User operation",
    log.String("user_id", "123"),
    log.String("action", "login"),
    log.Duration("duration", time.Since(start)))
```

### Available Field Types

```go
log.String("key", "value")          // String field
log.Int("count", 42)                // Integer field  
log.Int64("timestamp", time.Now().Unix())  // 64-bit integer
log.Float64("score", 98.5)          // Float field
log.Bool("success", true)           // Boolean field
log.Duration("elapsed", duration)    // Time duration
log.Time("created_at", time.Now())  // Time field
log.Err(err)                        // Error field (key="error")
log.Any("data", complexObject)      // Any type
```

## Logger Instances

Create logger instances with additional context:

```go
// Create a logger with predefined attributes
logger := log.With(log.WithAttrs(map[string]any{
    "service": "user-service",
    "version": "1.0.0",
}))

logger.Info("Service started")
// Output includes: {"service": "user-service", "version": "1.0.0", "message": "Service started"}
```

### Multiple Options

You can combine multiple options:

```go
logger := log.With(
    log.WithName("user-service"),
    log.WithAttrs(map[string]any{"component": "api"}),
    log.AddCaller(),
)

logger.Info("Processing request")
```

## Usage Examples

### HTTP Handler Logging

```go
func (h *UserHandler) CreateUser(c *gin.Context) {
    logger := log.With(log.WithAttrs(map[string]any{
        "handler": "CreateUser",
        "trace_id": c.GetString("trace_id"),
    }))
    
    logger.Info("Request received")
    
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        logger.Error("Invalid request", log.Err(err))
        c.JSON(400, gin.H{"error": "Invalid request"})
        return
    }
    
    user, err := h.userService.Create(c.Request.Context(), &req)
    if err != nil {
        logger.Error("Failed to create user", 
            log.Err(err),
            log.String("email", req.Email))
        c.JSON(500, gin.H{"error": "Internal error"})
        return
    }
    
    logger.Info("User created", 
        log.String("user_id", user.ID),
        log.String("email", user.Email))
    
    c.JSON(201, user)
}
```

### Business Logic Logging

```go
func (s *UserService) Create(ctx context.Context, req *CreateUserRequest) (*User, error) {
    logger := log.With(log.WithAttrs(map[string]any{
        "service": "UserService",
        "method": "Create",
    }))
    
    logger.Debug("Validating input", log.String("email", req.Email))
    
    if err := s.validateEmail(req.Email); err != nil {
        logger.Warn("Validation failed", log.Err(err))
        return nil, err
    }
    
    user, err := s.repo.Create(ctx, req)
    if err != nil {
        logger.Error("Database error", log.Err(err))
        return nil, err
    }
    
    logger.Info("User created", log.String("user_id", user.ID))
    return user, nil
}
```

## Best Practices

### Log Levels
- **Debug**: Detailed diagnostic information
- **Info**: General application flow
- **Warn**: Potentially harmful situations  
- **Error**: Error events that don't stop the application

```go
log.Debug("Processing input", log.Int("size", len(data)))
log.Info("Server started", log.String("port", ":8080"))
log.Warn("Rate limit approaching", log.String("user", userID))
log.Error("Database error", log.Err(err))
```

### Error Handling
Always include error context:

```go
if err != nil {
    log.Error("Operation failed",
        log.Err(err),
        log.String("operation", "user_creation"),
        log.String("user_id", userID))
    return err
}
```

### Performance Tips
- Use appropriate log levels for each environment
- Avoid logging sensitive data (passwords, tokens)
- Use structured fields instead of string formatting

```go
// Good - structured
log.Info("User login", 
    log.String("user_id", userID),
    log.String("ip", clientIP))

// Avoid - formatted strings  
log.Info(fmt.Sprintf("User %s logged in from %s", userID, clientIP))
```

### Context Propagation
Include relevant context in logs:

```go
logger := log.With(log.WithAttrs(map[string]any{
    "request_id": getRequestID(ctx),
    "user_id": getUserID(ctx),
}))

logger.Info("Processing request")
```

## Log Collection

### File Rotation
Configure automatic log rotation:

```json
{
  "log": {
    "file": {
      "file_name": "app.log",
      "max_size": 100,
      "max_backups": 3,
      "max_age": 28
    }
  }
}
```

### Simple Log Viewing with Logdy
For development, use [Logdy](https://github.com/logdyhq/logdy-core) for real-time log viewing:

```bash
tail -f app.log | logdy
# or
logdy follow app.log
```

### Production Setup with Grafana Loki
For production, use [Grafana Loki](https://grafana.com/oss/loki/) with Docker Compose:

```yaml
version: "3.8"
services:
  loki:
    image: grafana/loki:latest
    ports:
      - "3100:3100"
    
  promtail:
    image: grafana/promtail:latest
    volumes:
      - ./app.log:/var/log/app.log:ro
      - ./promtail-config.yml:/etc/promtail/config.yml
    depends_on:
      - loki
      
  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
```

## Advanced Features

### Initialization with Attributes
Set global attributes for all logs:

```go
func main() {
    config := log.NewDefaultConfig()
    config.Level = "debug"
    
    attrs := map[string]any{
        "service": "user-api",
        "version": "1.0.0",
        "environment": "production",
    }
    
    log.Init(config, attrs)
    defer log.Sync() // Flush logs before exit
    
    log.Info("Application started")
}
```

### Logger Options
Available options for customizing loggers:

```go
logger := log.With(
    log.WithName("component-name"),           // Set logger name
    log.AddCaller(),                         // Include caller info
    log.WithAttrs(map[string]any{            // Add attributes
        "module": "auth",
    }),
)
```
