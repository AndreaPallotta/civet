# Changelog

All notable changes to Civet will be documented in this file.

The format is based on Keep a Changelog, and this project adheres to Semantic Versioning.

## [0.1.1] - 2026-08-17

### Security
- Pin `softprops/action-gh-release` to immutable full commit SHA in release workflow.

### Documentation
- Add Civet project logo and badges to README.

## [0.1.0] - 2026-08-16

### Added
- Initial release of Civet CI/CD pipeline analyzer and assistant.
- Dual platform support for GitLab CI (`.gitlab-ci.yml`) and GitHub Actions (`.github/workflows/*.yml`).
- Core rules engine supporting 30 built-in rules across Universal, GitLab, and GitHub Actions domains.
- Scoring system evaluating pipelines across 6 dimensions: Performance, Security, Reliability, Maintainability, Observability, and Compliance.
- Multiple output formats: ANSI colored terminal output, structured JSON, GitHub Flavored Markdown, and token-dense LLM context (`--llm`).
- Dedicated `--llm` flag generating structured pipeline diagnostics, health grades, and actionable prompt directives for AI coding agents (Claude, Codex, Antigravity).
- Directory and repository recursive scanner (`civet scan`).
- Interactive and flagged pipeline scaffolding engine (`civet new`) with embedded best-practice templates for Go, Node.js, Python, Rust, and Docker.
- Pluggable AI provider integration supporting Anthropic Claude, OpenAI GPT, Google Gemini, and local Ollama.
- CI gate enforcement via `--fail-under` flag.
- Configuration management with `.civet.yml` supporting rule disabling and severity overrides.
