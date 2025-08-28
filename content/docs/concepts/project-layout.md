---
title: Project Layout
---

Sphere projects follow a pragmatic layout optimized for fast iteration while staying maintainable. The standard layout (`sphere-layout`) looks like this:

Top-Level
- `api/`: Go files generated from `.proto` definitions
- `assets/`: Static assets and templates
- `cmd/`: Application entry points (e.g., `cmd/app`), and developer tools (`cmd/tools`)
- `devops/`: Docker, compose files, CI/CD bits
- `internal/`: Private application code
- `proto/`: Protobuf source files
- `scripts/`: Helper scripts (e.g., Swagger TypeScript generation)
- `swagger/`: Generated OpenAPI/Swagger artifacts

Internal Breakdown
- `internal/config/`: Config definitions and loading
- `internal/pkg/`: Shared internal packages (DB, rendering, auth adapters)
- `internal/pkg/database/ent/`: Ent client and generated code
- `internal/server/`: HTTP server setup, middleware, docs/debug servers
- `internal/service/`: Service implementations that satisfy generated interfaces

Make Targets
- `gen/db`: Generate Ent code (and optionally related proto)
- `gen/proto`: Generate proto artifacts and run protoc plugins
- `gen/docs`: Generate Swagger docs from registered routers
- `gen/wire`: Generate dependency injection
- `gen/all`: Clean, then run end-to-end generation
- `gen/dts`: Generate TypeScript SDKs from Swagger
- `run`, `run/swag`: Run main server and Swagger server

When to Customize
- You can adapt directories to your team’s preferences, but keeping the above separation helps codegen and tooling work smoothly.

