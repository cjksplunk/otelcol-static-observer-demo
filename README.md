# otelcol-static-observer-demo

## What this proves

Operators who want to collect metrics from multiple instances of the same
service — and tag each instance's telemetry with distinct resource attributes
like `service.name` or `deployment.environment` — today have no clean solution.
Options require either one pipeline per instance, a custom processor, or changes
to the receiver itself.

This repo demonstrates a solution using only config: the new `static_observer`
extension paired with `receiver_creator`. Two `hostmetrics` instances run side
by side, each stamping completely different `service.name` and
`deployment.environment` resource attributes, with **zero changes to the
hostmetrics receiver**.

### How it works

1. **`static_observer`** — a new contrib extension
   ([`extension/observer/staticobserver`](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/47174))
   that fires a single synthetic endpoint of type `"static"` on startup.
   Unlike k8s or docker observers it performs no dynamic discovery.

2. **`receiver_creator`** — matches the synthetic endpoint via
   `rule: type == "static"` and instantiates each subreceiver template once.
   Each template has its own `resource_attributes` block; the `receiver_creator`
   stamps those key/value pairs onto all telemetry from that instance.

3. **`hostmetrics`** receiver — unchanged. It has no knowledge of the
   `resource_attributes`. The wrapping happens entirely in `receiver_creator`.

The result: a single config file can express per-instance identity for any
existing receiver, today, without touching any receiver code.

---

Implemented in
[cjksplunk/opentelemetry-collector-contrib@mysql-add-service-resource-attributes-clean](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/47174).

## What's wired up

Two `hostmetrics` subreceiver instances via `receiver_creator`:

| Instance | `service.name` | `deployment.environment` |
|---|---|---|
| `hostmetrics/prod` | `myapp-prod` | `production` |
| `hostmetrics/staging` | `myapp-staging` | `staging` |

Both scrape `cpu` and `memory` every 10 seconds and export to the `debug`
exporter (verbosity: detailed).

## Prerequisites

- Go 1.22+

## Run

The `static_observer` extension lives in a fork branch not yet in the public Go
checksum database, so `GONOSUMDB` is required.

```bash
# Fetch dependencies (first run only)
GONOSUMDB="github.com/cjksplunk/*" go mod tidy

# Run the collector
GONOSUMDB="github.com/cjksplunk/*" go run . --config config.yaml
```

Or via make:

```bash
make tidy   # first run only
make run
```

## Expected output

After startup you should see two sets of metrics in the debug output, each
with distinct resource attributes:

```
Resource SchemaURL:
Resource attributes:
     -> service.name: Str(myapp-prod)
     -> deployment.environment: Str(production)
...
Resource attributes:
     -> service.name: Str(myapp-staging)
     -> deployment.environment: Str(staging)
```

Both instances collect the same metrics from the same host — the only difference
is the resource attributes, set entirely in config.

## Dependency pinning

`go.mod` uses `replace` directives to pin all
`github.com/open-telemetry/opentelemetry-collector-contrib` packages to the
`mysql-add-service-resource-attributes-clean` branch of
[cjksplunk/opentelemetry-collector-contrib](https://github.com/cjksplunk/opentelemetry-collector-contrib).
All `go.opentelemetry.io/collector/*` core modules resolve from upstream.

Because the fork is not in the public Go checksum database, `GONOSUMDB` must be
set when downloading dependencies:

```bash
GONOSUMDB="github.com/cjksplunk/*" go mod tidy
```
