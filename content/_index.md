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
  Protobuf-first Go service framework&nbsp;<br class="hx:sm:block hx:hidden" />for definition-driven development
{{< /hextra/hero-headline >}}
</div>

<div class="hx:mb-12">
{{< hextra/hero-subtitle >}}
  Start monolithic, scale to microservices. Define once, generate everything with rapid tooling.
{{< /hextra/hero-subtitle >}}
</div>

<div class="hx:mb-6">
{{< hextra/hero-button text="Get Started" link="docs/getting-started/quickstart" >}}
</div>

{{< hextra/feature-grid >}}
  {{< hextra/feature-card
    title="Protocol-First Design"
    subtitle="Define once in Protobuf, generate everywhere. Get Go handlers, HTTP routing, client SDKs, and OpenAPI docs from a single source of truth."
  >}}
  {{< hextra/feature-card
    title="Pragmatic Monolith Template"
    subtitle="Start simple with Gin + Wire + Ent in a single binary. Clean architecture that scales from MVP to microservices when needed."
  >}}
  {{< hextra/feature-card
    title="Complete Code Generation"
    subtitle="Automated toolchain with protoc-gen-sphere ecosystem: server stubs, HTTP routing, field binding, typed errors, and validation."
  >}}
  {{< hextra/feature-card
    title="Structured Error Handling"
    subtitle="Define error enums in protobuf with automatic HTTP status mapping. Get consistent JSON responses with code, reason, and message."
  >}}
  {{< hextra/feature-card
    title="Full-Stack Development"
    subtitle="Generate Swagger documentation, TypeScript SDKs, and validation schemas. Bridge backend and frontend with type safety."
  >}}
  {{< hextra/feature-card
    title="Developer Experience"
    subtitle="sphere-cli for project scaffolding, Makefile workflows, and clean project structure. Focus on business logic, not boilerplate."
  >}}
{{< /hextra/feature-grid >}}
