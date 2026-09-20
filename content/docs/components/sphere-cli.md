---
title: sphere-cli
weight: 31
---

Sphere CLI (`sphere-cli`) is a small bootstrap tool for [`sphere`](https://github.com/go-sphere/sphere) projects. It creates projects from templates, records the exact upstream template revision, lists available templates, renames module paths, and provides lightweight service skeleton helpers.

It is intentionally not the primary build, deploy, or runtime orchestration tool. After a project is created, the generated Makefile, Buf, Go, Wire, Swag, Docker, and project-local tools own the day-to-day workflow.

## Installation

To install `sphere-cli`, ensure you have Go installed and run the following command:

```shell
go install github.com/go-sphere/sphere-cli@latest
```

## Interactive Mode

When a command is run **without its required flags and both stdin/stdout are terminals**, `sphere-cli` launches an interactive wizard:

- `sphere-cli create` walks through template selection (fetched from the remote catalog, with a built-in fallback), the project name, the Go module path, and optional steps (git init, dependency installation).
- `sphere-cli service proto` / `service golang` list the Ent schemas detected in the current project so the entity can be picked instead of typed, offer package/module defaults (module read from `go.mod`), show a preview of the generated code, and let you choose between printing to stdout and writing to a suggested path.
- `sphere-cli rename` reads the current module from `go.mod` so only the new module path has to be entered, then confirms before rewriting imports.

All commands keep their flags for scripted use. When stdout is not a terminal (CI, pipes), the wizards are skipped automatically and the required flags apply.

## Usage

The general syntax for `sphere-cli` is:

```shell
sphere-cli [command] [flags]
```

For detailed information on any command, you can use the `--help` flag:

```shell
sphere-cli [command] --help
```

## Commands

Here is an overview of the available commands.

## Scope

`sphere-cli` is responsible for:

- creating projects from official or custom templates;
- recording a Git-backed layout's exact source revision in `.sphere/layout.lock.json`;
- listing available templates;
- renaming Go module paths;
- generating small service skeletons when useful.

`sphere-cli` is not responsible for:

- building binaries;
- running test and lint workflows;
- generating every project artifact;
- managing Docker or Kubernetes deployment;
- replacing `make`, `buf`, `go`, `wire`, `swag`, or Docker.

### `create`

Initializes a new Sphere project with a default template.

**Usage:**
```shell
sphere-cli create --name <project-name> [--module <go-module-name>] [--layout <layout-name-or-uri>] [--no-git] [--no-deps]
```

**Flags:**
- `--name string`: Required in scripted mode. The name for the new Sphere project.
- `--module string`: (Optional) The Go module path for the project. Defaults to the project name when omitted.
- `--layout string`: (Optional) Official layout name or custom layout JSON URI.
- `--no-git`: Skip git repository initialization and the initial commit.
- `--no-deps`: Skip dependency installation (`make init` + `go mod tidy`).

Official layout names are `standard` (the default), `simple`, `bun`, and
`telegram`. Official and custom Git-backed layouts define `source`, `ref`, and
`mod`; legacy custom ZIP definitions using `uri`, `mod`, and `path` remain
compatible but cannot record an exact synchronization baseline.

**Example:**
```shell
sphere-cli create --name myproject --module github.com/myorg/myproject
```

This command creates a new project directory with the [`sphere-layout`](https://github.com/go-sphere/sphere-layout) template, including:
- Makefile for build automation
- buf configuration for protobuf management
- Standard directory structure
- Example configurations
- `.sphere/layout.lock.json` with the source repository, configured ref, and exact commit SHA

After creation, initialize and run the project through `make`:

```shell
cd myproject
make init
make run
```

### `create list`

Lists the available project templates.

```shell
sphere-cli create list
```

### Updating a generated project

There is intentionally no `sphere-cli update` command. An AI agent should read
the project's `.sphere/layout.json` ownership rules and
`.sphere/layout.lock.json` baseline, obtain the base and target upstream
snapshots, normalize their module paths, and perform a three-way merge. It must
leave project-owned files alone, regenerate generated files, and advance the
lock only after formatting, tests, lint, and build all pass. The normative
protocol is `docs/LAYOUT_CONTRACT.md` in `sphere-layout`.

For an older project without a lock, the agent should match the initial Git
tree against upstream history. If that does not identify one commit uniquely,
the user must choose the baseline; the agent must not guess.

### `service`

Generates small service skeletons, including Protobuf definitions and Go service implementations. This is a convenience helper, not the main repeatable generation pipeline.

For normal project regeneration, use the Makefile:

```shell
make gen/db
make gen/proto
make gen/docs
make gen/wire
```

This command has two subcommands: `proto` and `golang`.

#### `service proto`

Generates a `.proto` file for a new service.

**Usage:**
```shell
sphere-cli service proto [--name <service-name>] [--package <package-name>] [--out <file>]
```

**Flags:**
- `--name string`: Required in scripted mode. The name of the service.
- `--package string`: The package name for the generated `.proto` file (default: `dash.v1`).
- `--out string`: Write the generated code to this file instead of stdout.

The generated proto references `entpb.<Entity>` messages, so the entity must already exist as an Ent schema and be annotated for `entproto` generation. Run the command inside the project so `--name` is matched against the real schema types.

**Example:**
```shell
sphere-cli service proto --name UserService --package api.v1
```

#### `service golang`

Generates the Go implementation for a service from its definition.

**Usage:**
```shell
sphere-cli service golang [--name <service-name>] [--package <package-name>] [--mod <go-module-path>] [--out <file>]
```

**Flags:**
- `--name string`: Required in scripted mode. The name of the service.
- `--package string`: The package name for the generated Go code (default: `dash.v1`).
- `--mod string`: The Go module path for the generated code (default: `github.com/go-sphere/sphere-layout`).
- `--out string`: Write the generated code to this file instead of stdout.

The generated skeleton calls APIs produced by the project's own generators (`entbind.CreateXxx`, `ent.Xxx.Create`, `s.render.Xxx`), so the entity must exist as an Ent schema and `make gen/proto` must have run at least once. Run the command inside the project directory so the name is matched against the real schema types; outside a project, pass multi-word entities in separated form (`key_value_store`).

**Example:**
```shell
sphere-cli service golang --name UserService --package api.v1 --mod github.com/myorg/myproject
```

### `rename`

Renames the Go module path across the entire repository. This is useful when you need to change the module path after creating a project.

**Usage:**
```shell
sphere-cli rename [--old <old-module-path>] --new <new-module-path> [--target <directory>]
```

**Flags:**
- `--old string`: Current Go module path. Optional: detected from `go.mod` in the target directory when omitted.
- `--new string`: (Required) The new Go module path.
- `--target string`: Root directory of the project to rename (default: `.`).

## Common Workflows

- Quick project setup and generation: ../../getting-started/quickstart
- Day-to-day development loop: ../../getting-started/workflow

## Best Practices

1. **Use consistent naming**: Follow Go naming conventions for services and packages
2. **Organize proto files**: Keep related messages and services in appropriate packages
3. **Version your APIs**: Use versioned packages (e.g., `api.v1`, `api.v2`) for backward compatibility
4. **Document your protos**: Add comments to your proto files for better generated documentation
5. **Run generation regularly**: Use `make gen/all` frequently to keep generated code up to date

## Common Issues and Solutions

### Module Path Mismatches

If you need to change the module path after project creation:

```shell
sphere-cli rename --old github.com/old/path --new github.com/new/path
```

### Missing Dependencies

If you encounter import errors, ensure all dependencies are properly installed:

```shell
make install  # Install required tools
make init     # Initialize dependencies
```

### Stale Generated Code

If generated code seems out of sync:

```shell
make clean    # Clean generated files
make gen/all  # Regenerate everything
```
