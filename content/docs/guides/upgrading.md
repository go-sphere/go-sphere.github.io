---
title: Upgrading to v0.0.5 / v0.0.6
weight: 45
---

`sphere` v0.0.6 and `httpx` v0.0.5 are a breaking pair. `httpx` collapsed its two middleware forms into one, unified the adapter options that had drifted apart, and stopped validating bound structs; `sphere` migrated its own middleware to the surviving form and moved the official templates to the `stdx` (plain `net/http`) engine.

Generated code is unaffected: handlers use only `httpx.Context` / `Handler` / `Router`, `Bind*`, `Group` / `Handle` and the `httpz` wrappers, and none of those changed shape. Hand-written middleware, engine setup, custom binders, and adapter options are what need attention.

```bash
go get github.com/go-sphere/sphere@v0.0.6
go get github.com/go-sphere/httpx@v0.0.5
go get github.com/go-sphere/httpx/stdx@v0.0.5   # adapters are tagged per module
```

The authoritative lists live in the repositories:

- `httpx/CHANGELOG.md` — every v0.0.5 signature and behavior change, per adapter;
- `sphere/compat/api-incompatibilities.txt` — `apidiff`-visible signature changes;
- `sphere/compat/behavior-changes.md` — same signature, different runtime.

## One middleware form

v0.0.4 had two layer types. v0.0.5 ships one, and it is the composed form:

```go
// v0.0.4
type Middleware func(ctx httpx.Context) error
type Interceptor func(next httpx.Handler) httpx.Handler

// v0.0.5
type Middleware func(next httpx.Handler) httpx.Handler
```

A layer receives the rest of the chain and returns what runs in its place. Move the old body into an inner closure and return `next(ctx)` where it used to return `ctx.Next()`:

```go
// v0.0.4
func RequestID(ctx httpx.Context) error {
    ctx.SetContext(withRequestID(ctx.Context()))
    return ctx.Next()
}

// v0.0.5
func RequestID(next httpx.Handler) httpx.Handler {
    return func(ctx httpx.Context) error {
        ctx.SetContext(withRequestID(ctx.Context()))
        return next(ctx)
    }
}
```

| v0.0.4 | v0.0.5 |
| --- | --- |
| `httpx.Interceptor` | `httpx.Middleware` |
| `httpx.InterceptorScope` | `httpx.MiddlewareScope` |
| `scope.UseInterceptor(...)` | `scope.Use(...)` |
| `httpx.UseInterceptor(scope, ...)` | `scope.Use(...)` |
| `httpx.ComposeInterceptors` | `httpx.ComposeMiddleware` |
| `httpx.InterceptorChain` / `NewInterceptorChain` | `httpx.MiddlewareChain` / `NewMiddlewareChain` |
| `httpx.InterceptorFallback` / `NewInterceptorFallback` | `httpx.MiddlewareFallback` / `NewMiddlewareFallback` |
| `httpx.AsMiddleware`, `httpx.AsInterceptor` | removed — there is nothing to convert |
| `Context.Next()` | removed — call the `next` handler you were given |

The two signatures are unrelated, so every call site that passed the old form is a compile error, never a silent change. Four rules follow from the collapse:

- `Use` layers always run **inside** anything mounted with the adapter's `UseNative`, whatever the registration order. The `Adapt<Framework>Middleware` helpers were removed with it: all of a scope's layers now share one native handler slot, so a native middleware's own continuation would advance past it.
- An inner error is rendered at the route where the chain was composed, not at the failing layer. Read it from the value `next` returns, not only `ctx.StatusCode()`.
- An engine-scope chain also covers unmatched paths (404/405) on all five adapters; a group's chain does not.
- A late `engine.Use` reaches routes an existing group registers afterwards (previously echox/fiberx only).

`sphere/server/middleware` kept its constructor names; their return type is the new form. From v0.0.6 these all return `httpx.Middleware`: `auth.NewAuthMiddleware`, `auth.NewPermissionMiddleware`, `cors.NewCORS`, `(*online.Online).Middleware`, `ratelimiter.NewRateLimiter`, `ratelimiter.NewRateLimiterByClientIP`, `selector.NewSelectorMiddleware`, `logger.Log`, `logger.RecoveryLog`. The `*Interceptor` variants exposed in sphere v0.0.5 were removed — the plain names are the composed form.

## Adapter option renames

| v0.0.4 | v0.0.5 |
| --- | --- |
| `ginx.WithHTTPXErrorHandler(httpx.ErrorHandler)` | `ginx.WithErrorHandler(...)` |
| `hertzx.WithHTTPXErrorHandler(httpx.ErrorHandler)` | `hertzx.WithErrorHandler(...)` |
| `ginx.WithErrorHandler(ginx.ErrorHandler)` | `ginx.WithNativeErrorHandler(...)` |
| `hertzx.WithErrorHandler(hertzx.ErrorHandler)` | `hertzx.WithNativeErrorHandler(...)` |
| `echox.DefaultHTTPErrorHandler` | `echox.DefaultErrorHandler` |
| `ginx.WithServerAddr`, `echox.WithServerAddr` | `WithAddr` |
| `ginx.QueryBinding` | unexported |

`WithErrorHandler` is now `httpx.ErrorHandler` on every adapter and is the portable surface. `DefaultErrorHandler` is each adapter's own default in its **native** shape, so the five signatures differ on purpose and are not interchangeable.

## Bind* no longer validates

Adapters used to run go-playground/validator over the `binding` tag after a successful decode. That made the multi-source binding generated code uses impossible:

```go
struct {
    Name string `json:"name"`
    ID   string `uri:"id" binding:"required"`
}
```

The generated sequence `BindJSON` → `BindHeader` → `BindQuery` → `BindURI` failed with 400 at `BindJSON`, validating a `uri` field nothing had populated yet. `Bind*` now only decodes; a `binding` tag means nothing to `httpx`. Validation belongs above the transport — `protoc-gen-sphere` emits `protovalidate.Validate(&in)` — and a decode failure is still a 400 through `WrapBindError`.

## Removed symbols

| Removed | Instead |
| --- | --- |
| `httpx.WithJson`, `httpx.H` | `sphere/server/httpz.WithJson` |
| `httpx.AsResponseInfo` | `ResponseInfo` is part of `Context`; read `ctx.StatusCode()` |
| `httpx.ListenAndAutoShutdown` | `httpx.Start` / `httpx.Close` |
| `AdaptGinMiddleware`, `AdaptEchoMiddleware`, `AdaptFiberMiddleware`, `AdaptHertzMiddleware` | the adapter's `UseNative` |
| `httpx.FixWildcardPathIfNeed` | register the named wildcard (`/files/*name`) directly |
| `Responder.DataFromReader` with `size int` | `size int64` |

The anonymous wildcard `/files/*` is now rejected at registration on every adapter with httpx's own error; only the named form (`/files/*name`) is accepted, and its result must not be passed to `Handle`.

## Errors no longer leak raw text

`httpx.ParseError` no longer falls back to `err.Error()`: an error carrying no `httpx.MessageError` yields an empty message, and rendering substitutes the generic status text. That closes the driver/SQL/panic-text leak at the source, so the old filter in `AbortWithJsonError` (which dropped a parser message equal to `err.Error()`) is gone.

The consequence for custom parsers: a message returned by `httpz.SetDefaultErrorParser` is now authoritative and reaches the client verbatim. Audit any parser for errors it does not classify and return an empty message for those.

## What changed in sphere v0.0.6

- Middleware constructors return the composed form (see above).
- `fileserver.RegisterFileDownloader` registers `/*filename` directly; the `FixWildcardPathIfNeed` fix-then-register dance is both redundant and now wrong.
- `httpz.EndpointsToMatches` indexes only the path the caller registered; code that looked up the old anonymous spelling (`/files/*`) gets nothing back.
- `ratelimiter.NewRateLimiter` re-reads the cache inside the singleflight, so a request can no longer mint a second limiter (and a fresh full burst) after a concurrent flight's write.
- On `stdx`, `ClientIP` uses the direct peer address and ignores `X-Forwarded-For` / `X-Real-IP` unless `stdx.WithTrustedProxies` is set. Moving from gin/echo/hertz (which trust every peer by default) to `stdx` therefore tightens rate limiting; behind a real proxy without trusted proxies configured, every request shares one bucket.

New in v0.0.6, not migration work: eager-commit SSE streams (`httpz.WithSSEEagerCommit`), `log/logbuffer`, `server/middleware/logger`, and `task.NewFunc`. See [Logging](logging), [Infrastructure](infrastructure), and [Server Streaming](server-streaming).

## Templates and stdx

Official templates (`sphere-layout`, `sphere-simple-layout`, `sphere-bun-layout`, `sphere-telegram-layout`) now depend on `sphere` v0.0.6 and `httpx/stdx` v0.0.5, and already apply every change above. Existing projects created from older templates still need this page. A template registers its engine like this:

```go
engine := stdx.New(
    stdx.WithServer(httpServer),
    stdx.WithErrorHandler(httpz.AbortWithJsonError),
)
engine.Use(logger.Log(lg), logger.RecoveryLog(lg, true))
```

## Related

- [HTTP Runtime](http-runtime) — the current middleware and adapter contracts
- [Customizing the Stack](customizing-stack)
- [Infrastructure](infrastructure) — TTL, lifecycle, and `ClientIP` behavior
