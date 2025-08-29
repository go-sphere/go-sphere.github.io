---
title: Logging
weight: 30
---

Sphere provides a flexible and powerful logging system built on top of `go.uber.org/zap`, a high-performance structured logger. It supports multiple outputs, including console and file, and is designed for both local development and production environments.

The core logging functionality can be found in the [`log`](https://github.com/go-sphere/sphere/tree/main/log) directory of the Sphere framework.

## Configuration

The logger is configured through the main `config.json` file in your project root. The configuration allows you to define log levels and outputs.

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

### Configuration Fields

* `level`: The minimum log level to record (e.g., `debug`, `info`, `warn`, `error`).
* `console`: Console logging options.
    * `disable`: Set to `true` to turn off console output (default: `false`)
* `file`: File logging options. If this section is omitted, file logging is disabled.
    * `file_name`: The path to the log file.
    * `max_size`: The maximum size in megabytes of the log file before it gets rotated.
    * `max_backups`: The maximum number of old log files to retain.
    * `max_age`: The maximum number of days to retain old log files.

## Basic Usage

Sphere provides a global logger that can be used anywhere in your application.

### Standard Logging

You can use the global functions for standard logging:

```go
package main

import "github.com/go-sphere/sphere/log"

func main() {
	log.Debug("This is a debug message")
	log.Info("This is an info message", log.String("user", "test"))
	log.Warn("This is a warning")
	log.Error("This is an error", log.Err(fmt.Errorf("an error occurred")))
}
```

### Structured Logging

To add structured context to your logs, you can use `log.With` to create a new logger instance with predefined fields.

```go
logger := log.With(log.String("service", "UserService"), log.String("traceId", "xyz-123"))

logger.Info("User lookup successful")
// Output will include {"service": "UserService", "traceId": "xyz-123", "message": "User lookup successful"}
```

### Common Log Fields

Sphere provides convenient helper functions for common log fields:

```go
// User-related fields
log.Info("User created", 
    log.String("userId", "123"),
    log.String("email", "user@example.com"))

// Error logging
log.Error("Database connection failed", 
    log.Err(err),
    log.String("database", "postgresql"),
    log.Duration("timeout", 5*time.Second))

// HTTP request logging
log.Info("HTTP request processed",
    log.String("method", "POST"),
    log.String("path", "/api/users"),
    log.Int("status", 201),
    log.Duration("duration", 150*time.Millisecond))
```

## Service Integration

### In HTTP Handlers

Use structured logging in your HTTP handlers to track requests:

```go
func (s *UserService) CreateUser(c *gin.Context) {
    logger := log.With(
        log.String("handler", "CreateUser"),
        log.String("requestId", c.GetHeader("X-Request-ID")),
    )
    
    logger.Info("Creating user request received")
    
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        logger.Error("Failed to bind request", log.Err(err))
        c.JSON(400, gin.H{"error": "Invalid request"})
        return
    }
    
    user, err := s.userBiz.CreateUser(c.Request.Context(), &req)
    if err != nil {
        logger.Error("Failed to create user", 
            log.Err(err),
            log.String("email", req.Email))
        c.JSON(500, gin.H{"error": "Internal server error"})
        return
    }
    
    logger.Info("User created successfully", 
        log.String("userId", user.ID),
        log.String("email", user.Email))
    
    c.JSON(201, user)
}
```

### In Business Logic

Add contextual logging to your business logic:

```go
func (b *UserBiz) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    logger := log.With(log.String("method", "CreateUser"))
    
    logger.Debug("Validating user input", log.String("email", req.Email))
    
    if err := b.validateEmail(req.Email); err != nil {
        logger.Warn("Email validation failed", 
            log.Err(err),
            log.String("email", req.Email))
        return nil, fmt.Errorf("invalid email: %w", err)
    }
    
    logger.Debug("Checking if user exists")
    existing, err := b.userRepo.GetByEmail(ctx, req.Email)
    if err != nil {
        logger.Error("Database query failed", log.Err(err))
        return nil, fmt.Errorf("failed to check existing user: %w", err)
    }
    
    if existing != nil {
        logger.Warn("User already exists", log.String("email", req.Email))
        return nil, ErrUserAlreadyExists
    }
    
    logger.Info("Creating new user", log.String("email", req.Email))
    user, err := b.userRepo.Create(ctx, req)
    if err != nil {
        logger.Error("Failed to create user in database", log.Err(err))
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    logger.Info("User created successfully", 
        log.String("userId", user.ID),
        log.String("email", user.Email))
    
    return user, nil
}
```

## Log Levels

Sphere supports the following log levels:

### Debug
Use for detailed diagnostic information, typically only of interest when diagnosing problems.

```go
log.Debug("Processing user input", 
    log.String("input", userInput),
    log.Int("inputLength", len(userInput)))
```

### Info
Use for general informational messages that highlight the progress of the application.

```go
log.Info("Server started", 
    log.String("address", ":8080"),
    log.String("environment", "production"))
```

### Warn
Use for potentially harmful situations that are not necessarily errors.

```go
log.Warn("Rate limit approaching", 
    log.String("userId", userID),
    log.Int("currentRequests", currentReqs),
    log.Int("limit", rateLimit))
```

### Error
Use for error events that do not stop the application from running.

```go
log.Error("Failed to send email", 
    log.Err(err),
    log.String("recipient", email),
    log.String("template", "welcome"))
```

## Performance Considerations

### Use Appropriate Log Levels
Set your log level appropriately for each environment:
- **Development**: `debug` for detailed information
- **Staging**: `info` for general application flow
- **Production**: `warn` or `error` for minimal overhead

### Structured Fields
Use structured fields instead of string formatting:

```go
// Good - structured logging
log.Info("User logged in", 
    log.String("userId", userID),
    log.String("ip", clientIP))

// Avoid - string formatting
log.Info(fmt.Sprintf("User %s logged in from %s", userID, clientIP))
```

### Conditional Logging
For expensive log operations, check log level first:

```go
if log.IsDebugEnabled() {
    expensiveData := generateExpensiveLogData()
    log.Debug("Expensive debug info", log.Any("data", expensiveData))
}
```

## Log Collection and Monitoring

### File-based Collection

For simple log collection, you can use file output with log rotation:

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

### Simple Collection with Logdy

For simple, real-time log viewing without a complex setup, you can use [Logdy](https://github.com/logdyhq/logdy-core). It's a lightweight tool that can stream logs directly to your browser.

If your application is writing logs to `app.log`, you can start Logdy with the following command:

```bash
tail -f app.log | logdy
# or
logdy follow app.log
```

This is ideal for debugging during development or monitoring a single instance.

### Advanced Collection with Grafana Loki

For a production-grade, scalable log aggregation system, Sphere recommends using [Grafana Loki](https://grafana.com/oss/loki/). Loki is a horizontally scalable, multi-tenant log aggregation system inspired by Prometheus.

You can set up a local Loki stack using Docker Compose:

```yaml
version: "3.8"

services:
  loki:
    image: grafana/loki:latest
    ports:
      - "3100:3100"
    command: -config.file=/etc/loki/local-config.yaml
    
  promtail:
    image: grafana/promtail:latest
    volumes:
      - ./app.log:/var/log/app.log:ro
      - ./promtail-config.yml:/etc/promtail/config.yml
    command: -config.file=/etc/promtail/config.yml
    depends_on:
      - loki
      
  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    depends_on:
      - loki
```

### JSON Output for Structured Logging

For production environments, consider using JSON output format for better integration with log aggregation systems:

```json
{
  "log": {
    "level": "info",
    "console": {
      "disable": false,
      "format": "json"
    }
  }
}
```

This produces machine-readable logs:

```json
{
  "level": "info",
  "time": "2023-08-29T10:30:00Z",
  "caller": "service/user.go:45",
  "message": "User created successfully",
  "userId": "123",
  "email": "user@example.com"
}
```

## Best Practices

### Context Propagation
Always include relevant context in your logs:

```go
func (s *Service) ProcessOrder(ctx context.Context, orderID string) error {
    logger := log.With(
        log.String("orderID", orderID),
        log.String("userID", getUserIDFromContext(ctx)),
    )
    
    logger.Info("Processing order")
    // ... processing logic
}
```

### Error Logging
When logging errors, include the error and relevant context:

```go
if err != nil {
    log.Error("Payment processing failed",
        log.Err(err),
        log.String("orderID", order.ID),
        log.String("paymentMethod", order.PaymentMethod),
        log.Float64("amount", order.Amount))
    return err
}
```

### Sensitive Data
Never log sensitive information like passwords, API keys, or personal data:

```go
// Bad
log.Info("User login", log.String("password", user.Password))

// Good
log.Info("User login attempt", 
    log.String("email", user.Email),
    log.Bool("success", loginSuccess))
```

### Performance Monitoring
Use logs to track performance metrics:

```go
start := time.Now()
defer func() {
    log.Info("Operation completed",
        log.String("operation", "processPayment"),
        log.Duration("duration", time.Since(start)))
}()
```

## Integration with Monitoring

### Metrics from Logs
You can extract metrics from structured logs using tools like Prometheus with log-based alerting:

```go
// Log metrics-friendly data
log.Info("HTTP request completed",
    log.String("method", "POST"),
    log.String("endpoint", "/api/users"),
    log.Int("status", 201),
    log.Duration("duration", duration),
    log.String("user_agent", userAgent))
```

### Alerting
Set up alerts based on log patterns:

```go
// Log errors with severity levels
log.Error("Database connection failed",
    log.Err(err),
    log.String("severity", "critical"),
    log.String("component", "database"))
```

### Distributed Tracing
Include trace IDs in logs for distributed tracing:

```go
traceID := getTraceIDFromContext(ctx)
logger := log.With(log.String("traceID", traceID))

logger.Info("Processing request")
```

## Troubleshooting

### Common Issues

**High CPU usage from logging**: Reduce log level in production or use asynchronous logging.

**Disk space issues**: Ensure log rotation is configured properly.

**Missing context**: Always use structured logging with relevant fields.

### Debugging Tips

1. **Use appropriate log levels** during development
2. **Include request IDs** for tracing requests across services
3. **Log entry and exit points** of important functions
4. **Use structured fields** consistently for better searchability

This logging approach provides observability into your Sphere applications while maintaining performance and enabling effective monitoring and debugging.
