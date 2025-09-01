---
title: Quick Start
weight: 11
---

Sphere is a pragmatic Go backend toolkit centered on a clean monolithic template and a small toolchain that automates schema, API contracts, server stubs, Swagger, and even TypeScript clients.

What you build:
- Define entities with Ent and APIs with Protobuf  
- Generate Go handlers, Swagger, error types, and client SDKs
- Compose services with Gin + Wire; deploy as a single binary

## Prerequisites

- Go 1.24+
- Docker + Docker Compose  (Optional)
- Node.js + npm (for TypeScript clients)

## Install Tooling

**Install CLI:**
```bash
go install github.com/go-sphere/sphere-cli@latest
```

**Verify Installation:**
```bash
sphere-cli --version || sphere-cli -h
```

> **Note**: You don't need to manually install protoc plugins. After creating a project, run `make init` to automatically install all required dependencies including protoc plugins.

## Create Your First Project

### 1. Bootstrap Project

```bash
# Create a new project using the template
# Replace 'myproject' with your project name and update the Go module path
sphere-cli create --name myproject --mod github.com/yourusername/myproject
cd myproject
```

This generates a new project with a clean structure based on [sphere-layout](https://github.com/go-sphere/sphere-layout) and automatically installs all required protoc plugins and dependencies.

### 2. Define Database Schema (Ent)

Create your database entities in `internal/pkg/database/schema/`. For example, create `internal/pkg/database/schema/user.go`:

```go
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.String("email").Unique(),
		field.Int("age").Positive(),
	}
}
```

Generate the database code:
```bash
make gen/db
```

### 3. Define API Contracts (Protobuf)

Define shared messages in `proto/shared/v1/user.proto`:

```protobuf
syntax = "proto3";

package shared.v1;

message User {
  int64 id = 1;
  string name = 2;
  string email = 3;
  int32 age = 4;
}
```

Define your service API in `proto/api/v1/user.proto`:

```protobuf
syntax = "proto3";

package api.v1;

import "buf/validate/validate.proto";
import "google/api/annotations.proto";
import "shared/v1/user.proto";

service UserService {
  rpc GetUser(GetUserRequest) returns (shared.v1.User) {
    option (google.api.http) = {
      get: "/v1/users/{id}"
    };
  }
  
  rpc CreateUser(CreateUserRequest) returns (shared.v1.User) {
    option (google.api.http) = {
      post: "/v1/users"
      body: "*"
    };
  }
}

message GetUserRequest {
  int64 id = 1 [(buf.validate.field).int64.gt = 0];
}

message CreateUserRequest {
  string name = 1 [(buf.validate.field).required = true];
  string email = 2 [(buf.validate.field).required = true];
  int32 age = 3 [(buf.validate.field).int32.gt = 0];
}
```

### 4. Generate Code

Generate Go handlers, routers, and documentation:
```bash
make gen/proto
make gen/docs
```

### 5. Implement Business Logic

Implement service methods in `internal/service/**` using generated interfaces and your Ent client. Keep logic simple in the service layer.

### 6. Wire and Run

Generate dependency injection wiring:
```bash
make gen/wire
```

Start the server:
```bash
make run
```

Serve Swagger UI:
```bash
make run/swag
```

### 7. Generate TypeScript Client (Optional)

Generate a typed client from Swagger for frontend integration:
```bash
make gen/dts
```

## Development Workflow

1. **Model and Storage**: Define Ent schemas → `make gen/db`
2. **API Contract**: Create `.proto` files → `make gen/proto` → `make gen/docs`  
3. **Business Logic**: Implement service methods in `internal/service/**`
4. **Wire and Run**: `make gen/wire` → `make run`
5. **Client SDKs**: `make gen/dts` for TypeScript clients

## Next Steps

- [Project Structure](../concepts/project-structure) - Understand the directory layout
- [API Definitions](../guides/api-definitions) - Learn HTTP mapping and binding rules  
- [Error Handling](../guides/error-handling) - Define typed, standardized errors
- [Development Workflow](workflow) - Complete development process

