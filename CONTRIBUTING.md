# Contributing

Use Conventional Commits (`feat:`, `fix:`, `docs:`, `test:`, `chore:`). Pull
requests must pass `go test ./...`, `go vet ./...`, and `golangci-lint`.

Keep Kubernetes mutations explicit and test command contracts with fakes before
adding integration coverage. Do not include credentials, cluster dumps, or
provider payloads in logs or fixtures.
