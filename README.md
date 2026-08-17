<p align="center">
  <img src="assets/logo.png" alt="Civet Logo" width="120" />
</p>

# Civet

Civet is a fast, deterministic CI/CD pipeline analyzer and assistant. It evaluates pipeline configurations against industry best practices, scores health across key architectural dimensions, and provides actionable recommendations.

Civet supports both **GitLab CI** (`.gitlab-ci.yml`) and **GitHub Actions** (`.github/workflows/*.yml`). It operates fully offline with a built-in rules engine, and supports optional AI provider integrations (Anthropic Claude, OpenAI, Google Gemini, Ollama) for contextual architectural reviews.

---

## Features

- **Platform-Agnostic Rules Engine**: Over 30 built-in rules covering caching, DAG optimization, image pinning, security, secret exposure, retries, and timeouts.
- **Dimensional Scoring**: Scores pipelines on a 0 to 100 scale across 6 categories: Performance, Security, Reliability, Maintainability, Observability, and Compliance.
- **Multiple Output Formats**: Human-readable terminal output with score bars, structured JSON for AI coding agents and CI automation, and Markdown for pull request summaries.
- **Repository Scanning**: Recursively discovers and audits all CI/CD pipelines across repositories in one command.
- **Pipeline Scaffolding**: Detects project languages (Go, Node.js, Python, Rust, Docker) and generates hardened, best-practice pipeline templates.
- **Optional AI Enhancement**: Pluggable AI provider layer for deeper semantic analysis and risk assessment.

---

## Installation

### Via Go Install

```bash
go install github.com/AndreaPallotta/civet@latest
```

### From Source

```bash
git clone https://github.com/AndreaPallotta/civet.git
cd civet
go build -o civet.exe .
```

---

## Quickstart

### Audit a Pipeline

Run an audit on a single pipeline file:

```bash
civet audit .gitlab-ci.yml
civet audit .github/workflows/ci.yml
```

Output dense, token-optimized context tailored for LLM coding agents (Claude, Codex, Antigravity):

```bash
civet audit .gitlab-ci.yml --llm
civet scan . --llm
```

Output structured JSON or Markdown:

```bash
civet audit .gitlab-ci.yml --format json
civet audit .gitlab-ci.yml --format markdown
```

Fail in CI if the pipeline score is below a defined threshold:

```bash
civet audit .gitlab-ci.yml --fail-under 80
```

### Scan a Repository or Directory

Recursively locate and audit every pipeline:

```bash
civet scan .
civet scan /path/to/projects --fail-under 75
```

### Scaffold a New Pipeline

Generate a best-practice pipeline for your project:

```bash
civet new --platform gitlab --lang go
civet new --platform github --lang node --with docker
```

### Initialize Configuration

Generate a `.civet.yml` configuration file:

```bash
civet init
```

---

## Scoring Dimensions

| Category | Key Checks |
| :--- | :--- |
| **Performance** | Caching strategies, DAG optimization (`needs`), parallel execution, artifact retention |
| **Security** | Immutable image digests, secret exposure, runner permissions, approval gates |
| **Reliability** | Job timeouts, transient failure retries, concurrency controls, resource groups |
| **Maintainability** | Template includes, composite actions, DRY job configurations |
| **Observability** | Coverage reporting regexes, environment declarations, deployment tracking |
| **Compliance** | Deprecated syntax (`only`/`except`), broad triggers (`pull_request_target`) |

---

## AI Provider Configuration

Civet functions completely offline without any API keys. To enable optional AI analysis, configure `.civet.yml` or set environment variables:

```yaml
ai:
  provider: claude # Options: claude, openai, gemini, ollama
  model: claude-sonnet-4-20250514
  api_key_env: ANTHROPIC_API_KEY
```

Then run:

```bash
civet audit .gitlab-ci.yml --ai claude
```

---

## License

MIT License. Copyright (c) 2026 Andrea Pallotta.
