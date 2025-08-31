---
title: Development Workflow
weight: 12
aliases:
  - /docs/start/workflow/
  - /docs/start/workflow
---

A typical day-to-day development cycle with Sphere follows this pattern:

## Core Development Loop

1. **Model Your Data**
   - Define Ent schemas in `internal/pkg/database/schema/`
   - Run `make gen/db` to generate database code

2. **Define Your API**
   - Create `.proto` files in `proto/`
   - Use `google.api.http` annotations for HTTP endpoints
   - Run `make gen/proto` to generate Go handlers

3. **Implement Business Logic**
   - Write service implementations in `internal/service/`
   - Keep logic focused and testable

4. **Wire and Test**
   - Run `make gen/wire` for dependency injection
   - Start with `make run`
   - Test endpoints via Swagger UI (`make run/swag`)

## Adding New Features

**For a new entity:**
```bash
# 1. Add Ent schema
# 2. Regenerate database code
make gen/db

# 3. Add proto messages and services
# 4. Regenerate API code
make gen/proto
make gen/docs

# 5. Implement service methods
# 6. Wire and run
make gen/wire
make run
```

**For a new project setup:**
```bash
# 1. Create project
sphere-cli create --name myproject --mod github.com/user/myproject
cd myproject

# 2. Initialize all dependencies and tools
make init

# 3. Start development
make run
```

**For API changes:**
```bash
# 1. Update .proto files
# 2. Regenerate
make gen/proto
make gen/docs

# 3. Update service implementations
# 4. Test
make run
```

## Generated Artifacts

- `make gen/db` → Ent client and schemas
- `make gen/proto` → Go handlers, types, and routing
- `make gen/docs` → OpenAPI/Swagger documentation  
- `make gen/wire` → Dependency injection wiring
- `make gen/dts` → TypeScript client SDKs

## Best Practices

- **Keep services thin**: Move complex business logic to `internal/biz/`
- **Test early**: Use generated Swagger UI to verify endpoints
- **Iterate fast**: Small changes → quick regeneration cycles
- **Version APIs**: Use proto package versioning for breaking changes

See [Quick Start](quickstart) for the complete setup process.

