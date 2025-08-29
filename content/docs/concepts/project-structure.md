---
title: Project Structure
weight: 10
aliases:
  - /docs/concepts/project-code-structure/
  - /docs/concepts/project-layout/
---

A pragmatic structure that keeps code generation, server code, and business logic cleanly separated while staying fast to iterate on.

Highlights
- Clear separation: `internal/service`, `internal/biz`, `internal/pkg`, `cmd/server`
- Data layer with Ent in `internal/pkg/database/ent`
- API contracts in `proto/**` with buf and Make targets

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
- Adapt directories to your team’s preferences, but keep the separation above so codegen and tooling remain smooth.

