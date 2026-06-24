---
title: Sphere
layout: hextra-home
---

{{< hextra/hero-container image="images/sphere.png" imageClass="hx:hidden hx:sm:block" >}}
{{< hextra/hero-badge link="https://github.com/TBXark" >}}
  <div class="hx:w-2 hx:h-2 hx:rounded-full hx:bg-primary-400"></div>
  <span>Made&nbsp;with&nbsp;love&nbsp;by&nbsp;TBXark</span>
  {{< icon name="arrow-circle-right" attributes="height=14" >}}
{{< /hextra/hero-badge >}}

<div class="hx:mt-6 hx:mb-6">
{{< hextra/hero-headline >}}
  Thin integration layer&nbsp;
  <br class="hx:sm:block hx:hidden" />for definition-driven development
{{< /hextra/hero-headline >}}
</div>

<div class="hx:mb-12">
{{< hextra/hero-subtitle >}}
  Compose mature Go tools with Protobuf contracts, small adapters, code generators, templates, and Makefile-driven workflows.
{{< /hextra/hero-subtitle >}}
</div>

<div class="hx:mb-6">
{{< hextra/hero-button text="Get Started" link="docs/getting-started/quickstart" >}}
</div>
{{< /hextra/hero-container >}}

<div class="hx:mt-12">
{{< hextra/feature-grid>}}
  {{< hextra/feature-card
    title="Protocol-First Design"
    subtitle="Use Protobuf as the service contract. Generate transport glue, binding metadata, typed errors, and documentation from the same source."
    link="docs/concepts/protocol-and-codegen"
  >}}
  {{< hextra/feature-card
    title="Replaceable Defaults"
    subtitle="Official templates pick Gin, Ent or Bun, Wire, Buf, and Swagger, but those choices remain visible and replaceable."
    link="docs/concepts/philosophy"
  >}}
  {{< hextra/feature-card
    title="Makefile Workflow"
    subtitle="The CLI creates a project; the generated Makefile owns init, generation, formatting, linting, running, and builds."
    link="docs/concepts/makefile-contract"
  >}}
  {{< hextra/feature-card
    title="Focused Code Generation"
    subtitle="The protoc-gen-sphere ecosystem creates repeatable plumbing while business logic stays in ordinary Go packages."
    link="docs/components"
  >}}
  {{< hextra/feature-card
    title="Structured Error Handling"
    subtitle="Define error enums in protobuf with automatic HTTP status mapping. Get consistent JSON responses with code, reason, and message."
    link="docs/guides/error-handling"
  >}}
  {{< hextra/feature-card
    title="Customizable Stack"
    subtitle="Swap routers, persistence, dependency injection, documentation, or deployment flow while keeping the same project contracts."
    link="docs/guides/customizing-stack"
  >}}
{{< /hextra/feature-grid >}}
</div>
