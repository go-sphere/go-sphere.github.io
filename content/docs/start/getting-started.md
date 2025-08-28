---
title: Getting Started
weight: 10
---

This guide walks you through creating, generating, and running a new Sphere-based service end-to-end.

Prerequisites
- Go 1.24+
- Docker + Docker Compose
- Node.js + npm (for TypeScript clients)

Install CLI and Plugins
- `go install github.com/go-sphere/sphere-cli@latest`
- `go install github.com/go-sphere/protoc-gen-sphere@latest`
- `go install github.com/go-sphere/protoc-gen-route@latest`
- `go install github.com/go-sphere/protoc-gen-sphere-binding@latest`
- `go install github.com/go-sphere/protoc-gen-sphere-errors@latest`

Create a Project
- `sphere-cli create --name myproject --mod github.com/you/myproject`
- This bootstraps a project based on `sphere-layout` with Makefile, buf configs, and a clean directory structure.

Define the Schema (Ent)
- Add schemas under `internal/pkg/database/ent/schema/*.go`
- Run `make gen/db` to generate Ent code and optional proto messages (if configured in your layout).

Define the API (Protobuf)
- Place `.proto` files under `proto/**` (e.g., `proto/shared/v1/user.proto`, `proto/api/v1/user.proto`)
- Use `google.api.http` annotations to expose HTTP endpoints and Sphere binding options to control path/query/body mapping.
- Run `make gen/proto` then `make gen/docs` to generate Go handlers and Swagger.

Implement Business Logic
- Implement service methods in `internal/service/**` using generated interfaces and your Ent client.
- Keep logic simple in service layer; move complex orchestration into `internal/biz` if needed.

Wire and Run
- `make gen/wire` to generate DI wiring
- `make run` to start the server
- `make run/swag` to serve Swagger UI (from generated specs under `swagger/`)

Generate TypeScript Client (optional)
- `make gen/dts` to produce a typed client from Swagger for frontend integration

Next Steps
- Read Concepts → Project Layout for directory details
- See Guides → API Definitions for HTTP mapping and binding rules
- See Guides → Error Handling for typed, standardized errors

