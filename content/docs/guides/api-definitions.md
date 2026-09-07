---
title: API Definitions
weight: 41
---

Sphere allows you to define HTTP interfaces for your services using standard Protobuf definitions and `google.api.http` annotations. This document outlines the rules and conventions for mapping your gRPC methods to RESTful HTTP endpoints.

## Getting Started: A Basic Example

To expose a gRPC method as an HTTP endpoint, you need to define it in a `.proto` file and add an HTTP annotation.

Here is a basic example of a `TestService` that defines a simple `RunTest` method, exposed as an HTTP `POST` request.

```protobuf
syntax = "proto3";

package your.service.v1;

import "google/api/annotations.proto";
import "sphere/binding/binding.proto";

// The Test service definition.
service TestService {
  // RunTest method
  rpc RunTest(RunTestRequest) returns (RunTestResponse) {
    option (google.api.http) = {
      post: "/v1/test/{path_test1}"
      body: "*"
    };
  }
}

// The request message for the RunTest RPC.
message RunTestRequest {
  // URI path parameter
  string path_test1 = 1 [(sphere.binding.location) = BINDING_LOCATION_URI];
  // Request body field
  string field_test1 = 2;
  // Query parameter
  string query_test1 = 3 [(sphere.binding.location) = BINDING_LOCATION_QUERY];
}

// The response message for the RunTest RPC.
message RunTestResponse {
  string field_test1 = 1;
  string query_test1 = 3;
}
```

### Key Components

1. **`import "google/api/annotations.proto";`**: This import is required to use HTTP annotations.
2. **`import "sphere/binding/binding.proto";`**: This import is required for binding annotations.
3. **`service TestService { ... }`**: Defines your gRPC service.
4. **`rpc RunTest(...) returns (...)`**: Defines a method within the service.
5. **`option (google.api.http) = { ... };`**: This is the core of the HTTP mapping.
   - **`post: "/v1/test/{path_test1}"`**: This specifies that the `RunTest` method should be exposed as an HTTP `POST`
   - **`body: "*"`**: Indicates that the entire request message (except URI params) should be sent as JSON body
6. **`[(sphere.binding.location) = ...]`**: This annotation specifies where the field should be bound from in the HTTP request.

Sphere uses these definitions to automatically generate server-side stubs and routing information.

## URL Path Mapping

Sphere converts gRPC-Gateway style URL paths from your `.proto` definitions into `httpx` routes (`:param` and `*wildcard`). Official templates serve those routes through a Gin adapter; Fiber, Echo, and Hertz adapters accept the same patterns.

The following table shows how Protobuf URL paths are translated:

| Protobuf Path Template                           | Generated httpx Route                       |
|--------------------------------------------------|---------------------------------------------|
| `/users/{user_id}`                               | `/users/:user_id`                           |
| `/users/{user_id}/posts/{post_id}`               | `/users/:user_id/posts/:post_id`            |
| `/files/{file_path=**}`                          | `/files/*file_path`                         |
| `/files/{name=*}`                                | `/files/:name`                              |
| `/static/{path=assets/*}`                        | `/static/assets/:path`                      |
| `/static/{path=assets/**}`                       | `/static/assets/*path`                      |
| `/projects/{project_id}/locations/{location=**}` | `/projects/:project_id/locations/*location` |
| `/api/{version=v1}/users`                        | `/api/v1/users`                             |
| `/users/{user_id}/posts/{post_id=drafts}`        | `/users/:user_id/posts/drafts`              |
| `/docs/{path=guides/**}`                         | `/docs/guides/*path`                        |

## HTTP Methods and Body Binding

### Request Binding Behavior

Different HTTP methods have different default binding behaviors:

#### GET, HEAD, DELETE, and OPTIONS Requests
- Path parameters must be top-level fields marked with `BINDING_LOCATION_URI`. Nested templates such as `{user.id}` are not supported.
- Remaining fields must be marked `BINDING_LOCATION_QUERY`, `BINDING_LOCATION_HEADER`, or `BINDING_LOCATION_FORM` (or the message `default_location`). Unmarked fields fail generation.
- No request body is expected. JSON-bound fields on these methods fail generation.
- If `google.api.http` still declares `body`, `protoc-gen-sphere` warns (or fails with `fail_on_warn`) and generates the handler **without** `BindJSON`

#### POST, PUT, and PATCH Requests
- Path parameters are bound to fields marked with `BINDING_LOCATION_URI`
- By default, all other fields are expected in the JSON request body
- You can override this with explicit binding locations

### Field Binding Locations

Use the `sphere.binding.location` annotation to control where each field is bound from:

```protobuf
message GetUserRequest {
  // URI path parameter
  int64 user_id = 1 [(sphere.binding.location) = BINDING_LOCATION_URI];
  // Query parameter
  repeated string fields = 2 [(sphere.binding.location) = BINDING_LOCATION_QUERY];
  // Header value
  string auth_token = 3 [(sphere.binding.location) = BINDING_LOCATION_HEADER];
}

message UpdateUserRequest {
  // URI path parameter
  int64 user_id = 1 [(sphere.binding.location) = BINDING_LOCATION_URI];
  // JSON body (default for POST/PUT/PATCH)
  User user = 2;
}
```

Available binding locations:
- `BINDING_LOCATION_URI`: Path parameters
- `BINDING_LOCATION_QUERY`: Query string parameters
- `BINDING_LOCATION_JSON`: JSON request body (default for non-GET methods)
- `BINDING_LOCATION_HEADER`: HTTP headers
- `BINDING_LOCATION_FORM`: Form data

## Advanced Binding Configuration

### Message-Level Defaults

You can set default binding behavior for entire messages:

```protobuf
message SearchUsersRequest {
  option (sphere.binding.default_location) = BINDING_LOCATION_QUERY;
  option (sphere.binding.default_auto_tags) = "form";
  
  string name = 1;        // Will be bound from query by default
  int32 age = 2;          // Will be bound from query by default
  string email = 3;       // Will be bound from query by default
}
```

### Custom Struct Tags

Add custom Go struct tags using the `auto_tags` annotation:

```protobuf
message DatabaseModel {
  option (sphere.binding.default_auto_tags) = "db";
  
  string name = 1;     // Generated: `db:"name" json:"name"`
  string email = 2;    // Generated: `db:"email" json:"email"`
}
```

## Request Body Patterns

### Full Body Binding
Most common for create/update operations:

```protobuf
rpc CreateUser(CreateUserRequest) returns (User) {
  option (google.api.http) = {
    post: "/v1/users"
    body: "*"  // Entire request message as JSON body
  };
}
```

### Specific Field as Body
When you want only one field as the body:

```protobuf
rpc UpdateUserProfile(UpdateUserProfileRequest) returns (User) {
  option (google.api.http) = {
    put: "/v1/users/{user_id}/profile"
    body: "profile"  // Only the 'profile' field as JSON body
  };
}

message UpdateUserProfileRequest {
  int64 user_id = 1 [(sphere.binding.location) = BINDING_LOCATION_URI];
  UserProfile profile = 2;  // This becomes the JSON body
}
```

## Response Body Patterns

### Default Response
By default, the entire response message is returned as JSON:

```protobuf
rpc GetUser(GetUserRequest) returns (User) {
  option (google.api.http) = { get: "/v1/users/{id}" };
}
// Returns: {"id": 1, "name": "John", "email": "john@example.com"}
```

### Specific Field as Response Body
You can return only a specific field:

```protobuf
rpc GetUserName(GetUserNameRequest) returns (GetUserNameResponse) {
  option (google.api.http) = {
    get: "/v1/users/{id}/name"
    response_body: "name"
  };
}

message GetUserNameResponse {
  string name = 1;  // Only this field is returned
}
// Returns: "John Doe" (just the string, not wrapped in JSON object)
```

## Server-Streaming Responses

Add `stream` to the reply type to expose a server-streaming RPC as SSE:

```protobuf
rpc Watch(WatchRequest) returns (stream WatchResponse) {
  option (google.api.http) = { get: "/v1/watch/{topic}" };
}
```

`protoc-gen-sphere` generates a push-style service method:

```go
Watch(context.Context, *WatchRequest, func(*WatchResponse) error) error
```

Request fields use the same URI, query, header, form, and JSON binding rules as
unary methods. Each `send` becomes one SSE JSON event. The default runtime ends
a successful stream with a `done` event and a committed failure with an `error`
event. See [Server Streaming](server-streaming) for implementation, wire format,
resume, and operational guidance.

Only server-streaming is supported by the HTTP generator. Client-streaming and
bidirectional RPCs are skipped with a warning. Do not use `response_body` on a
stream: the generator warns and sends each whole reply message.

## Best Practices

1. **Use meaningful field names**: Field names become tag values, so use clear, descriptive names
2. **Choose appropriate binding locations**:
   - `BINDING_LOCATION_URI`: For resource identifiers in the path
   - `BINDING_LOCATION_QUERY`: For optional filters and pagination
   - `BINDING_LOCATION_JSON`: For complex data structures and create/update operations
3. **Be consistent with HTTP method semantics**:
   - GET: Retrieve data (no body, use query params for filters)
   - POST: Create new resources (use body for data)
   - PUT: Replace entire resources (use body for new data)
   - PATCH: Partial updates (use body for changes)
   - DELETE: Remove resources (no body, use path params for ID)

4. **Avoid overly broad wildcards** in paths to prevent ambiguous routing
5. **Prefer explicit body field** (`body: "fieldName"`) when payloads are nested
6. **Prefer not to use `oneof`** in HTTP-exposed request/response messages. Tags land on wrapper structs (`Message_Field`); generated handlers bind the parent request, so QUERY/URI/HEADER oneof members are not filled by `BindQuery`/`BindURI`/`BindHeader`. JSON codecs also handle oneof awkwardly on the wire.
7. **Use only server response streams for SSE**. Model uploads and bidirectional conversations as unary HTTP operations, WebSockets, or a separate transport.
8. **Treat stream termination as part of the contract**. Clients should understand `done`, `error`, and an interrupted connection without either terminal event.


## Integration with buf

Add the required dependencies to your `buf.yaml`:

```yaml
version: v2
deps:
  - buf.build/googleapis/googleapis
  - buf.build/go-sphere/binding
```

Configure HTTP generation in `buf.gen.yaml`:

```yaml
version: v2
managed:
  enabled: true
plugins:
  - local: protoc-gen-sphere
    out: api
    opt:
      - paths=source_relative
```

Run `protoc-gen-sphere-binding` from a second template after `protoc-gen-go` has written the structs. Official templates use `buf.binding.yaml`:

```yaml
version: v2
managed:
  enabled: true
plugins:
  - local: protoc-gen-sphere-binding
    out: api
    opt:
      - paths=source_relative
      - out=api
```
