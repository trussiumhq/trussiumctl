# trussiumctl

`trussiumctl` is the standalone, public Go CLI for Trussium platform
operations. It will own Kubernetes and Helm workflows while the Python
`trussium` CLI continues to own runtime-local commands.

## Current status

The repository starts with a versioned command boundary. The next milestones
add read-only cluster inspection, diagnostics, and then guarded install,
upgrade, and rollback workflows against the published Helm chart and Operator.

```console
$ trussiumctl version
dev
$ trussiumctl runtime status --url http://127.0.0.1:9000
{"status":"ready"}
$ trussiumctl operator status --namespace trussium-system
$ trussiumctl helm status --namespace trussium-system --release trussium
$ trussiumctl compatibility check --runtime 1.22.0 --chart 1.3.0 --operator 1.0.2
$ trussiumctl diagnostics cluster --namespace trussium-system --runtime-version 1.22.0 --chart-version 1.3.0 --operator-version 1.0.2
$ trussiumctl install --dry-run --namespace trussium-system --chart trussium/trussium
$ trussiumctl install --dry-run --server-dry-run --namespace trussium-system --chart trussium/trussium
$ trussiumctl upgrade --dry-run --namespace trussium-system --current-runtime 1.22.0 --current-chart 1.3.0 --current-operator 1.0.2 --target-runtime 1.23.0 --target-chart 1.3.0 --target-operator 1.0.2
$ trussiumctl rollback --dry-run --namespace trussium-system --target-runtime 1.22.0 --target-chart 1.3.0 --target-operator 1.0.2
```

The CLI does not embed runtime execution logic, provider SDKs, or credentials.
Inspection commands invoke only read-only `kubectl get` and `helm status`
operations and return bounded JSON suitable for automation.
Compatibility checks are local and read-only; they fail closed when versions
are missing, malformed, or below the current supported baseline.
`diagnostics cluster` composes the four read-only checks and preserves bounded
section results when one dependency is unavailable.
`install` currently requires `--dry-run` and invokes only `helm template`; it
cannot change a cluster.
With `--server-dry-run`, the rendered manifest is additionally sent to
`kubectl apply --dry-run=server`; the API server validates it without storing
resources.
`upgrade` also requires `--dry-run`, validates both version sets, and renders
the target chart without invoking `helm upgrade`.
`rollback` requires `--dry-run`, validates the target release, and renders the
rollback chart without invoking `helm rollback`.
Rollback reports include an explicit verification result; dry-run verification
is always marked as not performed. Future mutating commands will require the
exact confirmation token `TRUSSIUM`.

## Development

Requires Go 1.25 or newer.

```console
go test ./...
go vet ./...
go run ./cmd/trussiumctl version
```

Releases use Conventional Commits and semantic-release tags. Release assets
are built for Linux and macOS on amd64 and arm64; package publication is
deferred until a registry credential is intentionally configured.

## License

Apache-2.0. See [LICENSE](LICENSE).
