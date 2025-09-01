---
title: sphere-cli
weight: 31
---

Sphere CLI (`sphere-cli`) is a command-line tool designed to streamline the development of [Sphere](https://github.com/go-sphere/sphere) projects. It helps you create new projects, generate service code, manage Protobuf definitions, and perform other common development tasks.

## Installation

To install `sphere-cli`, ensure you have Go installed and run the following command:

```shell
go install github.com/go-sphere/sphere-cli@latest
```

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

### `create`

Initializes a new Sphere project with a default template.

**Usage:**
```shell
sphere-cli create --name <project-name> [--module <go-module-name>] [--layout <template-uri>]
```

**Flags:**
- `--name string`: (Required) The name for the new Sphere project.
- `--module string`: (Optional) The Go module path for the project.
- `--layout string`: (Optional) Custom template layout URI.

**Example:**
```shell
sphere-cli create --name myproject --module github.com/myorg/myproject
```

This command creates a new project directory with the [sphere-layout](https://github.com/go-sphere/sphere-layout) template, including:
- Makefile for build automation
- buf configuration for protobuf management
- Standard directory structure
- Example configurations

### `entproto`

Converts Ent schemas into Protobuf (`.proto`) definitions. This command reads your Ent schema files and generates corresponding `.proto` files.

**Usage:**
```shell
sphere-cli entproto [flags]
```

**Flags:**
- `--path string`: Path to the Ent schema directory (default: `./schema`).
- `--proto string`: Output directory for the generated `.proto` files (default: `./proto`).
- `--all_fields_required`: Treat all fields as required, ignoring `Optional()` (default: `true`).
- `--auto_annotation`: Automatically add `@entproto` annotations to the schema (default: `true`).
- `--enum_raw_type`: Use `string` as the type for enums in Protobuf (default: `true`).
- `--skip_unsupported`: Skip fields with types that are not supported (default: `true`).
- `--time_proto_type string`: Protobuf type to use for `time.Time` fields. Options: `int64`, `string`, `google.protobuf.Timestamp` (default: `int64`).
- `--uuid_proto_type string`: Protobuf type to use for `uuid.UUID` fields. Options: `string`, `bytes` (default: `string`).
- `--unsupported_proto_type string`: Protobuf type to use for unsupported fields. Options: `google.protobuf.Any`, `google.protobuf.Struct`, `bytes` (default: `google.protobuf.Any`).
- `--import_proto string`: Define external Protobuf imports. Format: `path1,package1,type1;path2,package2,type2` (default: `google/protobuf/any.proto,google.protobuf,Any;`).

**Example:**
```shell
sphere-cli entproto --path ./internal/pkg/database/ent/schema --proto ./proto/entpb
```

### `service`

Generates service code, including both Protobuf definitions and Go service implementations.

This command has two subcommands: `proto` and `golang`.

#### `service proto`

Generates a `.proto` file for a new service.

**Usage:**
```shell
sphere-cli service proto --name <service-name> [--package <package-name>]
```

**Flags:**
- `--name string`: (Required) The name of the service.
- `--package string`: The package name for the generated `.proto` file (default: `dash.v1`).

**Example:**
```shell
sphere-cli service proto --name UserService --package api.v1
```

#### `service golang`

Generates the Go implementation for a service from its definition.

**Usage:**
```shell
sphere-cli service golang --name <service-name> [--package <package-name>] [--mod <go-module-path>]
```

**Flags:**
- `--name string`: (Required) The name of the service.
- `--package string`: The package name for the generated Go code (default: `dash.v1`).
- `--mod string`: The Go module path for the generated code (default: `github.com/go-sphere/sphere-layout`).

**Example:**
```shell
sphere-cli service golang --name UserService --package api.v1 --mod github.com/myorg/myproject
```

### `rename`

Renames the Go module path across the entire repository. This is useful when you need to change the module path after creating a project.

**Usage:**
```shell
sphere-cli rename --old <old-module-path> --new <new-module-path>
```

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