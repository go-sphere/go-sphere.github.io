---
title: Sphere
layout: hextra-home
---

{{< hextra/hero-badge link="https://github.com/TBXark" >}}
  <div class="hx:w-2 hx:h-2 hx:rounded-full hx:bg-primary-400"></div>
  <span>Made&nbsp;with&nbsp;love&nbsp;by&nbsp;TBXark</span>
  {{< icon name="arrow-circle-right" attributes="height=14" >}}
{{< /hextra/hero-badge >}}

<div class="hx:mt-6 hx:mb-6">
{{< hextra/hero-headline >}}
  Build pragmatic Go backends&nbsp;<br class="hx:sm:block hx:hidden" />with Sphere
{{< /hextra/hero-headline >}}
</div>

<div class="hx:mb-12">
{{< hextra/hero-subtitle >}}
  Monolithic-first toolkit with codegen for contracts, errors, stubs, and clients
{{< /hextra/hero-subtitle >}}
</div>

<div class="hx:mb-6">
{{< hextra/hero-button text="Get Started" link="docs/getting-started/quickstart" >}}
</div>

{{< hextra/feature-grid >}}
  {{< hextra/feature-card
    title="Contracts First"
    subtitle="Define APIs in Protobuf and entities in Ent; generate Go handlers, routers, bindings, and clients."
  >}}
  {{< hextra/feature-card
    title="Monolith-First Template"
    subtitle="Start with a single binary using Gin + Wire; evolve to multi-service when needed."
  >}}
  {{< hextra/feature-card
    title="Code Generation"
    subtitle="protoc-gen-sphere, protoc-gen-route, protoc-gen-sphere-binding, and protoc-gen-sphere-errors automate server, routing, tags, and typed errors."
  >}}
  {{< hextra/feature-card
    title="Typed Errors"
    subtitle="Define error enums in .proto; get consistent HTTP JSON with status, code, reason, and message."
  >}}
  {{< hextra/feature-card
    title="Swagger & Clients"
    subtitle="Generate OpenAPI from contracts and optional TypeScript SDKs for frontends."
  >}}
  {{< hextra/feature-card
    title="CLI & Layout"
    subtitle="sphere-cli and sphere-layout bootstrap projects with Makefile workflows and a sane project structure."
  >}}
{{< /hextra/feature-grid >}}
