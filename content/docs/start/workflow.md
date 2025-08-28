---
title: Development Workflow
weight: 20
---

A practical, end-to-end flow you can use from day one.

1) Bootstrap
- Install CLI and plugins
- `sphere-cli create --name <name> --mod <module>`
- Inspect Makefile and buf configs

2) Model and Storage
- Define Ent schemas under `internal/pkg/database/ent/schema` (e.g., `User`)
- `make gen/db` to generate Ent code

3) API Contract
- Create `.proto` under `proto/` for shared messages and service RPCs
- Add `google.api.http` and `sphere.binding` annotations for routing and field binding
- `make gen/proto` to run protoc plugins

4) Server + Docs
- `make gen/docs` to generate Swagger
- Implement service methods in `internal/service/**`
- Use Ent client directly for CRUD; keep orchestration simple

5) Errors
- Define error enums in `.proto` with `sphere/errors` options
- Regenerate (`make gen/proto`) and use typed errors: `MyError_FOO.Join(err)`

6) Wire and Run
- `make gen/wire` to build the DI graph
- `make run` to start the server
- `make run/swag` to browse API docs

7) Client SDKs (optional)
- `make gen/dts` to produce TypeScript clients from Swagger

Iterate Fast
- Add new entities → regenerate db
- Add/adjust RPCs → regenerate proto/server/docs
- Keep service logic clean and focused

