---
title: protoc-gen-sphere-binding
weight: 34
---

`protoc-gen-sphere-binding` is a protoc plugin that generates Go struct tags for Sphere binding from `.proto` files. It is designed to inspect service definitions within your protobuf files and automatically generate corresponding Go struct tags based on sphere binding annotations.

Unlike other protoc plugins that generate new files, this plugin modifies the generated Go structs by injecting appropriate tags for request binding in HTTP handlers.

## Installation

To install `protoc-gen-sphere-binding`, use the following command:

```bash
go install github.com/go-sphere/protoc-gen-sphere-binding@latest
```

## Prerequisites

You need to have the sphere binding proto definitions in your project. Add the following dependency to your `buf.yaml`:

```yaml
deps:
  - buf.build/go-sphere/binding
```

## Configuration Parameters

The behavior of `protoc-gen-sphere-binding` can be customized with the following parameters:

- **`version`**: Print the current plugin version and exit. (Default: `false`)
- **`out`**: The output directory for the modified `.proto` files. (Default: `api`)

## Usage with Buf

To use `protoc-gen-sphere-binding` with `buf`, you can configure it in your `buf.gen.yaml` file. **Note**: `protoc-gen-sphere-binding` cannot be used with the standard `buf.gen.yaml` because it does not generate Go code, but rather modifies the `.proto` files to include Sphere binding tags.

Here is an example configuration:

```yaml
version: v2
managed:
  enabled: true
  override:
    - file_option: go_package_prefix
      value: github.com/go-sphere/sphere-layout/api
plugins:
  - local: protoc-gen-go
    out: api
    opt:
      - paths=source_relative
  - local: protoc-gen-sphere-binding
    out: api
    opt:
      - paths=source_relative
```

## How It Works

`protoc-gen-sphere-binding` processes protobuf files and adds Go struct tags based on sphere binding annotations. The plugin works in conjunction with `protoc-gen-go` and should be run after the standard Go code generation to add the binding tags to the generated structs.

## Binding Locations

The plugin supports the following binding locations through the `sphere.binding.location` annotation:

- `BINDING_LOCATION_BODY`: Fields bound to JSON request body (default behavior)
- `BINDING_LOCATION_QUERY`: Fields bound to query parameters (adds `form` tag)
- `BINDING_LOCATION_URI`: Fields bound to URI path parameters (adds `uri` tag, removes `json` tag)
- `BINDING_LOCATION_HEADER`: Fields bound to HTTP headers (adds `header` tag)
- `BINDING_LOCATION_FORM`: Fields bound to form data (adds `form` tag)

## Proto Definition Example

Here's a comprehensive example showing different binding locations:

```protobuf
syntax = "proto3";

package shared.v1;

import "buf/validate/validate.proto";
import "google/api/annotations.proto";
import "sphere/binding/binding.proto";

service TestService {
  rpc RunTest(RunTestRequest) returns (RunTestResponse) {
    option (google.api.http) = {
      post: "/v1/test/{path_test1}"
      body: "*"
    };
  }
}

message RunTestRequest {
  // URI path parameter
  string path_test1 = 1 [(sphere.binding.location) = BINDING_LOCATION_URI];
  // Fields without binding annotation go to JSON body
  string field_test1 = 2;
  // Query parameter
  string query_test1 = 3 [(sphere.binding.location) = BINDING_LOCATION_QUERY];
  // Header value
  string auth_token = 8 [(sphere.binding.location) = BINDING_LOCATION_HEADER];
  // Custom tags
  repeated int32 ids = 9 [
    (sphere.binding.location) = BINDING_LOCATION_QUERY,
    (sphere.binding.auto_tags) = "custom:\"ids\""
  ];
}

// Message with default auto tags
message BodyPathTestRequest {
  option (sphere.binding.default_auto_tags) = "db";
  
  message Request {
    string name = 1;  // Will get: db:"name" json:"name"
    string email = 2; // Will get: db:"email" json:"email"
  }
  Request request = 1;
}

enum TestEnum {
  TEST_ENUM_UNSPECIFIED = 0;
  TEST_ENUM_VALUE1 = 1;
  TEST_ENUM_VALUE2 = 2;
}
```

## Generated Code

After running `protoc-gen-go` followed by `protoc-gen-sphere-binding`, the generated Go struct will have appropriate tags:

```go
type RunTestRequest struct {
    state         protoimpl.MessageState `protogen:"open.v1"`
    sizeCache     protoimpl.SizeCache
    unknownFields protoimpl.UnknownFields

    PathTest1  string  `protobuf:"bytes,1,opt,name=path_test1,json=pathTest1,proto3" json:"-" uri:"path_test1"`
    FieldTest1 string  `protobuf:"bytes,2,opt,name=field_test1,json=fieldTest1,proto3" json:"field_test1"`
    QueryTest1 string  `protobuf:"bytes,3,opt,name=query_test1,json=queryTest1,proto3" json:"-" form:"query_test1"`
    AuthToken  string  `protobuf:"bytes,8,opt,name=auth_token,json=authToken,proto3" json:"-" header:"auth_token"`
    Ids        []int32 `protobuf:"varint,9,rep,packed,name=ids,proto3" json:"-" form:"ids" custom:"ids"`
}

type BodyPathTestRequest_Request struct {
    state         protoimpl.MessageState `protogen:"open.v1"`
    sizeCache     protoimpl.SizeCache
    unknownFields protoimpl.UnknownFields

    Name  string `protobuf:"bytes,1,opt,name=name,proto3" json:"name" db:"name"`
    Email string `protobuf:"bytes,2,opt,name=email,proto3" json:"email" db:"email"`
}
```

## Advanced Features

### Message-Level Configuration

You can set default binding behavior for entire messages:

```protobuf
message SearchRequest {
  option (sphere.binding.default_location) = BINDING_LOCATION_QUERY;
  option (sphere.binding.default_auto_tags) = "form";
  
  string name = 1;        // Will be bound from query with form tag
  int32 age = 2;          // Will be bound from query with form tag
  string email = 3;       // Will be bound from query with form tag
}
```

### Oneof Support

The plugin also supports oneof fields with default configurations:

```protobuf
message TestRequest {
  oneof test_oneof {
    option (sphere.binding.default_oneof_location) = BINDING_LOCATION_QUERY;
    option (sphere.binding.default_oneof_auto_tags) = "form";
    
    string option_a = 1;
    string option_b = 2;
  }
}
```

### Custom Tags

Add custom Go struct tags using the `auto_tags` annotation:

```protobuf
message DatabaseModel {
  option (sphere.binding.default_auto_tags) = "db";
  
  string name = 1;     // Generated: `db:"name" json:"name"`
  string email = 2 [(sphere.binding.auto_tags) = "validate:\"email\""];
  // Generated: `db:"email" json:"email" validate:"email"`
}
```

## Usage in HTTP Handlers

The generated tags work seamlessly with Gin's binding functions:

```go
func (s *TestService) RunTest(c *gin.Context) {
    var req RunTestRequest
    
    // Bind URI parameters
    if err := c.ShouldBindUri(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // Bind query parameters
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // Bind JSON body
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // Bind headers
    if err := c.ShouldBindHeader(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // Process request...
}
```

## Integration with Other Generators

`protoc-gen-sphere-binding` works best when combined with other Sphere generators:

1. **protoc-gen-go**: Generates the base Go structs
2. **protoc-gen-sphere-binding**: Adds binding tags to the structs
3. **protoc-gen-sphere**: Generates HTTP handlers that use the tagged structs

## Best Practices

1. **Use consistent binding locations**: Establish patterns for where different types of data should be bound from
2. **Leverage message-level defaults**: Use `default_location` and `default_auto_tags` to reduce repetition
3. **Be explicit when needed**: Override defaults with field-level annotations when necessary
4. **Test your bindings**: Verify that the generated tags work correctly with your HTTP framework
5. **Keep it simple**: Avoid overly complex binding patterns that might confuse API consumers

## Troubleshooting

### Common Issues

1. **Tags not appearing**: Make sure you're running `protoc-gen-sphere-binding` after `protoc-gen-go`
2. **Binding not working**: Verify that your HTTP framework supports the generated tag format
3. **Conflicts with existing tags**: Check for conflicts between generated tags and existing protobuf tags

### Debugging

You can verify the generated tags by examining the generated Go files or using reflection:

```go
import (
    "fmt"
    "reflect"
)

func inspectTags() {
    t := reflect.TypeOf(RunTestRequest{})
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        fmt.Printf("Field: %s, Tags: %s\n", field.Name, field.Tag)
    }
}
```
