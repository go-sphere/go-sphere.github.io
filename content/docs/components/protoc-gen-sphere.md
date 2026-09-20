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
| `omitempty` | Skip methods without a `google.api.http` rule instead of synthesizing a default `POST` route for them; a file whose services all lack a rule emits nothing | `true` |
| `omitempty_prefix` | Path prefix for synthesized default routes (`<prefix>/<fully.qualified.Service>/<Method>`, also used when a rule declares no path) | `""` |
| `fail_on_warn` | Treat generation warnings as errors (skipped client/bidirectional streams, ignored streaming `response_body`, invalid body declarations) | `false` |
| `template_file` | Custom Go text/template; empty uses the embedded default | `""` |
| `swagger_auth_header` | Comment injected as the authorization header in generated Swagger docs | `// @Param Authorization header string false "Bearer token"` |
| `router_type` | Router type | `github.com/go-sphere/httpx;Router` |
| `context_type` | Request context type | `github.com/go-sphere/httpx;Context` |
| `handler_type` | Type returned by each generated handler | `github.com/go-sphere/httpx;Handler` |
| `context_load_func` | Expression appended to the context value to obtain a `context.Context` | `.Context()` |
| `data_resp_type` | Success envelope; must support generics | `github.com/go-sphere/sphere/server/httpz;DataResponse` |
| `error_resp_type` | Error envelope | `github.com/go-sphere/sphere/server/httpz;ErrorResponse` |
| `server_handler_func` | Wrapper that adapts the generated handler to the response model; must support generics | `github.com/go-sphere/sphere/server/httpz;WithJson` |
| `stream_handler_func` | Wrapper for server-streaming SSE handlers; must support generics | `github.com/go-sphere/sphere/server/httpz;WithSSE` |
| `stream_type` | Generic stream type returned by the streaming prepare phase | `github.com/go-sphere/sphere/server/httpz;SSEStream` |

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
- A server-streaming method (`returns (stream Reply)`) generates an SSE handler with `httpz.WithSSE`; its service interface receives `send func(*Reply) error`.
- Client-streaming and bidirectional methods are skipped with a warning. On a server stream, `response_body` is ignored because each event carries the whole reply message.
- GET, HEAD, DELETE, and OPTIONS never emit `BindJSON`. If a proto still declares `body` on those methods, the plugin warns (or fails with `fail_on_warn`) and generates the handler without a body bind.
- Custom verbs: a `:verb` suffix on the last path segment (e.g. `post: "/v1/reports:generate"`) stays a literal part of the route; only a segment-leading `:` is treated as a path parameter.
- `QUERY`, `URI` and `HEADER` bind a value from a single string token, so only scalar fields (and the well-known scalar wrappers `Timestamp`, `Duration` and `wrapperspb.*Value`) may use them. A `map`, `bytes` or arbitrary `message` field with one of these locations is a generation-time error; use `JSON` (or `FORM` for `bytes`/files) instead.
- Proto3 `optional` fields are documented as `required=false` in Swagger unless `(buf.validate.field).required = true` overrides it.
- Pair with [`protoc-gen-sphere-binding`](https://github.com/go-sphere/protoc-gen-sphere-binding) to inject binding tags into generated structs.
- Official templates run on `httpx/stdx`, the plain net/http adapter; Gin, Fiber, Echo, and Hertz adapters are also available. Changing `router_type` / `context_type` is how you retarget the generated code.
- See [Server Streaming](../guides/server-streaming) for the generated contract, wire format, and lifecycle rules.
