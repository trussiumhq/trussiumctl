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
```

The CLI does not embed runtime execution logic, provider SDKs, or credentials.
Inspection commands invoke only read-only `kubectl get` and `helm status`
operations and return bounded JSON suitable for automation.

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
