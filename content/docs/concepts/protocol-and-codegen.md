---
title: Protocol & Codegen
weight: 20
aliases:
  - /docs/concepts/protobuf-protocol/
  - /docs/concepts/code-generation-engine/
---

Define APIs once in Protobuf, generate the rest.

Protobuf
- Define services and messages in `proto/**`
- Use `google.api.http` annotations for REST mapping
- Sphere binding options map path, query, and body fields

Generated Outputs
- Go server code and handlers
- Router and adapter glue
- Swagger/OpenAPI docs
- Optional client SDKs (e.g., TypeScript)

Generators
- `protoc-gen-sphere`: HTTP server code
- `protoc-gen-route`: router and adapter glue
- `protoc-gen-sphere-binding`: Go struct tags from binding options
- `protoc-gen-sphere-errors`: typed errors from enums

See Components → Generators for details and flags.

