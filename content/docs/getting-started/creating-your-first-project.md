---
title: Creating Your First Project
weight: 30
---

Create a project:
- `sphere-cli create --name myproject --mod github.com/you/myproject`
- Bootstraps `sphere-layout` with Makefile, buf configs, and clean structure

Define the schema (Ent):
- Add schemas under `internal/pkg/database/ent/schema/*.go`
- Run `make gen/db`

Define the API (Protobuf):
- Place `.proto` under `proto/**`
- Use `google.api.http` annotations and Sphere binding options
- Run `make gen/proto` then `make gen/docs`

Implement business logic:
- Implement service methods in `internal/service/**`

Wire and run:
- `make gen/wire` then `make run`
- Optional Swagger UI: `make run/swag`

Optional: Generate TypeScript client
- `make gen/dts`

Next: Explore Core Concepts.

