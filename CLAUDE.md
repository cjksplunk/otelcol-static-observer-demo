# CLAUDE.md — otelcol-static-observer-demo

## Repo context

Single-contributor demo repo. Commit directly to `main` — no branches or PRs needed.

## Architecture

- `static_observer` extension fires a single synthetic endpoint on startup
- `receiver_creator` matches it (`type == "static"`) and starts two `hostmetrics` subreceivers
- Each subreceiver gets distinct `resource_attributes` (`service.name`, `deployment.environment`)
- Output goes to the `debug` exporter (verbosity: detailed)

The fork modules (`github.com/cjksplunk/opentelemetry-collector-contrib/...`) are pinned to tagged
releases on the cjksplunk fork. The fork's tags have been force-pushed in the past, causing
`go.sum` hash mismatches.

## Known environment issues

### TMPDIR misconfiguration
The local shell has `TMPDIR` set to the project directory. Go treats that path as a system temp
root and ignores `go.mod` there. The Makefile sets `TMPDIR=/tmp` explicitly to guard against this.

### proxy.golang.org stale cache
`proxy.golang.org` has cached stale zips for the cjksplunk fork modules. Using the proxy causes
checksum mismatches even when `go.sum` is correct. The Makefile sets `GOPROXY=off,direct` to
bypass the proxy entirely.

### Port 8888 conflict
The otelcol Prometheus metrics endpoint binds to `localhost:8888`. If a previous run didn't exit
cleanly, `make run` will fail with `bind: address already in use`. Kill the stale process first:

```bash
lsof -ti tcp:8888 | xargs kill -9
```

## Running

```bash
make run
```

All required env vars (`GOPROXY`, `GONOSUMDB`, `TMPDIR`) are set by the Makefile. No manual setup needed.
