Sphere
---

**Sphere** is a pragmatic Go backend toolkit centered around a clean monolithic template and a small toolchain that automates the boring parts: schema, API contracts, server stubs, Swagger, and even TypeScript clients. Start simple, scale when needed.

- Core: `sphere` 
- Template: `sphere-layout`
- Toolchain: 
    - `sphere-cli`
    - `protoc-gen-sphere`
    - `protoc-gen-route`
    - `protoc-gen-sphere-binding`
    - `protoc-gen-sphere-errors`

What you build:
- Define entities with Ent and APIs with Protobuf
- Generate Go handlers, Swagger, error types, and client SDKs
- Compose services with Gin + Wire; deploy as a single binary

[Get started in the docs](./docs/)
