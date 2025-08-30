---
title: protoc-gen-sphere
weight: 32
---

Generates HTTP server code from `.proto` service definitions, using `google.api.http` annotations and Sphere server utilities (Gin-based).

Install
- `go install github.com/go-sphere/protoc-gen-sphere@latest`

Key Flags
- `version`: print version
- `omitempty`: skip files without `google.api.http` (default true)
- `omitempty_prefix`: apply omitempty only to prefixed paths
- `template_file`: custom go text/template
- `swagger_auth_header`: auth header comment for Swagger
- `router_type`: router type (default `github.com/gin-gonic/gin;IRouter`)
- `context_type`: context type (default `github.com/gin-gonic/gin;Context`)
- `data_resp_type`: data model with generics (default `github.com/go-sphere/sphere/server/ginx;DataResponse`)
- `error_resp_type`: error model (default `github.com/go-sphere/sphere/server/ginx;ErrorResponse`)
- `server_handler_func`: wrapper (default `github.com/go-sphere/sphere/server/ginx;WithJson`)
- `parse_json_func`, `parse_uri_func`, `parse_form_func`: request parsing hooks

Buf Example
```yaml
version: v2
managed:
  enabled: true
  disable:
    - file_option: go_package_prefix
      module: buf.build/googleapis/googleapis
    - file_option: go_package_prefix
      module: buf.build/bufbuild/protovalidate
  override:
    - file_option: go_package_prefix
      value: github.com/go-sphere/sphere-layout/api
plugins:
  - local: protoc-gen-sphere
    out: api
    opt:
      - paths=source_relative
      - swagger_auth_header=// @Security ApiKeyAuth
```

Notes
- Works hand-in-hand with Sphere’s Gin helpers: `WithJson`, `ShouldBindJSON`, `ShouldBindUri`, `ShouldBindQuery`
- Pair with `protoc-gen-sphere-binding` to inject binding tags into generated structs

## See Also

- Guides: ../../guides/api-definitions
- Component: ./protoc-gen-sphere-binding
- Component: ./protoc-gen-sphere-errors
- Concepts: ../../concepts/protocol-and-codegen
