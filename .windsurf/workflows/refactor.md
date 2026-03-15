---
description: Apply safe refactoring transformations from the codequality catalog
---

# Refactoring Workflow

Apply mechanical, safe refactoring transformations from `.codequality/transformations.yaml`.

## Prerequisites

// turbo

1. Verify clean starting state: `make ci`

2. Read `.codequality/transformations.yaml` to load available transformations.

3. Read `.codequality/agent_config.yaml` to understand quality targets.

## Phase 1: Identify Candidates

4. For each transformation in the catalog, search the codebase for matching triggers:
   - T-SEC-*: Security hardening (highest priority)
   - T-ERR-*: Error handling improvements
   - T-PERF-*: Performance optimizations
   - T-CONC-*: Concurrency safety
   - T-QUAL-*: Code quality
   - T-TEST-*: Testing improvements

## Phase 2: Apply Transformations (iterate)

5. Pick the highest-priority applicable transformation.

6. Apply the transformation to one location at a time.

7. After each change, validate:

// turbo

- `go build -o gh-app-auth .`

// turbo

- `go test ./...`

8. If validation passes, continue to next location or transformation.
   If validation fails, revert and skip.

## Phase 3: Batch Validation

// turbo

9. Run full CI: `make ci`

10. If CI passes, commit with appropriate conventional commit message:
    - `refactor(<scope>): <description>` for structural changes
    - `perf(<scope>): <description>` for performance improvements
    - `fix(<scope>): <description>` for bug fixes

## Phase 4: Report

11. Summarize applied transformations:
    - Number of transformations applied per category
    - Files modified
    - Any skipped transformations with reasons
