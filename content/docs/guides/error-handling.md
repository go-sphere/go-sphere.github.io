---
title: Error Handling
weight: 42
---

Sphere provides a powerful mechanism for generating typed, consistent error-handling code directly from your `.proto` definitions. By defining your errors as enums, you can ensure that error codes, HTTP statuses, and messages are standardized across your application.

This process is handled by `protoc-gen-sphere-errors`, a `protoc` plugin that inspects your `.proto` files and generates Go error-handling code.

## Installation

To install `protoc-gen-sphere-errors`, use the following command:

```bash
go install github.com/go-sphere/protoc-gen-sphere-errors@latest
```

## Configuration with Buf

To integrate the generator with `buf`, add the plugin to your `buf.gen.yaml` file. This configuration tells `buf` how to execute the plugin and where to place the generated files.

```yaml
version: v2
managed:
  enabled: true
plugins:
  - local: protoc-gen-sphere-errors
    out: api
    opt:
      - paths=source_relative
```

## Defining Errors in `.proto`

Errors are defined as `enum` types in your `.proto` files. You can use custom options from `sphere/errors/errors.proto` to attach metadata like HTTP status codes and default messages to each error.

First, import the necessary definitions in your `.proto` file:

```protobuf
import "sphere/errors/errors.proto";
```

Next, define an `enum` for your errors.

### Example: Basic Error Enum

Here is an example of an error enum:

```protobuf
syntax = "proto3";

package shared.v1;

import "sphere/errors/errors.proto";

enum UserError {
  option (sphere.errors.default_status) = 500;  // Default status for all values
  
  USER_ERROR_UNSPECIFIED = 0;
  USER_ERROR_NOT_FOUND = 1001 [(sphere.errors.options) = {
    status: 404
    message: "User not found"
  }];
  USER_ERROR_INVALID_EMAIL = 1002 [(sphere.errors.options) = {
    status: 400
    reason: "INVALID_EMAIL"
    message: "Invalid email format"
  }];
  USER_ERROR_ALREADY_EXISTS = 1003 [(sphere.errors.options) = {
    status: 409
    reason: "USER_EXISTS"
    message: "User already exists"
  }];
}
```

### Advanced Example with Reasons

```protobuf
enum PaymentError {
  option (sphere.errors.default_status) = 500;
  
  PAYMENT_ERROR_UNSPECIFIED = 0;
  PAYMENT_ERROR_INSUFFICIENT_FUNDS = 2001 [(sphere.errors.options) = {
    status: 402
    reason: "INSUFFICIENT_FUNDS"
    message: "Insufficient funds in account"
  }];
  PAYMENT_ERROR_CARD_DECLINED = 2002 [(sphere.errors.options) = {
    status: 402
    reason: "CARD_DECLINED"
    message: "Payment card was declined"
  }];
  PAYMENT_ERROR_INVALID_AMOUNT = 2003 [(sphere.errors.options) = {
    status: 400
    reason: "INVALID_AMOUNT"
    message: "Payment amount must be positive"
  }];
}
```

### Annotation Reference

- `(sphere.errors.default_status)`: An enum-level option that sets the default HTTP status code for all values. If an error value does not have a specific status, this one will be used.
- `(sphere.errors.options)`: A value-level option to customize a specific error.
  - `status`: The HTTP status code (e.g., `400`, `404`, `500`).
  - `reason`: A machine-readable reason code for programmatic error handling.
  - `message`: A user-facing default error message.

## Using the Generated Code

After running `buf generate`, the plugin will create a file named `{proto_name}_errors.pb.go` (e.g., `user_errors.pb.go`). This file contains a Go enum and several helper methods that allow you to use it as a standard Go error.

### Generated Methods

For each `enum UserError`, the following methods are generated:

- `Error() string`: Returns the error reason, making the type compatible with Go's `error` interface. If `reason` is not set, it returns a string representation of the enum value.
- `GetCode() int32`: Returns the numeric enum value (e.g., `1001`).
- `GetStatus() int32`: Returns the configured HTTP status code.
- `GetMessage() string`: Returns the default error message.
- `GetReason() string`: Returns the error reason (if specified).
- `Join(errs ...error) error`: Wraps one or more source errors, returning a `statuserr.Error` that includes the code, status, and message from the enum. This is the recommended way to return an error while preserving the original cause.
- `JoinWithMessage(msg string, errs ...error) error`: Similar to `Join`, but allows you to provide a custom, dynamic message at runtime.

### Example: Returning an Error in Go

In your service implementation, you can now return one of the generated errors.

```go
package service

import (
    "context"
    "fmt"
    sharedv1 "myproject/api/shared/v1" // Import the generated package
)

func (s *UserService) GetUser(ctx context.Context, req *GetUserRequest) (*User, error) {
    if req.Id <= 0 {
        return nil, sharedv1.UserError_USER_ERROR_INVALID_ID.Join(
            fmt.Errorf("user ID must be positive, got: %d", req.Id))
    }
    
    user, err := s.userRepo.GetByID(ctx, req.Id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, sharedv1.UserError_USER_ERROR_NOT_FOUND.Join(err)
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    if !isValidEmail(req.Email) {
        return nil, sharedv1.UserError_USER_ERROR_INVALID_EMAIL.JoinWithMessage(
            fmt.Sprintf("email '%s' is not valid", req.Email), nil)
    }
    
    // Check if user already exists
    existing, _ := s.userRepo.GetByEmail(ctx, req.Email)
    if existing != nil {
        return nil, sharedv1.UserError_USER_ERROR_ALREADY_EXISTS.Join(
            fmt.Errorf("user with email %s already exists", req.Email))
    }
    
    user, err := s.userRepo.Create(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    return user, nil
}
```

### HTTP Error Response

When this error is handled by Sphere's server layer, it will automatically be converted into an HTTP response with the appropriate status code and JSON body:

```json
{
  "status": 404,
  "code": 1001,
  "error": "USER_NOT_FOUND",
  "message": "User not found"
}
```

## Error Configuration Options

### Enum Level Options

- `default_status`: Sets the default HTTP status code for all enum values that don't specify their own status

### Enum Value Options

- `status`: HTTP status code (overrides default_status)
- `reason`: Optional machine-readable reason code
- `message`: Human-readable error message for client display

## Best Practices

1. **Use meaningful error codes**: Choose enum values that clearly indicate the error type
2. **Set appropriate HTTP status codes**: Use standard HTTP status codes (400, 401, 403, 404, 500, etc.)
3. **Provide clear messages**: Write user-friendly error messages in the appropriate language
4. **Use reasons for API consumers**: Include reason strings for programmatic error handling
5. **Group related errors**: Keep related errors in the same enum for better organization
6. **Preserve original errors**: Always use `.Join()` to wrap underlying errors for better debugging

## Common HTTP Status Codes

- `400`: Bad Request - Client error, invalid input
- `401`: Unauthorized - Authentication required
- `403`: Forbidden - Permission denied
- `404`: Not Found - Resource doesn't exist
- `409`: Conflict - Resource conflict
- `422`: Unprocessable Entity - Validation failed
- `429`: Too Many Requests - Rate limiting
- `500`: Internal Server Error - Server-side error
- `502`: Bad Gateway - External service error
- `503`: Service Unavailable - Service temporarily down

## Integration with buf

Add the required dependency to your `buf.yaml`:

```yaml
version: v2
deps:
  - buf.build/go-sphere/errors
```

Configure the plugin in your `buf.gen.yaml`:

```yaml
version: v2
managed:
  enabled: true
plugins:
  - local: protoc-gen-sphere-errors
    out: api
    opt:
      - paths=source_relative
```

## Error Composition

You can compose multiple errors using the generated methods:

```go
// Simple error with context
return nil, UserError_USER_ERROR_NOT_FOUND.Join(err)

// Error with custom message
return nil, UserError_USER_ERROR_INVALID_EMAIL.JoinWithMessage(
    fmt.Sprintf("Invalid email format: %s", email), validationErr)

// Multiple errors can be joined
return nil, UserError_USER_ERROR_VALIDATION_FAILED.Join(
    emailErr, passwordErr, ageErr)
```

The generated error types integrate seamlessly with Sphere's HTTP server utilities to provide consistent error responses across your entire API.
