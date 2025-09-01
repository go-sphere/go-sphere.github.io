---
title: Protocol & Codegen
weight: 23
---

Sphere follows a "protocol-first" approach where you define your APIs once in Protobuf and generate everything else from those definitions. This ensures consistency across your entire stack and reduces boilerplate code.

## Core Philosophy

The fundamental principle is: **Define once, generate everywhere**.

Instead of writing HTTP handlers, request/response structs, validation code, and documentation separately, you:

1. **Define services and messages** in `.proto` files
2. **Annotate with HTTP mappings** using `google.api.http`
3. **Configure field binding** with Sphere binding options
4. **Generate everything else** using protoc plugins

This approach provides:
- **Consistency**: All layers use the same contracts
- **Type Safety**: Compile-time guarantees across the stack
- **Documentation**: API docs generated from source of truth
- **Client SDKs**: Automatically generated for multiple languages
- **Reduced Boilerplate**: No manual HTTP handler writing

## Protocol as Contract

### Single Source of Truth

Your `.proto` files serve as the authoritative definition of:
- **Data structures** (messages)
- **API operations** (services and methods)
- **Error conditions** (enums with metadata)
- **HTTP mapping** (via annotations)
- **Field constraints** (via validation rules)

### Service Definitions

Services declare operations and HTTP mappings via annotations. See ../guides/api-definitions for end‑to‑end examples.

### Message Definitions

Messages define the request/response contracts. Field binding is controlled with sphere/binding options; examples live in ../guides/api-definitions.

### Error Definitions

Define error enums with HTTP status/reason/message using sphere/errors. Usage and patterns are covered in ../guides/error-handling.

## Code Generation Pipeline

### Generator Chain

The code generation happens in a specific order:

1. **protoc-gen-go**: Generate base Go types
2. **protoc-gen-sphere-binding**: Add struct tags for binding
3. **protoc-gen-sphere**: Generate HTTP handlers and routing
4. **protoc-gen-sphere-errors**: Generate error types and handling
5. **protoc-gen-route**: Generate custom routing (optional)

### What Gets Generated

From your proto definitions, you automatically get:

#### Server-side Code
- **Service interfaces** to implement
- **HTTP handlers** with proper routing
- **Request binding** with validation
- **Response marshaling** with proper headers
- **Error handling** with consistent formatting

#### Client-side Code
- **OpenAPI/Swagger documentation**
- **TypeScript SDKs** (optional)
- **Go client stubs** (if needed)
- **Validation schemas** for frontend use

#### Developer Tools
- **Interactive documentation** via Swagger UI
- **API testing** endpoints
- **Type definitions** for IDE support

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

## Evolution and Maintenance

### Adding New Features
1. Define new messages/services in `.proto`
2. Run code generation
3. Implement service methods
4. Tests and documentation are automatically updated

### API Versioning
- Create new version packages for breaking changes
- Maintain multiple versions simultaneously
- Gradual migration paths for clients
- Automatic deprecation warnings

### Team Collaboration
- Proto files as API contracts between teams
- Code review focuses on API design
- Generated code handles implementation details
- Consistent patterns across all services

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

### Error Handling
1. **Comprehensive coverage**: Define errors for all failure modes
2. **Appropriate status codes**: Use correct HTTP status codes
3. **Clear messages**: Write user-friendly error messages
4. **Structured data**: Include relevant context in errors

## See Also

- [API Definitions Guide](../guides/api-definitions) - Detailed HTTP mapping rules and examples
- [Error Handling Guide](../guides/error-handling) - Comprehensive error patterns and implementation
- [Components Overview](../components/) - Individual generator documentation and configuration
- [Project Structure](project-structure) - How generated code fits into the overall project layout
