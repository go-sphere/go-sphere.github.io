---
title: protoc-gen-sphere-binding
weight: 30
---

Injects Go struct tags into generated message types based on Sphere binding annotations. This replaces ad-hoc tag injection and keeps tags defined at the Protobuf layer.

Install
- `go install github.com/go-sphere/protoc-gen-sphere-binding@latest`

Key Flags
- `version`: print version
- `out`: output directory for modified `.proto` (default `api`)

Buf Example (separate template)
```yaml
version: v2
managed:
  enabled: true
  override:
    - file_option: go_package_prefix
      value: github.com/go-sphere/sphere-layout/api
plugins:
  - local: protoc-gen-sphere-binding
    out: api
    opt:
      - paths=source_relative
      - out=api
```

Annotating Fields
```protobuf
import "sphere/binding/binding.proto";
message RunTestRequest {
  string path = 1 [(sphere.binding.location)=BINDING_LOCATION_URI];
  string q = 2 [(sphere.binding.location)=BINDING_LOCATION_QUERY];
  repeated int32 ids = 3 [
    (sphere.binding.location)=BINDING_LOCATION_QUERY,
    (sphere.binding.tags) = "sphere:\"ids\""
  ];
}
```

Notes
- Works alongside `protoc-gen-sphere` which binds via Gin helpers
- Prefer this over legacy `sphere-cli retags`

