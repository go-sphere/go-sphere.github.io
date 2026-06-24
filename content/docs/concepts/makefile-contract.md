---
title: Makefile Contract
weight: 22
---

Sphere templates use `make` as the day-to-day workflow entrypoint. The CLI creates a project; the project Makefile owns local development, generation, formatting, linting, running, and build commands.

This keeps Sphere aligned with common Go practices and leaves each tool visible.

## Standard Targets

Official templates should keep these target names stable where the capability exists:

| Target | Purpose |
| --- | --- |
| `make init` | Download modules, install required local tools, and prepare the project. |
| `make install` | Install local development tools such as Buf, Wire, Swag, protoc plugins, and linters. |
| `make gen/conf` | Generate example configuration. |
| `make gen/db` | Generate persistence code for the template's selected ORM. |
| `make gen/proto` | Run Protobuf generation and Sphere protoc plugins. |
| `make gen/docs` | Generate Swagger/OpenAPI documentation. |
| `make gen/wire` | Generate dependency injection wiring. |
| `make gen/all` | Clean and regenerate the full generated surface. |
| `make gen/dts` | Generate TypeScript clients when the template supports it. |
| `make run` | Run the application locally. |
| `make run/swag` | Serve generated Swagger UI when available. |
| `make build` | Build a binary for the current platform. |
| `make build/docker` | Build a Docker image when the template provides Docker support. |
| `make lint` | Run Go, Buf, and linter checks. |
| `make fmt` | Format Go, Protobuf, and import layout. |
| `make clean` | Remove generated artifacts and build outputs. |

Templates can add more targets, but these names should not change without a migration note.

## Tool Ownership

The Makefile should call mature tools directly:

- `go` for modules, tests, builds, and local execution;
- `buf` for proto dependency and generation workflows;
- `wire` for dependency injection generation;
- `swag` for Swagger documentation generation;
- `docker` or `docker buildx` for image builds;
- project-local Go tools under `cmd/tools/**` for template-specific generation.

Sphere CLI should not duplicate these responsibilities.

## Generated Code Boundaries

Generated outputs should be easy to clean and regenerate. In the default layout, this includes:

- `api/**`;
- `swagger/**`;
- generated Ent packages;
- Wire output;
- generated conversion, mapping, binding, and CRUD helpers.

Handwritten code should live outside generated paths, especially in:

- `proto/**` for API contracts;
- `internal/service/**` for generated interface implementations;
- `internal/biz/**` for business logic;
- `internal/pkg/**` for project infrastructure and adapters.

The Makefile is the place where these ownership rules become repeatable.
