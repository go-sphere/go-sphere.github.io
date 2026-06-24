---
title: Design Philosophy
weight: 21
---

Sphere is designed as a "frameworkless framework" for Go services. The core idea is to keep the framework thin and make the integration points explicit.

Sphere does not try to replace the Go toolchain or mature ecosystem tools. It gives them a consistent project shape and a set of contracts so generated code, handwritten business logic, HTTP adapters, middleware, persistence, documentation, and deployment scripts can work together.

## Goals

- Use Protobuf as the source of truth for service contracts.
- Generate repeatable transport and binding glue from those contracts.
- Keep business logic in ordinary Go packages.
- Make common infrastructure adapters available without forcing one vendor or library.
- Keep generated projects understandable to Go developers who already know `go`, `make`, `buf`, Docker, and their chosen libraries.

## Non-Goals

Sphere is not intended to be a full-stack platform CLI.

It does not try to own:

- binary builds;
- test orchestration;
- deployment orchestration;
- database modeling for every database;
- Kubernetes resource management;
- observability backend selection;
- runtime service mesh behavior.

Those concerns should stay in focused tools and project templates. Sphere can provide examples, adapters, and conventions, but it should not hide the underlying tools.

## Default, Not Mandatory

Official templates choose a default stack to make new projects immediately useful:

- Buf for Protobuf dependency and generation management;
- Gin-compatible HTTP handlers;
- Ent or Bun for persistence;
- Wire for dependency injection;
- Swagger/OpenAPI and TypeScript client generation;
- Makefile targets for local workflow.

These are defaults. They are not the framework boundary. A Sphere project should be able to replace the ORM, HTTP router, deployment flow, or generated-client pipeline without rewriting its core service contracts.

## Where Sphere Adds Value

Sphere is valuable when it keeps integration boring:

- generated service interfaces are predictable;
- request binding metadata is explicit in proto files;
- errors are generated from typed proto enums;
- HTTP routers are reached through `httpx` adapters;
- infrastructure packages expose small interfaces and replaceable implementations;
- template Makefiles document the workflow instead of a hidden CLI doing everything.

This makes the project feel structured without making the runtime feel trapped inside a framework.
