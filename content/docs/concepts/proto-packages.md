---
title: Proto Packages
weight: 30
---

Sphere provides several core proto packages that extend Protobuf with specialized annotations for HTTP binding, error handling, and method options. These packages enable the code generation pipeline and provide the foundation for Sphere's contract-first approach.

## Overview

The core proto packages are:

- **sphere/binding**: HTTP request binding annotations
- **sphere/errors**: Structured error handling with HTTP status codes  
- **sphere/options**: Method-level metadata and routing options

These packages are published to the Buf Schema Registry and can be imported into your projects as dependencies.

## sphere/binding

The binding package provides annotations for controlling how protobuf message fields are bound to HTTP request components.

### Installation

Add to your `buf.yaml`:

```yaml
deps:
  - buf.build/go-sphere/binding
```

### Features

- Field-level binding location control (URI, query, body, header, form)
- Custom struct tag generation
- Message-level and oneof-level defaults
- Integration with Go HTTP frameworks like Gin

### Basic Usage

```protobuf
syntax = "proto3";

import "sphere/binding/binding.proto";

message GetUserRequest {
  // Bind from URI path parameter
  int64 user_id = 1 [(sphere.binding.location) = BINDING_LOCATION_URI];
  
  // Bind from query parameter
  repeated string fields = 2 [(sphere.binding.location) = BINDING_LOCATION_QUERY];
  
  // Bind from request header
  string auth_token = 3 [(sphere.binding.location) = BINDING_LOCATION_HEADER];
}

message CreateUserRequest {
  // Fields without annotation default to JSON body for POST/PUT/PATCH
  string name = 1;
  string email = 2;
}
```

### Message-Level Defaults

Set default binding behavior for entire messages:

```protobuf
message SearchRequest {
  option (sphere.binding.default_location) = BINDING_LOCATION_QUERY;
  option (sphere.binding.default_auto_tags) = "form";
  
  string query = 1;     // Will be bound from query parameters
  int32 limit = 2;      // Will be bound from query parameters  
  int32 offset = 3;     // Will be bound from query parameters
}
```

### Custom Struct Tags

Add custom Go struct tags for integration with validation, database, or other libraries:

```protobuf
message User {
  option (sphere.binding.default_auto_tags) = "db";
  
  string name = 1;   // Generated: `db:"name" json:"name"`
  string email = 2 [(sphere.binding.auto_tags) = "validate:\"email\""];
  // Generated: `db:"email" json:"email" validate:"email"`
}
```

### Available Binding Locations

- `BINDING_LOCATION_URI`: Path parameters (`{id}` in URL)
- `BINDING_LOCATION_QUERY`: Query string parameters (`?param=value`)
- `BINDING_LOCATION_BODY`: JSON request body (default for non-GET methods)
- `BINDING_LOCATION_HEADER`: HTTP headers
- `BINDING_LOCATION_FORM`: Form data

## sphere/errors

The errors package enables defining structured, typed errors directly in protobuf with automatic Go code generation.

### Installation

Add to your `buf.yaml`:

```yaml
deps:
  - buf.build/go-sphere/errors
```

### Features

- Enum-based error definitions
- HTTP status code mapping
- Custom error messages and reason codes
- Generated Go helpers for error wrapping
- Consistent JSON error responses

### Basic Usage

```protobuf
syntax = "proto3";

import "sphere/errors/errors.proto";

enum UserError {
  option (sphere.errors.default_status) = 500;  // Default for all values
  
  USER_ERROR_UNSPECIFIED = 0;
  USER_ERROR_NOT_FOUND = 1001 [(sphere.errors.options) = {
    status: 404
    reason: "USER_NOT_FOUND"  
    message: "User not found"
  }];
  USER_ERROR_INVALID_EMAIL = 1002 [(sphere.errors.options) = {
    status: 400
    reason: "INVALID_EMAIL"
    message: "Invalid email format"
  }];
}
```

### Generated Go Usage

The plugin generates helpful methods for each error enum:

```go
// Return an error
return nil, UserError_USER_ERROR_NOT_FOUND.Join(dbErr)

// Return with custom message
return nil, UserError_USER_ERROR_INVALID_EMAIL.JoinWithMessage(
    fmt.Sprintf("Email '%s' is invalid", email), validationErr)

// Check error properties
err := UserError_USER_ERROR_NOT_FOUND
fmt.Println(err.GetStatus())   // 404
fmt.Println(err.GetReason())   // "USER_NOT_FOUND"
fmt.Println(err.GetMessage())  // "User not found"
```

### HTTP Response Format

Errors are automatically converted to consistent JSON responses:

```json
{
  "status": 404,
  "code": 1001,
  "error": "USER_NOT_FOUND",
  "message": "User not found"
}
```

### Error Organization

Group related errors in domain-specific enums:

```protobuf
enum AuthError {
  option (sphere.errors.default_status) = 401;
  
  AUTH_ERROR_INVALID_TOKEN = 2001 [(sphere.errors.options) = {
    message: "Invalid authentication token"
  }];
  AUTH_ERROR_TOKEN_EXPIRED = 2002 [(sphere.errors.options) = {
    message: "Authentication token has expired"
  }];
}

enum ValidationError {
  option (sphere.errors.default_status) = 400;
  
  VALIDATION_ERROR_REQUIRED_FIELD = 3001 [(sphere.errors.options) = {
    message: "Required field is missing"
  }];
  VALIDATION_ERROR_INVALID_FORMAT = 3002 [(sphere.errors.options) = {
    message: "Field format is invalid"
  }];
}
```

## sphere/options

The options package provides method-level metadata for custom routing and code generation patterns.

### Installation

Add to your `buf.yaml`:

```yaml
deps:
  - buf.build/go-sphere/options
```

### Features

- Method-level key-value metadata
- Support for different value types (boolean, string, number)
- Extensible extra data mapping
- Integration with custom code generators
- Multiple transport protocol support

### Basic Usage

```protobuf
syntax = "proto3";

import "sphere/options/options.proto";

service BotService {
  rpc HandleStart(StartRequest) returns (StartResponse) {
    option (sphere.options.options) = [
      {
        key: "command"
        text: "/start"
      },
      {
        key: "callback_query"  
        text: "start_.*"
      }
    ];
  }
  
  rpc HandleMenu(MenuRequest) returns (MenuResponse) {
    option (sphere.options.options) = [
      {
        key: "callback_query"
        text: "menu_.*"
      },
      {
        key: "admin_only"
        flag: true
      }
    ];
  }
}
```

### Value Types

The options support multiple value types:

```protobuf
option (sphere.options.options) = [
  {
    key: "feature_flag"
    flag: true                    // Boolean value
  },
  {
    key: "route_pattern"
    text: "/api/v1/users"         // String value
  },
  {
    key: "rate_limit"
    number: 100                   // Integer value
  },
  {
    key: "metadata"
    extra: {                      // Map of additional strings
      "priority": "high"
      "category": "admin"
    }
  }
];
```

### Integration with Code Generation

The options are used by `protoc-gen-route` to generate routing and adapter code:

```go
// Generated constants
const OperationBotServiceHandleStart = "/bot.v1.BotService/HandleStart"

// Generated metadata
var ExtraDataBotServiceHandleStart = map[string]string{
    "command": "/start",
    "callback_query": "start_.*",
}

// Generated helper functions
func GetExtraDataByOperation(operation string) map[string]string {
    switch operation {
    case OperationBotServiceHandleStart:
        return ExtraDataBotServiceHandleStart
    // ... other cases
    }
}
```

### Use Cases

#### Bot Command Routing

```protobuf
rpc HandleCommand(CommandRequest) returns (CommandResponse) {
  option (sphere.options.options) = [
    {
      key: "telegram_command"
      text: "/users"
    },
    {
      key: "slack_command"  
      text: "/list-users"
    }
  ];
}
```

#### Feature Flags

```protobuf
rpc ExperimentalFeature(FeatureRequest) returns (FeatureResponse) {
  option (sphere.options.options) = [
    {
      key: "experimental"
      flag: true
    },
    {
      key: "min_version"
      text: "v2.1.0"
    }
  ];
}
```

#### API Versioning

```protobuf
rpc GetUserV2(GetUserV2Request) returns (UserV2) {
  option (sphere.options.options) = [
    {
      key: "api_version"
      text: "v2"
    },
    {
      key: "deprecated_in"
      text: "v3"
    }
  ];
}
```

## Integration Patterns

### Combined Usage

These packages work together to provide comprehensive API definitions:

```protobuf
syntax = "proto3";

import "google/api/annotations.proto";
import "sphere/binding/binding.proto";
import "sphere/errors/errors.proto";
import "sphere/options/options.proto";

// Error definitions
enum UserServiceError {
  option (sphere.errors.default_status) = 500;
  
  USER_SERVICE_ERROR_NOT_FOUND = 1001 [(sphere.errors.options) = {
    status: 404
    message: "User not found"
  }];
}

// Service with HTTP and custom routing
service UserService {
  rpc GetUser(GetUserRequest) returns (User) {
    option (google.api.http) = {
      get: "/v1/users/{id}"
    };
    option (sphere.options.options) = [
      {
        key: "cache_ttl"
        number: 300
      }
    ];
  }
}

// Request with field binding
message GetUserRequest {
  int64 id = 1 [(sphere.binding.location) = BINDING_LOCATION_URI];
  repeated string fields = 2 [(sphere.binding.location) = BINDING_LOCATION_QUERY];
}
```

### Buf Configuration

Example `buf.gen.yaml` for using all packages:

```yaml
version: v2
managed:
  enabled: true
  disable:
    - file_option: go_package_prefix
      module: buf.build/go-sphere/binding
    - file_option: go_package_prefix  
      module: buf.build/go-sphere/errors
    - file_option: go_package_prefix
      module: buf.build/go-sphere/options
plugins:
  - local: protoc-gen-go
    out: api
    opt:
      - paths=source_relative
  - local: protoc-gen-sphere-binding
    out: api
    opt:
      - paths=source_relative
  - local: protoc-gen-sphere
    out: api
    opt:
      - paths=source_relative
  - local: protoc-gen-sphere-errors
    out: api
    opt:
      - paths=source_relative
  - local: protoc-gen-route
    out: api
    opt:
      - paths=source_relative
      - options_key=cache
```

## Best Practices

### Binding Annotations
1. **Be explicit**: Use binding annotations for clarity even when defaults would work
2. **Group logically**: Keep related query parameters together
3. **Validate early**: Use binding to ensure data comes from expected sources
4. **Security**: Never bind sensitive data from query parameters

### Error Definitions  
1. **Domain grouping**: Create separate error enums for different domains
2. **Meaningful codes**: Use descriptive error codes and reasons
3. **User-friendly messages**: Write clear messages for end users
4. **HTTP semantics**: Use appropriate HTTP status codes

### Options Usage
1. **Consistent keys**: Establish conventions for option key names
2. **Documentation**: Comment the purpose of custom options
3. **Validation**: Validate option values in generators
4. **Versioning**: Consider compatibility when changing options

## Migration Guide

### From Manual HTTP Handlers

If migrating from manual HTTP handlers:

1. **Extract bindings**: Identify how request data is currently bound
2. **Add annotations**: Use binding annotations to replicate the behavior
3. **Define errors**: Convert error responses to error enums
4. **Generate code**: Use protoc-gen-sphere to generate new handlers

### From Other Frameworks

When migrating from other RPC frameworks:

1. **Map concepts**: Identify equivalent patterns in your current framework
2. **Proto-first**: Rewrite service definitions in protobuf
3. **Incremental**: Migrate one service at a time
4. **Testing**: Ensure compatibility with existing clients

These proto packages form the foundation of Sphere's code generation approach, enabling consistent, type-safe APIs with minimal boilerplate code.
