---
title: protoc-gen-route
weight: 33
---

`protoc-gen-route` is a protoc plugin that generates routing code from `.proto` files. It is designed to inspect service definitions within your protobuf files and automatically generate corresponding route handlers based on a specified template. This plugin creates Go code that provides structured routing with operation constants, extra data handling, server interfaces, and codec interfaces for seamless integration with various transport protocols.

## Features

- Generates operation constants for each service method
- Creates extra data mappings from proto options
- Provides server and codec interfaces for type-safe implementations
- Supports custom request and response models
- Generates handler functions with automatic request/response conversion
- Integrates with the sphere options framework
- Supports flexible template customization

## Installation

To install `protoc-gen-route`, use the following command:

```bash
go install github.com/go-sphere/protoc-gen-route@latest
```

## Prerequisites

You need to have the sphere options proto definitions in your project. Add the following dependency to your `buf.yaml`:

```yaml
deps:
  - buf.build/go-sphere/options
```

## Configuration Parameters

The behavior of `protoc-gen-route` can be customized with the following parameters:

- **`version`**: Print the current plugin version and exit. (Default: `false`)
- **`options_key`**: The key for the option extension in your proto file that contains routing information. (Default: `route`)
- **`file_suffix`**: The suffix for the generated files. (Default: `_route.pb.go`)
- **`template_file`**: Path to a custom Go template file. If not provided, the default internal template is used.
- **`request_model`**: (Required) The fully qualified Go type for the request model (e.g., `github.com/gin-gonic/gin.Context`).
- **`response_model`**: (Required) The fully qualified Go type for the response model.
- **`extra_data_model`**: The fully qualified Go type for an additional data model to be used in the template.
- **`extra_data_constructor`**: A function that constructs and returns a pointer to the `extra_data_model`. (Required if `extra_data_model` is set).

## Usage with Buf

To use `protoc-gen-route` with `buf`, you can configure it in your `buf.gen.yaml` file. Here is an example configuration:

```yaml
version: v2
managed:
  enabled: true
  disable:
    - file_option: go_package_prefix
      module: buf.build/go-sphere/options
plugins:
  - local: protoc-gen-go
    out: api
    opt:
      - paths=source_relative
  - local: protoc-gen-route
    out: api
    opt:
      - paths=source_relative
      - options_key=bot
      - request_model=github.com/go-sphere/sphere/social/telegram;Update
      - response_model=github.com/go-sphere/sphere/social/telegram;Message
      - extra_data_model=github.com/go-sphere/sphere/social/telegram;MethodExtraData
      - extra_data_constructor=github.com/go-sphere/sphere/social/telegram;NewMethodExtraData
```

## Proto Definition Example

Here's how to define services with routing options in your `.proto` files:

```protobuf
syntax = "proto3";

package bot.v1;

import "sphere/options/options.proto";

service MenuService {
  // UpdateCount handles count update operations
  rpc UpdateCount(UpdateCountRequest) returns (UpdateCountResponse) {
    option (sphere.options.options) = [
      {
        key: "callback_query"
        text: "start"
      },
      {
        key: "command"
        text: "start"
      }
    ];
  }
  
  rpc ProcessMenu(ProcessMenuRequest) returns (ProcessMenuResponse) {
    option (sphere.options.options) = [
      {
        key: "callback_query"
        text: "menu_.*"
      }
    ];
  }
}

message UpdateCountRequest {
  int64 value = 1;
  int64 offset = 2;
}

message UpdateCountResponse {
  int64 value = 1;
}

message ProcessMenuRequest {
  string menu_id = 1;
  string action = 2;
}

message ProcessMenuResponse {
  string result = 1;
}
```

## Generated Code

The plugin generates Go code with the following components for each service:

### Operation Constants

```go
const OperationBotMenuServiceUpdateCount = "/bot.v1.MenuService/UpdateCount"
const OperationBotMenuServiceProcessMenu = "/bot.v1.MenuService/ProcessMenu"
```

### Extra Data Variables

```go
var ExtraBotDataMenuServiceUpdateCount = telegram.NewMethodExtraData(map[string]string{
    "callback_query": "start",
    "command":        "start",
})

var ExtraBotDataMenuServiceProcessMenu = telegram.NewMethodExtraData(map[string]string{
    "callback_query": "menu_.*",
})
```

### Helper Functions

```go
func GetExtraBotDataByMenuServiceOperation(operation string) *telegram.MethodExtraData {
    switch operation {
    case OperationBotMenuServiceUpdateCount:
        return ExtraBotDataMenuServiceUpdateCount
    case OperationBotMenuServiceProcessMenu:
        return ExtraBotDataMenuServiceProcessMenu
    default:
        return nil
    }
}
```

### Server Interface

```go
type MenuServiceBotServer interface {
    UpdateCount(context.Context, *UpdateCountRequest) (*UpdateCountResponse, error)
    ProcessMenu(context.Context, *ProcessMenuRequest) (*ProcessMenuResponse, error)
}
```

### Codec Interface

```go
type MenuServiceBotCodec interface {
    DecodeUpdateCountRequest(*telegram.Update) (*UpdateCountRequest, error)
    EncodeUpdateCountResponse(*UpdateCountResponse) (*telegram.Message, error)
    DecodeProcessMenuRequest(*telegram.Update) (*ProcessMenuRequest, error)
    EncodeProcessMenuResponse(*ProcessMenuResponse) (*telegram.Message, error)
}
```

### Registration Function

```go
func RegisterMenuServiceBotServer(server MenuServiceBotServer, codec MenuServiceBotCodec) map[string]telegram.Handler {
    handlers := make(map[string]telegram.Handler)
    
    handlers[OperationBotMenuServiceUpdateCount] = func(ctx context.Context, update *telegram.Update) (*telegram.Message, error) {
        req, err := codec.DecodeUpdateCountRequest(update)
        if err != nil {
            return nil, err
        }
        
        resp, err := server.UpdateCount(ctx, req)
        if err != nil {
            return nil, err
        }
        
        return codec.EncodeUpdateCountResponse(resp)
    }
    
    handlers[OperationBotMenuServiceProcessMenu] = func(ctx context.Context, update *telegram.Update) (*telegram.Message, error) {
        req, err := codec.DecodeProcessMenuRequest(update)
        if err != nil {
            return nil, err
        }
        
        resp, err := server.ProcessMenu(ctx, req)
        if err != nil {
            return nil, err
        }
        
        return codec.EncodeProcessMenuResponse(resp)
    }
    
    return handlers
}
```

## Common Use Cases

### Bot Command Routing

The plugin is commonly used for routing bot commands in messaging platforms:

```protobuf
service BotService {
  rpc Start(StartRequest) returns (StartResponse) {
    option (sphere.options.options) = [
      {
        key: "command"
        text: "/start"
      },
      {
        key: "callback_query"
        text: "start_.*"
      }
    ];
  }
}
```

### Event Routing

You can also use it for general event routing:

```protobuf
service EventService {
  rpc HandleUserCreated(UserCreatedRequest) returns (UserCreatedResponse) {
    option (sphere.options.options) = [
      {
        key: "event_type"
        text: "user.created"
      }
    ];
  }
}
```

## Integration Example

Here's how you might integrate the generated code in a bot application:

```go
package main

import (
    "context"
    "log"
    
    botv1 "myproject/api/bot/v1"
    "github.com/go-sphere/sphere/social/telegram"
)

// Implement the server interface
type MenuService struct{}

func (s *MenuService) UpdateCount(ctx context.Context, req *botv1.UpdateCountRequest) (*botv1.UpdateCountResponse, error) {
    // Your business logic here
    newValue := req.Value + req.Offset
    return &botv1.UpdateCountResponse{
        Value: newValue,
    }, nil
}

func (s *MenuService) ProcessMenu(ctx context.Context, req *botv1.ProcessMenuRequest) (*botv1.ProcessMenuResponse, error) {
    // Your business logic here
    return &botv1.ProcessMenuResponse{
        Result: fmt.Sprintf("Processed menu %s with action %s", req.MenuId, req.Action),
    }, nil
}

// Implement the codec interface
type MenuServiceCodec struct{}

func (c *MenuServiceCodec) DecodeUpdateCountRequest(update *telegram.Update) (*botv1.UpdateCountRequest, error) {
    // Parse the telegram update into your request struct
    // This would typically extract data from callback data or command parameters
}

func (c *MenuServiceCodec) EncodeUpdateCountResponse(resp *botv1.UpdateCountResponse) (*telegram.Message, error) {
    // Convert your response into a telegram message
    return &telegram.Message{
        Text: fmt.Sprintf("Count updated to: %d", resp.Value),
    }, nil
}

// Similar implementations for other methods...

func main() {
    server := &MenuService{}
    codec := &MenuServiceCodec{}
    
    // Register handlers
    handlers := botv1.RegisterMenuServiceBotServer(server, codec)
    
    // Use with your bot framework
    for operation, handler := range handlers {
        extraData := botv1.GetExtraBotDataByMenuServiceOperation(operation)
        // Register with your bot router using the operation and extra data
        log.Printf("Registered handler for %s with extra data: %v", operation, extraData)
    }
}
```

## Advanced Configuration

### Custom Templates

You can provide your own template file for custom code generation:

```yaml
plugins:
  - local: protoc-gen-route
    out: api
    opt:
      - paths=source_relative
      - template_file=./templates/custom_route.tmpl
      - options_key=custom
```

### Multiple Option Keys

You can generate routing code for different transport protocols by using different option keys:

```protobuf
service UserService {
  rpc GetUser(GetUserRequest) returns (GetUserResponse) {
    option (sphere.options.options) = [
      {
        key: "http"
        text: "GET /users/{id}"
      },
      {
        key: "bot"
        text: "user_info"
      },
      {
        key: "event"
        text: "user.get"
      }
    ];
  }
}
```

## Best Practices

1. **Use meaningful operation keys**: Choose keys that clearly indicate the transport or routing mechanism
2. **Group related operations**: Keep related service methods in the same service for better organization
3. **Implement proper error handling**: Ensure your codec implementations handle errors gracefully
4. **Use consistent naming**: Follow consistent naming patterns for your services and methods
5. **Document your options**: Include comments explaining the routing options for each method

## Troubleshooting

### Common Issues

1. **Missing option annotations**: Ensure all methods that should generate routes have the appropriate `sphere.options.options` annotations
2. **Template errors**: Check your custom template syntax if using a custom template file
3. **Missing dependencies**: Verify that all required model types are available and properly imported

## Integration with Other Generators

`protoc-gen-route` works well alongside other Sphere generators:

- Use with `protoc-gen-sphere` for HTTP routing
- Combine with `protoc-gen-sphere-errors` for consistent error handling
- Works with any transport layer that can use the generated interfaces

## See Also

- Concepts: ../../concepts/protocol-and-codegen
- Component: ./protoc-gen-sphere
