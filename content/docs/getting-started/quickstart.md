---
title: Quick Start
weight: 11
aliases:
  - /docs/start/getting-started/
  - /docs/start/getting-started
  - /docs/getting-started/quickstart
  - /docs/getting-started/introduction
  - /docs/getting-started/creating-your-first-project
---

Sphere is a pragmatic Go backend toolkit centered on a clean monolithic template and a small toolchain that automates schema, API contracts, server stubs, Swagger, and even TypeScript clients.

What you build:
- Define entities with Ent and APIs with Protobuf  
- Generate Go handlers, Swagger, error types, and client SDKs
- Compose services with Gin + Wire; deploy as a single binary

## Prerequisites

- Go 1.24+
- Docker + Docker Compose  
- Node.js + npm (for TypeScript clients)

## Install Tooling

**Install CLI:**
```bash
go install github.com/go-sphere/sphere-cli@latest
```

**Install Protoc Plugins:**
```bash
go install github.com/go-sphere/protoc-gen-sphere@latest
go install github.com/go-sphere/protoc-gen-route@latest
go install github.com/go-sphere/protoc-gen-sphere-binding@latest
go install github.com/go-sphere/protoc-gen-sphere-errors@latest
```

**Verify Installation:**
```bash
sphere-cli --version || sphere-cli -h
protoc-gen-sphere --version
protoc-gen-route --version
protoc-gen-sphere-binding --version
protoc-gen-sphere-errors --version
```

## Create Your First Project

### 1. Bootstrap Project

```bash
# Create a new project using the template
# Replace 'myproject' with your project name and update the Go module path
sphere-cli create --name myproject --mod github.com/yourusername/myproject
cd myproject
```

This generates a new project with a clean structure based on [sphere-layout](https://github.com/go-sphere/sphere-layout).

### 2. Define Database Schema (Ent)

Create your database entities in `/internal/pkg/database/ent/schema`. For example, create `internal/pkg/database/ent/schema/user.go`:

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

option go_package = "myproject/api/shared/v1;sharedv1";

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

option go_package = "myproject/api/api/v1;apiv1";

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
  int64 id = 1;
}

message CreateUserRequest {
  string name = 1;
  string email = 2;
  int32 age = 3;
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

