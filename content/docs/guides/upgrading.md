---
title: Upgrading to v0.0.4
weight: 45
---

`sphere` v0.0.4 is a breaking runtime release. Generated HTTP code already targeted `httpx` / `httpz` in v0.0.3; this release tightens contracts (cache TTL, `Close` ownership, error envelopes, boot/task lifecycle) and fixes security bugs that cannot keep the old signatures.

```bash
go get github.com/go-sphere/sphere@v0.0.4
```

Then apply the call-site changes below. Official templates (`sphere-layout`, `sphere-simple-layout`, `sphere-bun-layout`) have been bumped to v0.0.4. Existing projects created from older templates still need this page.

The authoritative lists live in the sphere module:

- `compat/api-incompatibilities.txt` — signature changes `apidiff` sees against v0.0.3
- `compat/behavior-changes.md` — same signature, different runtime

## Compile-time

| Change | Migration |
| --- | --- |
| `utils/secure.CryptPassword(pwd) string` → `(string, error)` | Handle the error. bcrypt failures no longer fall back to plaintext. |
| `boot.WithLoggerInit` | Still works, marked deprecated. Prefer `boot.WithLoggerBackend(zapx.NewBackend(conf))`. |
| `fileserver.WithCreateFileKey` callback gained a `ttl time.Duration` argument | Add the parameter. |
| `cors.NewCORS(...) httpx.Middleware` → `(httpx.Middleware, error)` | `"*"` plus credentials now returns `cors.ErrWildcardWithCredentials`. |
| `online.NewOnline()` → `NewOnline(...Option)` | The tracker is a `task.Task`; `Start` it so expired entries are swept. |
| `ratelimiter.NewNewRateLimiterByClientIP` / `WithSetTTL` removed | Use `NewRateLimiterByClientIP` / `NewRateLimiter`. |
| `log.Init` removed | Use `log.InitWithBackends`. |
| `core/task.ErrGroupNotStarted` removed | `Stop` before `Start` is recorded and honored. |

`AlphabetBase32` changed to Crockford base32 (32 symbols, no `L`). Values encoded with v0.0.3 do not decode. Re-encode persisted IDs or keep the old alphabet as a migration table.

## Silent behavior

These compile. They change what is stored or returned.

### HTTP errors

Unclassified errors no longer put `err.Error()` in `ErrorResponse.Message`. `code` is `0` unless the error implements `httpx.CodeError`. Enable `httpz.SetDebugMode(true)` only in development.

### Cache TTL

Across every driver: `expiration > 0` expires, `0` never expires **and clears any existing TTL**, `< 0` returns `cache.ErrInvalidTTL`.

v0.0.3 differences this removes:

- Redis `Set` used `KEEPTTL`
- mcache / badgerdb treated `0` as already expired
- mcache treated `-1` as never expire

`Close()` on a cache closes only what that constructor created. Wrappers (`CodecCache`, `NSCache`) never close the injected backend.

### Auth, CORS, downloads

- JWT claims without `uid` authenticate as 401, not as user `0`
- CORS rejects wildcard origins with credentials
- File downloads send `Content-Disposition: attachment` unless `WithInlineDownload()` is set
- Reverse-proxy cache no longer stores private / credentialed responses

### Boot and tasks

- `WithShutdownTimeout` default is 30s. Non-positive means unbounded, not an already-expired context.
- `Group.Stop` bounds member `Stop`, not only the caller's wait.
- `Stop` before `Start` does not launch members.
- Scheduler `Start` returns when the run context ends; call `Stop` to drain. Periodic jobs require a single replica.
- Redis `infra/redis.NewClient` no longer `Ping`s at construction. Probe in `boot.AddBeforeStart` if you want fail-fast startup.

### Storage keys

`storage.NormalizeKey` is applied by every driver. Persist the key `UploadFile` returns, not the argument you passed in. `DeleteFile` is idempotent.

## Logging

- Construct a backend (`zapx.NewBackend` or `log.NewStdioBackend`) and pass it to `InitWithBackends` / `WithLoggerBackend`.
- `WithStackAt` only attaches stacks. Use `WithMinLevel` to filter `StdioBackend`.
- `InitWithBackends` with no usable backend keeps the current logger.

## Templates

Official templates now depend on `sphere` v0.0.4 and already apply these call-site changes (`WithLoggerBackend`, `cors.NewCORS` error return, `NewRateLimiterByClientIP`). Existing projects created before that bump still need the table above.
