---
title: protoc-gen-route
weight: 20
---

Generates routing/adapter glue from `.proto` using a flexible options-based template. Useful when mapping RPCs to non-HTTP transports or custom dispatch layers.

Install
- `go install github.com/go-sphere/protoc-gen-route@latest`

Concept
- Adds a custom `sphere.options.options` extension to `MethodOptions` (from `sphere/options/options.proto`)
- You place key–value pairs in `.proto` under a named key (e.g., `route`, `bot`)
- The plugin emits operation constants, decode/encode interfaces, and handler factories per method

Key Flags
- `version`: print version
- `options_key`: extension key to read (default `route`)
- `template_file`: custom template path
- `request_model`: Go ident `path;Type` representing request envelope
- `response_model`: Go ident `path;Type` representing response envelope
- `extra_data_model`: optional Go ident used to hold per-method extras
- `extra_data_constructor`: constructor that returns a pointer to `extra_data_model`

Buf Example
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
    opt: paths=source_relative
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

Proto Annotation Sketch
```protobuf
import "sphere/options/options.proto";
service CounterService {
  rpc Start(StartRequest) returns (StartResponse) {
    option (sphere.options.options) = {
      options: { key: "bot", text: "start" }
    };
  }
}
```

Output Highlights
- Operation constants per RPC
- `Register<Service><Key>Server` to register handlers
- Codec interface to decode envelope into RPC request and encode response

