---
title: protoc-gen-sphere-errors
weight: 35
---

`protoc-gen-sphere-errors` is a protoc plugin that generates error handling code from `.proto` files. It is designed to inspect enum definitions within your protobuf files and automatically generate corresponding error handling code based on the sphere errors framework. This plugin creates Go code that provides structured error handling with HTTP status codes, error codes, and customizable messages.

This code is inspired by [protoc-gen-go-errors](https://github.com/go-kratos/kratos/tree/main/cmd/protoc-gen-go-errors) but is specifically designed for the go-sphere framework.

## Features

- Generates error structs with HTTP status codes
- Supports custom error messages and reasons
- Provides `Join` and `JoinWithMessage` methods for error composition
- Integrates with the sphere error handling framework
- Supports default status codes for enum types
- Individual error value customization through options

## Installation

To install `protoc-gen-sphere-errors`, use the following command:

```bash
go install github.com/go-sphere/protoc-gen-sphere-errors@latest
```

## Prerequisites

You need to have the sphere errors proto definitions in your project. Add the following dependency to your `buf.yaml`:

```yaml
deps:
  - buf.build/go-sphere/errors
```

## Usage with Buf

To use `protoc-gen-sphere-errors` with `buf`, you can configure it in your `buf.gen.yaml` file. Here is an example configuration:

```yaml
version: v2
managed:
  enabled: true
  disable:
    - file_option: go_package_prefix
      module: buf.build/go-sphere/errors
  override:
    - file_option: go_package_prefix
      value: github.com/go-sphere/sphere-layout/api
plugins:
  - local: protoc-gen-sphere-errors
    out: api
    opt:
      - paths=source_relative
```

## Proto Definition Example

Here's how to define error enums in your `.proto` files:

```protobuf
syntax = "proto3";

package shared.v1;

import "sphere/errors/errors.proto";

enum TestError {
  option (sphere.errors.default_status) = 500;
  
  TEST_ERROR_UNSPECIFIED = 0;
  TEST_ERROR_INVALID_PATH_TEST2 = 1001 [(sphere.errors.options) = {
    status: 400
    reason: "INVALID_PATH"
    message: "Invalid path parameter"
  }];
  TEST_ERROR_MISSING_FIELD = 1002 [(sphere.errors.options) = {
    status: 400
    reason: "MISSING_FIELD"
    message: "Required field is missing"
  }];
}

enum UserError {
  option (sphere.errors.default_status) = 500;
  
  USER_ERROR_UNSPECIFIED = 0;
  USER_ERROR_NOT_FOUND = 2001 [(sphere.errors.options) = {
    status: 404
    reason: "USER_NOT_FOUND"
    message: "User not found"
  }];
  USER_ERROR_INVALID_EMAIL = 2002 [(sphere.errors.options) = {
    status: 400
    reason: "INVALID_EMAIL"
    message: "Invalid email format"
  }];
  USER_ERROR_ALREADY_EXISTS = 2003 [(sphere.errors.options) = {
    status: 409
    reason: "USER_EXISTS"
    message: "User already exists"
  }];
}
```

## Generated Code

The plugin generates Go code with the following methods for each error enum:

- `Error() string` - Returns a string representation of the error
- `GetCode() int32` - Returns the error code (enum value)
- `GetStatus() int32` - Returns the HTTP status code
- `GetMessage() string` - Returns the custom error message
- `GetReason() string` - Returns the error reason (if specified)
- `Join(errs ...error) error` - Wraps the error with additional errors
- `JoinWithMessage(msg string, errs ...error) error` - Wraps with custom message

## Usage in Code

### Direct Error Returns

```go
func (s *service) ValidateField(field string) error {
    if field == "" {
        return sharedv1.TestError_TEST_ERROR_MISSING_FIELD
    }
    return nil
}
```

### Error Handling in HTTP Handlers

```go
func (s *service) RunTest(ctx context.Context, req *sharedv1.RunTestRequest) (*sharedv1.RunTestResponse, error) {
    if req.FieldTest1 == "" {
        return nil, sharedv1.TestError_TEST_ERROR_MISSING_FIELD.Join(
            fmt.Errorf("field_test1 cannot be empty"))
    }
    
    // Business logic...
    user, err := s.userRepo.GetByID(ctx, req.UserId)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, sharedv1.UserError_USER_ERROR_NOT_FOUND.Join(err)
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    return &sharedv1.RunTestResponse{
        FieldTest1: req.FieldTest1,
    }, nil
}
```

### Custom Messages

```go
func (s *service) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    if !isValidEmail(req.Email) {
        return nil, sharedv1.UserError_USER_ERROR_INVALID_EMAIL.JoinWithMessage(
            fmt.Sprintf("Email '%s' does not match required format", req.Email), nil)
    }
    
    // Check if user exists
    existing, _ := s.userRepo.GetByEmail(ctx, req.Email)
    if existing != nil {
        return nil, sharedv1.UserError_USER_ERROR_ALREADY_EXISTS.JoinWithMessage(
            fmt.Sprintf("User with email '%s' already exists", req.Email), nil)
    }
    
    // Create user...
}
```

## HTTP Response Integration

When used with Sphere's HTTP server utilities, these errors are automatically converted to structured JSON responses:

```json
{
  "status": 404,
  "code": 2001,
  "error": "USER_NOT_FOUND",
  "message": "User not found"
}
```

## Best Practices

1. **Use meaningful error codes**: Choose enum values that clearly indicate the error type
2. **Set appropriate HTTP status codes**: Use standard HTTP status codes (400, 401, 403, 404, 500, etc.)
3. **Provide clear messages**: Write user-friendly error messages that can be displayed to end users
4. **Use reasons for programmatic handling**: Include reason strings that client applications can use for conditional logic
5. **Group related errors**: Keep related errors in the same enum for better organization
6. **Preserve context with Join**: Always use `.Join()` to wrap underlying errors for better debugging

## Integration Notes

- Sphere's Gin layer maps these to structured JSON with correct HTTP status
- Pair with a global error parser if you need to merge validation/notfound/custom errors
- The generated errors integrate seamlessly with Sphere's server utilities for consistent API responses

## See Also

- Guides: ../../guides/error-handling
- Concepts: ../../concepts/protocol-and-codegen
