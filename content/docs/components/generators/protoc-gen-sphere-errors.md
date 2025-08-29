---
title: protoc-gen-sphere-errors
weight: 40
---

Generates Go error types and helpers from `.proto` enums, standardizing error codes, statuses, and messages across services.

Install
- `go install github.com/go-sphere/protoc-gen-sphere-errors@latest`

Buf Example
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
    opt: paths=source_relative
```

Using in Protobuf
```protobuf
import "sphere/errors/errors.proto";
enum AdminError {
  option (sphere.errors.default_status) = 500;
  ADMIN_ERROR_CANNOT_DELETE_SELF = 1001 [(sphere.errors.options) = {
    status: 400
    message: "cannot delete self"
  }];
}
```

Using in Go
```go
return sharedv1.AdminError_ADMIN_ERROR_CANNOT_DELETE_SELF.Join(err)
```

Integration Notes
- Sphere’s Gin layer maps these to structured JSON with correct HTTP status
- Pair with a global error parser if you need to merge validation/notfound/custom errors

