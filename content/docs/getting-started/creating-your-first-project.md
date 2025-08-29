---
title: Creating Your First Project
weight: 30
---

This guide provides a step-by-step walkthrough for creating, building, and running a new application using the Sphere framework.

## Prerequisites

Before you begin, ensure you have the following installed:
- Go (version 1.24 or later)
- Docker and Docker Compose
- Node.js and npm (for TypeScript client generation)

## 1. Create a New Project

The first step is to install the `sphere-cli` and use it to bootstrap your new application. The CLI will create a new project directory with the recommended layout.

```bash
# Install the command-line tool
go install github.com/go-sphere/sphere-cli@latest

# Create a new project using the template
# Replace 'myproject' with your project name and update the Go module path
sphere-cli create --name myproject --mod github.com/TBXark/myproject
```

This command generates a new project with a clean structure. `sphere` is designed to be flexible, so feel free to modify the directory layout to suit your project's needs.

## 2. Define Database Schema with Ent

Sphere uses `ent` to manage the database schema. Define your database entities (tables) in the `/internal/pkg/database/ent/schema` directory.

For example, to create a `User` schema, you can create a file `internal/pkg/database/ent/schema/user.go`:

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

For a detailed guide on writing `ent` schemas, refer to the [official ent documentation](https://entgo.io/docs/getting-started).

After defining your schemas, run the following `make` command:

```bash
# This command generates:
# - Ent client and schema code in /internal/pkg/database/ent
# - Corresponding Protobuf message definitions in /proto/entpb
make gen/db
```

## 3. Define API Interfaces with Protobuf

Your service's API endpoints are defined in `.proto` files located in the `/proto` directory. Sphere utilizes [gRPC-Gateway](https://grpc-ecosystem.github.io/grpc-gateway/) style annotations to define how gRPC services map to HTTP/JSON endpoints.

It's a good practice to define shared messages in the `proto/shared` directory. For example, `proto/shared/v1/user.proto`:

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

Then, you can define a `UserService` in `proto/api/v1/user.proto` that uses the shared `User` message:

```protobuf
syntax = "proto3";

package api.v1;

option go_package = "myproject/api/api/v1;apiv1";

import "google/api/annotations.proto";
import "sphere/binding/binding.proto";
import "shared/v1/user.proto";

message GetUserRequest {
  int64 id = 1 [(sphere.binding.location) = BINDING_LOCATION_URI];
}

message CreateUserRequest {
  string name = 1;
  string email = 2;
  int32 age = 3;
}

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
```

For more details on writing `.proto` files and using transcoding, see:
- [Protobuf Language Guide (proto3)](https://developers.google.com/protocol-buffers/docs/proto3)
- [gRPC Transcoding](https://cloud.google.com/endpoints/docs/grpc/transcoding)
- [API Definition Rules](../guides/api-definitions) for Sphere-specific conventions.

Once your API is defined, generate the server code, client stubs, and Swagger/OpenAPI documentation:

```bash
# This command generates:
# - Go server stubs in the /api directory
# - Swagger/OpenAPI v2 specifications in the /swagger directory
make gen/docs
```

You can also generate a TypeScript client for your frontend application:
```bash
make gen/dts
```

## 4. Implement Business Logic

With your API defined and code generated, it's time to implement the actual business logic. This happens in two main layers:

### Service Layer
The service layer (`internal/service/`) implements the generated service interfaces. This is where you handle HTTP-specific concerns and coordinate with the business logic layer.

### Business Logic Layer
The business logic layer (`internal/biz/`) contains the core application logic, independent of transport protocols.

Example service implementation:

```go
package service

import (
	"context"
	apiv1 "myproject/api/api/v1"
	sharedv1 "myproject/api/shared/v1"
	"myproject/internal/biz"
)

type UserService struct {
	userBiz *biz.UserBiz
}

func NewUserService(userBiz *biz.UserBiz) *UserService {
	return &UserService{
		userBiz: userBiz,
	}
}

func (s *UserService) GetUser(ctx context.Context, req *apiv1.GetUserRequest) (*sharedv1.User, error) {
	user, err := s.userBiz.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *apiv1.CreateUserRequest) (*sharedv1.User, error) {
	user, err := s.userBiz.CreateUser(ctx, req.Name, req.Email, req.Age)
	if err != nil {
		return nil, err
	}
	return user, nil
}
```

## 5. Wire Dependencies

Sphere uses [Wire](https://github.com/google/wire) for dependency injection. After implementing your services, generate the wiring code:

```bash
make gen/wire
```

## 6. Run Your Application

Now you can build and run your application:

```bash
# Run the application
make run

# Or run with Swagger UI (in a separate terminal)
make run/swag
```

Your API will be available at `http://localhost:8080`, and if you started the Swagger UI, you can explore your API at `http://localhost:8081`.

## 7. Test Your API

You can test your API using curl or any HTTP client:

```bash
# Create a user
curl -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com", "age": 30}'

# Get a user
curl http://localhost:8080/v1/users/1
```

## Next Steps

- Read [Core Concepts](../concepts) for understanding the architecture
- Explore [API Definition Rules](../guides/api-definitions) for advanced HTTP mapping
- Learn about [Error Handling](../guides/error-handling) for typed, consistent errors
- Check out the [Components](../components) documentation for detailed tool usage

