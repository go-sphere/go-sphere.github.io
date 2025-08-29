---
title: API Definitions
---

Sphere exposes RESTful HTTP endpoints using Protobuf (`proto3`) and `google.api.http` annotations. Binding options from `sphere/binding/binding.proto` control where fields are read from: path, query, or JSON body.

Minimal Example
```protobuf
syntax = "proto3";
package api.v1;
import "google/api/annotations.proto";
import "sphere/binding/binding.proto";

service UserService {
  rpc GetUser(GetUserRequest) returns (User) {
    option (google.api.http) = { get: "/v1/users/{id}" };
  }
}

message GetUserRequest {
  int64 id = 1 [(sphere.binding.location) = BINDING_LOCATION_URI];
}

message User {
  int64 id = 1;
  string name = 2;
}
```

Path Mapping Rules
- `/users/{user_id}` → `/users/:user_id`
- `/files/{file_path=**}` → `/files/*file_path`
- `/static/{path=assets/*}` → `/static/assets/:path`
- Nested placeholders like `/v1/users/{user.id}` normalize to `/v1/users/:user_id`

HTTP and Field Binding
- GET/DELETE: non-path fields become query params by default
- POST/PUT/PATCH: non-path fields come from JSON body by default (with `body: "*"`)
- To force query binding, add `(sphere.binding.location)=BINDING_LOCATION_QUERY` on fields
- To bind a single field as body: set `body: "fieldName"`; for response, use `response_body: "fieldName"`

Field Tags and Defaults
- `(sphere.binding.tags)`: inject custom tags (e.g., `json:\"userId\"`)
- `(sphere.binding.default_auto_tags)`: a default tag key for all fields in a message (e.g., `db`)

Practical Tips
- Avoid overly broad wildcards in paths to prevent ambiguous routing
- Prefer explicit body field (`body: "item"`) when payloads are nested
- Avoid `oneof` in exposed HTTP request/response messages due to JSON codec limitations

See Also
- Components → Generators → protoc-gen-sphere for generating HTTP handlers
- Components → Generators → protoc-gen-sphere-binding for tag injection
