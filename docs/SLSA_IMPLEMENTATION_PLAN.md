# SLSA L3 Immutable Releases Implementation Plan

## Plan Date: February 23, 2026

## Target: gh-app-auth SLSA L3 compliance with immutable releases

---

## Overview

This plan implements SLSA Level 3 compliant releases with artifact attestations and immutable releases for the gh-app-auth project. The implementation follows the release-please workflow for automated versioning and CHANGELOG generation.

---

## Phase 1: Foundation - Release Please Setup

### Iteration 1.1: Create Release Please Workflow

**Goal**: Implement automated release PR generation

**Changes Required**:

1. Create `.github/workflows/release-please.yml`
   - Trigger on push to main
   - Use `googleapis/release-please-action@v4`
   - Configure for Go strategy (simple approach with version.txt)

**Review Criteria**:

- [ ] Workflow triggers correctly on main branch pushes
- [ ] Release Please config exists at `.github/release-please-config.json`
- [ ] Manifest at `.github/release-please-manifest.json`
- [ ] Action has correct permissions (contents: write, pull-requests: write)

**Dependencies**: None

---

### Iteration 1.2: Prepare Repository for Release Please

**Goal**: Add versioning support and verify conventional commits

**Changes Required**:

1. Create `version.txt` at root with current version
2. Verify `.gitmessage` follows conventional commit format (already exists)
3. Update `.commitlintrc.yml` if needed

**Review Criteria**:

- [ ] version.txt exists with format `x.y.z`
- [ ] Commitlint config enforces conventional commits
- [ ] CHANGELOG.md will be auto-generated on first release

**Dependencies**: Iteration 1.1

---

### Iteration 1.3: Test Release Please (Dry Run)

**Goal**: Verify Release Please creates PRs correctly

**Steps**:

1. Merge workflow to main (via this feature branch)
2. Create a test PR with a `fix:` or `feat:` commit
3. Verify Release Please creates/updates release PR

**Review Criteria**:

- [ ] Release PR is created automatically
- [ ] PR contains changelog updates
- [ ] PR shows correct version bump

**Dependencies**: Iteration 1.2

---

## Phase 2: Artifact Attestations (SLSA L2)

### Iteration 2.1: Add Attestations to Existing Release Workflow

**Goal**: Generate build provenance for release artifacts

**Changes Required**:

1. Update `.github/workflows/release.yml`:
   - Add required permissions (id-token: write, attestations: write)
   - Add `actions/attest-build-provenance@v3` step after artifact creation
   - Attest all built binaries

**Review Criteria**:

- [ ] Workflow has required permissions
- [ ] Attestation step runs after precompile action
- [ ] Attestations are created for all platform binaries
- [ ] Attestations appear in workflow summary

**Dependencies**: Phase 1 complete

---

### Iteration 2.2: Verify SLSA L2 Compliance

**Goal**: Confirm L2 requirements are met

**Verification Steps**:

1. Run release workflow
2. Verify attestation is created and signed
3. Check attestation is stored in Sigstore format
4. Verify with `gh attestation verify` command

**Review Criteria**:

- [ ] Attestation is digitally signed
- [ ] Attestation links to build platform (GitHub-hosted runner)
- [ ] Provenance contains build info (commit SHA, workflow ref)
- [ ] Verification command succeeds

**Dependencies**: Iteration 2.1

---

## Phase 3: SLSA L3 Compliance

### Iteration 3.1: Create Reusable Attestation Workflow

**Goal**: Centralize attestation process for SLSA L3

**Changes Required**:

1. Create `.github/workflows/attest-release.yml`:
   - Workflow triggered via `workflow_call`
   - Inputs: artifact-path, artifact-digest
   - Uses `actions/attest-build-provenance@v3`
   - Runs on GitHub-hosted runner (ephemeral)

**Review Criteria**:

- [ ] Workflow is callable via workflow_call
- [ ] Has required permissions (id-token, attestations, contents: read)
- [ ] Separates signing from build process
- [ ] Runs on ephemeral GitHub-hosted runner

**Dependencies**: Phase 2 complete

---

### Iteration 3.2: Update Release Workflow to Use Reusable Workflow

**Goal**: Implement SLSA L3 separation of concerns

**Changes Required**:

1. Modify `.github/workflows/release.yml`:
   - Build binaries in one job
   - Call reusable workflow for attestation
   - Ensure signing happens separately from build

**Review Criteria**:

- [ ] Build job creates artifacts
- [ ] Separate attest job calls reusable workflow
- [ ] Attest job has no access to build secrets
- [ ] Workflow graph shows proper separation

**Dependencies**: Iteration 3.1

---

### Iteration 3.3: Verify SLSA L3 Compliance

**Goal**: Confirm L3 requirements are met

**Verification Steps**:

1. Verify reusable workflow was used
2. Check signing happens on ephemeral machine
3. Confirm no cross-job influence possible
4. Test with `gh artifact verify --signer-workflow`

**Review Criteria**:

- [ ] Attestation created by reusable workflow
- [ ] Signer workflow is isolated from build
- [ `gh artifact verify --signer-workflow` succeeds
- [ ] All L3 requirements documented

**Dependencies**: Iteration 3.2

---

## Phase 4: Immutable Releases

### Iteration 4.1: Enable Immutable Releases (Settings)

**Goal**: Protect releases from tampering

**Changes Required**:

1. Repository settings (manual step, documented):
   - Enable "Immutable releases" in repository settings
   - This is a repository-level setting, not a workflow change

**Review Criteria**:

- [ ] Setting enabled in repository
- [ ] New releases are marked immutable
- [ ] Tags protected from modification

**Dependencies**: Phase 3 complete

---

### Iteration 4.2: Update Release Workflow for Immutability

**Goal**: Ensure releases are created as immutable

**Changes Required**:

1. Update `.github/workflows/release.yml`:
   - Remove draft release approach (immutable releases should be direct)
   - OR use GitHub API to mark release as immutable
   - Ensure release is published correctly

**Review Criteria**:

- [ ] Releases created with immutability enabled
- [ ] Assets cannot be modified after publication
- [ ] Tag protection is enforced

**Dependencies**: Iteration 4.1

---

## Phase 5: Documentation and Verification

### Iteration 5.1: Document SLSA Compliance

**Goal**: Provide transparency to users

**Changes Required**:

1. Create `docs/SLSA_COMPLIANCE.md`:
   - Current SLSA level achieved
   - How to verify releases
   - Example verification commands

**Review Criteria**:

- [ ] Document explains current SLSA level
- [ ] Verification commands are tested
- [ ] Clear instructions for users

**Dependencies**: Phase 4 complete

---

### Iteration 5.2: Update README with Badges

**Goal**: Display SLSA compliance status

**Changes Required**:

1. Add SLSA level badge to README.md
2. Add "Verified Reproducible Builds" section
3. Link to SLSA_COMPLIANCE.md

**Review Criteria**:

- [ ] Badge shows current SLSA level
- [ ] Section explains verification process
- [ ] Links work correctly

**Dependencies**: Iteration 5.1

---

## Phase 6: Cleanup and Integration

### Iteration 6.1: Deprecate Controlled Release Workflow

**Goal**: Remove old controlled-release workflow if not needed

**Changes Required**:

1. Review `.github/workflows/controlled-release.yml`
2. Decide: integrate into new flow or remove
3. Document decision

**Review Criteria**:

- [ ] Decision documented
- [ ] Old workflow either integrated or removed
- [ ] No orphaned workflows

**Dependencies**: All previous phases

---

### Iteration 6.2: Final Integration Test

**Goal**: End-to-end test of complete workflow

**Steps**:

1. Create test feature branch with a `feat:` commit
2. Merge to main
3. Verify Release Please creates PR
4. Merge release PR
5. Verify release created with:
   - Correct version
   - Attestations attached
   - Immutable status
   - All binaries attested

**Review Criteria**:

- [ ] Full workflow completes without errors
- [ ] Release has attestations
- [ ] Release is immutable
- [ ] Verification commands work

**Dependencies**: All previous phases

---

## Implementation Schedule

| Phase | Iteration | Estimated Time | Status |
|-------|-----------|----------------|--------|
| 1 | 1.1 | 30 min | 🔲 |
| 1 | 1.2 | 20 min | 🔲 |
| 1 | 1.3 | 15 min (post-merge) | 🔲 |
| 2 | 2.1 | 45 min | 🔲 |
| 2 | 2.2 | 30 min | 🔲 |
| 3 | 3.1 | 40 min | 🔲 |
| 3 | 3.2 | 30 min | 🔲 |
| 3 | 3.3 | 30 min | 🔲 |
| 4 | 4.1 | 10 min (settings) | 🔲 |
| 4 | 4.2 | 30 min | 🔲 |
| 5 | 5.1 | 45 min | 🔲 |
| 5 | 5.2 | 20 min | 🔲 |
| 6 | 6.1 | 20 min | 🔲 |
| 6 | 6.2 | 30 min | 🔲 |

**Total Estimated Time**: ~5-6 hours of active development

---

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Release Please conflicts with existing tag workflow | Medium | High | Test on fork first, coordinate timing |
| Attestation fails for large binaries | Low | Medium | Use subject-digest approach |
| Immutable releases conflict with draft workflow | Medium | Medium | Remove draft approach, go direct |
| Permissions issues with reusable workflow | Low | High | Test permissions thoroughly |

---

## Rollback Plan

If issues occur:

1. Disable immutable releases in settings (preserves existing)
2. Revert to previous release.yml version
3. Disable Release Please workflow
4. Return to manual tag-based releases

---

## Success Criteria

- [ ] Release Please creates release PRs automatically
- [ ] Releases are SLSA L3 compliant
- [ ] Immutable releases enabled and working
- [ ] Users can verify releases with `gh attestation verify`
- [ ] Documentation is clear and complete
- [ ] No regression in existing functionality

---

*Plan created: 2026-02-23*
*Plan version: 1.0*
