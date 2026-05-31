# Pulse

**Take the pulse of your codebase.**

Pulse is an open-source, multi-language code quality analyzer that tracks complexity, duplication, and maintainability — with time-series trends across your git history.

## Features

| Feature | Description |
|---------|-------------|
| Cyclomatic Complexity | Independent execution paths per function |
| Cognitive Complexity | Human readability score per function |
| Maintainability Index | Composite A-F grade per file |
| Duplication Detection | Copy-paste clone finder across your codebase |
| Time-Series Trends | Track metrics across git history |
| PR-Level Diffs | "This PR changed complexity by +X" |
| COCOMO Estimation | Person-months and cost estimates |
| Quality Gates | Fail CI when thresholds are exceeded |

## Quick Start

```bash
# Install
go install github.com/bobbydeveaux/pulse/app/cmd/pulse@latest

# Analyse your codebase
pulse check ./src

# See trends over your last 50 commits
pulse trend --last 50

# Check complexity diff for staged changes
pulse diff --staged

# Set quality gates
pulse gate --max-ccn 15 --max-duplication 5
```

## GitHub Action

```yaml
name: Code Quality
on: [pull_request]

jobs:
  pulse:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6
        with:
          fetch-depth: 0
      - uses: bobbydeveaux/pulse@main
        with:
          comment: true
          max_ccn: 15
          max_cognitive: 20
          max_duplication: 5
```

## Supported Languages

Go, TypeScript, JavaScript, Python, Java, Rust, C/C++, Ruby, PHP, Kotlin, Swift, and more.

## Part of the Newlife Ecosystem

Pulse works alongside [Guardian](https://guardian.stackramp.io) for security scanning. Guardian blocks insecure code; Pulse tracks code quality.

## GitHub Marketplace billing webhook

Pulse ships a stdlib-only HTTP server that receives `marketplace_purchase` events when a customer upgrades, downgrades, or cancels Pulse on the GitHub Marketplace. The binary lives at `app/cmd/webhook/` and is packaged as a distroless image via `Dockerfile.webhook` for Cloud Run (or any container runtime).

### Routes

| Method | Path | Purpose |
|--------|------|---------|
| `POST` | `/webhook/github-marketplace` | Verify `X-Hub-Signature-256` and log the event |
| `GET`  | `/health` | Liveness probe — returns `200 {"status":"ok"}` |

### Environment

| Variable | Required | Default | Purpose |
|----------|----------|---------|---------|
| `GITHUB_WEBHOOK_SECRET` | yes | — | HMAC-SHA256 shared secret configured in the Marketplace listing. The handler returns `500` on every event until this is set. |
| `PORT` | no | `8080` | Listen port (Cloud Run convention). |

### Accepted `marketplace_purchase` actions

`purchased`, `cancelled`, `changed`, `pending_change`, and `pending_change_cancelled`. Any other action returns `422` so misconfigurations surface loudly. The body is capped at 1 MiB; signature failures return `400`; malformed JSON returns `400`. GitHub `ping` deliveries are acknowledged with `200` so the handshake completes cleanly.

### Build and run locally

```bash
docker build -f Dockerfile.webhook -t pulse-webhook .

docker run \
  -e GITHUB_WEBHOOK_SECRET=replace-me \
  -e PORT=8080 \
  -p 8080:8080 \
  pulse-webhook
```

## License

MIT
