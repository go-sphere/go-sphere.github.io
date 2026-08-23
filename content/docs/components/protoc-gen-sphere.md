---
title: protoc-gen-sphere
weight: 33
---

Generates HTTP server code from `.proto` service definitions, using `google.api.http` annotations and Sphere's `httpx` / `httpz` runtime.

Install
- `go install github.com/go-sphere/protoc-gen-sphere@latest`

Key Flags

Type flags use the `import/path;Identifier` format.

| Flag | Description | Default |
| --- | --- | --- |
| `version` | Print version and exit | `false` |
| `omitempty` | Skip files whose methods have no `google.api.http` option | `true` |
| `omitempty_prefix` | When set, `omitempty` only applies to files with this prefix | `""` |
| `fail_on_warn` | Treat generation warnings as errors (skipped streaming methods, GET/DELETE declaring a body, missing body) | `false` |
| `template_file` | Custom Go text/template; empty uses the embedded default | `""` |
| `swagger_auth_header` | Comment injected as the authorization header in generated Swagger docs | `// @Param Authorization header string false "Bearer token"` |
| `router_type` | Router type | `github.com/go-sphere/httpx;Router` |
| `context_type` | Request context type | `github.com/go-sphere/httpx;Context` |
| `handler_type` | Type returned by each generated handler | `github.com/go-sphere/httpx;Handler` |
| `context_load_func` | Expression appended to the context value to obtain a `context.Context` | `.Context()` |
| `data_resp_type` | Success envelope; must support generics | `github.com/go-sphere/sphere/server/httpz;DataResponse` |
| `error_resp_type` | Error envelope | `github.com/go-sphere/sphere/server/httpz;ErrorResponse` |
| `server_handler_func` | Wrapper that adapts the generated handler to the response model; must support generics | `github.com/go-sphere/sphere/server/httpz;WithJson` |

Request binding no longer uses standalone `parse_*_func` flags. Binding is performed through methods on the configured `context_type` (`BindJSON` / `BindQuery` / `BindURI` / `BindHeader` / `BindForm`). Customize binding by changing `context_type`.

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
- Generated handlers take and return `httpx` types. The default wrapper is `httpz.WithJson`.
- GET, HEAD, DELETE, and OPTIONS never emit `BindJSON`. If a proto still declares `body` on those methods, the plugin warns (or fails with `fail_on_warn`) and generates the handler without a body bind.
- Pair with [`protoc-gen-sphere-binding`](https://github.com/go-sphere/protoc-gen-sphere-binding) to inject binding tags into generated structs.
- Official templates wrap Gin, Fiber, Echo, or Hertz behind `httpx` adapters. Changing `router_type` / `context_type` is how you retarget the generated code.
