---
title: Core
weight: 10
---

Sphere Core repositories and what they provide.

Core Repos
- `sphere`: multi-server application template with Gin utilities, docs/debug servers, middleware, and helpers
  - Repo: https://github.com/go-sphere/sphere
  - Provides `ginx` helpers like `WithJson`, `AbortWithJsonError`, `DataResponse`, `ErrorResponse`
- `sphere-layout`: standard project layout and Makefile-driven workflow
  - Repo: https://github.com/go-sphere/sphere-layout
  - Includes buf configs, Swagger generation, TypeScript SDK generation, and examples of `sphere` usage

Getting Started
- Install CLI: `go install github.com/go-sphere/sphere-cli@latest`
- Create a project from the layout template: `sphere-cli create --name <project-name> [--module <go-module-name>]`
- Extend with your schema and APIs; generate servers/bindings with the Sphere generators in this section
