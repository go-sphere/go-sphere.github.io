---
title: Runtime & Server
---

Sphere’s HTTP layer builds on Gin with small helpers to keep handlers concise and error responses consistent.

Key Pieces
- `WithJson[T]`: wraps a handler returning `(T, error)` and serializes responses to `DataResponse[T]`
- `AbortWithJsonError`: normalizes errors to `ErrorResponse` with proper HTTP status
- `ShouldBindJSON/Uri/Query`: thin wrappers to bind request payloads
- Docs/Debug servers: auxiliary HTTP servers for Swagger UI and diagnostics

Typical Flow
- Protobuf + `protoc-gen-sphere` generate handler plumbing
- Your service method returns data or a typed error
- Wrapper turns it into a stable JSON envelope for clients

Extensibility
- Customize router/context/response types via `protoc-gen-sphere` flags
- Provide your own error parser to merge validation or domain-specific errors

