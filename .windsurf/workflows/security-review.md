---
description: Perform a security-focused review of authentication code
---

# Security Review Workflow

Focused security audit for authentication-critical code paths.

## Prerequisites

// turbo

1. Verify the project builds: `go build -o gh-app-auth .`

## Phase 1: Secret Exposure Scan

2. Search all Go files for potential secret logging:
   - `grep -rn 'fmt.Print.*token' --include='*.go' | grep -v _test.go | grep -v HashToken`
   - `grep -rn 'fmt.Print.*password' --include='*.go' | grep -v _test.go`
   - `grep -rn 'fmt.Print.*[Kk]ey' --include='*.go' | grep -v _test.go | grep -v HashToken`
   - `grep -rn 'log.Print.*token' --include='*.go' | grep -v _test.go`

3. For any finding, apply transformation T-SEC-03 from `.codequality/transformations.yaml`.

## Phase 2: HTTP Client Security

4. Search for unprotected HTTP clients:
   - `grep -rn 'http.Client{' --include='*.go' | grep -v Timeout`
   - `grep -rn 'http.Get(' --include='*.go'`
   - `grep -rn 'http.Post(' --include='*.go'`

5. Verify all HTTP calls use `httpclient.Default()` or have explicit timeouts.

## Phase 3: File Permission Checks

6. Search for file operations and verify permission validation:
   - `grep -rn 'os.WriteFile' --include='*.go'`
   - `grep -rn 'os.MkdirAll' --include='*.go'`
   - Verify permissions are 0400/0600 for files, 0700 for directories.

## Phase 4: Input Validation

7. Check path traversal prevention in all user inputs:
   - `grep -rn '\.\.' --include='*.go' | grep -v _test.go | grep -v vendor`

8. Check URL validation for credential routing patterns.

## Phase 5: Token Lifecycle

9. Verify token storage is memory-only:
   - Search for any disk persistence of tokens
   - Verify cache TTL is 55 minutes or less
   - Confirm tokens are not included in config serialization

## Phase 6: SAST Tools

// turbo

10. Run gosec: `make security-gosec`

// turbo

11. Run govulncheck: `make security-vulncheck`

## Phase 7: Report

12. Produce findings categorized as:
    - **P0 Critical**: Must fix before release
    - **P1 High**: Should fix soon
    - **P2 Medium**: Improvement opportunity
    - **P3 Low**: Best practice suggestion
