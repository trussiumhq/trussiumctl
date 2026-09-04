# trussiumctl roadmap

The standalone binary is intentionally platform-focused. The Python runtime
continues to own request execution, provider adapters, and runtime-local
diagnostics.

## Priorities

1. **Foundation (current):** stable version/help contract, reproducible Go
   builds, tests, linting, vulnerability scanning, CodeQL, semantic tags, and
   signed release-ready archives.
2. **Read-only operations:** Kubernetes context discovery, runtime/operator
   status, Helm release inspection, compatibility checks, and bounded
   diagnostics with safe errors. The runtime health contract and local version
   compatibility preflight are delivered first (`runtime status` and
   `compatibility check`). The combined cluster report is now delivered as
   `diagnostics cluster`. Helm rendering is now available through the guarded
   `install --dry-run` command; mutating installs remain intentionally disabled.
   The equivalent `upgrade --dry-run` preflight now validates current and target
   versions before rendering the target chart.
3. **Guarded changes:** explicit install, upgrade, and rollback commands that
   reuse the published Helm chart and Operator artifacts, require confirmation
   for mutations, and support dry-run output.
4. **Operational safety:** preflight checks, timeout/deadline controls,
   structured command output, audit-friendly summaries, and rollback
   verification.

The CLI will not embed provider SDKs or duplicate runtime execution. Commands
will call Kubernetes/Helm APIs and the runtime's documented HTTP contracts.
