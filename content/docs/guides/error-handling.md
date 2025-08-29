---
title: Error Handling
---

Sphere generates typed, consistent errors from `.proto` enums using `protoc-gen-sphere-errors`. Each enum value can carry an HTTP status, a machine-readable reason, and a default message.

Define Errors in Protobuf
```protobuf
syntax = "proto3";
package shared.v1;
import "sphere/errors/errors.proto";

enum TestError {
  option (sphere.errors.default_status) = 500;
  TEST_ERROR_INVALID_FIELD = 1000 [(sphere.errors.options) = {
    status: 400
    reason: "INVALID_ARGUMENT"
    message: "invalid field"
  }];
}
```

Generate and Use in Go
- Add the plugin to `buf.gen.yaml` (see Components → Generators → protoc-gen-sphere-errors)
- After generation, use the typed helpers in your services:

```go
if badInput {
    return sharedv1.TestError_TEST_ERROR_INVALID_FIELD.Join(fmt.Errorf("field empty"))
}
```

HTTP Response Mapping
- Sphere server utilities convert typed errors to JSON consistently:
- `status`: mapped to HTTP status; `code`: enum numeric value; `error`: reason; `message`: default/custom message

Why This Works Well
- Centralized error semantics defined in `.proto`
- Strongly typed usage in Go with helpful helpers
- Stable mapping to consistent HTTP error responses for clients

- See Also
- Components → Generators → protoc-gen-sphere-errors for setup and flags
- Guides → API Definitions for endpoint conventions
