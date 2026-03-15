# Agentic Code Review Resources — Specification

## 1. Purpose

This specification defines a set of **machine-readable and agent-consumable resources** that
drive autonomous code review, refactoring, and quality enforcement in Go CLI projects.

It is derived from a 9-part specification (L1–L9) covering static analysis, architecture
enforcement, knowledge graphs, self-healing repositories, and refactoring transformation
libraries — adapted for **Go CLI extensions** (not microservices).

## 2. Applicability Filter

The source specs target cloud-native microservices. For a Go CLI tool like `gh-app-auth`,
the following categories **apply**:

| Category | Applicable | Rationale |
|----------|-----------|-----------|
| Static analysis rules | Yes | Universal Go quality |
| Security rules (SAST) | Yes | Critical — authentication project |
| Concurrency safety | Yes | Goroutines, channels, mutex usage |
| Performance patterns | Yes | String builders, preallocation |
| Error handling | Yes | CLI user experience |
| Architecture enforcement | Partial | cmd/pkg layering, not full hexagonal |
| Testing quality | Yes | Coverage, table-driven tests |
| Dependency governance | Yes | Supply chain security |
| Observability rules | Partial | Structured logging only (no metrics/tracing) |
| Cloud-native rules | No | Not a server (no healthz, graceful shutdown) |
| Database rules | No | No SQL/DB usage |
| Kubernetes rules | No | Not an operator |
| gRPC/HTTP server rules | No | CLI tool, not a service |

## 3. Architecture Model

### Layer Model (CLI-adapted)

```text
cmd/          → CLI entry points (Cobra commands)
 ↓
pkg/auth/     → Core authentication logic
pkg/config/   → Configuration management
pkg/jwt/      → Token generation
 ↓
pkg/secrets/  → Secure storage (keyring/filesystem)
pkg/cache/    → In-memory token caching
pkg/matcher/  → URL pattern matching
pkg/logger/   → Diagnostic logging
pkg/scope/    → Credential scope detection
```

### Architecture Rules

1. **Layer isolation**: `pkg/` packages must not import `cmd/`
2. **Domain purity**: `pkg/auth/`, `pkg/jwt/` must not import CLI framework (`cobra`)
3. **No cyclic dependencies**: Package graph must remain a DAG
4. **Interface locality**: Interfaces defined where consumed
5. **Config isolation**: Business logic must not read env vars directly

## 4. Agentic Resource Types

### 4.1 Architecture Rules (YAML)

Machine-readable constraints evaluated against the codebase.
File: `.codequality/architecture_rules.yaml`

### 4.2 Review Checklist (YAML)

Structured evaluation criteria organized by category and severity.
File: `.codequality/review_checklist.yaml`

### 4.3 Transformation Catalog (YAML)

Safe, mechanical refactoring patterns the agent can apply.
File: `.codequality/transformations.yaml`

### 4.4 Agent Configuration (YAML)

Agent behavior, quality targets, and project metadata.
File: `.codequality/agent_config.yaml`

### 4.5 Windsurf Workflows (Markdown)

Step-by-step procedures for common agent tasks.
Files: `.windsurf/workflows/*.md`

## 5. Agent Workflow

The autonomous review loop:

```text
1. Load agent_config.yaml → understand project context
2. Load architecture_rules.yaml → know constraints
3. Run static analysis (make quality) → collect findings
4. Load review_checklist.yaml → evaluate systematically
5. Load transformations.yaml → match findings to fixes
6. Apply fixes in small iterations
7. Validate after each change (make ci)
8. Commit working state only
```

## 6. Quality Targets

| Metric | Target | Enforcement |
|--------|--------|-------------|
| Test coverage (overall) | ≥70% | CI gate |
| Lint errors | 0 | CI gate |
| Security vulnerabilities | 0 critical | CI gate |
| Cyclomatic complexity | <15 per function | Advisory |
| Package size | <2000 LOC | Advisory |
| Dependency count | <30 direct | Advisory |
| Markdown lint errors | 0 | CI gate |

## 7. Safety Guardrails

Before applying any transformation:

- `go build` must succeed
- `go test ./...` must pass
- `make lint` must pass
- Coverage must not decrease
- Only commit fully working state
- Security-critical code changes require explicit review

## 8. Severity Model

| Level | Meaning | Action |
|-------|---------|--------|
| P0 | Security vulnerability | Fix immediately |
| P1 | Architecture violation | Fix in current iteration |
| P2 | Code correctness bug | Fix in current iteration |
| P3 | Performance issue | Fix if low-risk |
| P4 | Style/maintainability | Fix opportunistically |

## 9. Files Produced

```text
.codequality/
├── AGENTIC_RESOURCES_SPEC.md    ← This file
├── agent_config.yaml            ← Project metadata and targets
├── architecture_rules.yaml      ← Layer and dependency constraints
├── review_checklist.yaml        ← Structured evaluation criteria
└── transformations.yaml         ← Safe refactoring patterns

.windsurf/workflows/
├── code-review.md               ← Full code review workflow
├── security-review.md           ← Security-focused review
└── refactor.md                  ← Refactoring workflow
```
