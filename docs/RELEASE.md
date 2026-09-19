# Release and verification guide

`trussiumctl` releases are created from `main` by the semantic-release
workflow. Contributors do not create version tags or upload archives by hand.

## Before merging a release-producing change

- Use a Conventional Commit (`feat:`, `fix:`, `docs:`, `test:`, or `chore:`).
- Confirm `go test ./...`, `go vet ./...`, and `golangci-lint` pass.
- Confirm CodeQL and `govulncheck` pass.
- Confirm mutating workflows retain explicit confirmation, server-side
  validation, timeouts, and post-operation verification.
- Update `README.md` and `docs/ROADMAP.md` when user-visible behavior changes.
- Review the generated release scope with `git log <last-tag>..main`.

The current published baseline is `v1.16.0`. Releases `v1.14.0` through
`v1.16.0` add
deployed-version discovery for diagnostics and upgrade preflights, fail-closed
baseline discovery for guarded upgrades, and end-to-end guarded upgrade and
rollback safety coverage.

## Automated release

On a push to `main`, `.github/workflows/semantic-release.yml` runs
`go-semantic-release` with GoReleaser. It creates the semantic tag and GitHub
release, builds Linux and macOS archives for amd64 and arm64, and publishes
`checksums.txt`. Documentation-only commits are excluded from archive
changelog entries by the GoReleaser configuration.

The workflow also generates a keyless GitHub artifact attestation for the
checksum file. A release is not considered complete until the release assets,
checksums, and attestation are all present.

## Post-release verification

Set `RELEASE_VERSION` to the published tag and verify the release from a clean
checkout:

```console
gh release view "$RELEASE_VERSION" --repo trussiumhq/trussiumctl
gh release download "$RELEASE_VERSION" --repo trussiumhq/trussiumctl --dir /tmp/trussiumctl-release
gh attestation verify /tmp/trussiumctl-release/checksums.txt --repo trussiumhq/trussiumctl
```

The token used by `gh attestation verify` must have GitHub's **Attestations:
read** permission. Without it, GitHub returns HTTP 404 for the attestation API
even when the public attestation record exists. The release workflow itself
uses `attestations: write` and reports the uploaded attestation URL in its
summary.

Verify that the release contains archives for `linux/amd64`, `linux/arm64`,
`darwin/amd64`, and `darwin/arm64`, then test the matching binary:

```console
./trussiumctl version
./trussiumctl --help
```

If an archive, checksum, or attestation is missing, treat the release as
incomplete and investigate the workflow before announcing it.
