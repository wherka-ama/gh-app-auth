---
description: Perform a systematic code review using agentic resources
---

# Code Review Workflow

Systematic code review driven by `.codequality/` resources.

## Prerequisites

// turbo

1. Verify the project builds: `go build -o gh-app-auth .`

// turbo

2. Verify tests pass: `go test ./...`

## Phase 1: Load Context

3. Read `.codequality/agent_config.yaml` to understand project structure and quality targets.

4. Read `.codequality/architecture_rules.yaml` to understand constraints.

5. Read `.codequality/review_checklist.yaml` to know what to evaluate.

## Phase 2: Static Analysis

// turbo

6. Run full linting: `make lint`

// turbo

7. Run security scan: `make security-scan`

// turbo

8. Check cyclomatic complexity: `go run github.com/fzipp/gocyclo/cmd/gocyclo@latest -over 15 .`

## Phase 3: Architecture Verification

9. For each rule in `architecture_rules.yaml`, verify compliance:
   - Check package imports respect layer boundaries (ARCH-001, ARCH-002, ARCH-003)
   - Search for forbidden patterns (SEC-001 through SEC-004)
   - Check code quality patterns (QUAL-001 through QUAL-003)

## Phase 4: Checklist Evaluation

10. Walk through each category in `review_checklist.yaml`:
    - Security (P0/P1 items first)
    - Error Handling
    - Concurrency Safety
    - Resource Management
    - Performance
    - Testing
    - Code Quality
    - Dependencies
    - Documentation

## Phase 5: Fix Generation

11. For each finding, consult `transformations.yaml` for matching safe transformations.

12. Apply fixes one at a time, smallest first.

13. After each fix, validate:

// turbo

- `go build -o gh-app-auth .`

// turbo

- `go test ./...`

## Phase 6: Final Validation

// turbo

14. Run full CI: `make ci`

15. Commit working state with conventional commit message.
