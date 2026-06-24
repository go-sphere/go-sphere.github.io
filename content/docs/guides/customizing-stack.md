---
title: Customizing the Stack
weight: 43
---

Sphere templates ship with a complete default stack, but the framework is not defined by that stack. You can replace parts of the template as long as the project keeps the same contracts: proto files remain the API source of truth, generated code stays regenerable, handwritten service logic stays outside generated outputs, and the Makefile remains the workflow entrypoint.

## What Is Safe To Replace

### HTTP Router

The default generated HTTP code is optimized for the standard Sphere template. The `httpx` package exists to keep the router boundary small. When replacing the router, preserve:

- generated service interfaces;
- request binding semantics from `sphere.binding`;
- response envelope behavior expected by your API;
- middleware order and error handling conventions.

### Persistence Layer

The default layout uses Ent. Other templates can use Bun or a different persistence layer.

When replacing persistence, update:

- `make gen/db`;
- any project-local generators under `cmd/tools/**`;
- repository or DAO packages under `internal/pkg/**`;
- service implementations that currently depend on Ent types.

The API contracts in `proto/**` do not need to depend on the ORM.

### Dependency Injection

Wire is a default choice, not a requirement. You can replace it with manual constructors or another composition approach.

If you replace Wire, update:

- `make gen/wire`;
- application bootstrapping under `cmd/app`;
- provider sets or constructor packages under `internal/**`.

### Documentation and Clients

Swagger/OpenAPI and TypeScript client generation are template features. You can keep them, replace them, or add another OpenAPI pipeline.

The important part is that API documentation can be regenerated from the same source contracts.

### Deployment

Sphere does not own deployment orchestration. Templates can provide Docker targets, Compose files, or deployment scripts, but production deployment should stay in your CI/CD platform, Kubernetes tooling, or infrastructure workflow.

## Keep These Contracts Stable

- `proto/**` remains the contract source.
- `buf.yaml` and `buf.gen.yaml` describe proto dependencies and generation.
- generated directories can be removed and recreated.
- handwritten service logic remains outside generated code.
- `make init`, `make gen/*`, `make run`, `make build`, `make fmt`, and `make lint` stay discoverable.

If a customization keeps those contracts intact, it is aligned with Sphere's design.
