---
title: Proto Packages
weight: 22
---

Sphere extends Protobuf with a few focused packages that keep HTTP binding, errors, and custom options declarative. These annotations power Sphere’s generators and help maintain consistent, type‑safe APIs.

## Overview

- sphere/binding: Field → HTTP mapping and tagging
- sphere/errors: Enum‑based typed errors with HTTP semantics
- sphere/options: Method‑level metadata for custom routing/generators

## sphere/binding

Purpose
- Declare where each field binds from (URI, query, body, header, form)
- Set message/oneof defaults and auto‑tags for generated structs
- Work seamlessly with Sphere’s Gin helpers for request binding

Use when
- You want explicit, generator‑driven request parsing rules
- You prefer consistent struct tags without hand editing

See also
- API details and examples: ../guides/api-definitions

## sphere/errors

Purpose
- Define typed error enums with HTTP status, reason, and message
- Generate helpers to wrap causes and produce uniform JSON errors

Use when
- You need consistent error shapes across services and clients
- You want programmatic access to status/code/reason/message

See also
- Concepts overview: protocol-and-codegen
- How to author and use errors: ../guides/error-handling

## sphere/options

Purpose
- Attach simple key/value metadata to RPC methods
- Let generators consume options for advanced routing or transports

Use when
- You build adapters beyond HTTP or need custom routing hints

