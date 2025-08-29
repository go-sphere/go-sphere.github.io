---
title: Development Workflow
weight: 40
aliases:
  - /docs/start/workflow/
  - /docs/start/workflow
---

An end-to-end flow you can iterate on from day one. This page links to detailed steps to avoid duplication.

1) Bootstrap
- Install the tooling and create a project
- See: Docs → Getting Started → Installation, Creating Your First Project

2) Model and Storage
- Define Ent schemas under `internal/pkg/database/ent/schema/**`
- See also: Concepts → Project Layout

3) API Contract
- Create `.proto` for messages and RPCs under `proto/**`
- Guidance: Guides → API Definitions

4) Server + Docs
- Implement service methods in `internal/service/**`; keep orchestration simple
- See: Docs → Getting Started → Quickstart (Wire and Run)

5) Errors
- Define error enums in `.proto` with `sphere/errors` options
- Guidance: Guides → Error Handling

6) Wire and Run
- Generate DI, run the server, and browse Swagger
- See: Docs → Getting Started → Quickstart (Wire and Run)

7) Client SDKs (optional)
- Generate TypeScript clients from Swagger
- See: Docs → Getting Started → Quickstart (Generate TypeScript Client)

Iterate Fast
- Add entities → regenerate db
- Add/adjust RPCs → regenerate proto/server/docs
- Keep service logic clean and focused

