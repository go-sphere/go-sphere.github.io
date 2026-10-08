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

Official templates (`sphere-layout`, `sphere-simple-layout`, `sphere-bun-layout`, `sphere-telegram-layout`) now depend on `sphere` v0.0.7, `httpx/stdx` v0.0.6 and `errors` v0.0.3, and already apply every change above — including `idgenerator.InitFromEnv`, the `httpz.ParseError` fallback and the engine-level body cap. Existing projects created from older templates still need this page. A template registers its engine like this:

```go
engine := stdx.New(
    stdx.WithServer(httpServer),
    stdx.WithErrorHandler(httpz.AbortWithJsonError),
)
engine.Use(logger.Log(lg), logger.RecoveryLog(lg, true))
```

## sphere v0.0.7

`sphere` v0.0.7 was released on 2026-10-08 with `httpx` v0.0.6. Most of its breaking changes compile cleanly and only change runtime behaviour. The authoritative list is the v0.0.7 section of `sphere/CHANGELOG.md`.

| Change | What to write instead |
| --- | --- |
| Importing `core/boot` no longer sets `time.Local` / `TZ` to `Asia/Shanghai`; without a call the process runs in the host zone (UTC in most images), which shifts log timestamps and cron/asynq schedules with an empty `Timezone`. | Call `boot.InitTimezone(boot.DefaultTimezone)` first thing in `main`. |
| `idgenerator` no longer reads `WORKER_ID` in package init; the first `NextId` does, and a malformed value panics there instead of at startup. | Call `idgenerator.InitFromEnv()` (or `idgenerator.Init(workerID)`) in `main` and fail on its error. |
| `boot.WithLoggerInit` is removed. | `boot.WithLoggerBackend(zapx.NewBackend(conf.Log, ...))` — see [Logging](logging). |
| `storageerr.ErrNotFound`, `ErrDestExists`, `ErrFileNameInvalid` no longer carry an HTTP status. The default parser `httpz.ParseError` maps them to 404/400, but a custom parser that falls back to `httpx.ParseError` renders them as 500. | Fall back to `httpz.ParseError` — see [Custom Error Parser](http-runtime#custom-error-parser). |
| `jwtauth.ParseToken` rejects tokens without `exp` (`jwt.ErrTokenRequiredClaimMissing`). | Set `ExpiresAt` on custom claims; `NewRBACClaims` already does. |
| `fileserver` stores upload tokens under a `sphere-upload-token:` prefix, so tokens issued before the upgrade stop working; `WithCreateFileKey` now takes `func(ctx context.Context) (string, error)` and only generates the token. | Re-issue pending upload tokens; return just the token from a custom `WithCreateFileKey`. |

The protoc plugins released in this batch (`protoc-gen-sphere` v0.0.6, `protoc-gen-sphere-binding` v0.0.6, `protoc-gen-sphere-errors` v0.0.4) also change generated output:

- `protoc-gen-sphere` now fails generation instead of emitting code that silently drops values: a request message whose JSON body contains a real `oneof`, a `Timestamp`/`Duration`/wrapper used as a query/uri/header parameter, a FORM field on a method without a body, and nested or missing `body`/`response_body` paths. Validation errors are rendered as 400 instead of 500, and a `body: "*"` request whose message has no JSON fields no longer binds a body. See [protoc-gen-sphere](../components/protoc-gen-sphere).
- `protoc-gen-sphere-errors` defaults `new_errors_func` to `github.com/go-sphere/errors/sphere/errors;NewError`, so generated errors no longer import `httpx` and the project needs `errors` v0.0.3 or newer. Set `new_errors_func=github.com/go-sphere/httpx;NewError` to keep the old constructor. See [protoc-gen-sphere-errors](../components/protoc-gen-sphere-errors).
- `protoc-gen-sphere-binding` no longer applies a message's `default_location` to its oneof members; they stay in the JSON body unless the oneof sets `default_oneof_location`. It also no longer treats the well-known types as bindable scalars for query/uri/header. See [protoc-gen-sphere-binding](../components/protoc-gen-sphere-binding#oneof-support).

## httpx v0.0.6

`httpx` v0.0.6 (with the adapter tags `ginx/v0.0.6`, `fiberx/v0.0.6`, `echox/v0.0.6`, `hertzx/v0.0.6`, `stdx/v0.0.6`) is breaking for `fiberx` and for third-party adapters:

| Change | What to write instead |
| --- | --- |
| A `fiberx` engine built by the adapter now routes like `stdx`: case sensitive, strict about a trailing slash, and a GET route no longer answers HEAD (`Allow` no longer lists it). | An app passed through `WithEngine` keeps its own settings; otherwise fix the routes or links that relied on the loose matching. |
| `fiberx` route precedence no longer depends on registration order: `/users/new` beats `/users/:id`, and a parameter beats a wildcard. A route that would have to move across a native middleware whose path could match it panics at registration. | Register the native middleware before the routes it wraps, as before. |
| Third-party adapters run through `httpxtest` must meet the new cases: `BodyLimit` requires wiring `Options.MaxBodySize` to your own limit; the routing, `File` and `Redirect` cases pin the behaviour below. `Caps.RawPathRouting` declares an engine that matches the still-encoded path. | Implement body-limit and `File` support, then run `httpxtest`. |
| `File` serves regular files only: a missing path or a directory writes nothing and returns a 404 rendered by the error handler (`hertzx` no longer lists a directory, `fiberx` no longer redirects to its absolute path). | Rely on the error handler; do not serve directories. |
| `ValidRedirectCode` accepts 300, 301, 302, 303, 307 and 308 only; 304, 305 and 306 return an error and write nothing. | Use one of the accepted codes. |
| `ParseError` classifies an error carrying a `*http.MaxBytesError` as 413 instead of 500, and `WrapBindError` reports it as 413 instead of 400. | Nothing — this is what lets `WithMaxBodySize` answer 413 on every adapter. |

New in v0.0.6: `WithMaxBodySize(n)` on all five adapters, `httpx.ResponseHeaderEditor` / `httpx.AsResponseHeaderEditor`, and `httpx.CheckServeFile` / `httpx.ServeFile`.

## Related

- [HTTP Runtime](http-runtime) — the current middleware and adapter contracts
- [Customizing the Stack](customizing-stack)
- [Infrastructure](infrastructure) — TTL, lifecycle, and `ClientIP` behavior
