---
title: Protocol & Codegen
weight: 20
aliases:
  - /docs/concepts/protobuf-protocol/
  - /docs/concepts/code-generation-engine/
---

Sphere follows a "protocol-first" approach where you define your APIs once in Protobuf and generate everything else from those definitions. This ensures consistency across your entire stack and reduces boilerplate code.

## Core Philosophy

The fundamental principle is: **Define once, generate everywhere**.

Instead of writing HTTP handlers, request/response structs, validation code, and documentation separately, you:

1. **Define services and messages** in `.proto` files
2. **Annotate with HTTP mappings** using `google.api.http`
3. **Configure field binding** with Sphere binding options
4. **Generate everything else** using protoc plugins

This approach provides:
- **Consistency**: All layers use the same contracts
- **Type Safety**: Compile-time guarantees across the stack
- **Documentation**: API docs generated from source of truth
- **Client SDKs**: Automatically generated for multiple languages
- **Reduced Boilerplate**: No manual HTTP handler writing

## Protobuf as the Contract

### Service Definitions

Services define the operations your API supports:

```protobuf
syntax = "proto3";

package api.v1;

import "google/api/annotations.proto";
import "sphere/binding/binding.proto";

service UserService {
  rpc GetUser(GetUserRequest) returns (User) {
    option (google.api.http) = {
      get: "/v1/users/{id}"
    };
  }
  
  rpc CreateUser(CreateUserRequest) returns (User) {
    option (google.api.http) = {
      post: "/v1/users"
      body: "*"
    };
  }
  
  rpc UpdateUser(UpdateUserRequest) returns (User) {
    option (google.api.http) = {
      put: "/v1/users/{id}"
      body: "user"
    };
  }
  
  rpc DeleteUser(DeleteUserRequest) returns (google.protobuf.Empty) {
    option (google.api.http) = {
      delete: "/v1/users/{id}"
    };
  }
}
```

### Message Definitions

Messages define the data structures for requests and responses:

```protobuf
message User {
  int64 id = 1;
  string name = 2;
  string email = 3;
  google.protobuf.Timestamp created_at = 4;
  google.protobuf.Timestamp updated_at = 5;
}

message GetUserRequest {
  int64 id = 1 [(sphere.binding.location) = BINDING_LOCATION_URI];
}

message CreateUserRequest {
  string name = 1 [(buf.validate.field).string.min_len = 1];
  string email = 2 [(buf.validate.field).string.pattern = "^[^@]+@[^@]+$"];
}

message UpdateUserRequest {
  int64 id = 1 [(sphere.binding.location) = BINDING_LOCATION_URI];
  User user = 2;
}

message DeleteUserRequest {
  int64 id = 1 [(sphere.binding.location) = BINDING_LOCATION_URI];
}
```

### Error Definitions

Errors are defined as enums with rich metadata:

```protobuf
import "sphere/errors/errors.proto";

enum UserError {
  option (sphere.errors.default_status) = 500;
  
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

## HTTP Mapping with Annotations

### google.api.http Annotations

These annotations define how gRPC methods map to HTTP endpoints:

```protobuf
rpc GetUser(GetUserRequest) returns (User) {
  option (google.api.http) = {
    get: "/v1/users/{id}"
  };
}
```

Key components:
- **HTTP method**: `get`, `post`, `put`, `patch`, `delete`
- **URL path**: With placeholder variables in `{braces}`
- **Body mapping**: `body: "*"` (entire message) or `body: "field_name"`
- **Response body**: `response_body: "field_name"` for partial responses

### Sphere Binding Options

Sphere extends HTTP mapping with field-level binding control:

```protobuf
message SearchUsersRequest {
  string query = 1 [(sphere.binding.location) = BINDING_LOCATION_QUERY];
  int32 limit = 2 [(sphere.binding.location) = BINDING_LOCATION_QUERY];
  string auth_token = 3 [(sphere.binding.location) = BINDING_LOCATION_HEADER];
}
```

Available locations:
- `BINDING_LOCATION_URI`: Path parameters
- `BINDING_LOCATION_QUERY`: Query string parameters  
- `BINDING_LOCATION_BODY`: JSON request body
- `BINDING_LOCATION_HEADER`: HTTP headers
- `BINDING_LOCATION_FORM`: Form data

## Code Generation Pipeline

### Generator Chain

The code generation happens in a specific order:

1. **protoc-gen-go**: Generate base Go types
2. **protoc-gen-sphere-binding**: Add struct tags for binding
3. **protoc-gen-sphere**: Generate HTTP handlers
4. **protoc-gen-sphere-errors**: Generate error types
5. **protoc-gen-route**: Generate custom routing (optional)

### Generated Outputs

From your proto definitions, you get:

#### Go Server Code
```go
// Generated service interface
type UserServiceServer interface {
    GetUser(context.Context, *GetUserRequest) (*User, error)
    CreateUser(context.Context, *CreateUserRequest) (*User, error)
    UpdateUser(context.Context, *UpdateUserRequest) (*User, error)
    DeleteUser(context.Context, *DeleteUserRequest) (*emptypb.Empty, error)
}

// Generated HTTP handlers
func RegisterUserServiceServer(r gin.IRouter, srv UserServiceServer) {
    r.GET("/v1/users/:id", _UserService_GetUser_Handler(srv))
    r.POST("/v1/users", _UserService_CreateUser_Handler(srv))
    r.PUT("/v1/users/:id", _UserService_UpdateUser_Handler(srv))
    r.DELETE("/v1/users/:id", _UserService_DeleteUser_Handler(srv))
}
```

#### Request/Response Types with Binding Tags
```go
type GetUserRequest struct {
    Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"-" uri:"id"`
}

type CreateUserRequest struct {
    Name  string `protobuf:"bytes,1,opt,name=name,proto3" json:"name"`
    Email string `protobuf:"bytes,2,opt,name=email,proto3" json:"email"`
}
```

#### Error Types
```go
func (e UserError) Error() string { /* ... */ }
func (e UserError) GetStatus() int32 { /* ... */ }
func (e UserError) Join(errs ...error) error { /* ... */ }
```

#### OpenAPI/Swagger Documentation
```yaml
paths:
  /v1/users/{id}:
    get:
      summary: GetUser
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
            format: int64
      responses:
        200:
          description: A successful response.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/User'
```

## Benefits of This Approach

### Type Safety
- Compile-time verification of API contracts
- No runtime surprises from mismatched types
- Automatic validation of required fields

### Consistency
- Single source of truth for API definitions
- Consistent naming across all generated code
- Uniform error handling patterns

### Developer Experience
- Faster iteration cycles
- Less boilerplate code to maintain
- Clear separation of concerns

### Documentation
- Always up-to-date API documentation
- Comprehensive error code reference
- Type definitions for client developers

## Customization Options

### Custom Templates
You can provide custom templates for code generation:

```yaml
plugins:
  - local: protoc-gen-sphere
    out: api
    opt:
      - template_file=./templates/custom_server.tmpl
```

### Generator Configuration
Most generators support extensive configuration:

```yaml
plugins:
  - local: protoc-gen-sphere
    out: api
    opt:
      - router_type=github.com/gin-gonic/gin;IRouter
      - context_type=github.com/gin-gonic/gin;Context
      - swagger_auth_header=// @Security ApiKeyAuth
```

### Multiple Output Formats
Generate for different platforms from the same definitions:

```yaml
plugins:
  - local: protoc-gen-go
    out: api/go
  - local: protoc-gen-ts
    out: api/typescript  
  - local: protoc-gen-openapi
    out: api/openapi
```

## Best Practices

### Proto Organization
1. **Version your APIs**: Use `v1`, `v2` packages
2. **Group by domain**: Keep related services together
3. **Shared types**: Use `shared/` packages for common messages
4. **Clear naming**: Use descriptive service and method names

### HTTP Mapping
1. **RESTful patterns**: Follow REST conventions where possible
2. **Consistent paths**: Use consistent URL patterns
3. **Appropriate methods**: Match HTTP methods to operations
4. **Body optimization**: Use `body: "field"` for cleaner APIs

### Error Handling
1. **Comprehensive coverage**: Define errors for all failure modes
2. **Appropriate status codes**: Use correct HTTP status codes
3. **Clear messages**: Write user-friendly error messages
4. **Grouping**: Keep related errors in the same enum

### Documentation
1. **Comment everything**: Add comments to services, methods, and fields
2. **Examples**: Include example values in comments
3. **Deprecation**: Mark deprecated APIs clearly
4. **Versioning**: Document breaking changes between versions

## See Also

- [API Definitions Guide](../guides/api-definitions) - Detailed HTTP mapping rules
- [Error Handling Guide](../guides/error-handling) - Comprehensive error patterns
- [Code Generators](../components/generators) - Individual generator documentation
- [Project Structure](project-structure) - How generated code fits into projects

