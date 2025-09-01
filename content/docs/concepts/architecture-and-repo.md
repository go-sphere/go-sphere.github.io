---
title: Architecture & Repository Strategy
weight: 25
---

Start monolith-first with clean boundaries, then split when needed. Use a single Go module with Makefile-driven code generation to keep contracts and runtime consistent.

Monolith-First
- Begin with one binary using the standard project structure
- Share contracts via `.proto` and generated code
- Extract services only when scaling requires it

Scaling Out
- Multiple binaries can reuse the same contracts and tooling
- Consistent errors, routing, and payloads via generators reduce drift

Repository
- Prefer a pragmatic single-repo setup at first
- Keep boundaries clear via directories and interfaces
- Automate with `gen/*` targets to maintain repeatability

