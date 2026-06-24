---
title: Protocol & Code Generation
weight: 23
---

Sphere follows a "protocol-first" approach where you define API contracts once in Protobuf and generate repeatable transport glue, binding metadata, error types, and documentation from those definitions. This keeps service boundaries consistent without hiding the underlying Go tools.

## Core Philosophy

The fundamental principle is: **Define once, generate the plumbing**.

Instead of writing HTTP handlers, request/response structs, validation code, and documentation separately, you:

1. **Define services and messages** in `.proto` files
2. **Annotate with HTTP mappings** using `google.api.http`
3. **Configure field binding** with Sphere binding options
4. **Generate transport and documentation glue** using protoc plugins

This approach provides:
- **Consistency**: All layers use the same contracts
- **Type Safety**: Compile-time guarantees across the stack
- **Documentation**: API docs generated from source of truth
- **Client SDKs**: Optional clients generated from OpenAPI or other dedicated tools
- **Reduced Boilerplate**: No manual HTTP handler writing

## Protocol as Contract

Your `.proto` files serve as the authoritative definition of:
- **Data structures** (messages)
- **API operations** (services and methods)
- **Error conditions** (enums with metadata)
- **HTTP mapping** (via annotations)
- **Field constraints** (via validation rules)

## Code Generation Pipeline

### Generator Chain

The code generation happens in a specific order:

1. **protoc-gen-go**: Generate base Go types
2. **[protoc-gen-sphere-binding](https://github.com/go-sphere/protoc-gen-sphere-binding)**: Add struct tags for binding
3. **[protoc-gen-sphere](https://github.com/go-sphere/protoc-gen-sphere)**: Generate HTTP handlers and routing
4. **[protoc-gen-sphere-errors](https://github.com/go-sphere/protoc-gen-sphere-errors)**: Generate error types and handling
5. **[protoc-gen-route](https://github.com/go-sphere/protoc-gen-route)**: Generate custom routing (optional)

### What Gets Generated

From your proto definitions, you automatically get:

**Server-side Code:**
- Service interfaces to implement
- HTTP handlers with proper routing
- Request binding with validation
- Response marshaling with proper headers
- Error handling with consistent formatting

**Client-side Code:**
- OpenAPI/Swagger documentation
- TypeScript SDKs (optional)
- Go client stubs (if needed)
- Validation schemas for frontend use

**Developer Tools:**
- Interactive documentation via Swagger UI
- API testing endpoints
- Type definitions for IDE support

## Benefits of This Approach

### Type Safety
- Compile-time verification of API contracts
- No runtime surprises from mismatched types
- Automatic validation of required fields
- IDE support with autocomplete and error checking

### Consistency
- Single source of truth for API definitions
- Consistent naming across all generated code
- Uniform error handling patterns
- Standardized HTTP response formats

### Developer Experience
- Faster iteration cycles
- Less boilerplate code to maintain
- Clear separation of concerns
- Automatic documentation updates

### Scalability
- Easy to add new services and methods
- Version management built-in
- Multiple output targets from one definition
- Team coordination through shared contracts

## Protocol Organization

### Recommended Structure
```
proto/
├── shared/v1/           # Common messages
│   ├── user.proto
│   └── common.proto
├── api/v1/              # Service definitions
│   ├── user_service.proto
│   └── auth_service.proto
└── errors/v1/           # Error definitions
    ├── user_errors.proto
    └── common_errors.proto
```

### Versioning Strategy
- Use explicit version packages (`v1`, `v2`)
- Keep shared types separate from services
- Plan for backwards compatibility
- Document breaking changes clearly

## Best Practices

### Proto Design
1. **Clear naming**: Use descriptive, consistent names
2. **Proper grouping**: Organize by domain and version
3. **Forward compatibility**: Design for future evolution
4. **Documentation**: Comment services, methods, and fields

### Code Generation
1. **Frequent regeneration**: Update generated code early and often
2. **Don't edit generated files**: All changes go in `.proto` files
3. **Version control**: Commit both `.proto` and generated files
4. **Automation**: Integrate generation into build process

## Related Guides

For detailed information on:
- **Defining HTTP APIs**: See [API Definitions](../guides/api-definitions)
- **Error handling**: See [Error Handling](../guides/error-handling)
- **Proto packages**: See [Proto Packages & Runtime](proto-packages-and-runtime)
