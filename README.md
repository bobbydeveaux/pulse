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

## License

MIT
