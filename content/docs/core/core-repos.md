---
title: Core Repos
weight: 10
---

Sphere Core
- `sphere`: multi-server application template with Gin utilities, docs/debug servers, middleware, and helpers
  - Repo: https://github.com/go-sphere/sphere
  - Provides `ginx` helpers like `WithJson`, `AbortWithJsonError`, `DataResponse`, `ErrorResponse`
- `sphere-layout`: standard project layout and Makefile-driven workflow
  - Repo: https://github.com/go-sphere/sphere-layout
  - Includes buf configs, Swagger generation, TypeScript SDK generation, and examples of `sphere` usage

Use `sphere-cli create` to start from the `sphere-layout` template, then extend with your schema and APIs.

