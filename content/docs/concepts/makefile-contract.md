---
title: Makefile Contract
weight: 22
---

Sphere Go repositories use `make` as the day-to-day workflow entrypoint. The CLI creates a project; the project Makefile additionally owns local development, generation, running, and build commands.

This keeps Sphere aligned with common Go practices and leaves each tool visible.

## Standard Targets

All Go repositories provide the dependency, formatting, test, lint, and check targets below. Official templates keep the additional project targets stable where the capability exists:

| Target | Purpose |
| --- | --- |
| `make deps-update` | Update direct Go dependencies and tidy every module in the repository. |
| `make tidy` | Tidy every Go module without upgrading dependencies. |
| `make fmt` | Format Go, Protobuf, and import layout. |
| `make test` | Run the repository's complete test workflow. |
| `make lint` | Run non-mutating format, Go, Buf, and linter checks. |
| `make check` | Verify dependencies, lint, and run tests. |
| `make init` | Download modules, install required local tools, generate artifacts, and write `config.json` if it is missing. |
| `make install` | Install local development tools such as Buf, Wire, Swag, protoc plugins, and linters. |
| `make gen/conf` | Generate example configuration (`config_gen.json`) and write `config.json` if it is missing. |
| `make gen/db` | Generate persistence code for the template's selected ORM. |
| `make gen/proto` | Run Protobuf generation and Sphere protoc plugins. |
| `make gen/docs` | Generate Swagger/OpenAPI documentation. |
| `make gen/wire` | Generate dependency injection wiring. |
| `make gen/all` | Clean and regenerate the full generated surface. |
| `make gen/dts` | Generate TypeScript clients when the template supports it. |
| `make run` | Run the application locally. |
| `make run/race` | Run locally with the race detector. Requires cgo; not the default `run` target. |
| `make run/swag` | Serve generated Swagger UI when available. |
| `make build` | Build a binary for the current platform. |
| `make build/docker` | Build a Docker image when the template provides Docker support. |
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
