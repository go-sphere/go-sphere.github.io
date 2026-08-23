---
title: Cache, Storage, Boot
weight: 46
---

Sphere ships small interfaces and default adapters for cache, object storage, messaging, scheduling, and process lifecycle. They are optional. Generated HTTP code does not depend on them.

## Boot and tasks

Everything long-running implements `core/task.Task` (`Identifier`, `Start`, `Stop`). `boot.Run` wraps a `task.Group` with OS signals, hooks, and a shutdown deadline.

```go
err := boot.Run(conf, func(c *Conf) (*boot.Application, error) {
    return boot.NewApplication(httpTask, consumerTask), nil
}, boot.WithLoggerBackend(backend))
```

- Default `WithShutdownTimeout` is 30s; non-positive means no limit.
- Members of a plain `NewApplication` stop concurrently. Use `NewStagedApplication` when one task's `Stop` tears down something another task still uses (last stage stops first).
- Close Wire-owned clients (`sql.DB`, Redis) in `AddAfterStop` or in the injector cleanup after `Run` returns.
- Constructors such as `infra/redis.NewClient` do not ping. Add `AddBeforeStart` if startup must fail when the backend is down.

## Cache

`cache.Cache[S]` is CRUD + batch + TTL + `Close`. Drivers: `memory` (ristretto), `mcache`, `redis`, `badgerdb`, `nocache`. `nscache` prefixes keys so several logical caches can share one backend.

TTL contract:

- `expiration > 0` — expires after that duration
- `expiration == 0` — never expires, and **clears** any existing TTL
- `expiration < 0` — `cache.ErrInvalidTTL`, no write

A constructor that opens a client/DB closes it; a constructor that receives one does not. `CodecCache` / `NSCache` `Close` is a no-op.

`DelAll` blast radius differs by driver (Redis `FLUSHDB` of the selected DB, memory/mcache the whole process cache, NSCache only its namespace). Do not share a Redis database with unrelated keys if you call `DelAll`.

## Storage

Drivers: `local`, `s3`, `qiniu`, `kvcache`. `fileserver` is an HTTP adapter with one-time upload tokens, not an S3 driver.

Every driver runs `storage.NormalizeKey` first. Persist the key `UploadFile` returns. `DeleteFile` succeeds when the object is already gone. Downloads are served as attachments by default.

## Scheduler and MQ

`scheduler/cron` and `scheduler/asynq` implement `task.Task`. `Start` blocks until the run context ends and does **not** drain; `Stop` is the cleanup half. Periodic jobs are not coordinated across processes — run one replica, or make handlers idempotent.

`mq` has memory and Redis drivers for queues and pubsub. `Close` waits for in-flight handlers.

## Related

- [Upgrading to v0.0.4](upgrading) — TTL, `Close` ownership, and boot timeout changes
- [Customizing the Stack](customizing-stack)
- [Logging](logging)
