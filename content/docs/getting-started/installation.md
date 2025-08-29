---
title: Installation
weight: 20
---

Prerequisites:
- Go 1.24+
- Docker + Docker Compose
- Node.js + npm (for TypeScript clients)

Install CLI
- `go install github.com/go-sphere/sphere-cli@latest`

Install Protoc Plugins
- `go install github.com/go-sphere/protoc-gen-sphere@latest`
- `go install github.com/go-sphere/protoc-gen-route@latest`
- `go install github.com/go-sphere/protoc-gen-sphere-binding@latest`
- `go install github.com/go-sphere/protoc-gen-sphere-errors@latest`

Verify Installation
```bash
sphere-cli --version || sphere-cli -h
protoc-gen-sphere --version
protoc-gen-route --version
protoc-gen-sphere-binding --version
protoc-gen-sphere-errors --version
```

Optional: TypeScript SDK Generation
- Ensure Node.js is installed for generating SDKs from Swagger (`make gen/dts` in templates)

Next: Creating Your First Project.
