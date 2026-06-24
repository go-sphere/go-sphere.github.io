---
title: Development Workflow
weight: 12
---

Sphere projects use the generated Makefile as the day-to-day workflow entrypoint. The CLI creates the project; `make`, `go`, `buf`, Wire, Swag, Docker, and project-local tools do the ongoing work.

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

## Tool Boundaries

- `sphere-cli`: create projects, list templates, rename modules, and optionally generate small skeleton files.
- `make`: run the project workflow.
- `buf`: manage proto dependencies and generation.
- `go`: run, test, build, and manage modules.
- `wire`: generate dependency injection code.
- `swag`: generate Swagger/OpenAPI files.
- Docker or CI/CD: build and deploy images when the template supports it.

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
sphere-cli create --name myproject --module github.com/user/myproject
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

Generated artifacts should be safe to clean and recreate. Keep handwritten business logic in `internal/service/**` and `internal/biz/**`, and keep API contracts in `proto/**`.

## Best Practices

- **Keep services thin**: Move complex business logic to `internal/biz/`
- **Test early**: Use generated Swagger UI to verify endpoints
- **Iterate fast**: Small changes → quick regeneration cycles
- **Version APIs**: Use proto package versioning for breaking changes
- **Keep tool ownership clear**: Add project workflow to the Makefile instead of expanding the CLI into a build or deployment platform
