---
title: Sphere CLI
weight: 10
---

Sphere CLI streamlines common project tasks: create projects, generate services, convert Ent schemas to proto, retag, and rename modules.

Install
- `go install github.com/go-sphere/sphere-cli@latest`

Usage
- `sphere-cli [command] [flags]`
- `sphere-cli [command] --help` for details

Key Commands
- `create`: bootstrap a new project
  - `sphere-cli create --name <project-name> [--module <go-module-name>]`
- `entproto`: convert Ent schemas into `.proto`
  - `--path`: ent schema dir; `--proto`: output dir
  - flags: `--all_fields_required`, `--auto_annotation`, `--enum_raw_type`, `--skip_unsupported`, `--time_proto_type`, `--uuid_proto_type`, `--unsupported_proto_type`, `--import_proto`
- `service proto`: scaffold a service `.proto`
  - `sphere-cli service proto --name <service-name> [--package <pkg>]`
- `service golang`: scaffold Go service implementation
  - `sphere-cli service golang --name <service> [--package <pkg>] [--mod <module>]`
- `retags` (deprecated): inject struct tags into generated `.pb.go` (use `protoc-gen-sphere-binding` instead)
- `rename`: rewrite Go module path across repo

When to Use
- First setup with `create`
- Evolve schema with `entproto`
- Bootstrap new APIs with `service`

